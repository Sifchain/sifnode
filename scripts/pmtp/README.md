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
