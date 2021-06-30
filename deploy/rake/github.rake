namespace :github do
    desc "Create Github Release."
    namespace :release_by_branch do
      desc "Create Github Release."
      task :create, [:branch, :release, :env, :token] do |t, args|
        require 'rest-client'
        require 'json'
        begin
          release_hash = { "devnet" => "DevNet", "testnet" =>"TestNet", "betanet" =>"MainNet" }
          release_target = { "devnet" => "develop", "testnet" =>"testnet", "betanet" =>"master" }
          release_name = release_hash[args[:env]]
          if "#{args[:app_env]}" == "betanet"
            headers = {content_type: :json, "Accept": "application/vnd.github.v3+json", "Authorization":"token #{args[:token]}"}
            payload = {"tag_name"  =>  "mainnet-#{args[:release]}", "target_commitish"  =>  args[:branch], "name"  =>  "#{release_name} v#{args[:release]}","body"  => "Sifchain MainNet Release v#{args[:release]}","prerelease"  =>  true}.to_json
            response = RestClient.post 'https://api.github.com/repos/Sifchain/sifnode/releases', payload, headers
            json_response_job_object = JSON.parse response.body
            puts json_response_job_object
          else
            headers = {content_type: :json, "Accept": "application/vnd.github.v3+json", "Authorization":"token #{args[:token]}"}
            payload = {"tag_name"  =>  "#{args[:env]}-#{args[:release]}", "target_commitish"  =>  args[:branch], "name"  =>  "#{release_name} v#{args[:release]}","body"  => "Sifchain #{args[:env]} Release v#{args[:release]}","prerelease"  =>  true}.to_json
            response = RestClient.post 'https://api.github.com/repos/Sifchain/sifnode/releases', payload, headers
            json_response_job_object = JSON.parse response.body
            puts json_response_job_object
          end
        rescue
          puts 'Release Already Exists'
        end
      end
    end

    desc "Create create_github_release_by_branch_and_repo."
    namespace :release_by_branch_and_repo do
      desc "Create create_github_release_by_branch_and_repo."
      task :create, [:branch, :release, :env, :token, :repo] do |t, args|
        require 'rest-client'
        require 'json'
          release_hash = { "develop" => "DevNet", "testnet" =>"TestNet", "master" =>"MainNet" }
          release_target = { "devnet" => "develop", "testnet" =>"testnet", "betanet" =>"master" }
          puts release_hash
          puts args[:env]
          puts args[:repo]
          puts args[:branch]
          puts args[:release]
          release_name = release_hash[args[:env]]
          puts "Release Name #{release_name}"
          if "#{args[:app_env]}" == "betanet"
            headers = {content_type: :json, "Accept": "application/vnd.github.v3+json", "Authorization":"token #{args[:token]}"}
            payload = {"tag_name"  =>  "mainnet-#{args[:release]}", "target_commitish"  =>  args[:branch], "name"  =>  "#{release_name} v#{args[:release]}","body"  => "#{args[:repo]} MainNet Release v#{args[:release]}","prerelease"  =>  true}.to_json
            url = "https://api.github.com/repos/Sifchain/#{args[:repo]}/releases"
            puts "github api url #{url}"
            response = RestClient.post url, payload, headers
            json_response_job_object = JSON.parse response.body
            puts json_response_job_object
          else
            headers = {content_type: :json, "Accept": "application/vnd.github.v3+json", "Authorization":"token #{args[:token]}"}
            payload = {"tag_name"  =>  "#{args[:env]}-#{args[:release]}", "target_commitish"  =>  args[:branch], "name"  =>  "#{release_name} v#{args[:release]}","body"  => "#{args[:repo]} #{args[:env]} Release v#{args[:release]}","prerelease"  =>  true}.to_json
            url = "https://api.github.com/repos/Sifchain/#{args[:repo]}/releases"
            puts "github api url #{url}"
            response = RestClient.post url, payload, headers
            json_response_job_object = JSON.parse response.body
            puts json_response_job_object
          end
      end
    end

end