cd $BASEDIR/smart-contracts
goldmaster=$BASEDIR/ui/core/src/assets.ethereum.mainnet.json
whitelistfile=$BASEDIR/ui/core/src/tokenwhitelist.${DEPLOYMENT_NAME}.json
yarn -s --cwd $BASEDIR/smart-contracts integrationtest:whitelistedTokens --bridgebank_address $BRIDGE_BANK_ADDRESS --json_path $BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME --ethereum_network $ETHEREUM_NETWORK | grep '^\[' | jq > $whitelistfile
for chain in sifchain ethereum
do
  node scripts/test/updateAddresses.js $whitelistfile $goldmaster | jq .$chain > $BASEDIR/ui/core/src/assets.${chain}.${DEPLOYMENT_NAME}.json
done
