desc "Manage sifnode deploy, upgrade, etc processes"
namespace :standalone do

  namespace :deploy do
    desc "Deploy a single standalone sifnode on to your cluster"
    task :standalone, [:cluster, :chainnet, :provider, :namespace, :image, :image_tag, :moniker, :mnemonic, :admin_clp_addresses, :admin_oracle_address, :minimum_gas_prices] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade sifnode #{cwd}/../../deploy/helm/sifnode \
        --set sifnode.env.chainnet=#{args[:chainnet]} \
        --set sifnode.env.moniker=#{args[:moniker]} \
        --set sifnode.args.mnemonic=#{args[:mnemonic]} \
        --set sifnode.args.adminCLPAddresses=#{args[:admin_clp_addresses]} \
        --set sifnode.args.adminOracleAddress=#{args[:admin_oracle_address]} \
        --set sifnode.args.minimumGasPrices=#{args[:minimum_gas_prices]} \
        --install -n #{ns(args)} --create-namespace \
        --set image.tag=#{image_tag(args)} \
        --set image.repository=#{image_repository(args)}
      }

      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end

    desc "Deploy a single network-aware sifnode on to your cluster"
    task :peer, [:cluster, :chainnet, :provider, :namespace, :image, :image_tag, :moniker, :mnemonic, :peer_address, :genesis_url] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade sifnode #{cwd}/../../deploy/helm/sifnode \
        --install -n #{ns(args)} --create-namespace \
        --set sifnode.env.chainnet=#{args[:chainnet]} \
        --set sifnode.env.moniker=#{args[:moniker]} \
        --set sifnode.args.mnemonic=#{args[:mnemonic]} \
        --set sifnode.args.peerAddress=#{args[:peer_address]} \
        --set sifnode.args.genesisURL=#{args[:genesis_url]} \
        --set image.tag=#{image_tag(args)} \
        --set image.repository=#{image_repository(args)}
      }

      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end

    desc "Deploy the sifnode API to your cluster"
    task :api, [:cluster, :chainnet, :provider, :namespace, :image, :image_tag, :node_host] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade sifnode-api #{cwd}/../../deploy/helm/sifnode-api \
        --install -n #{ns(args)} --create-namespace \
        --set sifnodeApi.args.chainnet=#{args[:chainnet]} \
        --set sifnodeApi.args.nodeHost=#{args[:node_host]} \
        --set image.tag=#{image_tag(args)} \
        --set image.repository=#{image_repository(args)}
      }

      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end

  desc "Sifnode with Vault"
  namespace :sifnode_vault do
    desc "Deploy a single standalone sifnode on to your cluster"
    task :standalone_vault, [:namespace, :image, :image_tag, :helm_values_file] do |t, args|
      #variable_template_replace(args[:template_file_name], args[:final_file_name])
      cmd = %Q{helm upgrade sifnode deploy/helm/sifnode-vault \
        --install -n #{args[:namespace]} --create-namespace \
        --set image.tag=#{args[:image_tag]} \
        --set image.repository=#{args[:image]} \
        -f #{args[:helm_values_file]} --kubeconfig=./kubeconfig
      }
      system(cmd) or exit 1
    end

    # desc "Deploy a single network-aware sifnode on to your cluster"
    # task :peer_vault, [:namespace, :image, :image_tag, :peer_address, :template_file_name, :final_file_name] do |t, args|
    #   variable_template_replace(args[:template_file_name], args[:final_file_name])
    #   cmd = %Q{helm upgrade sifnode deploy/helm/sifnode-vault \
    #     --install -n #{args[:namespace]} --create-namespace \
    #     --set sifnode.args.peerAddress=#{args[:peer_address]} \
    #     --set image.tag=#{args[:image_tag]} \
    #     --set image.repository=#{args[:image]} \
    #     -f #{args[:final_file_name]} --kubeconfig=./kubeconfig
    #   }
    #   system(cmd) or exit 1
    #   #:namespace, :image, :image_tag, :peer_address, :template_file_name, :final_file_name,:app_region, :app_env
    #   #rake "sifnode:standalone:peer_vault['sifnode', 'sifchain/sifnoded', 'testnet-genesis', '1b02f2eb065031426d37186efff75df268bb9097@54.164.57.141:26656', './deploy/helm/sifnode-vault/template.values.yaml', './deploy/helm/sifnode-vault/generated.values.yaml']"
    # end

    desc "Deploy a new sifnode vault to a new cluster"
    task :standalone, [:namespace, :image, :image_tag, :helm_values_file] do |t, args|
        cmd = %Q{helm upgrade sifnode deploy/helm/sifnode-vault \
          --install -n #{args[:namespace]} --create-namespace \
          --set image.tag=#{args[:image_tag]} \
          --set image.repository=#{args[:image]} \
          -f #{args[:helm_values_file]} --kubeconfig=./kubeconfig
        }
        system(cmd) or exit 1
    end

    desc "Deploy the sifnode API to your cluster"
    task :sifnode_api, [:chainnet, :namespace, :image, :image_tag, :node_host] do |t, args|
        cmd = %Q{helm upgrade sifnode-api deploy/helm/sifnode-api \
          --install -n #{args[:namespace]} --create-namespace \
          --set sifnodeApi.args.chainnet=#{args[:chainnet]} \
          --set sifnodeApi.args.nodeHost=#{args[:node_host]} \
          --set image.tag=#{args[:image_tag]} \
          --set image.repository=#{args[:image]} --kubeconfig=./kubeconfig
        }
        system(cmd) or exit 1
      end

    desc "Deploy a single network-aware sifnode on to your cluster"
    task :peer_vault, [:namespace, :image, :image_tag, :helm_values_file, :peer_address] do |t, args|
        cmd = %Q{helm upgrade sifnode deploy/helm/sifnode-vault \
        --install -n #{args[:namespace]} --create-namespace \
        --set sifnode.args.peerAddress=#{args[:peer_address]} \
        --set image.tag=#{args[:image_tag]} \
        --set image.repository=#{args[:image]} \
        -f #{args[:helm_values_file]} --kubeconfig=./kubeconfig
        }
        system(cmd) or exit 1
    end
  end
end
