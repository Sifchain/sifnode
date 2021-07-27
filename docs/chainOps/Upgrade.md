# Upgrading to Cosmos 0.42 on BetaNet (External Validator Operators)

## Standalone

_These instructions assume that you're running our docker standalone solution_.

### Pre-upgrade

Sifchain will publish the time at which your node should halt, via Discord. 

1. Stop your node.

2. Checkout the latest code:

```bash
git checkout master && git pull
```

If you previously launched your standalone validator with the `rake "genesis:sifnode:boot..."` command, you'll notice that when you pull down the latest updates, your `./deploy` folder is now empty (apart from your validator config). This is intentional, as we have moved the deployment helm charts and rake tasks to another repository. 

3. To update your node so that it halts at the correct time, run (from the root of the sifnode repository):

```bash
wget -O ./deploy/docker/mainnet/docker-compose.yml https://raw.githubusercontent.com/Sifchain/sifchain-deploy-public/master/docker/mainnet/docker-compose.yml
mkdir -p ./deploy/docker/scripts
wget -O ./deploy/docker/scripts/entrypoint.sh https://raw.githubusercontent.com/Sifchain/sifchain-deploy-public/master/docker/scripts/entrypoint.sh
chmod +x ./deploy/docker/scripts/entrypoint.sh

TIMESTAMP=<timestamp> GAS_PRICE="0.5rowan" docker-compose -f ./deloy/docker/mainnet/docker-compose.yml up
```

Replace `<timestamp>` with the timestamp that is announced on Discord. The above will also replace your existing `docker-compose.yml` and `entrypoint.sh` with slightly modified versions, to allow your node to halt and be upgraded. You will see a number of warnings about values being empty when you restart your node; you can safely ignore these.

### Upgrade

Once the upgrade time has been reached, your node will automatically halt. If you attempt to restart it, it'll remain idle.

To upgrade your node, run (from the root of the sifnode repository):

```bash
CHAINNET=sifchain-1 \
UPGRADE_NODE=true \
INITIAL_HEIGHT=<height> \
COSMOS_SDK_VERSION="0.40" \
GENESIS_TIME=<time> \
VERSION="0.9.0" \
CURRENT_VERSION="0.8.8" \
DATA_MIGRATE_VERSION="0.9" \
GAS_PRICE="0.5rowan" \
docker-compose -f ./deploy/docker/mainnet/docker-compose.yml up
```

Replace `<height>` and `<time>` with what is announced on Discord.

3. Your node will upgrade itself and connect back into the network.

## k8s

### Pre-upgrade

Like the standalone solution above, Sifchain will publish the time at which your node should halt, via Discord. 

1. Checkout the latest code:

```bash
git checkout master && git pull
```

2. Checkout the new `sifnode-deploy-public` repository (from the root of where you checked out the sifnode repository to), as these rake tasks and helm charts have been moved to a new repository:

```bash
git clone ssh://git@github.com/Sifchain/sifchain-deploy-public ./deploy
```

3. To update your node so that it halts at the correct time, run (from the root of the sifnode repository):

```bash
export CLUSTER_NAME=<cluster_name>
export KUBECONFIG=./.live/${CLUSTER_NAME}/kubeconfig_${CLUSTER_NAME}

helm upgrade sifnode ./deploy/helm/standalone/sifnode \
--set sifnode.env.chainnet=sifchain \
--install -n sifnode --create-namespace \
--set sifnode.args.enableRpc="true" \
--set sifnode.args.enableExternalRpc="false" \
--set sifnode.upgrade.timestamp="<timestamp>" \
--set image.tag=mainnet-genesis \
--set image.repository=sifchain/sifnoded
```

Replace `<cluster_name>` with the full name of your cluster and `<timestamp>` with the value provided by Sifchain in Discord.

### Upgrade

1. Once the timestamp has been reached, your node will halt and sit idle (this will also be announced on Discord). You will then need to shut down the node completely:

```bash
export CLUSTER_NAME=<cluster_name>
export KUBECONFIG=./.live/${CLUSTER_NAME}/kubeconfig_${CLUSTER_NAME}

kubectl scale deployment sifnode --replicas=0 -n sifnode
```

As above, replace `<cluster_name>` with the name of your cluster.

2. To check your node has terminated, run:

```bash
kubectl get pods -n sifnode
```

3. Once terminated, you can perform the upgrade by running:

```bash
helm upgrade sifnode ./deploy/helm/standalone/sifnode \
--install -n sifnode --create-namespace \
--set sifnode.args.enableRpc="true" \
--set sifnode.args.enableExternalRpc="false" \
--set sifnode.args.upgradeNode="true" \
--set sifnode.upgrade.initialHeight="<initial_height>" \
--set sifnode.upgrade.chainnet="sifchain-1" \
--set sifnode.upgrade.cosmosSDKVersion="0.40" \
--set sifnode.upgrade.genesisTime="<genesis_time>" \
--set sifnode.upgrade.version="0.9.0" \
--set sifnode.upgrade.currentVersion="0.8.8" \
--set sifnode.upgrade.dataMigrateVersion="0.9" \
--set image.tag=mainnet-0.9.0 \
--set image.repository=sifchain/sifnoded
```

Replace `<initial_height>` and `<genesis_time>` with the values provided by Sifchain in Discord.

4. The node should then perform the upgrade (data export, migration and restart) and reconnect back into the network.
