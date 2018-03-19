// Copyright (c) 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// golint lints the Go source files named on its command line.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/maximebedard/pikeman"
)

var (
	minConfidence = flag.Float64("min_confidence", 0.8, "minimum confidence of a problem to print it")
	setExitStatus = flag.Bool("set_exit_status", false, "set exit status to 1 if any issues are found")
	formatterType = flag.String("format", "text", "set the format. Available: text, json.")
	configPath    = flag.String("config_path", "", "set the configuration file.")
	suggestions   int
	formatter     problemFormatter
	config        *lint.Config
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tgolint [flags] # runs on package in current directory\n")
	fmt.Fprintf(os.Stderr, "\tgolint [flags] [packages]\n")
	fmt.Fprintf(os.Stderr, "\tgolint [flags] [directories] # where a '/...' suffix includes all sub-directories\n")
	fmt.Fprintf(os.Stderr, "\tgolint [flags] [files] # all must belong to a single package\n")
	fmt.Fprintf(os.Stderr, "Version: %s\n", VERSION)
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	switch *formatterType {
	case "json":
		formatter = &jsonFormatter{}
	default:
		formatter = &textFormatter{}
	}

	var err error
	if *configPath != "" {
		config, err = lint.ReadConfig(*configPath)
	} else {
		config, err = lint.ReadConfigFromWorkingDir()
	}

	if err != nil {
		writeError(1, "unable to read configuration file", err.Error())
	}

	if flag.NArg() == 0 {
		lintDir(".")
	} else {
		// dirsRun, filesRun, and pkgsRun indicate whether golint is applied to
		// directory, file or package targets. The distinction affects which
		// checks are run. It is no valid to mix target types.
		var dirsRun, filesRun, pkgsRun int
		var args []string
		for _, arg := range flag.Args() {
			if strings.HasSuffix(arg, "/...") && isDir(arg[:len(arg)-len("/...")]) {
				dirsRun = 1
				for _, dirname := range allPackagesInFS(arg) {
					args = append(args, dirname)
				}
			} else if isDir(arg) {
				dirsRun = 1
				args = append(args, arg)
			} else if exists(arg) {
				filesRun = 1
				args = append(args, arg)
			} else {
				pkgsRun = 1
				args = append(args, arg)
			}
		}

		if dirsRun+filesRun+pkgsRun != 1 {
			usage()
			os.Exit(2)
		}
		switch {
		case dirsRun == 1:
			for _, dir := range args {
				lintDir(dir)
			}
		case filesRun == 1:
			lintFiles(args...)
		case pkgsRun == 1:
			for _, pkg := range importPaths(args) {
				lintPackage(pkg)
			}
		}
	}

	if *setExitStatus && suggestions > 0 {
		writeError(1, "Found %d lint suggestions; failing.\n", suggestions)
	}
}

func writeError(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	os.Exit(1)
}

type textFormatter struct{}

func (tf *textFormatter) Write(ps []lint.Problem) {
	for _, p := range ps {
		fmt.Printf("%v: %s\n", p.Position, p.Text)
	}
}

type problemFormatter interface {
	Write([]lint.Problem)
}

type jsonFormatter struct{}

type jsonProblem struct {
	Filename   string  `json:"filename"`
	Line       int     `json:"line"`
	Column     int     `json:"column"`
	Text       string  `json:"text"`
	Link       string  `json:"link"`
	Confidence float64 `json:"confidence"`
	LineText   string  `json:"linetext"`
	Category   string  `json:"category"`
}

func (jf *jsonFormatter) Write(ps []lint.Problem) {
	problems := make([]jsonProblem, len(ps))

	for i, p := range ps {
		problems[i] = jsonProblem{
			Filename:   p.Position.Filename,
			Line:       p.Position.Line,
			Column:     p.Position.Column,
			Text:       p.Text,
			Link:       p.Link,
			Confidence: p.Confidence,
			LineText:   p.LineText,
			Category:   p.Category,
		}
	}

	b, _ := json.Marshal(problems)
	os.Stdout.Write(b)
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func lintFiles(filenames ...string) {
	files := make(map[string][]byte)
	for _, filename := range filenames {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		files[filename] = src
	}

	l := &lint.Linter{Config: config}
	ps, err := l.LintFiles(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	// TODO: maybe write + calculate suggestions in the same loop, but for now who cares.
	for _, p := range ps {
		if p.Confidence >= *minConfidence {
			suggestions++
		}
	}

	formatter.Write(ps)
}

func lintDir(dirname string) {
	pkg, err := build.ImportDir(dirname, 0)
	lintImportedPackage(pkg, err)
}

func lintPackage(pkgname string) {
	pkg, err := build.Import(pkgname, ".", 0)
	lintImportedPackage(pkg, err)
}

func lintImportedPackage(pkg *build.Package, err error) {
	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			// Don't complain if the failure is due to no Go source files.
			return
		}
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.TestGoFiles...)
	if pkg.Dir != "." {
		for i, f := range files {
			files[i] = filepath.Join(pkg.Dir, f)
		}
	}
	// TODO(dsymonds): Do foo_test too (pkg.XTestGoFiles)

	lintFiles(files...)
}
