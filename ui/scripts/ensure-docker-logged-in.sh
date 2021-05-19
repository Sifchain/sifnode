if  [[ ! -f  ~/.docker/config.json || $(cat ~/.docker/config.json  | jq '.auths["ghcr.io"].auth') == 'null' ]]; then
  echo "In order to run this script and push a new container to the github registry you need to create a personal access token and use it to login to ghcr with docker"
  echo ""
  echo "  echo \$MY_PAT | docker login ghcr.io -u USERNAME --password-stdin"
  echo ""
  echo "For more information see https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry"
  echo ""
  echo "Create a personal access token and log into docker using the above link then try running this script again."
  exit 1
fi
