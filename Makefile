NAME=pikeman
GO_PACKAGE=github.com/maximebedard/pikeman
RUBY_MODULE=Pikeman
VERSION=$(shell cat VERSION)
GOFILES=$(shell find . -type f -name '*.go')
GEM=$(NAME)-$(VERSION).gem

.PHONY: all

default: release
release: $(GEM)
test: gotest rbtest
binaries: build/linux-amd64/pikeman build/darwin-amd64/pikeman

build/linux-amd64/pikeman: $(GOFILES) golint/version.go
	GOOS=linux GOARCH=amd64 go build -o "$@" ./golint

build/darwin-amd64/pikeman: $(GOFILES) golint/version.go
	GOOS=darwin GOARCH=amd64 go build -o "$@" ./golint

gotest:
	go test -race -v ./...

rbtest:
	bundle exec rake test

$(GEM): \
	lib/$(NAME)/version.rb \
	build/linux-amd64/pikeman \
	build/darwin-amd64/pikeman
	gem build pikeman.gemspec && gem push $@ && rm $@

golint/version.go: VERSION
	mkdir -p $(@D)
	echo 'package main\n\nconst VERSION string = "$(VERSION)"' > $@

lib/$(NAME)/version.rb: VERSION
	mkdir -p $(@D)
	echo 'module $(RUBY_MODULE)\n  VERSION = "$(VERSION)"\nend' > $@
