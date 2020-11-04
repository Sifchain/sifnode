desc "faucet operations"
namespace :faucet do
  desc "validator operations"
  namespace :validator do
    desc "send funds to a validator"
    task :send, [:chainnet, :address] do |t, args|
      config = YAML.load_file(network_config(args[:chainnet]))

      cmd = %Q{printf "#{config[0]['password']}\n#{config[0]['password']}\n" | \
               sifnodecli tx send #{config[0]['address']} #{args[:address]} 10000000trwn,100000000000clink,100000000000chot,100000000000cusdt,100000000000cusdc \
               --home networks/validators/#{args[:chainnet]}/#{config[0]['moniker']}/.sifnodecli -y
             }

      system(cmd)
    end
  end
end
