#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
K8S_DIR="$ROOT_DIR/kubernetes"
ENV_FILE="${SIFNODE_GKE_ENV:-$K8S_DIR/sifnode.env}"
NAMESPACE="${SIFNODE_NAMESPACE:-sifnode}"
CLUSTER_NAME="${SIFNODE_GKE_CLUSTER:-sifnode-gke}"
GCP_REGION="${SIFNODE_GKE_REGION:-${GOOGLE_REGION:-}}"
GCP_PROJECT="${SIFNODE_GKE_PROJECT:-${GOOGLE_CLOUD_PROJECT:-}}"
ACTION="${1:-status}"

need() {
  command -v "$1" >/dev/null 2>&1 || { echo "missing required command: $1" >&2; exit 2; }
}

load_env() {
  [[ -f "$ENV_FILE" ]] || { echo "missing $ENV_FILE (copy sifnode.env.example)" >&2; exit 2; }
  # shellcheck disable=SC1090
  set -a; source "$ENV_FILE"; set +a
  : "${CHAINNET:?CHAINNET is required}"
  : "${MONIKER:?MONIKER is required}"
  : "${EXTERNAL_IP:?EXTERNAL_IP is required; use a reserved GCP address or the service address}"
  : "${GENESIS_URL:?GENESIS_URL is required}"
  : "${MNEMONIC:?MNEMONIC is required}"
  : "${PEER_ADDRESS:=}"
  [[ "$EXTERNAL_IP" != "CHANGE_ME" ]] || { echo "EXTERNAL_IP must be configured" >&2; exit 2; }
  [[ "$MNEMONIC" != "CHANGE_ME" ]] || { echo "MNEMONIC must be configured" >&2; exit 2; }
}

connect_cluster() {
  [[ -n "$GCP_REGION" ]] || { echo "SIFNODE_GKE_REGION/GOOGLE_REGION is required" >&2; exit 2; }
  [[ -n "$GCP_PROJECT" ]] || { echo "SIFNODE_GKE_PROJECT/GOOGLE_CLOUD_PROJECT is required" >&2; exit 2; }
  gcloud container clusters get-credentials "$CLUSTER_NAME" \
    --region "$GCP_REGION" \
    --project "$GCP_PROJECT" >/dev/null
}

validate() {
  need kubectl
  need gcloud
  need terraform
  load_env
  kubectl version --client >/dev/null
  gcloud auth list --filter=status:ACTIVE --format='value(account)' | grep -q . || {
    echo "no active gcloud account; run gcloud auth login" >&2; exit 2;
  }
  echo "GKE deployment prerequisites: OK"
}

apply() {
  validate
  connect_cluster
  kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

  kubectl create configmap sifnode-config \
    --namespace "$NAMESPACE" \
    --from-literal=chainnet="$CHAINNET" \
    --from-literal=moniker="$MONIKER" \
    --from-literal=external_ip="$EXTERNAL_IP" \
    --from-literal=peer_address="$PEER_ADDRESS" \
    --from-literal=genesis_url="$GENESIS_URL" \
    --dry-run=client -o yaml | kubectl apply -f -

  kubectl create secret generic sifnode \
    --namespace "$NAMESPACE" \
    --from-literal=mnemonic="$MNEMONIC" \
    --dry-run=client -o yaml | kubectl apply -f -

  kubectl apply -n "$NAMESPACE" -f "$K8S_DIR/sifnode.yaml"
  kubectl -n "$NAMESPACE" rollout status deployment/sifnode --timeout=15m
  kubectl -n "$NAMESPACE" wait --for=condition=available deployment/sifnode --timeout=2m
  echo "Sifnode is running on GKE"
}

status() {
  need kubectl
  kubectl -n "$NAMESPACE" get deployment,pod,svc,pvc -l app=sifnode -o wide
}

destroy() {
  need kubectl
  if [[ "${CONFIRM_DESTROY:-}" != "yes" ]]; then
    echo "Refusing destructive action. Re-run with CONFIRM_DESTROY=yes." >&2
    exit 2
  fi
  kubectl delete namespace "$NAMESPACE" --ignore-not-found
}

case "$ACTION" in
  validate) validate ;;
  apply) apply ;;
  status) status ;;
  destroy) destroy ;;
  *) echo "usage: $0 {validate|apply|status|destroy}" >&2; exit 2 ;;
esac
