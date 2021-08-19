#!/usr/bin/env bash


helpFunction()
{
   echo ""
   echo "Usage: $0 -u UpgradeName -c CurrentBinary -n NewBinary"
   echo -e "\t-u Name on the upgrade [Must match a handler defined in setup-handlers.go in NewBinary]"
   echo -e "\t-c Download link for current Binary"
   echo -e "\t-n Download link for new Binary"
   exit 1 # Exit script after printing help
}

while getopts "u:c:n:" opt
do
   case "$opt" in
      u ) UpgradeName="$OPTARG" ;;
      c ) CurrentBinary="$OPTARG" ;;
      n ) NewBinary="$OPTARG" ;;
      ? ) helpFunction ;; # Print helpFunction in case parameter is non-existent
   esac
done

if [ -z "$UpgradeName" ] || [ -z "$CurrentBinary" ] || [ -z "$NewBinary" ]
then
   echo "Some or all of the parameters are empty";
   helpFunction
fi

echo "$UpgradeName"
echo "$CurrentBinary"
echo "$NewBinary"

export DAEMON_HOME=$HOME/.sifnoded
export DAEMON_NAME=sifnoded
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true

make clean
wget $CurrentBinary -P $GOPATH/bin
chmod +x $GOPATH/bin/sifnoded
sifnoded init test --chain-id=localnet -o

echo "Generating deterministic account - sif"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnoded keys add sif --recover --keyring-backend=test

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnoded keys add akasha --recover --keyring-backend=test


sifnoded keys add mkey --multisig sif,akasha --multisig-threshold 2 --keyring-backend=test

sifnoded add-genesis-account $(sifnoded keys show sif -a --keyring-backend=test) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink,90000000000000000000ibc/96D7172B711F7F925DFC7579C6CCC3C80B762187215ABD082CDE99F81153DC80 --keyring-backend=test
sifnoded add-genesis-account $(sifnoded keys show akasha -a --keyring-backend=test) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test

sifnoded add-genesis-clp-admin $(sifnoded keys show sif -a --keyring-backend=test) --keyring-backend=test
sifnoded add-genesis-clp-admin $(sifnoded keys show akasha -a --keyring-backend=test) --keyring-backend=test

sifnoded set-genesis-whitelister-admin sif --keyring-backend=test
sifnoded set-gen-denom-whitelist scripts/denoms.json

sifnoded add-genesis-validators $(sifnoded keys show sif -a --bech val --keyring-backend=test) --keyring-backend=test

sifnoded gentx sif 1000000000000000000000000stake --chain-id=localnet --keyring-backend=test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis



mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
mkdir -p $DAEMON_HOME/cosmovisor/upgrades/$UpgradeName/bin


#wget $CurrentBinary -P $DAEMON_HOME/cosmovisor/genesis/bin
cp $GOPATH/bin/sifnoded $DAEMON_HOME/cosmovisor/genesis/bin
wget $NewBinary -P $DAEMON_HOME/cosmovisor/upgrades/$UpgradeName/bin/
chmod +x $DAEMON_HOME/cosmovisor/genesis/bin/sifnoded
chmod +x $DAEMON_HOME/cosmovisor/upgrades/$UpgradeName/bin/sifnoded

contents="$(jq '.app_state.gov.voting_params.voting_period = "10s"' $DAEMON_HOME/config/genesis.json)" && \
echo "${contents}" > $DAEMON_HOME/config/genesis.json

# Add state data here if required

cosmovisor start >> sifnode.log 2>&1  &
sleep 10
yes Y | sifnoded tx gov submit-proposal software-upgrade $UpgradeName --from sif --deposit 100000000stake --upgrade-height 20 --title $UpgradeName --description $UpgradeName --keyring-backend test --chain-id localnet
sleep 5
yes Y | sifnoded tx gov vote 1 yes --from sif --keyring-backend test --chain-id localnet
clear
sleep 5
sifnoded query gov proposal 1

tail -f sifnode.log

#Upgrade : 0.9.3
#Old : https://filebin.net/ez0139j5zzh76uba/sifnoded
#New : https://filebin.net/q67j0mjnckxgmn90/sifnoded