desc "Faucet operations"
namespace :faucet do
  desc "Validator operations"
  namespace :validator do
    desc "Send funds"
    task :send, [:from_address, :to_address, :amount, :node_address] do |t, args|
      cmd = %Q{sifnodecli tx send #{args[:from_address]} #{args[:to_address]} #{args[:amount]} \
              --node #{node_address(args)} \
              --keyring-backend file -y
      }

      system(cmd)
    end
  end
end
