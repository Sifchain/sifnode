desc ""
namespace :genesis do
  desc "network operations"
  namespace :network do
    desc "Scaffold a new genesis network for use in docker-compose"
    task :scaffold, [:chainnet] do |t, args|
      if args[:chainnet] == nil
        puts "Please provide chainnet argument. ie testnet, mainnet"
        exit 1
      end

      network_create({chainnet: args[:chainnet], validator_count: 4,
                      build_dir: "networks", seed_ip_address: "192.168.2.1",network_config: network_config(args[:chainnet])})
    end

    desc "Boot the new scaffolded network in docker-compose"
    task :boot, [:chainnet] do |t, args|
      if args[:chainnet] == nil
        puts "Please provide chainnet argument. ie testnet, mainnet"
        exit(1)
      end

      if !File.file?(network_config(args[:chainnet]))
        puts "the file #{network_config(args[:chainnet])} does not exist!"
        exit(1)
      end

      build_docker_image(args[:chainnet])
      boot_docker_network({chainnet: args[:chainnet], seed_network_address: "192.168.2.0/24"})
    end

    desc "Expose local seed node to the outside world"
    task :expose, [:chainnet] do |t, args|
      puts "Build me!" # TODO use something like ngrok to expose local ports of seed node to the world
    end
  end

  desc "node operations"
  namespace :sifnode do
    desc "Scaffold a new local node and configure it to connect to an existing network"
    task :scaffold, [:chainnet, :peer_address, :genesis_url] do |t, args|
      system("sifgen node create #{args[:chainnet]} #{args[:peer_address]} #{args[:genesis_url]}")
    end

    desc "boot scaffolded node and connect to existing network"
    task :boot do
      system("sifnoded start --p2p.laddr tcp://0.0.0.0:26658 ")
    end
  end
end

# Creates the config for a new network.
def network_create(chainnet:, validator_count:, build_dir:, seed_ip_address:, network_config:)
  system("sifgen network create #{chainnet} #{validator_count} #{build_dir} #{seed_ip_address} #{network_config}")
end

# Boot the new network.
def boot_docker_network(chainnet:, seed_network_address:)
  network = YAML.load_file(network_config(chainnet))

  cmd = "CHAINNET=#{chainnet} "
  network.each_with_index do |node, idx|
    cmd += "MONIKER#{idx+1}=#{node['moniker']} IPV4_ADDRESS#{idx+1}=#{node['ipv4_address']} "
  end

  cmd += "IPV4_SUBNET=#{seed_network_address} docker-compose -f ./docker-compose.yml up"
  system(cmd)
end

# Build docker image for the new network
def build_docker_image(chainnet)
  system("cd .. && docker build -f ./cmd/sifnoded/Dockerfile -t sifchain/sifnoded:#{chainnet} .")
end
