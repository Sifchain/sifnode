## Validation performed for this change

- `bash -n deploy/gke/deploy_gke.sh` passes locally.
- The pull-request workflow runs `terraform fmt -check`, `terraform init -backend=false`, and `terraform validate` on GitHub Actions.
- A live GKE deployment was not claimed in this environment because no Google Cloud project credentials are available to the execution environment.

This is intentionally recorded rather than claiming a live-cloud result that was not observed.
