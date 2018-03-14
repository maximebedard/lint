lib = File.expand_path("../lib", __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require "pikeman/version"

Gem::Specification.new do |spec|
  spec.name          = "pikeman"
  spec.version       = Pikeman::VERSION
  spec.authors       = ["Maxime Bedard"]
  spec.email         = ["maxime.bedard@shopify.com"]

  spec.summary       = %q{Very opinionated version of google/golint}
  spec.description   = %q{Very opinionated version of google/golint}
  spec.homepage      = "https://github.com/maximebedard/pikeman"
  spec.license       = "MIT"

  spec.files         = `git ls-files -z`.split("\x0").reject do |f|
    f.match(%r{^(test|spec|features)/})
  end

  spec.executables   = ["pikeman"]
  spec.require_paths = ["lib"]

  spec.add_dependency "cli-kit"
  spec.add_development_dependency "bundler", "~> 1.16"
  spec.add_development_dependency "rake", "~> 10.0"
  spec.add_development_dependency "minitest", "~> 5.0"
  spec.add_development_dependency "pry-byebug"
end
