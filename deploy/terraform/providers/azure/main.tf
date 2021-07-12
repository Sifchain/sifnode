terraform {
	required_providers{
		azurerm = {
			source = "hashicorp/azurerm"
			version = ">=2.65.0"
		}

		kubernetes = {
	      source  = "hashicorp/kubernetes"
	      version = ">= 2.0"
    	}
    
	    helm = {
	      source  = "hashicorp/helm"
	      version = ">= 2.0"
	    }

	}
}

provider "azurerm"{
	features{

	}
}

provider "kubernetes" {
	host = azurerm_kubernetes_cluster.sifchain_cluster.kube_config.0.host
	username = azurerm_kubernetes_cluster.sifchain_cluster.kube_config.0.username
	password = azurerm_kubernetes_cluster.sifchain_cluster.kube_config.0.password
	client_certificate = base64decode(azurerm_kubernetes_cluster.sifchain_cluster.kube_config.0.client_certificate)
	client_key = base64decode(azurerm_kubernetes_cluster.sifchain_cluster.kube_config.0.client_key)
	cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.sifchain_cluster.kube_config.0.cluster_ca_certificate)

}

resource "azurerm_resource_group" "resource_group_1" {
	name = var.cluster_resource_group
	location = var.cluster_location
}

resource "azurerm_kubernetes_cluster" "sifchain_cluster" {
	# name = "kubernetes_cluster_1"
  name = var.cluster_name
	location = azurerm_resource_group.resource_group_1.location
	resource_group_name = azurerm_resource_group.resource_group_1.name
	dns_prefix = var.dns_prefix

	default_node_pool {
		name = "default"
		node_count = var.node_count
		#2 CPU cores, 7GB RAM 100 GiB, 2 NIC, High Bandwidth
		vm_size = var.vm_size
		os_disk_size_gb = var.disk_size
	}


	identity {
		type = "SystemAssigned"
	}

	tags = {
    Project = var.project_name
    Name = var.cluster_name
    Cluster_Version = var.cluster_version
		Environment = var.environment
	}

}

resource "azurerm_virtual_network" "sifchain_vnetwork" {
  name                = "${var.cluster_name}-${var.vnet_name}"
  resource_group_name = azurerm_resource_group.resource_group_1.name
  address_space       = [var.vpc_cidr]
  location            = azurerm_resource_group.resource_group_1.location

  depends_on = [azurerm_resource_group.resource_group_1]
}


resource "azurerm_subnet" "sifchain_subnet" {
  name = "${var.cluster_name}-${var.subnet_name}"
  resource_group_name  = azurerm_resource_group.resource_group_1.name
  virtual_network_name = azurerm_virtual_network.sifchain_vnetwork.name
  address_prefixes     = [var.vpc_cidr]

  depends_on = [azurerm_virtual_network.sifchain_vnetwork, azurerm_resource_group.resource_group_1]

}

resource "local_file" "kube_config_file"{
  filename = "kubeconfig_${azurerm_kubernetes_cluster.sifchain_cluster.name}"
  file_permission = "0711"
  content = azurerm_kubernetes_cluster.sifchain_cluster.kube_config_raw
  depends_on = [azurerm_kubernetes_cluster.sifchain_cluster]
}


# module "aks" {
#   source = "Azure/aks/azurerm"
#   resource_group_name = azurerm_resource_group.resource_group_1.name
#   cluster_name = azurerm_kubernetes_cluster.sifchain_cluster.name
#   vnet_subnet_id = azurerm_subnet.sifchain_subnet.address_prefixes[0]

#   prefix = azurerm_kubernetes_cluster.sifchain_cluster.name

#   tags = merge({ "Name" = var.cluster_name }, var.tags)
  
#   depends_on = [azurerm_kubernetes_cluster.sifchain_cluster]
# }