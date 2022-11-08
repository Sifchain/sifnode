# How to: Upgrade existing Peggy

## Setup

Modify the .env file to include:

```
MAINNET_URL=<mainnet_url>
MAINNET_PROXY_ADMIN_PRIVATE_KEY=<private_key>
DEPLOYMENT_NAME="sifchain-1"  
```

Where:

|Item|Description|
|----|-----------|
|`<mainnet_url>`|Replace with the Mainnet URL|
|`<private_key>`|Replace with the ETH Private Key|

## Execution

### BetaNet

1. Pull down the latest code:

```bash
git checkout develop && git pull
```

2. Switch to the `smart-contracts` directory:

```bash
cd smart-contracts
```

3. Run the upgrade:

```bash
scripts/upgrade_contracts.sh sifchain-1 mainnet
```

4. Copy the artifacts into the deployment folder:

```bash
cp -r .openzeppelin deployments/sifchain-1/
```

5. Add the files to git and push into the `develop` branch:

```bash
git add deployments/
git commit -m 'Updated contracts.'
git push origin develop
```
