#!/usr/bin/env bash
# =============================================================
# sifnode GKE Deployment Script
# =============================================================
# Usage:
#   ./deploy_gke.sh create              # Full deployment
#   ./deploy_gke.sh destroy             # Tear down
#   ./deploy_gke.sh update              # Update Helm release
#   ./deploy_gke.sh status              # Check deployment status
#   ./deploy_gke.sh -p <project> ...    # Specify GCP project
#   ./deploy_gke.sh -r <region> ...     # Specify region
# =============================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# ── Defaults ──────────────────────────────────────────────────
PROJECT_ID="${GCP_PROJECT_ID:-}"
REGION="${GCP_REGION:-us-central1}"
CLUSTER_NAME="${CLUSTER_NAME:-sifnode-gke}"
NAMESPACE="${NAMESPACE:-sifnode}"
DEPLOY_ENV="${DEPLOY_ENV:-gke}"
TERRAFORM_DIR="${PROJECT_DIR}/deploy/terraform/providers/gke"
HELM_CHART="${PROJECT_DIR}/deploy/helm/sifnode"
HELM_VALUES="${PROJECT_DIR}/deploy/helm/sifnode/values_gke.yaml"

# ── Colors ────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; NC='\033[0m'
info() { echo -e "${CYAN}[INFO]${NC} $1"; }
ok()   { echo -e "${GREEN}[OK]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()  { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# ── Prerequisites Check ──────────────────────────────────────
check_prereqs() {
  info "Checking prerequisites..."
  command -v gcloud >/dev/null 2>&1 || err "gcloud CLI not installed. Install: https://cloud.google.com/sdk/docs/install"
  command -v terraform >/dev/null 2>&1 || err "Terraform not installed. Install: https://learn.hashicorp.com/tutorials/terraform/install-cli"
  command -v kubectl >/dev/null 2>&1 || err "kubectl not installed."
  command -v helm >/dev/null 2>&1 || err "Helm not installed."
  
  if [[ -z "$PROJECT_ID" ]]; then
    PROJECT_ID=$(gcloud config get-value project 2>/dev/null || true)
    if [[ -z "$PROJECT_ID" ]]; then
      err "GCP project ID not set. Use -p <project> or set GCP_PROJECT_ID env var."
    fi
  fi
  ok "Prerequisites satisfied. Project: ${PROJECT_ID}, Region: ${REGION}"
}

# ── Terraform ─────────────────────────────────────────────────
tf_init() {
  info "Initializing Terraform (GKE)..."
  cd "$TERRAFORM_DIR"
  terraform init -upgrade
  ok "Terraform initialized"
}

tf_apply() {
  info "Applying Terraform (GKE)..."
  cd "$TERRAFORM_DIR"
  terraform apply -auto-approve \
    -var="project_id=${PROJECT_ID}" \
    -var="region=${REGION}" \
    -var="cluster_name=${CLUSTER_NAME}"
  ok "Terraform applied successfully"
}

tf_destroy() {
  warn "Destroying GKE infrastructure..."
  cd "$TERRAFORM_DIR"
  terraform destroy -auto-approve \
    -var="project_id=${PROJECT_ID}" \
    -var="region=${REGION}" \
    -var="cluster_name=${CLUSTER_NAME}"
  ok "GKE infrastructure destroyed"
}

# ── Kubernetes Context ───────────────────────────────────────
setup_kubectl() {
  info "Configuring kubectl..."
  gcloud container clusters get-credentials "$CLUSTER_NAME" \
    --region "$REGION" \
    --project "$PROJECT_ID"
  kubectl config set-context --current --namespace="$NAMESPACE" 2>/dev/null || true
  ok "kubectl configured for cluster: ${CLUSTER_NAME}"
}

# ── Helm Deployment ───────────────────────────────────────────
helm_install() {
  info "Installing/Upgrading sifnode via Helm..."
  kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
  
  helm upgrade --install sifnode "$HELM_CHART" \
    --namespace "$NAMESPACE" \
    --create-namespace \
    --values "${HELM_VALUES}" \
    --set provider=gke \
    --set persistence.storageClass=gke-ssd \
    --wait \
    --timeout 15m
  ok "sifnode deployed to GKE"
}

helm_uninstall() {
  info "Uninstalling sifnode Helm release..."
  helm uninstall sifnode --namespace "$NAMESPACE" 2>/dev/null || true
  kubectl delete namespace "$NAMESPACE" --ignore-not-found
  ok "sifnode uninstalled"
}

# ── Status ────────────────────────────────────────────────────
check_status() {
  info "=== GKE Cluster Status ==="
  gcloud container clusters describe "$CLUSTER_NAME" \
    --region "$REGION" \
    --project "$PROJECT_ID" \
    --format="table(name, location, status, currentMasterVersion, currentNodeVersion, currentNodeCount)"
  
  info "=== Pod Status ==="
  kubectl get pods -n "$NAMESPACE" -o wide 2>/dev/null || echo "No pods found"
  
  info "=== Service Status ==="
  kubectl get svc -n "$NAMESPACE" 2>/dev/null || echo "No services found"
  
  info "=== Persistent Volume Claims ==="
  kubectl get pvc -n "$NAMESPACE" 2>/dev/null || echo "No PVCs found"
}

# ── CLI Argument Parsing ─────────────────────────────────────
while [[ $# -gt 0 ]]; do
  case "$1" in
    -p|--project) PROJECT_ID="$2"; shift 2 ;;
    -r|--region)  REGION="$2"; shift 2 ;;
    -c|--cluster) CLUSTER_NAME="$2"; shift 2 ;;
    create|deploy)
      check_prereqs
      tf_init && tf_apply
      setup_kubectl
      helm_install
      check_status
      info "✅ sifnode deployed to GKE successfully!"
      info "   Cluster: ${CLUSTER_NAME} (${REGION})"
      info "   Namespace: ${NAMESPACE}"
      info "   To check logs: kubectl logs -n ${NAMESPACE} -l app=sifnode"
      exit 0
      ;;
    destroy|teardown)
      check_prereqs
      helm_uninstall
      tf_destroy
      info "✅ GKE deployment destroyed."
      exit 0
      ;;
    update)
      check_prereqs
      setup_kubectl
      helm_install
      check_status
      exit 0
      ;;
    status)
      check_prereqs
      setup_kubectl
      check_status
      exit 0
      ;;
    --help|-h)
      head -20 "$0" | grep "^#"
      exit 0
      ;;
    *) err "Unknown option: $1. Use --help for usage." ;;
  esac
done

# Interactive mode: show menu
echo "sifnode GKE Deployment Script"
echo "1) Full deployment (create)"
echo "2) Update only (update)"
echo "3) Destroy (destroy)"
echo "4) Status (status)"
echo "Choose (1-4): "
read -r choice
case "$choice" in
  1) $0 create ;;
  2) $0 update ;;
  3) $0 destroy ;;
  4) $0 status ;;
  *) err "Invalid choice" ;;
esac
