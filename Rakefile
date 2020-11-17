require "rake"
require "yaml"
require "digest"

Dir["**/lib/*.rb"].each { |file| require_relative file }
Dir["**/*.rake"].each do |path|
  Rake.add_rakelib(path&.split("/")&.reverse&.drop(1)&.reverse&.join("/"))
end
