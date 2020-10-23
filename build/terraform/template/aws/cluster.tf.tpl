// Manage state file with terraform cloud
// terraform {
//        backend "remote" {
//        hostname = "app.terraform.io"
//    }
//}

// Manage state files with s3
// terraform {
//        backend "s3" {
//        bucket = ""
//        key    = ""
//        region = "us-west-1"
//    }
//}

// Sifchain terraform module
module sifchain {
    source                  = "../../build/terraform/providers/aws"
    region                  = "us-west-1"
    cluster_name            = "sifchain-aws-{{.chainnet}}"
    tags = {
        Terraform           = true
        Project             = "sifchain"
        Chainnet            = "{{.chainnet}}"
    }
}
