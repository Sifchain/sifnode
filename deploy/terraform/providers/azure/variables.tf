variable "name"         { type = string }
variable "location"     { type = string default = "eastus" }
variable "node_count"   { type = number default = 2 }
variable "node_size"    { type = string default = "Standard_D2s_v3" }
variable "kube_version" { type = string default = "1.30" }
variable "tags"         { type = map(string) default = {} }
