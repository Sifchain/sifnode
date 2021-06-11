namespace :openapi do
    namespace :deploy do
      desc "Deploy OpenAPI - Swagger documentation ui"
      task :swaggerui, [:chainnet, :provider, :namespace] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade swagger-ui \
                #{cwd}/../../deploy/helm/swagger-ui \
                --install -n #{ns(args)} --create-namespace \
                }

        system({"KUBECONFIG" => kubeconfig(args)}, cmd)
      end

      desc "Deploy OpenAPI - Prism Mock server "
      task :prism, [:chainnet, :provider, :namespace] do |t, args|
        check_args(args)

        cmd = %Q{helm upgrade prism \
                #{cwd}/../../deploy/helm/prism \
                --install -n #{ns(args)} --create-namespace \
                }

        system({"KUBECONFIG" => kubeconfig(args)}, cmd)
      end
    end
end
