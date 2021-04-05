# Distributes funds to these test wallets from ROWAN_SOURCE

test_wallets="sif1fpq67nw66thzmf2a5ng64cd8p8nxa5vl9d3cm4
sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
sif1hjkgsq0wcmwdh8pr3snhswx5xyy4zpgs833akh
sif1ypc5qcq5ha562xlak4xw3g6v352k39t6868jhx
sif1u7cp5e5kty8xwuu7k234ah4jsknvkzazqagvl6
sif1lj3rsayj4xtrhp2e3elv4nf7lazxty272zqegr
sif1cffgyxgvw80rr6n9pcwpzrm6v8cd6dax8x32f5"

for i in $test_wallets
do
  DESTINATION_ACCOUNT=$i python3 -m pytest --color=yes -x -olog_cli=true -olog_level=DEBUG -v -olog_file=vagrant/data/pytest.log src/py/token_distribution.py
done

