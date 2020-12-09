require "securerandom"

desc "management processes for the kube cluster and terraform commands"
namespace :cluster do
  desc "Scaffold new cluster environment configuration"
  task :scaffold, [:chainnet, :provider] do |t, args|
    check_args(args)

    # create path location
    system("mkdir -p #{cwd}/../../.live")
    system("mkdir #{path(args)}") or exit

    # create config from template
    system("go run github.com/belitre/gotpl #{cwd}/../terraform/template/aws/cluster.tf.tpl \
      --set chainnet=#{args[:chainnet]} > #{path(args)}/main.tf")

    system("go run github.com/belitre/gotpl #{cwd}/../terraform/template/aws/.envrc.tpl \
      --set chainnet=#{args[:chainnet]} > #{path(args)}/.envrc")

    # init terraform
    system("cd #{path(args)} && terraform init")

    puts "Cluster configuration scaffolding complete: #{path(args)}"
  end

  desc "Deploy a new cluster"
  task :deploy, [:chainnet, :provider] do |t, args|
    check_args(args)
    puts "Deploy cluster config: #{path(args)}"
    system("cd #{path(args)} && terraform apply -auto-approve") or exit 1
    system("chmod 600 #{path(args)}/kubeconfig_sifchain-#{args[:provider]}-#{args[:chainnet]}")
    puts "Cluster #{path(args)} created successfully"
  end

  desc "Status of your cluster"
  task :status, [:chainnet, :provider] do
    puts "Build me!"
  end

  desc "Backup your cluster"
  task :backup, [:chainnet, :provider] do
    puts "Build me!"
  end

  task :destroy, [:chainnet, :provider] do |t, args|
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

        system({"KUBECONFIG" => kubeconfig(args) }, cmd)
      end

      desc "Deploy OpenAPI - Prism Mock server "
      task :prism, [:chainnet, :provider, :namespace] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade prism #{cwd}/../../deploy/helm/prism \
          --install -n #{ns(args)} --create-namespace \
        }

        system({"KUBECONFIG" => kubeconfig(args) }, cmd)
      end
    end
  end

  desc "Manage sifnode deploy, upgrade, etc processes"
  namespace :sifnode do
    namespace :deploy do
      desc "Deploy a single standalone sifnode on to your cluster"
      task :standalone, [:chainnet, :provider, :namespace, :image, :image_tag, :moniker, :mnemonic] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade sifnode #{cwd}/../../deploy/helm/sifnode \
          --set sifnode.env.chainnet=#{args[:chainnet]} \
          --set sifnode.env.moniker=#{args[:moniker]} \
          --set sifnode.env.mnemonic=#{args[:mnemonic]} \
          --install -n #{ns(args)} --create-namespace \
          --set image.tag=#{image_tag(args)} \
          --set image.repository=#{image_repository(args)}
        }

        system({"KUBECONFIG" => kubeconfig(args) }, cmd)
      end

      desc "Deploy a single network-aware sifnode on to your cluster"
      task :peer, [:chainnet, :provider, :namespace, :image, :image_tag, :moniker, :mnemonic, :peer_address, :genesis_url] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade sifnode #{cwd}/../../deploy/helm/sifnode \
          --install -n #{ns(args)} --create-namespace \
          --set sifnode.env.chainnet=#{args[:chainnet]} \
          --set sifnode.env.moniker=#{args[:moniker]} \
          --set sifnode.env.mnemonic=#{args[:mnemonic]} \
          --set sifnode.env.peerAddress=#{args[:peer_address]} \
          --set sifnode.env.genesisURL=#{args[:genesis_url]} \
          --set image.tag=#{image_tag(args)} \
          --set image.repository=#{image_repository(args)}
        }

        system({"KUBECONFIG" => kubeconfig(args) }, cmd)
      end
    end

    desc "Destroy an existing namespace"
    task :destroy, [:chainnet, :provider, :namespace, :skip_prompt] do |t, args|
      check_args(args)
      are_you_sure(args)
      cmd = "kubectl delete namespace #{args[:namespace]}"
      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end

  desc "ebrelayer Operations"
  namespace :ebrelayer do
    desc "Deploy a new ebrelayer to an existing cluster"
    task :deploy, [:chainnet, :provider, :namespace, :image, :image_tag, :mnemonic, :eth_websocket_address, :eth_bridge_registry_address, :eth_private_key, :moniker] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade sifnode #{cwd}/../../deploy/helm/sifnode \
        --set sifnode.env.chainnet=#{args[:chainnet]} \
        --install -n #{ns(args)} \
        --set ebrelayer.image.repository=#{image_repository(args)} \
        --set ebrelayer.image.tag=#{image_tag(args)} \
        --set ebrelayer.enabled=true \
        --set ebrelayer.env.mnemonic=#{args[:mnemonic]} \
        --set ebrelayer.env.ethWebsocketAddress=#{args[:eth_websocket_address]} \
        --set ebrelayer.env.ethBridgeRegistryAddress=#{args[:eth_bridge_registry_address]} \
        --set ebrelayer.env.ethPrivateKey=#{args[:eth_private_key]} \
        --set ebrelayer.env.moniker=#{args[:moniker]}
      }

      system({"KUBECONFIG" => kubeconfig(args) }, cmd)
    end

    desc "Destroy a running ebrelayer on an existing cluster"
    task :destroy, [:chainnet, :provider, :namespace] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade sifnode #{cwd}/../../deploy/helm/sifnode \
        --set sifnode.env.chainnet=#{args[:chainnet]} \
        --install -n #{ns(args)} \
        --set ebrelayer.enabled=false
      }

      system({"KUBECONFIG" => kubeconfig(args) }, cmd)
    end
  end

  desc "Block Explorer"
  namespace :blockexplorer do
    desc "Deploy a Block Explorer to an existing cluster"
    task :deploy, [:chainnet, :provider] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade block-explorer #{cwd}/../../deploy/helm/block-explorer \
        --install -n block-explorer \
        --create-namespace
      }

      system({"KUBECONFIG" => kubeconfig(args) }, cmd)
    end

    desc "Destroy a running Block Explorer on an existing cluster"
    task :destroy, [:chainnet, :provider] do |t, args|
      check_args(args)

      cmd = %Q{helm delete block-explorer --namespace block-explorer && \
        kubectl delete ns block-explorer
      }

      system({"KUBECONFIG" => kubeconfig(args) }, cmd)
    end
  end

  desc "Manage eth full node deploy, upgrade, etc processes"
  namespace :ethnode do
    desc "Deploy a full eth node onto your cluster"
    task :deploy, [:chainnet, :provider:, :image, :image_tag, :network] do |t, args|
      check_args(args)
      eth_wallet(args)
      helmrepos(args)

      eth_wallet_private = `cat #{path(args)}/ethereum-generate-wallet/ethereum-wallet-generator.output \
        | grep -o -P '(?<=Private key: ).*' \
        | tr -d '[:space:]'`
      eth_wallet_public = `cat #{path(args)}/ethereum-generate-wallet/ethereum-wallet-generator.output \
        | grep -o -P '(?<=Public key: ).*' \
        | tr -d '[:space:]'`
      eth_wallet_address = `cat #{path(args)}/ethereum-generate-wallet/ethereum-wallet-generator.output \
        | grep -o -P '(?<=Address: ).*' \
        | tr -d '[:space:]'`
      eth_wallet_secret = SecureRandom.base64 20

      cmd = %Q{helm upgrade ethnode #{cwd}/../../deploy/helm/eth-full-node \
        --install -n ethnode \
        --set geth.image.tag=#{arg[:image]} \
        --set geth.image.tag=#{arg[:image_tag]} \
        --set geth.account.privateKey=#{eth_wallet_private} \
        --set geth.account.address=#{eth_wallet_address} \
        --set geth.account.secret=#{eth_wallet_secret} \
        --set geth.genesis.networkId=#{args[:network]} \
        --create-namespace \
        --debug
      }

      system({"KUBECONFIG" => kubeconfig(args) }, cmd)
    end

    desc "Uninstall ethnode"
    task :uninstall, [:chainnet, :provider] do |t, args|
      check_args(args)

      cmd = %Q{helm delete ethnode --namespace ethnode && \
        kubectl delete ns ethnode
      }

      system({"KUBECONFIG" => kubeconfig(args) }, cmd)
    end
  end
end

#
# Get the path to our terraform config based off the supplied rake args
#
# @param args Arguments passed to rake
#
def path(args)
  "#{cwd}/../../.live/sifchain-#{args[:provider]}-#{args[:chainnet]}"
end

#
# Get the path to our kubeconfig based off the supplied rake args
#
# @param args Arguments passed to rake
#
def kubeconfig(args)
  "#{path(args)}/kubeconfig_sifchain-#{args[:provider]}-#{args[:chainnet]}"
end

#
# Helm dependencies
#
# @param args Arguments passed to rake
#
def helmrepos(args)
  cmd = %Q{
    helm repo add stable https://charts.helm.sh/stable --force-update
  }
  system({"KUBECONFIG" => kubeconfig(args) }, cmd)
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

#
# Ethereum wallet
#
# @param args Arguments passed to rake
#
def eth_wallet(args)
  # Clone the generate-wallet repo
  system("cd #{path(args)} && [ ! -d 'ethereum-generate-wallet' ] &&  git clone https://github.com/vkobel/ethereum-generate-wallet")

  # Generate the wallet and export to file
  system("cd #{path(args)}/ethereum-generate-wallet && ./ethereum-wallet-generator.sh > ethereum-wallet-generator.output")
end
