# Production Environment Setups For Scripts

The following document will describe how to set up your environment in the `smart-contracts` folder to be able to talk to betanet effectively and run operational or maintenance actions.

1). Ensure that the `ETHEREUM_PRIVATE_KEY` variable in your .env file is set to the correct role for the action that you want to do. For example, if you want to pause the bridgebank, ensure that the private key in your env file has the pauser role capability.

2). Set the `INFURA_PROJECT_ID` variable in your .env file to the correct value of the actual infura id, not the one we give in the `.env.example` file.

This should be all you need in most cases. However, some scripts may require the bridgebank or other smart contract address.

3). Set the proper smart contract address required by your script in the command line as an env variable, or set that variable in the .env file.

4). Double check all the previous steps, ensure your data is correct and you are pointed to the correct network. Then run the script.
