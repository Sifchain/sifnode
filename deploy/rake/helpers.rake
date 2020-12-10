#
# Current working directory.
#
def cwd
  File.dirname(__FILE__)
end

#
# Check the supplied arguments
#
# @param args Arguments passed to rake
#
def check_args(args)
  if args[:chainnet] == nil
    puts "Please provider a chainnet argument E.g testnet, mainnet, etc"
    exit
  end

  case args[:provider]
  when "aws"
  when "az"
    puts "Build me!"
    exit
  when "gcp"
    puts "Build me!"
    exit
  when "do"
    puts "Build me!"
    exit
  else
    puts "Please provide a cloud host provider. E.g aws"
    exit
  end
end

#
# Network config
#
# @params chainnet Name or ID of the chain
#
def network_config(chainnet)
  if chainnet == 'localnet'
    "#{cwd}/../networks/network-definition.yml"
  else
    "#{cwd}/../networks/#{Digest::SHA256.hexdigest chainnet}.yml"
  end
end

#
# Generic prompt
#
# @param args Arguments passed to rake
#
def are_you_sure(args)
  if args[:skip_prompt].nil?
    STDOUT.puts "Are you sure? (y/n)"

    begin
      input = STDIN.gets.strip.downcase
    end until %w(y n).include?(input)

    exit(0) if input != 'y'
  end
end

#
# Node address
#
# @param args Arguments passed to rake
#
def node_address(args)
  args[:node].nil? ? "tcp://127.0.0.1:26657" : args[:node]
end

#
# Sifnode Pod name
#
# @param args Arguments passed to rake
#
def pod_name(args)
  cmd = %Q{kubectl get pods --selector=app.kubernetes.io/instance=sifnode \
          -n #{args[:namespace]} -o json | jq '.items[0].metadata.name'}
  system("export KUBECONFIG=#{kubeconfig(args)}")
  `#{cmd}`.strip
end

#
# Detect the O/S
#
def detect_os
  @os ||= (
  host_os = RbConfig::CONFIG['host_os']
  case host_os
  when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
    :windows
  when /darwin|mac os/
    :macosx
  when /linux/
    :linux
  when /solaris|bsd/
    :unix
  else
    raise Error::WebDriverError, "unknown os: #{host_os.inspect}"
  end
  )
end

def safe_system(cmd)
  if (!system(cmd))
    STDERR.puts("System cmd failed: #{cmd}")
  end
end