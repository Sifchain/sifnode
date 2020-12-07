desc "Faucet operations"
namespace :faucet do
  desc "Validator operations"
  namespace :validator do
    desc "Send funds"
    task :send, [:chainnet, :from_address, :to_address, :amount, :node] do |t, args|
      cmd = %Q{sifnodecli tx send #{args[:from_address]} #{args[:to_address]} #{args[:amount]} \
              --node #{node_address(args)} \
              --chain-id #{args[:chainnet]} \
              --keyring-backend file -y
      }

      system(cmd)
    end
  end
end
