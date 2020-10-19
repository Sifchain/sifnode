module sifchain {
    source                  = "../../build/terraform/providers/aws"
    region                  = "us-west-1"
    vpc_cidr                = "{{.aws.cidr}}"
    cluster_version         = {{.aws.cluster.version}}
    cluster_name            = "sifchain-aws-{{.chainnet}}"
    tags = {
        Terraform           = true
        Project             =  "sifchain"
        Chainnet            = "{{.chainnet}}"
    }
    node_group_settings = {
        ami_type            = "{{.aws.cluster.ami_type}}"
        desired_capacity    = {{.aws.cluster.desired_capacity}}
        max_capacity        = {{.aws.cluster.max_capacity}}
        min_capacity        = {{.aws.cluster.min_capacity}}
        instance_type       = "{{.aws.cluster.instance_type}}"
        disk_size           = {{.aws.cluster.disk_size}}
    }
}
