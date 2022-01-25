#!/bin/bash
# microtick and bitcanna contributed significantly here.
set -uxe

# set environment variables
export GOPATH=~/go
export PATH=$PATH:~/go/bin
export RPC=http://65.21.232.104:26657
export RPCN=htps://rpc.sifchain.finance:443,https://rpc-archive.sifchain.finance:443
export APPNAME=SIFNODED

# Install Gaia
go install -tags rocksdb ./...

# MAKE HOME FOLDER AND GET GENESIS
sifnoded init notional-sif-relays
wget -O ~/.sifnoded/config/genesis.json.gz https://cloudflare-ipfs.com/ipfs/QmeotEhwc67AnkHSYE53421DJAb1odHAsKLDUc7qBpXErA
cd ~/.sifnoded/config
gunzip -f genesis.json.gz
cd -


INTERVAL=1000

# GET TRUST HASH AND TRUST HEIGHT

LATEST_HEIGHT=$(curl -s $RPC/block | jq -r .result.block.header.height);
BLOCK_HEIGHT=$(($LATEST_HEIGHT-INTERVAL))
TRUST_HASH=$(curl -s "$RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)


# TELL USER WHAT WE ARE DOING
echo "TRUST HEIGHT: $BLOCK_HEIGHT"
echo "TRUST HASH: $TRUST_HASH"


# export state sync vars
export $(echo $APPNAME)_STATESYNC_ENABLE=true
export $(echo $APPNAME)_P2P_MAX_NUM_OUTBOUND_PEERS=500
export $(echo $APPNAME)_STATESYNC_RPC_SERVERS="$RPC,$RPCN"
export $(echo $APPNAME)_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export $(echo $APPNAME)_STATESYNC_TRUST_HASH=$TRUST_HASH
export $(echo $APPNAME)_P2P_PERSISTENT_PEERS="1208f890dbc1e6e40d1140ec5dbf47c2f7f745a2@35.83.206.211:26656,e2ffd994b3afe688ed789a8f582dc173818c202c@54.187.76.232:26656,612ab3973cd8566f45fc8a9929da374da54f06d1@35.81.10.98:26656,905adf4d4d515ef58fa320cb7557b8353486538e@35.160.36.158:26656,584dfa4fcf6ffb587371fc5c33ea834355b58486@34.208.32.73:26656,7e0cfb76afe681e391da5098ffa71c93088a52b4@35.83.13.80:26656,0d4981bdaf4d5d73bad00af3b1fa9d699e4d3bc0@44.235.108.41:26656,d3f068691f21c0b53a848f75a9d5479270c9eb00@34.214.39.89:26656,45022ff9b2f5f30f54601b0cdbeb770b4d79a05d@18.116.120.21:26656,663dec65b754aceef5fcccb864048305208e7eb2@34.248.110.88:26656,66cd3244dfe80e537c80024ed7a7df327e352839@34.251.129.118:26656,5f68649f085db69b32e0f66bee2f814c8b242797@52.17.217.188:26656,c74943e4882833b5d77bdd41049e03a4d602aff7@52.211.143.137:26656,bcc2d07a14a8a0b3aa202e9ac106dec0bef91fda@13.55.247.60:26656,5dec18976ca45c8b2b8a0546c874fc2468d8dc2d@176.57.184.122:26656,f33f1bb8894253daf9242cef96092686792cef80@161.97.157.60:26656,7b1ce449b4619136eb45991150a37a9962242ec6@168.119.9.103:26656,0c2d53932e7cc2bb2c44a29bc19e5aff9a167008@144.76.118.238:26656,0120f0a48e7e81cc98829ef4f5b39480f11ecd5a@52.76.185.17:26656,239baa0b73b33a1f8099c30a593539397e700240@62.171.143.207:26656,7c38e97c188c25f2d407af7b3c4d6af80360527f@85.237.192.105:26656,bd1a55a8a8bdfb747571e56ba9df518043fcc011@62.171.151.81:26656,26c306f0471032344bfa100f52b5b36071c301b6@91.230.111.86:26656,705376a038047d32a4830515139a18bb20b5b1ff@121.78.247.243:26656,ae74e0caaa799876434aeb872a75ff9d33c4db10@35.86.232.47:26656,6b37c01ae15298322273c3d7e4ec0c670962ac27@138.201.207.43:26656"
export $(echo $APPNAME)_P2P_SEEDS="4bf564ab479c860977759d050f4d42018f4bfbde@sif-seed.blockpane.com:26656"


sifnoded start --x-crisis-skip-assert-invariants --db_backend rocksdb 
