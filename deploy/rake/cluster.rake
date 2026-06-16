require "fileutils"

desc "management processes for the kube cluster and terraform commands"
namespace :cluster do
  desc "Scaffold new cluster environment configuration"
  task :scaffold, [:chainnet, :provider] do |_t, args|
    check_args(args)

    FileUtils.mkdir_p("#{cwd}/../../.live")
    FileUtils.mkdir_p(path(args))

    template_root = "#{cwd}/../terraform/template/#{template_provider_id(args[:provider])}"
    render_template(
      "#{template_root}/cluster.tf.tpl",
      "#{path(args)}/main.tf",
      { chainnet: args[:chainnet] }
    )

    envrc_template = "#{template_root}/.envrc.tpl"
    if File.exist?(envrc_template)
      render_template(
        envrc_template,
        "#{path(args)}/.envrc",
        { chainnet: args[:chainnet] }
      )
    end

    Dir.chdir(path(args)) do
      system("terraform init") or exit 1
    end

    puts "Cluster configuration scaffolding complete: #{path(args)}"
  end

  desc "Deploy a new cluster"
  task :deploy, [:chainnet, :provider] do |_t, args|
    check_args(args)
    puts "Deploy cluster config: #{path(args)}"

    Dir.chdir(path(args)) do
      system("terraform apply -auto-approve") or exit 1
    end

    write_kubeconfig(args)
    puts "Cluster #{path(args)} created successfully"
  end

  desc "Destroy a cluster"
  task :destroy, [:chainnet, :provider] do |_t, args|
    check_args(args)
    puts "Destroy running cluster: #{path(args)}"

    Dir.chdir(path(args)) do
      system("terraform destroy") or exit 1
    end

    File.delete(kubeconfig(args)) if File.exist?(kubeconfig(args))
    puts "Cluster #{path(args)} destroyed successfully"
  end

  desc "Manage sifnode deploy, upgrade, etc processes"
  namespace :sifnode do
    namespace :deploy do
      desc "Deploy a single standalone sifnode on to your cluster"
      task :standalone, [:chainnet, :provider, :namespace, :image, :image_tag, :moniker, :mnemonic] do |_t, args|
        check_args(args)

        cmd = %Q{helm upgrade sifnode #{cwd}/../../deploy/helm/sifnode \
          --install -n #{ns(args)} --create-namespace \
          --set sifnode.env.chainnet=#{args[:chainnet]} \
          --set sifnode.env.moniker=#{args[:moniker]} \
          --set sifnode.env.mnemonic="#{args[:mnemonic]}" \
          #{common_helm_flags(args)}
        }

        system({ "KUBECONFIG" => kubeconfig(args) }, cmd) or exit 1
      end

      desc "Deploy a single network-aware sifnode on to your cluster"
      task :peer, [:chainnet, :provider, :namespace, :image, :image_tag, :moniker, :mnemonic, :peer_address, :genesis_url] do |_t, args|
        check_args(args)

        cmd = %Q{helm upgrade sifnode #{cwd}/../../deploy/helm/sifnode \
          --install -n #{ns(args)} --create-namespace \
          --set sifnode.env.chainnet=#{args[:chainnet]} \
          --set sifnode.env.moniker=#{args[:moniker]} \
          --set sifnode.env.mnemonic="#{args[:mnemonic]}" \
          --set sifnode.env.peerAddress=#{args[:peer_address]} \
          --set sifnode.env.genesisURL=#{args[:genesis_url]} \
          #{common_helm_flags(args)}
        }

        system({ "KUBECONFIG" => kubeconfig(args) }, cmd) or exit 1
      end
    end

    desc "Destroy an existing namespace"
    task :destroy, [:chainnet, :provider, :namespace, :skip_prompt] do |_t, args|
      check_args(args)
      are_you_sure(args)
      cmd = "kubectl delete namespace #{ns(args)}"
      system({ "KUBECONFIG" => kubeconfig(args) }, cmd) or exit 1
    end
  end
end

#
# Get the path to our terraform config based off the supplied rake args
#
def path(args)
  "#{cwd}/../../.live/sifchain-#{provider_id(args[:provider])}-#{args[:chainnet]}"
end

#
# Get the path to our kubeconfig based off the supplied rake args
#
def kubeconfig(args)
  "#{path(args)}/kubeconfig_sifchain-#{provider_id(args[:provider])}-#{args[:chainnet]}"
end
