variable "name"         { type = string }
variable "region"       { type = string default = "us-central1" }
variable "node_count"   { type = number default = 2 }
variable "node_size"    { type = string default = "e2-standard-2" }
variable "kube_version" { type = string default = "1.30" }
variable "tags"         { type = map(string) default = {} }
