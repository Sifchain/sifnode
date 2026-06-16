#
# Current working directory.
#
def cwd
  File.dirname(__FILE__)
end

#
# Normalize provider aliases to the canonical names we use in paths and docs.
#
def provider_id(provider)
  case provider&.downcase
  when "aws"
    "aws"
  when "az", "azure"
    "azure"
  else
    provider
  end
end

#
# Terraform template folder name for the selected cloud.
#
def template_provider_id(provider)
  provider_id(provider)
end

#
# Check the supplied arguments
#
# @param args Arguments passed to rake
#
def check_args(args)
  if args[:chainnet].nil?
    puts "Please provide a chainnet argument, for example: testnet or mainnet"
    exit 1
  end

  case provider_id(args[:provider])
  when "aws", "azure"
  else
    puts "Please provide a cloud provider. Supported values: aws, azure"
    exit 1
  end
end

#
# Generic prompt
#
# @param args Arguments passed to rake
#
def are_you_sure(args)
  return unless args[:skip_prompt].nil?

  STDOUT.puts "Are you sure? (y/n)"

  begin
    input = STDIN.gets.strip.downcase
  end until %w[y n].include?(input)

  exit(0) if input != "y"
end

#
# k8s namespace
#
# @param args Arguments passed to rake
#
def ns(args)
  args[:namespace] ? args[:namespace].to_s : "sifnode"
end

#
# Image tag
#
# @param args Arguments passed to rake
#
def image_tag(args)
  args[:image_tag] ? args[:image_tag].to_s : "testnet"
end

#
# Image repository
#
# @param args Arguments passed to rake
#
def image_repository(args)
  args[:image] ? args[:image].to_s : "sifchain/sifnoded"
end

#
# Helm values shared across deployment tasks.
#
def common_helm_flags(args)
  [
    "--set global.cloudProvider=#{provider_id(args[:provider])}",
    "--set image.tag=#{image_tag(args)}",
    "--set image.repository=#{image_repository(args)}"
  ].join(" ")
end

#
# Render a gotpl template into a file.
#
def render_template(template_path, output_path, values)
  set_args = values.map { |key, value| "--set #{key}=#{value}" }.join(" ")
  system("go run github.com/belitre/gotpl #{template_path} #{set_args} > #{output_path}") or exit 1
end

#
# Capture terraform output and persist it as the kubeconfig used by later tasks.
#
def write_kubeconfig(args)
  Dir.chdir(path(args)) do
    kubeconfig_content = `terraform output -raw kubectl_config`
    raise "terraform output -raw kubectl_config failed" unless $?.success?

    File.write(kubeconfig(args), kubeconfig_content)
  end
end
