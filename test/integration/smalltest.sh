cd /sifnode/test/integration
. /tmp/testenv.sh
script foo -c "python3 -m pytest --color=yes -olog_cli=true -olog_level=DEBUG -olog_file=vagrant/data/pytest.log -v src/py/test_eth_transfers.py"
less -Ricer foo
