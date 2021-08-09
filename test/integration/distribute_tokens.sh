# Distributes funds to these test wallets from ROWAN_SOURCE

# make sure the limits for transferring ceth are high.  To use scripts/setTokenLockBurnLimit.js, we need
# the json artifacts for the current deployment put into $BASEDIR/smart-contracts/build/contracts first.

cp $BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME/BridgeBank.json $BASEDIR/smart-contracts/build/contracts
( cd $BASEDIR/smart-contracts && UPDATE_ADDRESS=0x0000000000000000000000000000000000000000 npx truffle exec scripts/setTokenLockBurnLimit.js --network ropsten 1000000000000000000000000000000 )

test_wallets="sif1fpq67nw66thzmf2a5ng64cd8p8nxa5vl9d3cm4
sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
sif1hjkgsq0wcmwdh8pr3snhswx5xyy4zpgs833akh
sif1ypc5qcq5ha562xlak4xw3g6v352k39t6868jhx
sif1u7cp5e5kty8xwuu7k234ah4jsknvkzazqagvl6
sif1lj3rsayj4xtrhp2e3elv4nf7lazxty272zqegr
sif1cffgyxgvw80rr6n9pcwpzrm6v8cd6dax8x32f5
sif1dlse3w2pxlmuvsj5eda344zp99fegual958qyr
sif1m7257566ehx7ya4ypeq7lj4y2h075z6u2xu79v
sif1qrxylp97p25wcqn4cs9nd02v672073ynpkt4yr
sif13rysrrdlhtmuc2pzve7jk0t4pwytwyxhaqcqcn
sif1shywxv2g8gvjcqknvkxu4p6lkqhfclwwj2qk6h"

for i in $test_wallets
do
  DESTINATION_ACCOUNT=$i python3 -m pytest --color=yes -x -olog_cli=true -olog_level=DEBUG -v -olog_file=vagrant/data/pytest.log src/py/token_distribution.py
done
