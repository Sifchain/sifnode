variable "name"        { type = string }
variable "provider"    { type = string  # aws | azure | gcp
  validation { condition = contains(["aws","azure","gcp"], var.provider)
    error_message = "provider must be aws, azure or gcp" } }
variable "region"      { type = string }
variable "node_count"  { type = number default = 2 }
variable "node_size"   { type = string  # e.g. t3.medium / Standard_D2s_v3 / e2-standard-2 }
variable "kube_version"{ type = string default = "1.30" }
variable "tags"        { type = map(string) default = {} }
