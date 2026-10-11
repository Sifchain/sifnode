#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
K8S_DIR="$ROOT_DIR/kubernetes"
ENV_FILE="${SIFNODE_GKE_ENV:-$K8S_DIR/sifnode.env}"
NAMESPACE="${SIFNODE_NAMESPACE:-sifnode}"
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
  : "${GENESIS_URL:?GENESIS_URL is required}"
  : "${MNEMONIC:?MNEMONIC is required}"
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
  kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
  kubectl apply -n "$NAMESPACE" -f "$K8S_DIR/sifnode.yaml"
  kubectl apply -n "$NAMESPACE" -f "$K8S_DIR/configmap.yaml"

  kubectl create secret generic sifnode \
    --namespace "$NAMESPACE" \
    --from-literal=mnemonic="$MNEMONIC" \
    --dry-run=client -o yaml | kubectl apply -f -

  kubectl -n "$NAMESPACE" rollout status deployment/sifnode --timeout=15m
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
