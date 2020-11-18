require "rake"
require "yaml"
require "digest"
require "rbconfig"

Dir["**/*.rake"].each do |path|
  Rake.add_rakelib(path&.split("/")&.reverse&.drop(1)&.reverse&.join("/"))
end
