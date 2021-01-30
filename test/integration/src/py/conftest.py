import logging

import pytest

import test_utilities


@pytest.fixture
def smart_contracts_dir():
    return test_utilities.get_required_env_var("SMART_CONTRACTS_DIR")


@pytest.fixture
def integration_dir():
    return test_utilities.get_required_env_var("TEST_INTEGRATION_DIR")


@pytest.fixture
def rowan_source():
    """A sifchain address or key that has rowan and can send that rowan to other address"""
    return test_utilities.get_required_env_var("ROWAN_SOURCE")


@pytest.fixture
def sifnodecli_node():
    return test_utilities.get_optional_env_var("SIFNODE", None)


@pytest.fixture
def chain_id(is_ropsten_testnet):
    id = "sandpit" if is_ropsten_testnet else 5777
    return test_utilities.get_optional_env_var("CHAIN_ID", id)


@pytest.fixture
def ropsten_wait_time():
    return 30 * 60


@pytest.fixture
def is_ropsten_testnet(sifnodecli_node):
    """if sifnode_clinode is set, we're talking to ropsten/sandpit"""
    return sifnodecli_node


@pytest.fixture
def source_ethereum_address(is_ropsten_testnet, smart_contracts_dir):
    # account with some starting eth that can be transferred out
    if is_ropsten_testnet:
        return test_utilities.get_required_env_var("ETHEREUM_ADDRESS")
    else:
        addr = test_utilities.get_optional_env_var("ETHEREUM_ADDRESS", "")
        return addr if addr else test_utilities.ganache_owner_account(smart_contracts_dir)


@pytest.fixture(scope="function")
def ganache_timed_blocks(integration_dir):
    logging.info("restart ganache with timed blocks (keeps existing database)")
    yield test_utilities.get_shell_output(f"{integration_dir}/ganache_start.sh 2")
    logging.info("restart ganache with instant mining (keeps existing database)")
    test_utilities.get_shell_output(f"{integration_dir}/ganache_start.sh")
