require "test_helper"

class PikemanTest < Minitest::Test
  def test_lint_file
    path = "./test/fixtures/shittymain.go"
    expected = Pikeman::Error.new(
      filename: "./test/fixtures/shittymain.go",
      line: 6,
      column: 6,
      text: "exported type Patate should have comment or be unexported",
      link: "https://golang.org/wiki/CodeReviewComments#doc-comments",
      confidence: 1,
      linetext: "type Patate struct{}\n",
      category: "comments"
    )
    actual, * = Pikeman.lint_file(path)

    assert_equal(expected, actual)
  end
end
