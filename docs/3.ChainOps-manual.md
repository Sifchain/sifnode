# Sifnode ChainOps

Sifnode is build on top of a kubernetes stack with the use of Ruby/Rake tasks to manage the needed tasks for managing the cluster

# Requirements

1. Ruby 2.6.x
2. Terraform 0.13.x
3. Helm 3

# Usage

1. `cd ./build`

2. `bundle install`

3. `rake -T # See the rake task options` 

4. `rake dependencies:install`

5. `rake cluster:scaffold[testnet,aws] # Scaffold a terraform config .live/` 

6. `vim ../.live/sifnode-aws-test # Edit config if required`

7. `rake cluster:deploy[testnet,aws] # Deploy cluster in aws` 

8. `rake cluster:sifnode:deploy:standalone[chainnet,provider,namespace,image,image_tag]`
