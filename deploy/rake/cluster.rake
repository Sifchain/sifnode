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

  #======================================= PIPELINE AUTOMATION RUBY CONVERSIONS =============================================#
  desc "Vault Login"
  namespace :vault do
    desc "Ensure vault-0 pod has been successfully logged into with token. "
    task :login, [] do |t, args|
      cluster_automation = %Q{kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault login ${VAULT_TOKEN} > /dev/null}
      system(cluster_automation) or exit 1
    end
  end

  desc "Sifchain Art."
  namespace :generate do
    desc "Sifchain Art."
    task :art, [] do |t, args|

      cluster_automation = %Q{
#!/usr/bin/env bash
set +x
echo '                       iiii     ffffffffffffffff                 hhhhhhh                                 iiii'
echo '                      i::::i   f::::::::::::::::f                h:::::h                                i::::i'
echo '                       iiii   f::::::::::::::::::f               h:::::h                                 iiii'
echo '                              f::::::fffffff:::::f               h:::::h'
echo '        ssssssssss   iiiiiii  f:::::f       ffffffcccccccccccccccch::::h hhhhh         aaaaaaaaaaaaa   iiiiiiinnnn  nnnnnnnn'
echo '      ss::::::::::s  i:::::i  f:::::f           cc:::::::::::::::ch::::hh:::::hhh      a::::::::::::a  i:::::in:::nn::::::::nn'
echo '    ss:::::::::::::s  i::::i f:::::::ffffff    c:::::::::::::::::ch::::::::::::::hh    aaaaaaaaa:::::a  i::::in::::::::::::::nn'
echo '    s::::::ssss:::::s i::::i f::::::::::::f   c:::::::cccccc:::::ch:::::::hhh::::::h            a::::a  i::::inn:::::::::::::::n'
echo '     s:::::s  ssssss  i::::i f::::::::::::f   c::::::c     ccccccch::::::h   h::::::h    aaaaaaa:::::a  i::::i  n:::::nnnn:::::n'
echo '       s::::::s       i::::i f:::::::ffffff   c:::::c             h:::::h     h:::::h  aa::::::::::::a  i::::i  n::::n    n::::n'
echo '          s::::::s    i::::i  f:::::f         c:::::c             h:::::h     h:::::h a::::aaaa::::::a  i::::i  n::::n    n::::n'
echo '    ssssss   s:::::s  i::::i  f:::::f         c::::::c     ccccccch:::::h     h:::::ha::::a    a:::::a  i::::i  n::::n    n::::n'
echo '    s:::::ssss::::::si::::::if:::::::f        c:::::::cccccc:::::ch:::::h     h:::::ha::::a    a:::::a i::::::i n::::n    n::::n'
echo '    s::::::::::::::s i::::::if:::::::f         c:::::::::::::::::ch:::::h     h:::::ha:::::aaaa::::::a i::::::i n::::n    n::::n'
echo '     s:::::::::::ss  i::::::if:::::::f          cc:::::::::::::::ch:::::h     h:::::h a::::::::::aa:::ai::::::i n::::n    n::::n'
echo '      sssssssssss    iiiiiiiifffffffff            cccccccccccccccchhhhhhh     hhhhhhh  aaaaaaaaaa  aaaaiiiiiiii nnnnnn    nnnnnn'
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Install Cert-Manager If Not Exists"
  namespace :certmanager do
    desc "Install Cert-Manager Into Kubernetes"
    task :install, [] do |t, args|

      check_namespace=`kubectl get namespaces --kubeconfig=./kubeconfig | grep cert-manager`
      puts "check namespace #{check_namespace}"
      if check_namespace.empty?
            create_namespace=`kubectl create namespace --kubeconfig=./kubeconfig cert-manager`
            puts "create namespace #{create_namespace}"
        else
            puts "Namespace exists"
      end

      check_helm_repo_installed = `helm repo list --kubeconfig=./kubeconfig | grep jetstack`
      puts "check helm repo installed #{check_helm_repo_installed}"
      if check_helm_repo_installed.empty?
            add_helm_repo=`helm repo add jetstack https://charts.jetstack.io --kubeconfig=./kubeconfig`
            puts "add helm repo #{add_helm_repo}"
            helm_repo_update=`helm repo update --kubeconfig=./kubeconfig`
            puts "helm repo update #{helm_repo_update}"
      else
            puts "helm repo already installed."
      end

      check_cert_manager_installed = `kubectl get deployment -n cert-manager --kubeconfig=./kubeconfig | grep cert-manager`
      if check_helm_repo_installed.empty?
            helm_install=`helm install cert-manager jetstack/cert-manager --namespace cert-manager --version v1.2.0 --kubeconfig=./kubeconfig --set installCRDs=true`
            puts "cert-manager install: #{helm_install}"
      else
            puts "cert-manager already installed."
      end

      rollout_status = `kubectl rollout status deployment/cert-manager -n cert-manager  --kubeconfig=./kubeconfig`

      puts rollout_status

    end
  end

  desc "Install Vault If Not Exists"
  namespace :vault do
    desc "Install Vault into Kubernetes Env Configured"
    task :install, [:env, :region, :path, :kmskey, :aws_role] do |t, args|
      require 'fileutils'
      require 'net/http'

      APP_NAME='vault'
      APP_NAMESPACE='vault'
      POD='vault-0'
      SERVICE='vault-internal'
      CSR_NAME='vault-csr'
      NAMESPACE='vault'
      SECRET_NAME="#{APP_NAME}-#{POD}-tls"
      TMPDIR='/tmp'

      check_namespace=`kubectl get namespaces --kubeconfig=./kubeconfig | grep vault`
      puts "check namespace #{check_namespace}"
      if check_namespace.empty?
            create_namespace=`kubectl create namespace --kubeconfig=./kubeconfig vault`
            puts "create namespace #{create_namespace}"
        else
            puts "Namespace exists"
      end

      delete_secret_if_exists = `kubectl delete secret -n vault vault-eks-creds --ignore-not-found=true`
      puts delete_secret_if_exists

      create_aws_secret=`kubectl create secret generic --kubeconfig=./kubeconfig vault-eks-creds --from-literal=AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" --from-literal=AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" -n vault`
      puts create_aws_secret

      check_vault_installed = `kubectl get pods -n vault --kubeconfig=./kubeconfig | grep vault`
      if check_vault_installed.empty?
        generate_certificate_key=`openssl genrsa -out #{TMPDIR}/vault.key 2048`

        csr_config = %Q{
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = #{SERVICE}
DNS.2 = #{SERVICE}.#{NAMESPACE}
DNS.3 = #{SERVICE}.#{NAMESPACE}.svc
DNS.4 = #{SERVICE}.#{NAMESPACE}.svc.cluster.local

DNS.5 = vault-0.#{SERVICE}
DNS.6 = vault-0.#{SERVICE}.#{NAMESPACE}
DNS.7 = vault-0.#{SERVICE}.#{NAMESPACE}.svc
DNS.8 = vault-0.#{SERVICE}.#{NAMESPACE}.svc.cluster.local

DNS.9 = vault-1.#{SERVICE}
DNS.10 = vault-1.#{SERVICE}.#{NAMESPACE}
DNS.11 = vault-1.#{SERVICE}.#{NAMESPACE}.svc
DNS.12 = vault-1.#{SERVICE}.#{NAMESPACE}.svc.cluster.local

DNS.13 = vault-2.#{SERVICE}
DNS.14 = vault-2.#{SERVICE}.#{NAMESPACE}
DNS.15 = vault-2.#{SERVICE}.#{NAMESPACE}.svc
DNS.16 = vault-2.#{SERVICE}.#{NAMESPACE}.svc.cluster.local

IP.1 = 127.0.0.1
}
        File.open("#{TMPDIR}/csr.conf", 'w') { |file| file.write(csr_config) }
        generate_server_csr = `openssl req -new -key #{TMPDIR}/vault.key -subj "/CN=#{SERVICE}.#{NAMESPACE}.svc" -config #{TMPDIR}/csr.conf -out #{TMPDIR}/server.csr`
        CSR_BASE64=`cat /tmp/server.csr | base64 | tr -d '\n'`

        certificate_request = %Q{
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: #{CSR_NAME}
  namespace: #{NAMESPACE}
spec:
  groups:
  - system:authenticated
  request: #{CSR_BASE64}
  usages:
  - digital signature
  - key encipherment
  - server auth
}

        File.open("#{TMPDIR}/csr.yaml", 'w') { |file| file.write(certificate_request) }

        certificate_request = %Q{
            kubectl delete --kubeconfig=./kubeconfig -f #{TMPDIR}/csr.yaml

            kubectl apply --kubeconfig=./kubeconfig -f #{TMPDIR}/csr.yaml

            kubectl certificate approve --kubeconfig=./kubeconfig #{CSR_NAME}

            serverCert=$(kubectl get csr --kubeconfig=./kubeconfig #{CSR_NAME} -o jsonpath='{.status.certificate}')

            echo "${serverCert}" | openssl base64 -d -A -out #{TMPDIR}/vault.crt

            kubectl config view --kubeconfig=./kubeconfig --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}' | base64 --decode > #{TMPDIR}/vault.ca

            vault_ca_base64=$(kubectl config view --kubeconfig=./kubeconfig --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')

            kubectl delete secret --kubeconfig=./kubeconfig #{SECRET_NAME} --namespace #{NAMESPACE} --ignore-not-found=true
            #apply_certificate_signing_requests = `kubectl apply --kubeconfig=./kubeconfig -f #{TMPDIR}/csr.yaml`
            #puts apply_certificate_signing_requests

            #approve_certificate_signing_requets = `kubectl certificate approve --kubeconfig=./kubeconfig #{CSR_NAME}`
            #puts approve_certificate_signing_requets

            #retrieve_server_certificate = `kubectl get csr --kubeconfig=./kubeconfig #{CSR_NAME} -o jsonpath='{.status.certificate}' | openssl base64 -d -A -out #{TMPDIR}/vault.crt`
            #puts retrieve_server_certificate

            #get_vault_ca = `kubectl config view --kubeconfig=./kubeconfig --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}' | base64 --decode > #{TMPDIR}/vault.ca`
            #puts get_vault_ca
            #kubectl get csr ${CSR_NAME} -o jsonpath='{.status.certificate}'

            #FileUtils.rm_rf("#{TMPDIR}/csr.conf")
            #FileUtils.rm_rf("#{TMPDIR}/csr.yaml")
            #FileUtils.rm_rf("#{TMPDIR}/vault.ca")
            #FileUtils.rm_rf("#{TMPDIR}/vault.key")
            #FileUtils.rm_rf("#{TMPDIR}/vault.crt")
            #FileUtils.rm_rf("#{TMPDIR}/vault.key")
        }
        system(certificate_request)

        puts "Create Kubernetes TLS Secret For Vault"
        cluster_automation = `kubectl create secret generic --kubeconfig=./kubeconfig #{SECRET_NAME} --namespace #{NAMESPACE} --from-file=vault.key=#{TMPDIR}/vault.key --from-file=vault.crt=#{TMPDIR}/vault.crt --from-file=vault.ca=#{TMPDIR}/vault.ca --from-file=vault.ca.key=#{TMPDIR}/vault.key`

        puts "Check if helm repo is installed if not install."
        check_helm_repo_setup = `helm repo list --kubeconfig=./kubeconfig | grep hashicorp`
        if check_helm_repo_setup.empty?
                add_helm_repo=`helm repo add hashicorp https://helm.releases.hashicorp.com --kubeconfig=./kubeconfig`
                puts "add helm repo #{add_helm_repo}"
                helm_repo_update = `helm repo update --kubeconfig=./kubeconfig`
                puts "helm repo update #{helm_repo_update}"
        else
                puts "Namespace exists"
        end
      end

      puts "Template the overrides file for vault."
      template_file_text = File.read("#{args[:path]}override-values.yaml").strip
      ENV.each_pair do |k, v|
          replace_string="-=#{k}=-"
          if replace_string == "-=aws_region=-"
            template_file_text.include?(k) ? (template_file_text.gsub! replace_string, "#{args[:region]}") : (puts 'env matching...')
          elsif replace_string == "-=kmskey=-"
            template_file_text.include?(k) ? (template_file_text.gsub! replace_string, "#{args[:kmskey]}") : (puts 'env matching...')
          elsif replace_string == "-=aws_role=-"
            template_file_text.include?(k) ? (template_file_text.gsub! replace_string, "#{args[:aws_role]}") : (puts 'env matching...')
          end
      end
      File.open("#{args[:path]}override-values.yaml", 'w') { |file| file.write(template_file_text) }

      puts "Check if deployment exists and install if it doesn't"
      check_vault_deployment_exist = `kubectl get statefulsets --kubeconfig=./kubeconfig -n vault | grep vault`
      if check_vault_deployment_exist.empty?
            helm_install = `helm install vault hashicorp/vault --namespace vault -f #{args[:path]}override-values.yaml --kubeconfig=./kubeconfig`
            puts "helm install #{helm_install}"
        else
            helm_upgrade = `helm upgrade vault hashicorp/vault --namespace vault -f #{args[:path]}override-values.yaml --kubeconfig=./kubeconfig`
            puts "helm upgrade #{helm_upgrade}"
      end

      puts "sleep for 300 seconds to wait for vault to start."
      sleep(300)

      puts "Ensure there is  avault pod that exists as extra mesure to ensure vault is up and running."
      check_vault_pod_exist = `kubectl get pod --kubeconfig=./kubeconfig -n vault | grep vault`
      if check_vault_pod_exist.empty?
            puts "Something went wrong no vault pods. #{check_vault_pod_exist}"
            exit 1
        else
            puts "Everything Looks Good. #{check_vault_pod_exist}"
      end

      ENV["VAULT_TOKEN"]=""
      puts "Check if vault init has been completed."
      check_vault_init = `kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault status | grep Initialized | grep true`

      if check_vault_init.empty?
            puts "vault not completed initalizing vault. "
            vault_init = %Q{
                vault_init_output=`kubectl exec -n vault  vault-0 -- vault operator init -n 1 -t 1`
                echo -e ${vault_init_output} > vault_output
                echo "sleep for 30 seconds to let vault init."
                sleep 30
                export VAULT_TOKEN=`echo $vault_init_output | cut -d ':' -f 7 | cut -d ' ' -f 2`

                ./vault_login > /dev/null
             }
            system(certificate_request)
          puts "Check if the vault token has been set."
          puts "vault token #{ENV["VAULT_TOKEN"]}"

          if ENV["VAULT_TOKEN"].to_s.empty?
                puts "no vault token there was an issue"
                exit 1
          end

          puts "uploading to s3"
          upload_to_s3 = `aws s3 cp ./vault_output s3://sifchain-vault-output-backup/#{args[:env]}/#{args[:region]}/vault-master-keys.$(date  | sed -e 's/ //g').backup --region us-west-2`
          puts upload_to_s3
          puts "ok now that the token is there lets login"
          vault_login = `kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault login ${VAULT_TOKEN} > /dev/null`
        else
            puts "Vault Already Inited."
      end

      check_kv_engine_enabled=`kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault secrets list | grep kv-v2`
      if check_kv_engine_enabled.empty?
           enable_kv_enagine=`kubectl exec --kubeconfig=./kubeconfig -n vault  vault-0 -- vault secrets enable kv-v2`
           puts "enable kv engine #{enable_kv_enagine}"
      else
           puts "kv engine already enabled."
      end

      create_test_secret = `kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault kv put kv-v2/staging/test username=test123 password=foobar123`
      puts create_test_secret

      puts "sleep for 30 seconds"
      sleep(30)

      get_test_secret = `kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault kv get kv-v2/staging/test | grep "test123"`
      if get_test_secret.empty?
           puts "Secret not found"
           exit 1
      else
           puts "Secret Found Vault Running Properly"
      end

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

  desc "Generate Temp Secrets For Application Path In Vault"
  namespace :vault do
    desc "Generate Temp Secrets For Application Path In Vault"
    task :pull_temp_secrets_file_app, [:app_name,:app_region,:app_env] do |t, args|
        require "json"
        secrets_json = `kubectl exec -n vault --kubeconfig=./kubeconfig -it vault-0 -- vault kv get -format json kv-v2/#{args[:app_region]}/#{args[:app_env]}/#{args[:app_name]}`
        data = JSON.parse(secrets_json)
        temp_secrets_string = ""
        data['data']['data'].each do |key, value|
          temp_secrets_string += "export #{key}='#{value}' \n"
        end
        File.open("tmp_secrets", 'w') { |file| file.write(temp_secrets_string) }
    end
  end

  desc "Generate Temp Secrets For Path"
  namespace :vault do
    desc "Generate Temp Secrets For Path"
    task :pull_temp_secrets_file, [:path] do |t, args|
        require "json"
        secrets_json = `kubectl exec -n vault --kubeconfig=./kubeconfig -it vault-0 -- vault kv get -format json #{args[:path]}`
        data = JSON.parse(secrets_json)
        temp_secrets_string = ""
        data['data']['data'].each do |key, value|
          temp_secrets_string += "export #{key}='#{value}' \n"
        end
        File.open("tmp_secrets", 'w') { |file| file.write(temp_secrets_string) }
    end
  end

  desc "CONFIGURE AWS PROFILE AND KUBECONFIG"
  namespace :automation do
    desc "Deploy a new ebrelayer to an existing cluster"
    task :configure_aws_credentials, [:APP_ENV, :AWS_ACCESS_KEY_ID, :AWS_SECRET_ACCESS_KEY, :AWS_REGION, :AWS_ROLE, :CLUSTER_NAME] do |t, args|
        require 'fileutils'
        require 'net/http'

        puts "Download aws-iam-authenticator"
        File.write("aws-iam-authenticator", Net::HTTP.get(URI.parse("https://amazon-eks.s3.us-west-2.amazonaws.com/1.19.6/2021-01-05/bin/linux/amd64/aws-iam-authenticator")))

        puts "Create AWS Directory!"
        FileUtils.mkdir_p("/home/runner/.aws")

        credential_file = %Q{
        [default]
        aws_access_key_id = #{args[:AWS_ACCESS_KEY_ID]}
        aws_secret_access_key = #{args[:AWS_SECRET_ACCESS_KEY]}
        region = #{args[:AWS_REGION]}

        [sifchain-base]
        aws_access_key_id = #{args[:AWS_ACCESS_KEY_ID]}
        aws_secret_access_key = #{args[:AWS_SECRET_ACCESS_KEY]}
        region = #{args[:AWS_REGION]}
        }

        config_file = %Q{
        [profile #{args[:APP_ENV]}]
        source_profile = sifchain-base
        role_arn = #{args[:AWS_ROLE]}
        color = 83000a
        role_session_name = elk_stack
        region = #{args[:AWS_REGION]}
        }

        if ENV["pipeline_debug"] == "true"
            puts "config file"
            puts config_file

            puts "credential file"
            puts credential_file
        end

        puts "Write AWS Config File."
        File.open("/home/runner/.aws/config", 'w') { |file| file.write(config_file) }

        puts "Write AWS Credential File"
        File.open("/home/runner/.aws/credentials", 'w') { |file| file.write(credential_file) }

        puts "Generate Kubernetes Config from configured profile"
        get_kubectl = %Q{
              export PATH=$(pwd):${PATH}
              aws eks update-kubeconfig --name #{args[:CLUSTER_NAME]} \
              --region #{args[:AWS_REGION]} \
              --role-arn #{args[:AWS_ROLE]} \
              --profile #{args[:APP_ENV]} \
              --kubeconfig ./kubeconfig
        }
        system(get_kubectl) or exit 1

        puts "Test Generated Kubernetes Profile"
        test_kubectl = %Q{
            kubectl get pods --all-namespaces --kubeconfig ./kubeconfig
        }
        system(test_kubectl) or exit 1
    end
  end


  desc "Utility for Doing Variable Replacement"
  namespace :utilities do
    desc "Utility for Doing Variable Replacement"
    task :template_variable_replace, [:template_file_name, :final_file_name] do |t, args|
        require 'fileutils'
        template_file_text = File.read("#{args[:template_file_name]}").strip
        ENV.each_pair do |k, v|
          replace_string="-=#{k}=-"
          template_file_text.include?(k) ? (template_file_text.gsub! replace_string, v) : (puts 'matching env vars for variable replacement...')
        end
        File.open("#{args[:final_file_name]}", 'w') { |file| file.write(template_file_text) }
    end
  end

  desc "Vault Create Policy"
  namespace :vault do
    desc "Create vault policy for application to read secrets."
    task :create_vault_policy, [:region, :app_namespace, :image, :image_tag, :env, :app_name] do |t, args|

        puts "Build Vault Policy File For Application #{args[:app_name]}"
        policy_file = %Q{
path "#{args[:region]}/#{args[:env]}/#{args[:app_name]}" { capabilities = ["read"] }
path "#{args[:region]}/#{args[:env]}/#{args[:app_name]}/*" { capabilities = ["read"] }
path "/#{args[:region]}/#{args[:env]}/#{args[:app_name]}" { capabilities = ["read"] }
path "/#{args[:region]}/#{args[:env]}/#{args[:app_name]}/*" { capabilities = ["read"] }
path "*" { capabilities = ["read"] }
        }
        File.open("#{args[:app_name]}-policy.hcl", 'w') { |file| file.write(policy_file) }

      puts "Copy Policy to the Vault Pod."
      copy_policy_file_to_pod = %Q{kubectl cp --kubeconfig=./kubeconfig #{args[:app_name]}-policy.hcl vault-0:/home/vault/#{args[:app_name]}-policy.hcl -n vault}
      system(copy_policy_file_to_pod) or exit 1

      puts "Delete Policy if it Exists for Update"
      delete_policy_if_exists = %Q{kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault policy delete #{args[:app_name]}}
      system(delete_policy_if_exists) or exit 1

      puts "Write Vault Policy Based on Copied File"
      write_policy = %Q{kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault policy write #{args[:app_name]} /home/vault/#{args[:app_name]}-policy.hcl}
      system(write_policy) or exit 1

      puts "Enable Policy"
      enable_policy = %Q{kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault write sys/internal/counters/config enabled=enable}
      system(enable_policy) or exit 1

      puts "Delete the Policy File and Cleanup After."
      File.delete("#{args[:app_name]}-policy.hcl") if File.exist?("#{args[:app_name]}-policy.hcl")

    end
  end

  desc "Vault Enable Kubernetes"
  namespace :vault do
    desc "Enable Application and Vault to Talk to Kubernetes."
    task :enable_kubernetes, [] do |t, args|
      check_kubernetes_enabled = `kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault auth list | grep kubernetes`
      if check_kubernetes_enabled.include?("kubernetes")
        puts "Kubernetes Already Enabled"
      else
        enable_kubernetes = %Q{kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault auth enable kubernetes}
        system(enable_kubernetes) or exit 1
      end
    end
  end

  desc "Vault Configure Kubernetes for Application"
  namespace :vault do
    desc "Setup Service Account, and Vault Security Connections for Application."
    task :configure_application, [:app_namespace, :image, :image_tag, :env, :app_name] do |t, args|
      service_account = %Q{
apiVersion: v1
kind: ServiceAccount
metadata:
  name: #{args[:app_name]}
  namespace: #{args[:app_namespace]}
  labels:
    app: #{args[:app_name]}
}
      puts "Create Service Account File."
      puts service_account
      File.open("service_account.yaml", 'w') { |file| file.write(service_account) }

      puts "Create Service Account If It Exists"
      create_service_account = `kubectl apply --kubeconfig=./kubeconfig -f service_account.yaml -n #{args[:app_namespace]}`
      puts create_service_account

      puts "Create Service Account If It Exists"
      create_service_account = `kubectl apply --kubeconfig=./kubeconfig -f service_account.yaml`
      puts create_service_account

      puts "Get the Token from Pod"
      token = `kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- cat /var/run/secrets/kubernetes.io/serviceaccount/token` or exit 1

      puts "Get the Kubernetes Cluster IP"
      kubernetes_cluster_ip = `kubectl exec --kubeconfig=./kubeconfig -it vault-0 -n vault -- printenv | grep KUBERNETES_PORT_443_TCP_ADDR | cut -d '=' -f 2 | tr -d '\\n' | tr -d '\\r'` or exit 1
      puts kubernetes_cluster_ip

      ENV["token"] = token
      ENV["kubernetes_cluster_ip"] = kubernetes_cluster_ip

      puts "Write Auth Config"
      write_config_auth = `kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault write auth/kubernetes/config token_reviewer_jwt="#{ENV["token"]}" kubernetes_host="https://#{ENV["kubernetes_cluster_ip"]}:443" kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt` or exit 1
      puts write_config_auth

      puts "Write Auth Role"
      write_auth_role = `kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault write auth/kubernetes/role/#{args[:app_name]} bound_service_account_names=#{args[:app_name]} bound_service_account_namespaces=#{args[:app_namespace]} policies=#{args[:app_name]} ttl=1h` or exit 1
      puts write_auth_role

      puts "Clean Up"
      remove_service_account = `rm -rf service_account.yaml`
      puts remove_service_account

    end
  end

  desc "Execute Anchore Security Image Scan"
  namespace :anchore do
    desc "Execute Anchore Security Image Scan"
    task :image_scan, [:image, :image_tag, :app_name] do |t, args|
      anchore_image_scan = %Q{curl -s https://ci-tools.anchore.io/inline_scan-latest | bash -s -- -f -r -d cmd/#{args[:app_name]}/Dockerfile -p "#{args[:image]}:#{args[:image_tag]}"}
      system(anchore_image_scan) or exit 1
    end
  end

  desc "Check Vault Secret Exists"
  namespace :vault do
    desc "Check Vault Secret Exists"
    task :check_application_configured, [:app_env, :region, :app_name] do |t, args|
      vault_secret_check = `kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault kv get kv-v2/#{args[:region]}/#{args[:app_env]}/#{args[:app_name]}`
      if vault_secret_check.include?("#No value found")
        puts "Application Not Configured Please Run https://github.com/Sifchain/chainOps/actions/workflows/setup_new_application_in_vault.yaml"
        exit 1
      else
        puts "Secret Exists"
      end
    end
  end

  desc "Deploy Helm Files"
  namespace :vault do
    desc "Deploy Helm Files"
    task :helm_deploy, [:app_namespace, :image, :image_tag, :env, :app_name] do |t, args|
      puts "Deploy the Helm Files."
      deoploy_helm = %Q{helm upgrade #{args[:app_name]} deploy/helm/#{args[:app_name]} --install -n #{args[:app_namespace]} --create-namespace --set image.repository=#{args[:image]} --set image.tag=#{args[:image_tag]} --kubeconfig=./kubeconfig}
      system(deoploy_helm) or exit 1

      puts "Use kubectl rollout to wait for pods to start."
      check_kubernetes_rollout_status = %Q{kubectl rollout status --kubeconfig=./kubeconfig deployment/#{args[:app_name]} -n #{args[:app_namespace]}}
      system(check_kubernetes_rollout_status) or exit 1
    end
  end


  desc "Check kubernetes pod for specific log entry to ensure valid deployment."
  namespace :kubernetes do
    desc "Check kubernetes pod for specific log entry to ensure valid deployment."
    task :log_validate, [:APP_NAME, :APP_NAMESPACE, :SEARCH_PATH] do |t, args|
        ENV["APP_NAMESPACE"] = "#{args[:APP_NAMESPACE]}"
        ENV["APP_NAME"] = "#{args[:APP_NAME]}"
        was_successful = false
        max_loops = 20
        loop_count = 0
        until was_successful == true
            pod_name = `kubectl get pods --kubeconfig=./kubeconfig -n #{ENV["APP_NAMESPACE"]} | grep #{ENV["APP_NAME"]} | cut -d ' ' -f 1`.strip
            puts "looking up logs fo #{pod_name}"
            pod_logs = `kubectl logs #{pod_name} --kubeconfig=./kubeconfig -c ebrelayer -n #{ENV["APP_NAMESPACE"]}`
            if pod_logs.include?(args[:SEARCH_PATH])
                #:SEARCH_PATH "new transaction witnessed in sifchain client."
                puts "Log Search Completed Container Running and Producing Valid Logs"
                was_successful = true
                break
            end
            loop_count += 1
            puts "On Loop #{loop_count} of #{max_loops}"
            if loop_count >= max_loops
                puts "Reached Max Loops"
                break
            end
            sleep(60)
        end
    end
  end


  desc "Wait for Release Pipeline to Finish."
  namespace :release do
    desc "Wait for Release Pipeline to Finish."
    task :wait_for_release_pipeline, [:APP_ENV, :RELEASE, :GIT_TOKEN] do |t, args|
        require 'rest-client'
        require 'json'
        job_succeeded = false
        max_loops = 20
        loop_count = 0
        until job_succeeded == true
            headers = {"Accept": "application/vnd.github.v3+json","Authorization":"token #{args[:GIT_TOKEN]}"}
            response = RestClient.get 'https://api.github.com/repos/Sifchain/sifnode/actions/workflows', headers
            find_release="#{args[:APP_ENV]}-#{args[:RELEASE]}"
            json_response_object = JSON.parse response.body
            json_response_object["workflows"].each do |child|
                if child["name"] == "Release"
                    workflow_id = child["id"]
                    response = RestClient.get "https://api.github.com/repos/Sifchain/sifnode/actions/workflows/#{workflow_id}/runs", headers
                    json_response_job_object = JSON.parse response.body
                    json_response_job_object["workflow_runs"].each do |job|
                        if job["head_branch"] == find_release
                            puts "Release Job: #{job["head_branch"]} finished with state: #{job["conclusion"]}"
                            puts job["head_branch"]
                            puts job["status"]
                            puts job["conclusion"]
                            job_succeeded = true
                            break
                        end
                    end
                end
            end
            loop_count += 1
            puts loop_count
            puts "On Loop #{loop_count} of #{max_loops}"
            if loop_count >= max_loops
                puts "Reached Max Loops"
                break
            end
            sleep(60)
        end
    end
  end

  desc "Import Key Ring"
  namespace :release do
    desc "Import Key Ring"
    task :import_keyring, [:moniker, :app_env] do |t, args|
        File.open("tmp_keyring_rendered", "w+") do |f|
            ENV["keyring_pem"]&.split("-=n=-")&.each { |line| f.puts(line) }
        end
       import_key_ring=`yes "${keyring_passphrase}" | go run ./cmd/sifnodecli keys import #{args[:moniker]} tmp_keyring_rendered --keyring-backend test`
       puts "import key ring"
       puts import_key_ring
    end
  end

  desc "Create Github Release."
  namespace :release do
    desc "Create Github Release."
    task :create_github_release, [:release, :env, :token] do |t, args|
        require 'rest-client'
        require 'json'
        begin
            headers = {content_type: :json, "Accept": "application/vnd.github.v3+json", "Authorization":"token #{args[:token]}"}
            payload = {"tag_name"  =>  "#{args[:env]}-#{args[:release]}","name"  =>  "#{args[:env]} v#{args[:release]}","body"  => "Sifchain #{args[:env]} Release v#{args[:release]}","prerelease"  =>  true}.to_json
            response = RestClient.post 'https://api.github.com/repos/Sifchain/sifnode/releases', payload, headers
            json_response_job_object = JSON.parse response.body
            puts json_response_job_object
        rescue
            puts 'Release Already Exists'
        end
    end
  end

  desc "Create Release Governance Request."
  namespace :release do
    desc "Create Release Governance Request."
    task :generate_governance_release_request, [:upgrade_hours, :block_time, :deposit, :rowan, :chainnet, :release_version, :from, :app_env, :token] do |t, args|
        require 'rest-client'
        require 'json'

        puts "Looking for the Release Handler"
        release_search = "release-#{args[:release_version]}"
        setupHandlers = File.read("app/setupHandlers.go").strip
        setupHandlers.include?(release_search) ? (puts 'Found') : (exit 1)

        puts "Calculating Upgrade Block Height"
        if "#{args[:app_env]}" == "mainnet"
            puts "Mainnet"
            response = RestClient.get "http://rpc.sifchain.finance/abci_info?"
            json_response_object = JSON.parse response.body
        else
            puts "Testnet"
            response = RestClient.get "http://rpc-#{args[:app_env]}.sifchain.finance/abci_info?"
            json_response_object = JSON.parse response.body
        end
        current_height = json_response_object["result"]["response"]["last_block_height"].to_f
        average_block_time = "#{args[:block_time]}".to_f
        average_time = 50 / average_block_time
        average_time = average_time * 60 * "#{args[:upgrade_hours]}".to_f
        future_block_height = current_height + average_time + 100
        block_height = future_block_height.round
        puts "Block Height #{block_height}"

        sha_token=""
        headers = {"Accept": "application/vnd.github.v3+json","Authorization":"token #{args[:token]}"}
        response = RestClient.get 'https://api.github.com/repos/Sifchain/sifnode/releases', headers
        json_response_job_object = JSON.parse response.body
        json_response_job_object.each do |release|
            if release["tag_name"].include?("#{args[:app_env]}-#{args[:release_version]}")
                release["assets"].each do |asset|
                    if asset["name"].include?(".sha256")
                        response = RestClient.get asset["browser_download_url"], headers
                        sha_token = response.body.strip
                    end
                end
            end
        end
        puts "Sha found #{sha_token}"

        if "#{args[:app_env]}" == "mainnet"
            governance_request = %Q{ yes "${keyring_passphrase}" | go run ./cmd/sifnodecli tx gov submit-proposal software-upgrade release-#{args[:release_version]} \
                --from #{args[:from]} \
                --deposit #{args[:deposit]} \
                --upgrade-height #{block_height} \
                --info '{"binaries":{"linux/amd64":"https://github.com/Sifchain/sifnode/releases/download/mainnet-#{args[:release_version]}/sifnoded-#{args[:app_env]}-#{args[:release_version]}-linux-amd64.zip?checksum='#{sha_token}'"}}' \
                --title release-#{args[:release_version]} \
                --description release-#{args[:release_version]} \
                --node tcp://rpc.sifchain.finance:80 \
                --keyring-backend test \
                -y \
                --chain-id #{args[:chainnet]} \
                --gas-prices "#{args[:rowan]}"
                sleep 60 }
            system(governance_request) or exit 1
        else
            puts "create dev net gov request #{sha_token}"
            governance_request = %Q{ yes "${keyring_passphrase}" | go run ./cmd/sifnodecli tx gov submit-proposal software-upgrade release-#{args[:release_version]} \
                --from #{args[:from]} \
                --deposit #{args[:deposit]} \
                --upgrade-height #{block_height} \
                --info '{"binaries":{"linux/amd64":"https://github.com/Sifchain/sifnode/releases/download/#{args[:app_env]}-#{args[:release_version]}/sifnoded-#{args[:app_env]}-#{args[:release_version]}-linux-amd64.zip?checksum='#{sha_token}'"}}' \
                --title release-#{args[:release_version]} \
                --description release-#{args[:release_version]} \
                --node tcp://rpc-#{args[:app_env]}.sifchain.finance:80 \
                --keyring-backend test \
                -y \
                --chain-id #{args[:chainnet]} \
                --gas-prices "#{args[:rowan]}"
                sleep 60 }
            system(governance_request) or exit 1
        end
    end
  end

  desc "Create Release Governance Request Vote."
  namespace :release do
    desc "Create Release Governance Request Vote."
    task :generate_vote, [:rowan, :chainnet, :from, :app_env] do |t, args|
        if "#{args[:app_env]}" == "mainnet"
            governance_request = %Q{
vote_id=$(go run ./cmd/sifnodecli q gov proposals --node tcp://rpc.sifchain.finance:80 --trust-node -o json | jq --raw-output 'last(.[]).id' --raw-output)
echo "vote_id $vote_id"
yes "${keyring_passphrase}" | go run ./cmd/sifnodecli tx gov vote ${vote_id} yes \
    --from #{args[:from]} \
    --keyring-backend test \
    --chain-id #{args[:chainnet]}  \
    --node tcp://rpc.sifchain.finance:80 \
    --gas-prices "#{args[:rowan]}" -y
sleep 15  }
            system(governance_request) or exit 1
        else
            governance_request = %Q{
vote_id=$(go run ./cmd/sifnodecli q gov proposals --node tcp://rpc-#{args[:app_env]}.sifchain.finance:80 --trust-node -o json | jq --raw-output 'last(.[]).id' --raw-output)
echo "vote_id $vote_id"
yes "${keyring_passphrase}" | go run ./cmd/sifnodecli tx gov vote ${vote_id} yes \
    --from #{args[:from]} \
    --keyring-backend test \
    --chain-id #{args[:chainnet]}  \
    --node tcp://rpc-#{args[:app_env]}.sifchain.finance:80 \
    --gas-prices "#{args[:rowan]}" -y
sleep 15 }
             system(governance_request) or exit 1
        end
    end
  end

  #=======================================RUBY CONVERSIONS END=============================================#

  desc "Block Explorer"
  namespace :blockexplorer do
    desc "Deploy a Block Explorer to an existing cluster"
    task :deploy, [:cluster, :chainnet, :provider, :namespace, :image, :image_tag, :root_url, :genesis_url, :rpc_url, :api_url, :mongo_password] do |t, args|
      check_args(args)

      cmd = %Q{helm upgrade block-explorer #{cwd}/../../deploy/helm/block-explorer \
        --install -n #{ns(args)} --create-namespace \
        --set image.repository=#{image_repository(args)} \
        --set image.tag=#{image_tag(args)} \
        --set blockExplorer.env.chainnet=#{args[:chainnet]} \
        --set blockExplorer.env.rootURL=#{args[:root_url]} \
        --set blockExplorer.env.genesisURL=#{args[:genesis_url]} \
        --set blockExplorer.env.remote.rpcURL=#{args[:rpc_url]} \
        --set blockExplorer.env.remote.apiURL=#{args[:api_url]} \
        --set blockExplorer.args.mongoPassword=#{args[:mongo_password]}
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
