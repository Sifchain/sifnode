# some sample shell commands to transfer eth to ceth
#

export ETHEREUM_ADDRESS=$operator_address

args="
--sifchain_address $USER1ADDR
--ethereum_address $operator_address
--ethereum_symbol eth
--sifchain_symbol ceth
--amount 7000
--smart_contracts_dir=$SMART_CONTRACTS_DIR
--ethereum_chain_id=5777
--chain_id localnet
--manual_block_advance
--sifnodecli_homedir ~/.sifnodecli
--from_key user1
--keyring_backend=test
--logfile logfile.txt
--loglevel debug
"

# python3 $TEST_INTEGRATION_PY_DIR/sifchain_to_ethereum.py $args

python3 $TEST_INTEGRATION_PY_DIR/ethereum_to_sifchain.py $args

# clear && python3 -m pytest -o log_cli=true -o log_cli_level=INFO  $TEST_INTEGRATION_DIR/test_new_account.py
