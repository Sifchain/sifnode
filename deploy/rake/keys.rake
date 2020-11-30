desc "key operations"
namespace :keys do
  namespace :generate do
    desc "Generate a new mnemonic phrase"
    task :mnemonic do
      mnemonic = `sifgen key generate`
      puts mnemonic
    end
  end
end
