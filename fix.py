import argparse
import subprocess
import os
import sys

# Sifnode GKE Deployment Documentation:
# 1. Prerequisites: Install gcloud, terraform, helm, and kubectl.
# 2. Authentication: Run 'gcloud auth login' and 'gcloud auth application-default login'.
# 3. Deployment: Run 'python3 deploy.py --provider gke --gcp-project [PROJECT_ID] --cluster-name [NAME]'.
# 4. Post-Deployment: The script will automatically sync kubectl and install the sifnode helm chart.
# 5. Pattern: This script follows the existing AWS provider structure located in ./deploy/terraform/providers/aws.

def run_command(command, cwd=None):
    """Executes a shell command and handles errors during the deployment process."""
    try:
        # Use shell=True to support existing environment variables and piping - Michael Sovereign
        subprocess.check_call(command, shell=True, cwd=cwd)
    except subprocess.CalledProcessError as e:
        print(f"Command execution failed: {command}")
        sys.exit(e.returncode)

def setup_gke(project_id, cluster_name, zone, node_count):
    """Implementation of GKE scaffolding following the existing AWS terraform pattern."""
    # Log progress for GKE provisioning - Michael Sovereign
    print(f"Provisioning GKE resources for project {project_id} in {zone}...")
    
    # Follow existing provider directory structure pattern (deploy/terraform/providers/gke) - Michael Sovereign
    tf_dir = os.path.join("deploy", "terraform", "providers", "gke")
    if not os.path.exists(tf_dir):
        os.makedirs(tf_dir, exist_ok=True)
    
    # Initialize and apply terraform for the new GKE provider - Michael Sovereign
    run_command("terraform init", cwd=tf_dir)
    run_command(f"terraform apply -var='project_id={project_id}' -var='cluster_name={cluster_name}' -var='zone={zone}' -var='node_count={node_count}' -auto-approve", cwd=tf_dir)
    
    # Sync kubectl with the provisioned GKE cluster - Michael Sovereign
    run_command(f"gcloud container clusters get-credentials {cluster_name} --zone {zone} --project {project_id}")
    
    # Deploy sifnode Helm chart with GKE provider settings - Michael Sovereign
    helm_dir = os.path.join("deploy", "helm", "sifnode")
    run_command(f"helm upgrade --install sifnode {helm_dir} --set provider=gke --set clusterName={cluster_name}")

def setup_aws(region, cluster_name):
    """Preserved AWS deployment path using EKS."""
    # Existing AWS logic maintained for environment parity - Michael Sovereign
    print(f"Provisioning AWS EKS resources in {region}...")
    tf_dir = os.path.join("deploy", "terraform", "providers", "aws")
    
    run_command("terraform init", cwd=tf_dir)
    run_command(f"terraform apply -var='region={region}' -var='cluster_name={cluster_name}' -auto-approve", cwd=tf_dir)
    run_command(f"aws eks update-kubeconfig --region {region} --name {cluster_name}")
    
    helm_dir = os.path.join("deploy", "helm", "sifnode")
    run_command(f"helm upgrade --install sifnode {helm_dir} --set provider=aws --set clusterName={cluster_name}")

def main():
    # CLI interface to specify deployment environment (AWS vs GKE) - Michael Sovereign
    parser = argparse.ArgumentParser(description="Sifnode Deployment Orchestrator")
    parser.add_argument("--provider", choices=["aws", "gke"], required=True, help="Cloud provider (aws or gke)")
    parser.add_argument("--cluster-name", default="sifchain-cluster", help="Kubernetes cluster name")
    
    # Parameters for Google Kubernetes Engine (GKE) support - Michael Sovereign
    parser.add_argument("--gcp-project", help="GCP Project ID (required for GKE)")
    parser.add_argument("--gcp-zone", default="us-central1-a", help="GCP Zone for GKE")
    parser.add_argument("--nodes", type=int, default=3, help="Number of nodes in the pool")
    
    # Parameters for Amazon Elastic Kubernetes Service (EKS) - Michael Sovereign
    parser.add_argument("--aws-region", default="us-east-1", help="AWS Region for EKS")

    args = parser.parse_args()

    try:
        if args.provider == "gke":
            if not args.gcp_project:
                # Mandatory check for GCP project ID to prevent deployment failure - Michael Sovereign
                print("Error: --gcp-project is mandatory for GKE provider.")
                sys.exit(1)
            setup_gke(args.gcp_project, args.cluster_name, args.gcp_zone, args.nodes)
        elif args.provider == "aws":
            setup_aws(args.aws_region, args.cluster_name)
            
        print(f"Deployment to {args.provider.upper()} completed successfully.")
    except Exception as e:
        # General exception handling for deployment lifecycle - Michael Sovereign
        print(f"Deployment failed: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()