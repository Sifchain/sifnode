variable "aws_region" {}
variable "cred_profile" {}
variable "cluster_subnet_id" {}
variable "cluster_security_group_id" {}
variable "profile" {}
variable "hsm_name" {}
variable "cluster_level_security_group_id" {}

provider "aws" {
  region     = var.aws_region
  shared_credentials_file = "~/.aws/credentials"
  profile                 = var.cred_profile
}

data "aws_subnet" "selected" {
  id = var.cluster_subnet_id
}

resource "aws_cloudhsm_v2_cluster" "cloudhsm_v2_cluster" {
  hsm_type   = "hsm1.medium"
  subnet_ids = data.aws_subnet.selected.*.id
  tags = {
    Name = var.hsm_name
  }
}

resource "aws_security_group_rule" "allow_all" {
  type              = "ingress"
  to_port           = 2223
  protocol          = "-1"
  from_port         = 2225
  security_group_id = aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.security_group_id
  source_security_group_id = var.cluster_security_group_id
}

resource "aws_security_group_rule" "allow_all_level" {
  type              = "ingress"
  to_port           = 2223
  protocol          = "-1"
  from_port         = 2225
  security_group_id = aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.security_group_id
  source_security_group_id = var.cluster_level_security_group_id
}

resource "aws_cloudhsm_v2_hsm" "cloudhsm_v2_hsm" {
  subnet_id  = data.aws_subnet.selected.id
  cluster_id = aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.cluster_id
  depends_on = ["aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster"]
}


resource "null_resource" "previous" {
    depends_on = ["aws_cloudhsm_v2_hsm.cloudhsm_v2_hsm"]
}

resource "time_sleep" "wait_30_seconds" {
  depends_on = [null_resource.previous]
  create_duration = "60s"
}

resource "null_resource" "echo_cluster_id" {
  provisioner "local-exec" {
    command = "echo \"${aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.cluster_id}\" > cluster_id"
  }
  depends_on = [time_sleep.wait_30_seconds]
}

resource "null_resource" "get_csr" {
  provisioner "local-exec" {
    command = "aws cloudhsmv2 describe-clusters --region ${var.aws_region} --filters clusterIds=${aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.cluster_id} --output text --query 'Clusters[].Certificates.ClusterCsr' > ${aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.cluster_id}_ClusterCsr.csr"
  }
  depends_on = [time_sleep.wait_30_seconds]
}

resource "null_resource" "generate_key" {
  provisioner "local-exec" {
    command = "openssl genrsa -aes256 -out customerCA.key 2048"
  }
  depends_on = ["null_resource.get_csr"]
}

resource "null_resource" "generate_csr" {
  provisioner "local-exec" {
    command = "openssl req -new -x509 -days 3652 -key customerCA.key -out customerCA.crt -subj '/C=US/ST=California/O=Sifchain/OU=chainOps/L=SanJose/CN=HSM'"
  }
  depends_on = ["null_resource.generate_key"]
}

resource "null_resource" "generate_cert" {
  provisioner "local-exec" {
    command = "openssl x509 -req -days 3652 -in ${aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.cluster_id}_ClusterCsr.csr -CA customerCA.crt -CAkey customerCA.key -CAcreateserial -out ${aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.cluster_id}_CustomerHsmCertificate.crt"
  }
  depends_on = ["null_resource.generate_csr"]
}

resource "null_resource" "initialize_cluster" {
  provisioner "local-exec" {
    command = "aws cloudhsmv2 initialize-cluster --region ${var.aws_region} --cluster-id ${aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.cluster_id} --signed-cert file://${aws_cloudhsm_v2_cluster.cloudhsm_v2_cluster.cluster_id}_CustomerHsmCertificate.crt --trust-anchor file://customerCA.crt"
  }
  depends_on = ["null_resource.generate_cert"]
}