desc "validator operations"
namespace :validator do
  desc "Stake a node so it can participate in consensus"
  task :stake, [:chainnet, :moniker, :amount, :gas, :pub_key, :node] do |t, args|
    node = if args[:node].nil?
             "tcp://127.0.0.1:26657"
           else
             args[:node]
           end

    cmd = %Q{sifnodecli tx staking create-validator \
            --commission-max-change-rate="0.1" \
            --commission-max-rate="0.1" \
            --commission-rate="0.1" \
            --amount=#{args[:amount]} \
            --pubkey=#{args[:pub_key]} \
            --chain-id=#{args[:chainnet]} \
            --min-self-delegation="1" \
            --gas-prices=#{args[:gas]} \
            --moniker=#{args[:moniker]} \
            --from=#{args[:moniker]} \
            --keyring-backend=file \
            --node=#{node}
    }

    system(cmd)
  end

  desc "Key operations"
  namespace :keys do
    desc "Print the validator public key"
    task :public, [:cluster, :provider, :namespace] do |t, args|
      pod_name = pod_name(args)
      if pod_name.nil?
        puts "Unable to find any pods!"
        exit(1)
      end

      cmd = %Q{kubectl exec --stdin --tty #{pod_name} -n #{args[:namespace]} -- cosmovisor tendermint show-validator}
      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end

  desc "Backup operations"
  namespace :backup do
    desc "Backup the validator config"
    task :config, [:cluster, :provider, :namespace, :save_path] do |t, args|
      pod_name = pod_name(args)
      if pod_name.nil?
        puts "Unable to find any pods!"
        exit(1)
      end

      cmd = %Q{kubectl cp #{args[:namespace]}/#{pod_name}:/root/.sifnoded/config #{args[:save_path]}}
      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end
end
