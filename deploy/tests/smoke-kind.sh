#!/bin/sh
# Destructive only to a new namespace on a disposable kind cluster. No cloud calls.
set -eu

: "${KUBECONFIG:?Set KUBECONFIG to a disposable kind cluster config}"
: "${SIFNODE_SMOKE_IMAGE:?Set SIFNODE_SMOKE_IMAGE to a kind-loaded image containing sifnoded and /bin/sh}"
case "$(kubectl config current-context)" in
    kind-*) ;;
    *) echo "This test requires a disposable kind context." >&2; exit 1 ;;
esac

chart=$(CDPATH='' cd -- "$(dirname -- "$0")/../helm/sifnode" && pwd)
namespace=sifnode-aks-smoke
temporary=$(mktemp -d)
kubectl create namespace "$namespace"
cleanup() {
    kubectl delete namespace "$namespace" --wait=false >/dev/null 2>&1 || true
    rm -rf "$temporary"
}
trap cleanup EXIT

# Create a random, unfunded, local-only test validator. Keys never leave this PVC.
cat > "$temporary/seed.yaml" <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: smoke-data
spec:
  accessModes: [ReadWriteOnce]
  storageClassName: standard
  resources:
    requests:
      storage: 1Gi
---
apiVersion: batch/v1
kind: Job
metadata:
  name: seed
spec:
  backoffLimit: 0
  template:
    spec:
      restartPolicy: Never
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
        fsGroup: 1000
      containers:
        - name: seed
          image: $SIFNODE_SMOKE_IMAGE
          imagePullPolicy: Never
          command: [/bin/sh, -ec]
          args:
            - |
              exec 3>&1
              stage=init
              trap 'echo "Local seed stopped at stage: \$stage" >&3' EXIT
              exec >/tmp/seed.log 2>&1
              sifnoded init smoke --chain-id aks-local-smoke --home /data
              stage=keys
              sifnoded keys add smoke --keyring-backend test --home /data
              stage=genesis
              address=\$(sifnoded keys show smoke -a --keyring-backend test --home /data)
              validator=\$(sifnoded keys show smoke -a --bech val --keyring-backend test --home /data)
              sifnoded add-genesis-account "\$address" 1000000000000000000000000stake,1000000000000000000000000rowan --home /data
              sifnoded add-genesis-clp-admin "\$address" --home /data
              sifnoded set-genesis-oracle-admin smoke --keyring-backend test --home /data
              sifnoded add-genesis-validators "\$validator" --home /data
              sifnoded set-genesis-whitelister-admin smoke --keyring-backend test --home /data
              stage=gentx
              sifnoded gentx smoke 100000000000000000000stake --chain-id aks-local-smoke --keyring-backend test --home /data
              sifnoded collect-gentxs --home /data
              sifnoded validate-genesis --home /data
              trap - EXIT
              cat /data/config/genesis.json >&3
          volumeMounts:
            - name: data
              mountPath: /data
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: smoke-data
EOF
kubectl -n "$namespace" apply -f "$temporary/seed.yaml"
kubectl -n "$namespace" wait --for=condition=complete job/seed --timeout=180s
kubectl -n "$namespace" logs job/seed > "$temporary/genesis.json"
python3 -c 'import json,sys; assert json.load(open(sys.argv[1]))["chain_id"] == "aks-local-smoke"' "$temporary/genesis.json"
kubectl -n "$namespace" create configmap smoke-genesis --from-file="genesis.json=$temporary/genesis.json"
kubectl -n "$namespace" delete job seed --wait=true

repository=${SIFNODE_SMOKE_IMAGE%:*}
tag=${SIFNODE_SMOKE_IMAGE##*:}
helm upgrade --install smoke "$chart" -n "$namespace" \
    --set-string "image.repository=$repository" --set-string "image.tag=$tag" \
    --set-string cloudProvider=azure --set-string chain.id=aks-local-smoke \
    --set-string genesis.existingConfigMap=smoke-genesis \
    --set-string persistence.existingClaim=smoke-data --set-string service.type=ClusterIP \
    --set-string resources.requests.cpu=100m --set-string resources.requests.memory=256Mi \
    --set-string resources.limits.memory=1Gi --wait --timeout=180s

height() {
    # Cosmos SDK's status command writes its JSON to stderr in some versions.
    kubectl -n "$namespace" exec -c sifnode smoke-sifnode-0 -- \
        sifnoded status --home /data --node tcp://127.0.0.1:26657 2>&1 |
        python3 -c 'import json,sys; d=json.load(sys.stdin); print(int(d.get("SyncInfo", d.get("sync_info", {}))["latest_block_height"]))'
}
first=$(height)
if [ "$first" -le 0 ]; then
    echo "The local chain did not produce a block." >&2
    exit 1
fi
checksum=$(kubectl -n "$namespace" exec smoke-sifnode-0 -- sha256sum /data/config/priv_validator_key.json)
kubectl -n "$namespace" delete pod smoke-sifnode-0 --wait=true
kubectl -n "$namespace" wait --for=condition=ready pod/smoke-sifnode-0 --timeout=180s
second=$(height)
after=$(kubectl -n "$namespace" exec smoke-sifnode-0 -- sha256sum /data/config/priv_validator_key.json)
test "$checksum" = "$after"
test "$second" -ge "$first"
kubectl -n "$namespace" exec smoke-sifnode-0 -- test -f /data/.initialized
echo "PASS: local sifnoded chain produced blocks (height $first -> $second) and retained its signing key across pod recreation."
helm uninstall smoke -n "$namespace" --wait
kubectl -n "$namespace" get pvc smoke-data -o name
echo "PASS: Helm uninstall retained the data PVC."
