def network_config(chainnet)
  "networks/#{Digest::SHA256.hexdigest chainnet}.yml"
end
