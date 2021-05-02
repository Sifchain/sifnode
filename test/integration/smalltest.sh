set -e 

cd /sifnode/test/integration
envfile=/tmp/tstenv.sh
python3 -u src/py/env_framework.py print_test_environment > $envfile
echo built $envfile
. $envfile
time script foo -c "python3 -m pytest --color=yes -olog_cli=false -olog_level=DEBUG -olog_file=vagrant/data/pytest.log -v src/py/test_eth_transfers.py"
less -Ricer foo
