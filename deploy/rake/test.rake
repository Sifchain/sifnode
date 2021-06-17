require "./deploy/rake/lib/sifchain/chainops/builder"
require "./deploy/rake/lib/sifchain/chainops/cli"
require "./deploy/rake/lib/sifchain/chainops/task"
require "./deploy/rake/lib/sifchain/chainops/test/testing"

namespace :test do
  desc "Test"
  task :testing, %i[arg1 arg2 arg3] do |t, args|
    puts ::Sifchain::Chainops::Task.new(task: t, args: args).build
  end
end
