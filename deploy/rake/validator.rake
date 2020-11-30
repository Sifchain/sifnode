desc "validator operations"
namespace :validator do
  desc "Stake a node so it can participate in consensus"
  task :stake, [:chainnet, :provider, :pod, :namespace] do |t, args|
    cmd = %()
    system({"KUBECONFIG" => kubeconfig(args)}, cmd)
  end
end
