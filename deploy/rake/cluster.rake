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
      task :standalone, [:cluster, :chainnet, :provider, :namespace, :image, :image_tag, :moniker, :mnemonic, :admin_clp_addresses, :admin_oracle_address, :minimum_gas_prices, :clp_config_url] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade sifnode #{cwd}/../../deploy/helm/sifnode \
          --set sifnode.env.chainnet=#{args[:chainnet]} \
          --set sifnode.env.moniker=#{args[:moniker]} \
          --set sifnode.args.mnemonic=#{args[:mnemonic]} \
          --set sifnode.args.adminCLPAddresses=#{args[:admin_clp_addresses]} \
          --set sifnode.args.adminOracleAddress=#{args[:admin_oracle_address]} \
          --set sifnode.args.minimumGasPrices=#{args[:minimum_gas_prices]} \
          --set sifnode.args.clpConfigURL=#{args[:clp_config_url]} \
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

  desc "Vault Login"
  namespace :vault do
    desc "Ensure vault-0 pod has been successfully logged into with token. "
    task :login, [] do |t, args|
      cluster_automation = %Q{
        set +x
        kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault login ${VAULT_TOKEN} > /dev/null
        echo "Vault Ready"
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Vault Create Policy"
  namespace :vault do
    desc "Create vault policy for application to read secrets."
    task :createpolicy, [:app_namespace, :image, :image_tag, :env, :app_name] do |t, args|
      cluster_automation = %Q{
        set +x
        echo "
path \\"#{args[:env]}/#{args[:app_name]}\\" {
    capabilities = [\\"create\\", \\"read\\", \\"update\\", \\"delete\\", \\"list\\"]
}
path \\"#{args[:env]}/#{args[:app_name]}/*\\" {
    capabilities = [\\"create\\", \\"read\\", \\"update\\", \\"delete\\", \\"list\\"]
}
path \\"/#{args[:env]}/#{args[:app_name]}\\" {
    capabilities = [\\"create\\", \\"read\\", \\"update\\", \\"delete\\", \\"list\\"]
}
path \\"/#{args[:env]}/#{args[:app_name]}/*\\" {
    capabilities = [\\"create\\", \\"read\\", \\"update\\", \\"delete\\", \\"list\\"]
}
path \\"*\\" {
    capabilities = [\\"create\\", \\"read\\", \\"update\\", \\"delete\\", \\"list\\"]
}
path \\"sys/internal/counters/activity\\" {
  capabilities = [\\"read\\"]
}
path \\"sys/internal/counters/config\\" {
  capabilities = [\\"read\\", \\"update\\"]
}
path \\"sys/namespaces\\" {
  capabilities = [\\"list\\", \\"read\\", \\"update\\"]
}
path \\"sys/internal/ui/namespaces\\" {
  capabilities = [\\"read\\", \\"list\\", \\"update\\", \\"sudo\\"]
}
path \\"sys/internal/ui/mounts\\" {
  capabilities = [\\"read\\", \\"sudo\\"]
}
path \\"+/sys/internal/counters/config\\" {
  capabilities = [\\"read\\", \\"update\\"]
}
path \\"+/sys/internal/counters/activity\\" {
  capabilities = [\\"read\\"]
}
        " > #{args[:app_name]}-policy.hcl

        cat #{args[:app_name]}-policy.hcl

        kubectl cp --kubeconfig=./kubeconfig #{args[:app_name]}-policy.hcl vault-0:/home/vault/#{args[:app_name]}-policy.hcl -n vault
        kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault policy delete #{args[:app_name]}
        kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault policy write #{args[:app_name]} /home/vault/#{args[:app_name]}-policy.hcl
        kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault write sys/internal/counters/config enabled=enable

        rm -rf #{args[:app_name]}-policy.hcl
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Vault Enable Kubernetes"
  namespace :vault do
    desc "Enable Application and Vault to Talk to Kubernetes."
    task :enablekubernetes, [] do |t, args|
      cluster_automation = %Q{
        set +x
        echo "APPLY VAULT AUTH ENABLE KUBERNETES"
        check_installed=`kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault auth list | grep kubernetes`
        [ -z "$check_installed" ] && kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault auth enable kubernetes || echo "Kubernetes Already Enabled"
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Vault Configure Kubernetes for Application"
  namespace :vault do
    desc "Setup Service Account, and Vault Security Connections for Application."
    task :configureapplication, [:app_namespace, :image, :image_tag, :env, :app_name] do |t, args|
      cluster_automation = %Q{
        set +x
        echo "
apiVersion: v1
kind: ServiceAccount
metadata:
  name: #{args[:app_name]}
  namespace: #{args[:app_namespace]}
  labels:
    app: #{args[:app_name]} " > service_account.yaml

        kubectl apply --kubeconfig=./kubeconfig -f service_account.yaml -n #{args[:app_namespace]}

        token=`kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- cat /var/run/secrets/kubernetes.io/serviceaccount/token`
        kubernetes_cluster_ip=`kubectl exec --kubeconfig=./kubeconfig -it vault-0 -n vault -- printenv | grep KUBERNETES_PORT_443_TCP_ADDR | cut -d '=' -f 2 | tr -d '\n' | tr -d '\r'`

        kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault write auth/kubernetes/config token_reviewer_jwt="$token" kubernetes_host="https://$kubernetes_cluster_ip:443" kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault write auth/kubernetes/role/#{args[:app_name]} bound_service_account_names=#{args[:app_name]} bound_service_account_namespaces=#{args[:app_namespace]} policies=#{args[:app_name]} ttl=1h

        rm -rf service_account.yaml
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Vault ebrelayer Operations"
  namespace :vault do
    desc "Deploy a new ebrelayer to an existing cluster"
    task :deploy, [:app_namespace, :image, :image_tag, :env, :app_name] do |t, args|
      cluster_automation = %Q{
        set +x
        helm upgrade #{args[:app_name]} deploy/helm/#{args[:app_name]} \
            --install -n #{args[:app_namespace]} \
            --create-namespace \
            --set image.repository=#{args[:image]} \
            --set image.tag=#{args[:image_tag]} \
            --kubeconfig=./kubeconfig

        kubectl rollout status \
            --kubeconfig=./kubeconfig deployment/#{args[:app_name]} \
            -n #{args[:app_namespace]}

      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Anchore Security Docker Vulnerability Scan"
  namespace :anchore do
    desc "Deploy a new ebrelayer to an existing cluster"
    task :scan, [:image, :image_tag, :app_name] do |t, args|
      cluster_automation = %Q{
        set +x
        curl -s https://ci-tools.anchore.io/inline_scan-latest | bash -s -- -f -r -d cmd/#{args[:app_name]}/Dockerfile -p "#{args[:image]}:#{args[:image_tag]}"
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Setup AWS Profile for Automation Pipelines"
  namespace :automation do
    desc "Deploy a new ebrelayer to an existing cluster"
    task :configure_aws_kube_profile, [:app_env, :aws_access_key_id, :aws_secret_access_key, :aws_region, :aws_role, :cluster_name] do |t, args|
      cluster_automation = %Q{
        set +x
          set +x
          curl -s -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.19.6/2021-01-05/bin/linux/amd64/aws-iam-authenticator
          chmod +x ./aws-iam-authenticator
          export PATH=$(pwd):${PATH}
          mkdir -p ~/.aws

          echo "[sifchain-base]" > ~/.aws/credentials
          echo "aws_access_key_id = #{args[:aws_access_key_id]}" >> ~/.aws/credentials
          echo "aws_secret_access_key = #{args[:aws_secret_access_key]}" >> ~/.aws/credentials
          echo "region = #{args[:aws_region]}" >> ~/.aws/credentials

          echo "[profile #{args[:app_env]}]" > ~/.aws/config
          echo "source_profile = sifchain-base" >> ~/.aws/config
          echo "role_arn = #{args[:aws_role]}" >> ~/.aws/config
          echo "color = 83000a" >> ~/.aws/config
          echo "role_session_name = elk_stack" >> ~/.aws/config
          echo "region = #{args[:aws_region]}" >> ~/.aws/config

          aws eks update-kubeconfig --name #{args[:cluster_name]} --region #{args[:aws_region]} --profile #{args[:app_env]} --kubeconfig ./kubeconfig
      }
      system(cluster_automation) or exit 1
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
