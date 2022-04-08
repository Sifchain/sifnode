# Running PMTP manual tests against localnet

From the root of the repo run the following commands:

```
make install
rm -rf ~/sifnoded
./scripts/init_w_prod_tokens.sh
./scripts/run.sh
```

then move to the following folder:

```
./scripts/pmtp/
```

Set the following variables:

```
# localnet
export ADMIN_ADDRESS="sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
export ADMIN_MNEMONIC="race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
export SIF_ACT=sif
export SIFNODE_CHAIN_ID=localnet
export SIFNODE_P2P_HOSTNAME=localhost
```

and now register the tokens and create 6 pools against those tokens by calling:

```
./register.sh && ./create-pools.sh
```

now calling the following commands should display all the available pools in localnet:

```
./pools.sh
```

# Test localnet using Sifgen

Setup the node with `sifgen`

```
sifgen node create sifchain-1 sifnode1 "connect rocket hat athlete kind fall auction measure wage father bridge tackle midnight athlete benefit faculty shove okay win entire reveal kit era truly" \
--admin-clp-addresses="sif1mxv2xmhm9r25cdxpwp4n43fd98t8xz97mg6xyt|sif1rkl3p87fanf8srn44lp9xrxx8smtux4mfjhwf2" \
--admin-oracle-address=sif1mxv2xmhm9r25cdxpwp4n43fd98t8xz97mg6xyt \
--standalone --with-cosmovisor
```

Setup cosmovisor:

```
export DAEMON_NAME=sifnoded
export DAEMON_HOME=$HOME/.sifnoded
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export DAEMON_RESTART_AFTER_UPGRADE=true
export UNSAFE_SKIP_BACKUP=true
```

Start the localnet chain:

```
cosmovisor start --rpc.laddr tcp://0.0.0.0:26657
```
