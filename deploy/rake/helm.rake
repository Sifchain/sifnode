namespace :deploy do
desc "Deploy helm chart with type and type_name arguments."
task :by_type, [:app_name, :namespace, :image, :image_tag, :type, :type_name, :values_file] do |t, args|
  cmd = %Q{helm upgrade #{args[:app_name]} deploy/helm/#{args[:app_name]}-vault \
    --set vanir.args.type=#{args[:type]} \
    --set vanir.args.type_name=#{args[:type_name]} \
    --install -n #{args[:namespace]} --create-namespace \
    --set image.tag=#{args[:image_tag]} \
    --set image.repository=#{args[:image]} \
    --kubeconfig=./kubeconfig \
    -f #{args[:values_file]}
  }
  system(cmd) or exit 1
end
