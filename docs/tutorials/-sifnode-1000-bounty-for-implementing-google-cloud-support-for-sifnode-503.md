# Tutorial: Implementing Google Cloud Support for Sifnode (GKE)

## Introduction
In this tutorial, we'll walk you through the process of implementing support for deploying sifnode to Google Cloud via Google Kubernetes Engine (GKE). This guide will provide detailed step-by-step instructions on how to create scaffolding and deployment scripts that enable seamless deployment. We’ll follow the existing patterns found in the AWS Terraform setup and ensure a user-friendly experience by allowing them to specify their deployment environment via CLI.

## Prerequisites
Before you begin, ensure you have the following tools installed:
- **GCloud SDK**: Install the Google Cloud SDK from [here](https://cloud.google.com/sdk/docs/install).
- **Kubernetes CLI (kubectl)**: Ensure it's properly configured and can communicate with your GKE cluster.
- **Terraform**: Version 1.0 or higher for compatibility with GKE.
- **Helm**: For deploying applications to the Kubernetes cluster.

## Step-by-Step Guide

### 1. Set Up Google Cloud Project
First, create a new project in the Google Cloud Console:
1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Select or create a project.
3. Enable billing for your project if you haven't already.

### 2. Authenticate GCloud SDK
Authenticate your local machine with the Google Cloud SDK using the following command:
```bash
gcloud auth application-default login
```

### 3. Create and Configure GKE Cluster
1. Install or update `kubectl` on your system.
2. Create a new GKE cluster using the following command:
   ```bash
   gcloud container clusters create my-cluster --zone=us-central1-a
   ```
3. Verify that the cluster was created successfully by listing all clusters:
   ```bash
   gcloud container clusters list
   ```

### 4. Set Up Terraform Configuration for GKE
Copy and paste the following template into a new file named `main.tf` in your project directory:
```hcl
provider "google" {
  credentials = file("<path_to_your_credentials_file>")
  project     = "<your_project_id>"
}

module "gke_cluster" {
  source         = "terraform-google-modules/gke/google"
  version        = "10.25.0"
  name           = "my-gke-cluster"
  location       = var.location
  num_nodes      = 3
  machine_type   = "e2-medium"
  min_cpu_platform = "automatic"
}
```
Replace `<path_to_your_credentials_file>` and `<your_project_id>` with your actual values.

### 5. Initialize Terraform and Deploy GKE Cluster
1. Initialize the Terraform configuration:
   ```bash
   terraform init
   ```
2. Plan the deployment of the cluster:
   ```bash
   terraform plan -out=tfplan
   ```
3. Apply the changes to create your GKE cluster:
   ```bash
   terraform apply tfplan
   ```

### 6. Create Helm Chart for sifnode
1. Navigate to the `sifnode` repository directory.
2. Copy the existing AWS Helm chart structure and modify it as needed for GKE.
3. Update the values.yaml file with appropriate configurations for GKE.

Example of a modified values.yaml snippet:
```yaml
gkeConfig:
  location: "us-central1-a"
  clusterName: "my-gke-cluster"
```

### 7. Deploy sifnode to GKE Using Helm
1. Install Helm if you haven't already:
   ```bash
   curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
   chmod 700 get_helm.sh
   ./get_helm.sh
   ```
2. Update the Helm repository and install sifnode:
   ```bash
   helm repo add sifchain https://charts.sifchain.finance
   helm repo update
   helm install my-sifnode sifchain/sifnode -f values.yaml
   ```

### 8. Verify Deployment
Use `kubectl` to check the status of your pods and services:
```bash
kubectl get all
```
You should see that all pods are running and ready.

## User Interface for CLI
To enable users to specify their deployment environment via a CLI, modify the sifnode script in `./scripts/deploy.sh`. Add logic to parse command-line arguments or use environment variables to set the deployment type.

Example of modified deploy.sh:
```bash
#!/bin/bash

# Default to AWS if no argument is provided
DEPLOY_ENV="AWS"

while getopts ":g" opt; do
  case ${opt} in
    g)
      DEPLOY_ENV="GKE"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done

if [ "$DEPLOY_ENV" == "GKE" ]; then
  # Run GKE deployment logic here
else
  # Run AWS deployment logic here
fi
```

## Documentation and Troubleshooting
### Basic Documentation
Provide a detailed guide on setting up sifnode on GKE, similar to the one for AWS:
- **Prerequisites**: List all necessary tools and dependencies.
- **Setup Instructions**: Follow each step outlined in this tutorial.
- **Configuration Details**: Explain how to modify values.yaml and other configuration files.

### Troubleshooting
Common issues and solutions when deploying sifnode on GKE include:
1. **Network Issues**:
   - Ensure your firewall rules allow necessary traffic.
2. **Resource Limitations**:
   - Check that you have sufficient resources (CPU, Memory) in your GKE cluster.
3. **Terraform Errors**:
   - Verify the path to your credentials file is correct and accessible.

## Conclusion
By following this comprehensive guide, you will be able to successfully deploy sifnode on Google Cloud using GKE. This process not only supports AWS but also enhances the flexibility of sifnode by adding support for another major cloud provider.