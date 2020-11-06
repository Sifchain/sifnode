desc ""
namespace :genesis do
  desc "network operations"
  namespace :network do
    desc "Scaffold a new genesis network for use in docker-compose"
    task :scaffold, [:chainnet] do |t, args|
      if args[:chainnet].nil?
        puts "Please provide chainnet argument. ie testnet, mainnet"
        exit 1
      end

      network_create({chainnet: args[:chainnet], validator_count: 4,
                      build_dir: "networks", seed_ip_address: "192.168.2.1",network_config: network_config(args[:chainnet])})
    end

    desc "Boot the new scaffolded network in docker-compose"
    task :boot, [:chainnet, :eth_address, :eth_keys, :eth_websocket] do |t, args|
      trap('SIGINT') { puts "Exiting..."; exit }

      if args[:chainnet].nil?
        puts "Please provide chainnet argument. ie testnet, mainnet"
        exit(1)
      end

      with_eth = eth_config(eth_address: args[:eth_address],
                            eth_keys: args[:eth_keys].split(" "),
                            eth_websocket: args[:eth_websocket])

      if !File.file?(network_config(args[:chainnet]))
        puts "the file #{network_config(args[:chainnet])} does not exist!"
        exit(1)
      end

      build_docker_image(args[:chainnet])
      boot_docker_network(chainnet: args[:chainnet], seed_network_address: "192.168.2.0/24", eth_config: with_eth)
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
def boot_docker_network(chainnet:, seed_network_address:, eth_config:)
  network = YAML.load_file(network_config(chainnet))

  cmd = "CHAINNET=#{chainnet} "
  network.each_with_index do |node, idx|
    cmd += "MONIKER#{idx+1}=#{node['moniker']} PASSWORD#{idx+1}=#{node['password']} IPV4_ADDRESS#{idx+1}=#{node['ipv4_address']} "
  end

  cmd += "IPV4_SUBNET=#{seed_network_address} #{eth_config} docker-compose -f ./genesis/docker-compose.yml up"
  system(cmd)
end

# Build docker image for the new network
def build_docker_image(chainnet)
  system("cd .. && docker build -f ./build/genesis/Dockerfile -t sifchain/sifnoded:#{chainnet} .")
end

# ethereum config
def eth_config(eth_address:, eth_keys:, eth_websocket:)
  config = "ETHEREUM_CONTRACT_ADDRESS=#{eth_address} "

  eth_keys.each_with_index do |address, idx|
    config += " ETHEREUM_PRIVATE_KEY#{idx+1}=#{eth_keys[idx]} "
  end

  config += "ETHEREUM_WEBSOCKET_ADDRESS=#{eth_websocket}"
  return config
end
