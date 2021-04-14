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

  desc "Install Cert-Manager If Not Exists"
  namespace :certmanager do
    desc "Install Cert-Manager Into Kubernetes"
    task :install, [] do |t, args|
      cluster_automation = %Q{
#!/usr/bin/env bash
set +x

echo "===================STAGE INIT - GLOBAL REQUIREMENT CHECKS==================="
check_created=`kubectl get namespaces --kubeconfig=./kubeconfig | grep cert-manager`
[ -z "$check_created" ] && kubectl create namespace --kubeconfig=./kubeconfig cert-manager || echo "Namespace Exists"

echo "===================STAGE 2 - SETUP & UPDATE HELM==================="
check_created=`helm repo list --kubeconfig=./kubeconfig | grep jetstack`
[ -z "$check_created" ] && helm repo add jetstack https://charts.jetstack.io --kubeconfig=./kubeconfig && helm repo update --kubeconfig=./kubeconfig || echo "Helm Repo Already Added For Cert-Manager"

echo "===================STAGE 3 - INSTALL CERT MANAGER==================="
echo "Install Cert Manager"
check_installed=`kubectl get deployment -n cert-manager --kubeconfig=./kubeconfig | grep cert-manager`
[ -z "$check_installed" ] && helm install cert-manager jetstack/cert-manager --namespace cert-manager --version v1.2.0 --kubeconfig=./kubeconfig --set installCRDs=true || echo "CERT-MANAGER already seems to be installed."

echo "===================STAGE 4 - CHECK CERT-MANAGER ROLLOUT STATUS ==================="
echo "Use KUBECTL roll out to check status"
kubectl rollout status deployment/cert-manager -n cert-manager  --kubeconfig=./kubeconfig

      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Install Vault If Not Exists"
  namespace :vault do
    desc "Install Vault into Kubernetes Env Configured"
    task :install, [:env, :region, :path, :kmskey, :aws_role] do |t, args|
      cluster_automation = %Q{
#!/usr/bin/env bash
set +x

echo "===================STAGE INIT - GLOBAL REQUIREMENT CHECKS==================="
APP_NAME=vault
APP_NAMESPACE=vault
POD=vault-0
SERVICE=vault-internal
export CSR_NAME=vault-csr
NAMESPACE=${APP_NAMESPACE}
SECRET_NAME=${APP_NAME}-${POD}-tls
TMPDIR=/tmp

echo "ENSURE NAMESPACE EXISTS"
check_secret=`kubectl get namespaces --kubeconfig=./kubeconfig | grep vault | grep -v grep`
[ -z "$check_secret" ] && kubectl create namespace --kubeconfig=./kubeconfig vault || echo "Namespace Exists"

echo "Check to see if VAULT AWS SECRET EXISTS IF NOT CREATE."
check_created=`kubectl get secret -n vault --kubeconfig=./kubeconfig | grep vault-eks-creds`
[ -z "$check_created" ] && kubectl create secret generic --kubeconfig=./kubeconfig vault-eks-creds --from-literal=AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" --from-literal=AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" -n vault || echo "Vault EKS Secret Already Created"

echo "===================STAGE 1 - GENERATE CA AND TLS KEY AND CERT==================="
openssl genrsa -out ${TMPDIR}/vault.key 2048

cat <<EOF >${TMPDIR}/csr.conf
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
DNS.1 = ${SERVICE}
DNS.2 = ${SERVICE}.${NAMESPACE}
DNS.3 = ${SERVICE}.${NAMESPACE}.svc
DNS.4 = ${SERVICE}.${NAMESPACE}.svc.cluster.local

DNS.5 = vault-0.${SERVICE}
DNS.6 = vault-0.${SERVICE}.${NAMESPACE}
DNS.7 = vault-0.${SERVICE}.${NAMESPACE}.svc
DNS.8 = vault-0.${SERVICE}.${NAMESPACE}.svc.cluster.local

DNS.9 = vault-1.${SERVICE}
DNS.10 = vault-1.${SERVICE}.${NAMESPACE}
DNS.11 = vault-1.${SERVICE}.${NAMESPACE}.svc
DNS.12 = vault-1.${SERVICE}.${NAMESPACE}.svc.cluster.local

DNS.13 = vault-2.${SERVICE}
DNS.14 = vault-2.${SERVICE}.${NAMESPACE}
DNS.15 = vault-2.${SERVICE}.${NAMESPACE}.svc
DNS.16 = vault-2.${SERVICE}.${NAMESPACE}.svc.cluster.local

IP.1 = 127.0.0.1
EOF

openssl req -new -key ${TMPDIR}/vault.key -subj "/CN=${SERVICE}.${NAMESPACE}.svc" -config ${TMPDIR}/csr.conf -out ${TMPDIR}/server.csr

cat <<EOF >${TMPDIR}/csr.yaml
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: ${CSR_NAME}
  namespace: ${NAMESPACE}
spec:
  groups:
  - system:authenticated
  request: $(cat ${TMPDIR}/server.csr | base64 | tr -d '\\n')
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF

kubectl apply --kubeconfig=./kubeconfig -f ${TMPDIR}/csr.yaml

kubectl certificate approve --kubeconfig=./kubeconfig ${CSR_NAME}

serverCert=$(kubectl get csr --kubeconfig=./kubeconfig ${CSR_NAME} -o jsonpath='{.status.certificate}')

echo "${serverCert}" | openssl base64 -d -A -out ${TMPDIR}/vault.crt

kubectl config view --kubeconfig=./kubeconfig --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}' | base64 --decode > ${TMPDIR}/vault.ca

vault_ca_base64=$(kubectl config view --kubeconfig=./kubeconfig --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')

kubectl create secret generic --kubeconfig=./kubeconfig ${SECRET_NAME} \
        --namespace ${NAMESPACE} \
        --from-file=vault.key=${TMPDIR}/vault.key \
        --from-file=vault.crt=${TMPDIR}/vault.crt \
        --from-file=vault.ca=${TMPDIR}/vault.ca \
        --from-file=vault.ca.key=${TMPDIR}/vault.key

echo "Clean up files"
rm -rf ${TMPDIR}/csr.conf
rm -rf ${TMPDIR}/csr.yaml
rm -rf ${TMPDIR}/vault.ca
rm -rf ${TMPDIR}/vault.key
rm -rf ${TMPDIR}/vault.crt
rm -rf ${TMPDIR}/vault.key

echo "===================STAGE 2 - SETUP and UPDATE VAULT REPO ==================="
check_created=`helm repo list --kubeconfig=./kubeconfig | grep hashicorp`
[ -z "$check_created" ] && helm repo add hashicorp https://helm.releases.hashicorp.com --kubeconfig=./kubeconfig && helm repo update --kubeconfig=./kubeconfig || echo "Helm Repo Already Added For Cert-Manager"

cat << EOF > helmvaulereplace.py
#!/usr/bin/env python
vaules_yaml = open("#{args[:path]}override-values.yaml", "r").read()
vaules_yaml = vaules_yaml.replace("-=region=-", "#{args[:region]}" )
vaules_yaml = vaules_yaml.replace("-=kmskey=-", "#{args[:kmskey]}" )
vaules_yaml = vaules_yaml.replace("-=role_arn=-", "#{args[:aws_role]}" )
open("#{args[:path]}override-values.yaml", "w").write(vaules_yaml)
EOF
python helmvaulereplace.py

echo "===================STAGE 3 - INSTALL VAULT ==================="
check_deployment=`kubectl get statefulsets --kubeconfig=./kubeconfig -n vault | grep vault`
[ -z "$check_deployment" ] && helm install vault hashicorp/vault --namespace vault -f #{args[:path]}override-values.yaml --kubeconfig=./kubeconfig || helm upgrade vault hashicorp/vault --namespace vault -f #{args[:path]}override-values.yaml --kubeconfig=./kubeconfig

echo "sleep for 2 min to let vault start up"
sleep 180

check_deployment=`kubectl get pod --kubeconfig=./kubeconfig -n vault | grep vault`
[ -z "$check_deployment" ] && echo "Something Went Wrong" || echo "Vault Deployed ${check_deployment}"

vault_init_output=`kubectl exec --kubeconfig=./kubeconfig -n vault vault-0 -- vault operator init -n 1 -t 1`
echo "sleep for 30 seconds to let vault init."
sleep 30

echo -e ${vault_init_output} > vault_output
export VAULT_TOKEN=$(echo $vault_init_output | cut -d ':' -f 7 | cut -d ' ' -f 2)

vault_output_wordcount=$(cat vault_output | wc | sed -e 's/ //g')

echo "vault output word count ${vault_output_wordcount}"

if [ "${vault_output_wordcount}" -ge "200" ]; then
    aws s3 cp ./vault_output s3://sifchain-vault-output-backup/#{args[:env]}/#{args[:region]}/vault-master-keys.$(date  | sed -e 's/ //g').backup --region us-west-2
fi

kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault login ${VAULT_TOKEN} > /dev/null

echo "create kv v2 engine"
kubectl exec --kubeconfig=./kubeconfig -n vault  vault-0 -- vault secrets enable kv-v2

echo "create test secret"
kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault kv put kv-v2/staging/test username=test123 password=foobar123

echo "validate secret made it in vault."
get_secrets=`kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault kv get kv-v2/staging/test | grep "test123"`
[ -z "$get_secrets" ] && echo "not present ${get_secrets}" && exit 1 || echo "Secre Present"

      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Vault Create Policy"
  namespace :vault do
    desc "Create vault policy for application to read secrets."
    task :createpolicy, [:region, :app_namespace, :image, :image_tag, :env, :app_name] do |t, args|
      cluster_automation = %Q{
        set +x
        echo "
path \\"#{args[:region]}/#{args[:env]}/#{args[:app_name]}\\" {
    capabilities = [\\"create\\", \\"read\\", \\"update\\", \\"delete\\", \\"list\\"]
}
path \\"#{args[:region]}/#{args[:env]}/#{args[:app_name]}/*\\" {
    capabilities = [\\"create\\", \\"read\\", \\"update\\", \\"delete\\", \\"list\\"]
}
path \\"/#{args[:region]}/#{args[:env]}/#{args[:app_name]}\\" {
    capabilities = [\\"create\\", \\"read\\", \\"update\\", \\"delete\\", \\"list\\"]
}
path \\"/#{args[:region]}/#{args[:env]}/#{args[:app_name]}/*\\" {
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

  desc "Generate tmp_secrets file with vault secrest to source and remove in your automations."
  namespace :vault do
    desc "Generate tmp_secrets file with vault secrest to source and remove in your automations."
    task :generate_vault_tmp_var_source_file, [:path] do |t, args|
      cluster_automation = %Q{
#!/usr/bin/env bash
set +x
cat << EOF > pyscript.py
#!/usr/bin/env python
import json
import urllib3
http = urllib3.PoolManager()
import subprocess
print("Starting to Pull Secrets")
result = subprocess.Popen(["kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault kv get -format json #{args[:path]}"], stdout=subprocess.PIPE, shell=True)
output,error = result.communicate()
vars_return = json.loads(output.decode('utf-8'))["data"]["data"]
print("Opening temporary secrets file for writing secrets")
temp_secrets = open("tmp_secrets", "w")
for var in vars_return:
    temp_secrets.write('export {key}=\\'{values}\\' \\n'.format(key=var, values=vars_return[var]))
temp_secrets.close()
print("secrets written.")
EOF
python pyscript.py
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
        kubectl delete --kubeconfig=./kubeconfig -f service_account.yaml -n #{args[:app_namespace]}
        kubectl create --kubeconfig=./kubeconfig -f service_account.yaml -n #{args[:app_namespace]}
        token=`kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- cat /var/run/secrets/kubernetes.io/serviceaccount/token`
        kubernetes_cluster_ip=`kubectl exec --kubeconfig=./kubeconfig -it vault-0 -n vault -- printenv | grep KUBERNETES_PORT_443_TCP_ADDR | cut -d '=' -f 2 | tr -d '\n' | tr -d '\r'`
        kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault write auth/kubernetes/config token_reviewer_jwt="$token" kubernetes_host="https://$kubernetes_cluster_ip:443" kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault write auth/kubernetes/role/#{args[:app_name]} bound_service_account_names=#{args[:app_name]} bound_service_account_namespaces=#{args[:app_namespace]} policies=#{args[:app_name]} ttl=1h
        rm -rf service_account.yaml
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Check ebrelayer logs for service running serach string basically use logs to ensure events are processed."
  namespace :ebrelayer do
    desc "Check ebrelayer logs for service running serach string basically use logs to ensure events are processed."
    task :check_deployment, [:app_name, :app_namespace, :search_string] do |t, args|
      cluster_automation = %Q{
#!/usr/bin/env bash
set +x
APP_NAMESPACE=#{args[:app_namespace]}
APP_NAME=#{args[:app_name]}
echo "get pod name"
pod_name=$(kubectl get pods --kubeconfig=./kubeconfig -n ${APP_NAMESPACE} | grep ${APP_NAME} | cut -d ' ' -f 1 | sed -e 's/ //g')
echo "POD NAME ${pod_name}"
echo "see if there is log output for the pod"
logs_check=$(kubectl logs --kubeconfig=./kubeconfig -n ${APP_NAME} ${pod_name} -c ${APP_NAME} | grep '#{args[:search_string]}')

echo "set the max check loop and current count"
max_check=50
check_count=0

echo "check if the logs output was empty"
if [ -z "${logs_check}" ]; then
    while true; do
        if [ "${max_check}" == "${check_count}" ]; then
            echo "max count reached"
            break
        fi
        echo "get pod name"
        pod_name=$(kubectl get pods --kubeconfig=./kubeconfig -n ${APP_NAMESPACE} | grep ${APP_NAME} | cut -d ' ' -f 1 | sed -e 's/ //g')
        echo "POD NAME ${pod_name}"
        echo "see if there is log output for the pod"
        logs_check_loop=$(kubectl logs --kubeconfig=./kubeconfig -n ${APP_NAME} ${pod_name} -c ${APP_NAME} | grep '#{args[:search_string]}')
        echo "see if log check had data meaning search string found"
        if [ -z "${logs_check_loop}" ]; then
            echo "sleep and wait for logs"
            sleep 5
        else
            echo "service successfully started."
            break
        fi
        check_count=$((check_count+1))
        echo "${check_count} of ${max_check}"
    done
else
    echo "service successfully started."
fi
      }
      system(cluster_automation) or exit 1
    end
  end


  desc "Create Release Governance Request."
  namespace :release do
    desc "Create Release Governance Request."
    task :generate_governance_release_request, [:upgrade_hours, :block_time, :deposit, :rowan, :chainnet, :release_version, :from, :app_env] do |t, args|

      cluster_automation = %Q{
#!/usr/bin/env bash

echo "Checking to see if the release version exist in app.go"
check_exist=$(cat app/setupHandlers.go | grep 'release-#{args[:release_version]}')
[ -z "$check_exist" ] && exit 1 || echo "release version exists."

set +x

env_check="#{args[:app_env]}"
if [ "${env_check}" == "prod" ]; then
    export CURRENT_HEIGHT=`curl -s http://rpc.sifchain.finance/abci_info? | jq --raw-output '.result.response.last_block_height'`
else
    export CURRENT_HEIGHT=`curl -s http://rpc-#{args[:app_env]}.sifchain.finance/abci_info? | jq --raw-output '.result.response.last_block_height'`
fi

cat << EOF > pyscript.py
#!/usr/bin/env python
import os
current_height = float(os.environ["CURRENT_HEIGHT"])
block_time=#{args[:block_time]}
average_time = 60 / block_time
average_time = average_time * 60 * #{args[:upgrade_hours]}
future_block_height = average_time + current_height + 100
print(int(future_block_height))
EOF
future_block_height=$(python pyscript.py)
echo ${future_block_height}


cat << EOF > pyscript.py
#!/usr/bin/env python
import requests
import urllib3
import sys
import os
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
headers = {"Accept": "application/vnd.github.v3+json","Authorization":"token " + os.environ["GITHUB_TOKEN"]}
releases_request = requests.get('https://api.github.com/repos/Sifchain/sifnode/releases', headers=headers, verify=False)
release_request_json = releases_request.json()
find_realease="#{args[:app_env]}-#{args[:release_version]}"
retrieved_sha = ""
for release in release_request_json:
    release_name = release["name"]
    release_body = release["body"]
    published_at = release["published_at"]
    if find_realease in release["tag_name"]:
        for asset in release["assets"]:
            if ".sha256" in asset["name"]:
                get_sha = requests.get(asset["browser_download_url"], headers=headers, verify=False)
                retrieved_sha = get_sha.text.replace("\\n", "")
                print(retrieved_sha)
                break
if not retrieved_sha:
    print("retrieved sha not found ",retrieved_sha)
    sys.exit(2)
EOF
retrieved_sha=$(python pyscript.py)
echo ${retrieved_sha}

if [ "${env_check}" == "prod" ]; then
    yes "${keyring_passphrase}" | go run ./cmd/sifnodecli tx gov submit-proposal software-upgrade release-#{args[:release_version]} \
        --from #{args[:from]} \
        --deposit #{args[:deposit]} \
        --upgrade-height ${future_block_height} \
        --info '{"binaries":{"linux/amd64":"https://github.com/Sifchain/sifnode/releases/download/mainnet-#{args[:release_version]}/sifnoded-#{args[:app_env]}-#{args[:release_version]}-linux-amd64.zip?checksum='${retrieved_sha}'"}}' \
        --title release-#{args[:release_version]} \
        --description release-#{args[:release_version]} \
        --node tcp://rpc.sifchain.finance:80 \
        --keyring-backend file \
        -y \
        --chain-id #{args[:chainnet]} \
        --gas-prices "#{args[:rowan]}"
else
    yes "${keyring_passphrase}" | go run ./cmd/sifnodecli tx gov submit-proposal software-upgrade release-#{args[:release_version]} \
        --from #{args[:from]} \
        --deposit #{args[:deposit]} \
        --upgrade-height ${future_block_height} \
        --info '{"binaries":{"linux/amd64":"https://github.com/Sifchain/sifnode/releases/download/#{args[:app_env]}-#{args[:release_version]}/sifnoded-#{args[:app_env]}-#{args[:release_version]}-linux-amd64.zip?checksum='${retrieved_sha}'"}}' \
        --title release-#{args[:release_version]} \
        --description release-#{args[:release_version]} \
        --node tcp://rpc-#{args[:app_env]}.sifchain.finance:80 \
        --keyring-backend file \
        -y \
        --chain-id #{args[:chainnet]} \
        --gas-prices "#{args[:rowan]}"
fi

echo "Sleeping for 1 minute to allow for proposal to finalize"
sleep 60
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Wait for Release.."
  namespace :release do
    desc "Wait for Release."
    task :wait_for_release, [:app_env, :release] do |t, args|

      cluster_automation = %Q{
#!/usr/bin/env bash
set +x
pip install requests

cat << EOF > pyscript.py
#!/usr/bin/env python
import requests
import time
import sys
import urllib3
import os
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
headers = {"Accept": "application/vnd.github.v3+json","Authorization":"token " + os.environ["GITHUB_TOKEN"]}
workflow_request = requests.get('https://api.github.com/repos/Sifchain/sifnode/actions/workflows', headers=headers, verify=False)
workflow_request_json = workflow_request.json()
find_realease="#{args[:app_env]}-#{args[:release]}"
print("Looking for release", find_realease)
max_loop = 20
loop_count = 0
while True:
    print("You are on attempt", loop_count, " of ", max_loop)
    if loop_count >= max_loop:
        sys.exit(1)
    for workflow in workflow_request_json["workflows"]:
        if workflow["name"] == "Release":
            release_workflow_id = workflow["id"]
            workflow_info_request = requests.get('https://api.github.com/repos/Sifchain/sifnode/actions/workflows/{workflow_id}/runs'.format(workflow_id=release_workflow_id), headers=headers, verify=False)
            workflow_info_request_json = workflow_info_request.json()
            for workflow_run in workflow_info_request_json["workflow_runs"]:
                if find_realease in workflow_run["head_branch"]:
                    print(find_realease)
                    print(workflow_run["head_branch"])
                    print("Found pipeline, lets see if its done running yet.")
                    print(workflow_run)
                    if workflow_run["status"] == "completed" and workflow_run["conclusion"] == "success":
                        print("Workflow run has completed good to create governance and begin release chain.")
                        sys.exit(0)
                    else:
                        print("Release hasn't finished yet going to sleep and loop until max loops is reached waiting for Release pipeline to finish.")
            print("release not found")
            print("Sleeping for 60 seconds.")
            loop_count += 1
            time.sleep(60)
EOF
python pyscript.py
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Create Github Release."
  namespace :release do
    desc "Create Github Release."
    task :create_release, [:release, :env, :token] do |t, args|

      cluster_automation = %Q{
#!/usr/bin/env bash
set +x
pip install requests

cat << EOF > pyscript.py
#!/usr/bin/env python
import requests
import time
import sys
import urllib3
import json
import os
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
data_payload = {
    "tag_name": "#{args[:env]}-#{args[:release]}",
    "name": "#{args[:env]} v#{args[:release]}",
    "body": "Sifchain #{args[:env]} Release v#{args[:release]}",
    "prerelease": True
}
print("sending payload")
print(data_payload)
headers = {"Accept": "application/vnd.github.v3+json","Authorization":"token #{args[:token]}"}
releases_request = requests.post('https://api.github.com/repos/Sifchain/sifnode/releases',data=json.dumps(data_payload),headers=headers,verify=False)
release_request_json = releases_request.json()
if releases_request.status_code == 201 or releases_request.status_code == 200:
    print("Release Published")
    print(str(release_request_json))
elif releases_request.status_code == 422:
    print(releases_request.content, releases_request.status_code)
    print("Release already exists.")
else:
    print(releases_request.content, releases_request.status_code)
    sys.exit(2)
EOF
python pyscript.py
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Generate Test Key Ring."
  namespace :release do
    desc "Generate Test Key Ring."
    task :generate_keyring, [:moniker] do |t, args|

      cluster_automation = %Q{
#!/usr/bin/env bash
set +x
echo -e "${keyring_pem}" > tmp_keyring
tail -c +4 tmp_keyring > tmp_keyring_rendered
cat tmp_keyring_rendered | wc
rm -rf tmp_keyring
echo "moniker #{args[:moniker]}"
yes "${keyring_passphrase}" | go run ./cmd/sifnodecli keys import #{args[:moniker]} tmp_keyring_rendered --keyring-backend file
      }
      system(cluster_automation) or exit 1
    end
  end

  desc "Create Release Governance Request Vote."
  namespace :release do
    desc "Create Release Governance Request Vote."
    task :generate_vote, [:rowan, :chainnet, :from, :env] do |t, args|

      cluster_automation = %Q{
#!/usr/bin/env bash
set +x

env_check="#{args[:env]}"
if [ "${env_check}" == "prod" ]; then
    vote_id=$(go run ./cmd/sifnodecli q gov proposals --node tcp://rpc.sifchain.finance:80 --trust-node -o json | jq --raw-output 'last(.[]).id' --raw-output)
    echo "vote_id $vote_id"
    yes "${keyring_passphrase}" | go run ./cmd/sifnodecli tx gov vote ${vote_id} yes \
        --from #{args[:from]} \
        --keyring-backend file \
        --chain-id #{args[:chainnet]}  \
        --node tcp://rpc.sifchain.finance:80 \
        --gas-prices "#{args[:rowan]}" -y

    sleep 15

else
    vote_id=$(go run ./cmd/sifnodecli q gov proposals --node tcp://rpc-#{args[:env]}.sifchain.finance:80 --trust-node -o json | jq --raw-output 'last(.[]).id' --raw-output)
    echo "vote_id $vote_id"
    yes "${keyring_passphrase}" | go run ./cmd/sifnodecli tx gov vote ${vote_id} yes \
        --from #{args[:from]} \
        --keyring-backend file \
        --chain-id #{args[:chainnet]}  \
        --node tcp://rpc-#{args[:env]}.sifchain.finance:80 \
        --gas-prices "#{args[:rowan]}" -y

    sleep 15
fi
      }
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


  desc "Update Dynamic Variables For Helm Values"
  namespace :ebrelayer do
    desc "Update Dynamic Variables For Helm Values"
    task :update_helm_values, [:region, :env, :app_name, :path] do |t, args|
      cluster_automation = %Q{
#!/usr/bin/env bash
set +x
cat << EOF > helmvaulereplace.py
#!/usr/bin/env python
vaules_yaml = open("#{args[:path]}values.yaml", "r").read()
vaules_yaml = vaules_yaml.replace("-=app_name=-", "#{args[:app_name]}" )
vaules_yaml = vaules_yaml.replace("-=region=-", "#{args[:region]}" )
vaules_yaml = vaules_yaml.replace("-=env=-", "#{args[:env]}" )
open("#{args[:path]}/values.yaml", "w").write(vaules_yaml)
EOF
python helmvaulereplace.py
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

  #=======================================RUBY CONVERSIONS START=============================================#

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
        template_file_text = File.read(:template_file_name).strip
        ENV.each_pair do |k, v|
          replace_string="-=#{k}=-"
          puts replace_string
          template_file_text.include?(k) ? (template_file_text.gsub! replace_string, v) : (puts 'env matching...')
        end
        File.open(:final_file_name, 'w') { |file| file.write(template_file_text) }
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
            pod_name = `kubectl get pods -n #{ENV["APP_NAMESPACE"]} | grep #{ENV["APP_NAME"]} | cut -d ' ' -f 1`.strip
            puts "looking up logs fo #{pod_name}"
            pod_logs = `kubectl logs #{pod_name} -c ebrelayer -n #{ENV["APP_NAMESPACE"]}`
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



  #=======================================RUBY CONVERSIONS END=============================================#

  desc "Setup AWS Profile for Automation Pipelines"
  namespace :automation do
    desc "Deploy a new ebrelayer to an existing cluster"
    task :configure_aws_kube_profile, [:app_env, :aws_access_key_id, :aws_secret_access_key, :aws_region, :aws_role, :cluster_name] do |t, args|
      cluster_automation = %Q{
          set +x
          curl -s -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.19.6/2021-01-05/bin/linux/amd64/aws-iam-authenticator
          chmod +x ./aws-iam-authenticator
          export PATH=$(pwd):${PATH}
          mkdir -p ~/.aws

          echo "[default]" > ~/.aws/credentials
          echo "aws_access_key_id = #{args[:aws_access_key_id]}" >> ~/.aws/credentials
          echo "aws_secret_access_key = #{args[:aws_secret_access_key]}" >> ~/.aws/credentials
          echo "region = #{args[:aws_region]}" >> ~/.aws/credentials

          echo "[sifchain-base]" >> ~/.aws/credentials
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
