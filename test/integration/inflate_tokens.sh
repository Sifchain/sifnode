# Runs through all the steps of creating and distributing tokens in a fresh environment

runcmd="python3 -m pytest -olog_cli=true -olog_level=DEBUG -v -olog_file=vagrant/data/pytest.log"
TOKENS_FILE=$BASEDIR/sifnode/ui/core/src/assets.sifchain.mainnet.json $runcmd src/py/token_setup.py
TOKEN_AMOUNT=120000000000 $runcmd src/py/token_refresh.py
bash ./distribute_tokens.sh
