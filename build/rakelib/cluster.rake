desc "management processes for the kube cluster and terraform commands"
namespace :cluster do
  desc "Scaffold new cluster environment configuration"
  task :scaffold, [:chainnet, :provider] do |t, args|
    check_args(args)

    # create path location
    system('mkdir -p ../.live')
    system("mkdir #{path(args)}") or exit

    # create config from template
    system("go run github.com/belitre/gotpl ./terraform/template/aws/cluster.tf.tpl \
      --set chainnet=#{args[:chainnet]} \
      > #{path(args)}/main.tf
    ")

    system("go run github.com/belitre/gotpl ./terraform/template/aws/.envrc.tpl \
      --set chainnet=#{args[:chainnet]} \
      > #{path(args)}/.envrc
    ")

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

  desc "Manage sifnode deploy, upgrade, etc processes"
  namespace :sifnode do
    namespace :deploy do
      desc "Deploy a single standalone sifnode on to your cluster"
      task :standalone, [:chainnet, :provider, :namespace, :image, :image_tag] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade #{ns(args)} ../build/helm/sifnode \
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

        cmd = %Q{helm upgrade #{ns(args)} ../build/helm/sifnode \
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

    task :uninstall, [:chainnet, :provider, :namespace] do |t, args|
      check_args(args)
      cmd = "helm delete #{ns(args)} -n #{ns(args)}"
      system({"KUBECONFIG" => kubeconfig(args) }, cmd)
    end
  end

  desc "Manage eth full node deploy, upgrade, etc processes"
  namespace :ethnode do
    desc "Deploy a full eth node on to your cluster"
    task :deploy do
      puts "Coming soon! "
    end
  end
end

# path returns the path of the terraform config that is generated as part of the scaffold task
def path(args)
  "../.live/sifchain-#{args[:provider]}-#{args[:chainnet]}"
end

# check_args checks to make sure the required args are passed in
def check_args(args)
  if args[:chainnet] == nil
    puts "Please provider a chainnet argument E.g testnet, mainnet, etc"
    exit
  end

  case args[:provider]
  when "aws"
  when "az"
    puts "Build me!"
    exit
  when "gcp"
    puts "Build me!"
    exit
  when "do"
    puts "Build me!"
    exit
  else
    puts "Please provide a cloud host provider. E.g aws"
    exit
  end
end

# kubeconfig returns the path to the kubeconfig file based on the args
def kubeconfig(args)
  "#{path(args)}/kubeconfig_sifchain-#{args[:provider]}-#{args[:chainnet]}"
end

# ns = namespace for kubes returns the arg with the namespace if set or the default setting
def ns(args)
  args[:namespace] ? "#{args[:namespace]}" : "sifnode"
end

# image_tag returns the arg for image_tag if set or the default tag setting
def image_tag(args)
  args[:image_tag] ? "#{args[:image_tag]}" : "testnet"
end

# image_repository returns the arg with a image if set or the default setting
def image_repository(args)
  args[:image] ? "#{args[:image]}" : "sifchain/sifnoded"
end
