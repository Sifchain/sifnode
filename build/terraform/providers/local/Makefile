# local installation of a microk8s cluster for development and testing.

/usr/local/bin/multipass:
	brew cask install multipass
	sudo multipass set local.driver=virtualbox

/usr/local/bin/microk8s: /usr/local/bin/multipass
	brew install ubuntu/microk8s/microk8s
	microk8s install

microk8s-cluster-start: /usr/local/bin/microk8s
	microk8s start

microk8s-cluster-stop: /usr/local/bin/microk8s
	microk8s stop

