require "pikeman/version"
require "cli/kit"
require "json"

module Pikeman
  Error = Struct.new(:filename, :line, :column, :text, :link, :confidence, :linetext, :category) do
    def initialize(filename:, line:, column:, text:, link:, confidence:, linetext:, category:)
      super(filename, line, column, text, link, confidence, linetext, category)
    end
  end

  def self.lint_file(absolute_path, _config_absolute_path = nil)
    out_and_err, _stat = CLI::Kit::System.capture2e(binary, "-format", "json", absolute_path)
    out_and_err.each_line.map do |line|
      data = JSON.parse(line)
      Error.new(
        filename: data.fetch("filename"),
        line: data.fetch("line"),
        column: data.fetch("column"),
        text: data.fetch("text"),
        link: data.fetch("link"),
        confidence: data.fetch("confidence"),
        linetext: data.fetch("linetext"),
        category: data.fetch("category")
      )
    end
  end

  private

  def self.binary
    File.join(File.dirname(__dir__), "bin", "pikeman")
  end
end
