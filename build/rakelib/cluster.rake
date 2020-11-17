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
    puts "Now run `rake cluster:create[#{args[:chainnet]},#{args[:provider]}]` to deploy your cluster"
  end

  desc "Deploy a cluster"
  task :deploy, [:chainnet, :provider] do |t, args|
    check_args(args)
    puts "Deploy cluster config: #{path(args)}"
    system("cd #{path(args)} && terraform apply -auto-approve") or exit 1
    puts "Cluster #{path(args)} created successfully"
    puts "Now run `rake sifnode:install[#{args[:chainnet]},#{args[:provider]}]` to deploy sifnode to your cluster"
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

        cmd = %Q{helm upgrade swagger-ui #{cwd}/../../build/helm/swagger-ui \
          --install -n #{ns(args)} --create-namespace \
        }

        system({"KUBECONFIG" => kubeconfig(args) }, cmd)
      end

      desc "Deploy OpenAPI - Prism Mock server "
      task :prism, [:chainnet, :provider, :namespace] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade prism #{cwd}/../../build/helm/prism \
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
      task :standalone, [:chainnet, :provider, :namespace, :image, :image_tag] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade sifnode #{cwd}/../../build/helm/sifnode \
          --set sifnode.env.chainnet=#{args[:chainnet]} \
          --install -n #{ns(args)} --create-namespace \
          --set image.tag=#{image_tag(args)} \
          --set image.repository=#{image_repository(args)}
        }

        system({"KUBECONFIG" => kubeconfig(args) }, cmd)
      end

      desc "Deploy a single network-aware sifnode on to your cluster"
      task :peer, [:chainnet, :provider, :namespace, :image, :image_tag, :peer_address, :genesis_url] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade sifnode #{cwd}/../../build/helm/sifnode \
          --install -n #{ns(args)} --create-namespace \
          --set sifnode.env.chainnet=#{args[:chainnet]} \
          --set sifnode.env.genesisURL=#{args[:genesis_url]} \
          --set sifnode.env.peerAddress=#{args[:peer_address]} \
          --set image.tag=#{image_tag(args)} \
          --set image.repository=#{image_repository(args)}
        }

        system({"KUBECONFIG" => kubeconfig(args) }, cmd)
      end
    end

    desc "delete a namespace"
    task :destroy, [:chainnet, :provider, :namespace, :skip_prompt] do |t, args|
      check_args(args)
      are_you_sure(args)
      cmd = "kubectl delete namespace #{args[:namespace]}"
      system({"KUBECONFIG" => kubeconfig(args)}, cmd)
    end
  end

  desc "ebrelayer operations"
  namespace :ebrelayer do
    desc "Install ebrelayer"
    task :deploy, [:chainnet, :provider, :namespace, :image, :image_tag, :eth_websocket_address, :eth_bridge_registry_address, :eth_private_key, :moniker] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade sifnode #{cwd}/../../build/helm/sifnode \
        --set sifnode.env.chainnet=#{args[:chainnet]} \
        --install -n #{ns(args)} \
        --set ebrelayer.image.repository=#{image_repository(args)} \
        --set ebrelayer.image.tag=#{image_tag(args)} \
        --set ebrelayer.enabled=true \
        --set ebrelayer.env.ethWebsocketAddress=#{args[:eth_websocket_address]} \
        --set ebrelayer.env.ethBridgeRegistryAddress=#{args[:eth_bridge_registry_address]} \
        --set ebrelayer.env.ethPrivateKey=#{args[:eth_private_key]} \
        --set ebrelayer.env.moniker=#{args[:moniker]}
      }

      system({"KUBECONFIG" => kubeconfig(args) }, cmd)
    end

    desc "Uninstall ebrelayer"
    task :uninstall, [:chainnet, :provider, :namespace] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade sifnode #{cwd}/../../build/helm/sifnode \
        --set sifnode.env.chainnet=#{args[:chainnet]} \
        --install -n #{ns(args)} \
        --set ebrelayer.enabled=false
      }

      system({"KUBECONFIG" => kubeconfig(args) }, cmd)
    end
  end

  desc "Manage eth full node deploy, upgrade, etc processes"
  namespace :ethnode do
    desc "Deploy a full eth node onto your cluster"
    task :deploy do
      puts "Coming soon! "
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
