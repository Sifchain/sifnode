require "securerandom"

desc "management processes for the kube cluster and terraform commands"
namespace :cluster do
  desc "Scaffold new cluster environment configuration"
  task :scaffold, [:cluster, :provider] do |t, args|
    check_args(args)

    # create path location
    system("mkdir -p #{cwd}/../../.live")
    system("mkdir #{path(args)}") or exit

    # create config from template
    system("go run github.com/belitre/gotpl #{cwd}/../terraform/template/aws/cluster.tf.tpl \
      --set chainnet=#{args[:cluster]} > #{path(args)}/main.tf")

    system("go run github.com/belitre/gotpl #{cwd}/../terraform/template/aws/.envrc.tpl \
      --set chainnet=#{args[:cluster]} > #{path(args)}/.envrc")

    # init terraform
    system("cd #{path(args)} && terraform init")

    puts "Cluster configuration scaffolding complete: #{path(args)}"
  end

  desc "Deploy a new cluster"
  task :deploy, [:cluster, :provider] do |t, args|
    check_args(args)
    puts "Deploy cluster config: #{path(args)}"
    system("cd #{path(args)} && terraform apply -auto-approve") or exit 1
    puts "Cluster #{path(args)} created successfully"
  end

  desc "Destroy a cluster"
  task :destroy, [:cluster, :provider] do |t, args|
    check_args(args)
    puts "Destroy running cluster: #{path(args)}"
    system("cd #{path(args)} && terraform destroy") or exit 1
    puts "Cluster #{path(args)} destroyed successfully"
  end

  namespace :openapi do
    namespace :deploy do
      desc "Deploy OpenAPI - Swagger documentation ui"
      task :swaggerui, [:chainnet, :provider, :namespace] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade swagger-ui #{cwd}/../../deploy/helm/swagger-ui \
          --install -n #{ns(args)} --create-namespace \
        }

        system({"KUBECONFIG" => kubeconfig(args)}, cmd)
      end

      desc "Deploy OpenAPI - Prism Mock server "
      task :prism, [:chainnet, :provider, :namespace] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade prism #{cwd}/../../deploy/helm/prism \
          --install -n #{ns(args)} --create-namespace \
        }

        system({"KUBECONFIG" => kubeconfig(args)}, cmd)
      end
    end
  end

  desc "Manage sifnode deploy, upgrade, etc processes"
  namespace :sifnode do
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
  end

  desc "ebrelayer Operations"
  namespace :ebrelayer do
    desc "Deploy a new ebrelayer to an existing cluster"
    task :deploy, [:cluster, :chainnet, :provider, :namespace, :image, :image_tag, :node_host, :eth_websocket_address, :eth_bridge_registry_address, :eth_private_key, :moniker, :mnemonic] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade ebrelayer #{cwd}/../../deploy/helm/ebrelayer \
        --install -n #{ns(args)} --create-namespace \
        --set image.repository=#{image_repository(args)} \
        --set image.tag=#{image_tag(args)} \
        --set ebrelayer.env.chainnet=#{args[:chainnet]} \
        --set ebrelayer.args.nodeHost=#{args[:node_host]} \
        --set ebrelayer.args.ethWebsocketAddress=#{args[:eth_websocket_address]} \
        --set ebrelayer.args.ethBridgeRegistryAddress=#{args[:eth_bridge_registry_address]} \
        --set ebrelayer.env.ethPrivateKey=#{args[:eth_private_key]} \
        --set ebrelayer.env.moniker=#{args[:moniker]} \
        --set ebrelayer.args.mnemonic=#{args[:mnemonic]}
      }

      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end

  desc "Block Explorer"
  namespace :blockexplorer do
    desc "Deploy a Block Explorer to an existing cluster"
    task :deploy, [:cluster, :chainnet, :provider, :namespace, :root_url, :genesis_url, :rpc_url, :lcd_url] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade block-explorer #{cwd}/../../deploy/helm/block-explorer \
        --install -n #{ns(args)} --create-namespace \
        --set blockExplorer.env.chainnet=#{args[:chainnet]} \
        --set blockExplorer.env.rootURL=#{args[:root_url]} \
        --set blockExplorer.env.genesisURL=#{args[:genesis_url]} \
        --set blockExplorer.env.remote.rpcURL=#{args[:rpc_url]} \
        --set blockExplorer.env.remote.lcdURL=#{args[:lcd_url]}
      }

      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end

  desc "eth operations"
  namespace :ethereum do
    desc "Deploy an ETH node"
    task :deploy, [:cluster, :provider, :namespace, :network] do |t, args|
      check_args(args)

      if args.has_key? :network
        network_id =  if args[:network] == "ropsten"
                        3
                      else
                        1
                      end
      end

      if args.has_key? :network
        cmd = %Q{helm upgrade ethereum #{cwd}/../../deploy/helm/ethereum \
            --install -n #{ns(args)} --create-namespace \
            --set geth.args.network='--#{args[:network]}' \
            --set geth.args.networkID=#{network_id} \
            --set ethstats.env.websocketSecret=#{SecureRandom.base64 20}
            }
      else
        cmd = %Q{helm upgrade ethereum #{cwd}/../../deploy/helm/ethereum \
            --install -n #{ns(args)} --create-namespace \
            --set ethstats.env.webSocketSecret=#{SecureRandom.base64 20}
            }
      end

      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end

  desc "logstash operations"
  namespace :logstash do
    desc "Deploy a logstash node"
    task :deploy, [:cluster, :provider, :namespace, :elasticsearch_username, :elasticsearch_password] do |t, args|
      cmd = %Q{helm upgrade logstash #{cwd}/../../deploy/helm/logstash \
            --install -n #{ns(args)} --create-namespace \
            --set logstash.args.cluster=#{args[:cluster]} \
            --set logstash.args.elasticsearchUsername=#{args[:elasticsearch_username]} \
            --set logstash.args.elasticsearchPassword=#{args[:elasticsearch_password]} \
      }

      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end

  desc "namespace operations"
  namespace :namespace do
    desc "Destroy an existing namespace"
    task :destroy, [:chainnet, :provider, :namespace, :skip_prompt] do |t, args|
      check_args(args)
      are_you_sure(args)
      cmd = "kubectl delete namespace #{args[:namespace]}"
      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end
end

#
# Get the path to our terraform config based off the supplied rake args
#
# @param args Arguments passed to rake
#
def path(args)
  return "#{cwd}/../../.live/sifchain-#{args[:provider]}-#{args[:cluster]}" if args.has_key? :cluster

  "#{cwd}/../../.live/sifchain-#{args[:provider]}-#{args[:chainnet]}"
end

#
# Get the path to our kubeconfig based off the supplied rake args
#
# @param args Arguments passed to rake
#
def kubeconfig(args)
  return "#{path(args)}/kubeconfig_sifchain-#{args[:provider]}-#{args[:cluster]}" if args.has_key? :cluster

  "#{path(args)}/kubeconfig_sifchain-#{args[:provider]}-#{args[:chainnet]}"
end

#
# k8s namespace
#
# @param args Arguments passed to rake
#
def ns(args)
  args[:namespace] ? "#{args[:namespace]}" : "sifnode"
end

#
# Image tag
#
# @param args Arguments passed to rake
#
def image_tag(args)
  args[:image_tag] ? "#{args[:image_tag]}" : "testnet"
end

#
# Image repository
#
# @param args Arguments passed to rake
#
def image_repository(args)
  args[:image] ? "#{args[:image]}" : "sifchain/sifnoded"
end
