desc "Install require dependencies"
namespace :dependencies do
  desc "Install dependencies for deploying to the cloud"
  task :install do
    case detect_os
    when :macosx
      puts "Installing Terraform"
      system("brew install terraform")

      puts "Installing aws-iam-authenticator"
      system("brew install aws-iam-authenticator")

      puts "Installing kubectl"
      system("brew install kubernetes-cli")

      puts "Installing Helm 3"
      system("brew install helm")
    when :linux
      # TODO complete....
      puts "Installing aws-iam-authenticator"
      system("curl -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.18.8/2020-09-18/bin/linux/amd64/aws-iam-authenticator")

      puts "Installing helm"
      system("curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash")
    else
      puts "Sorry not currently supported!"
    end
  end
end
