desc "key operations"
namespace :keys do
  desc "Generate"
  namespace :generate do
    desc "Generate a new mnemonic phrase"
    task :mnemonic do
      mnemonic = `sifgen key generate`
      puts mnemonic
    end
  end

  desc "Import (recover) a key, using the mnemonic"
  task :import, [:moniker] do |t, args|
    system("sifnoded keys add #{args[:moniker]} --recover --keyring-backend file")
  end
end
