import copy
import logging

import pytest

import test_utilities
from burn_lock_functions import decrease_log_level, force_log_level


@pytest.fixture
def smart_contracts_dir():
    return test_utilities.get_required_env_var("SMART_CONTRACTS_DIR")


@pytest.fixture
def validator_password():
    return test_utilities.get_optional_env_var("OWNER_PASSWORD", None)


@pytest.fixture
def validator_address():
    return test_utilities.get_optional_env_var("OWNER_ADDR", None)


@pytest.fixture
def integration_dir():
    return test_utilities.get_required_env_var("TEST_INTEGRATION_DIR")


@pytest.fixture
def bridgebank_address():
    return test_utilities.get_required_env_var("BRIDGE_BANK_ADDRESS")


@pytest.fixture
def bridgetoken_address():
    return test_utilities.get_required_env_var("BRIDGE_TOKEN_ADDRESS")


@pytest.fixture
def ethereum_network():
    return test_utilities.get_optional_env_var("ETHEREUM_NETWORK", "")


@pytest.fixture
def solidity_json_path(smart_contracts_dir):
    return test_utilities.get_optional_env_var("SOLIDITY_JSON_PATH", f"{smart_contracts_dir}/build/contracts")


@pytest.fixture
def sifnodecli_homedir(is_ropsten_testnet):
    if is_ropsten_testnet:
        base = test_utilities.get_required_env_var("HOME")
    else:
        base = test_utilities.get_required_env_var("CHAINDIR")
    result = f"""{base}/.sifnodecli"""
    return result


@pytest.fixture
def rowan_source(is_ropsten_testnet, validator_address):
    """A sifchain address or key that has rowan and can send that rowan to other address"""
    result = test_utilities.get_optional_env_var("ROWAN_SOURCE", None)
    if result:
        return result
    if is_ropsten_testnet:
        assert result
    else:
        assert validator_address
        return validator_address


@pytest.fixture
def rowan_source_key(is_ropsten_testnet):
    """A sifchain address or key that has rowan and can send that rowan to other address"""
    result = test_utilities.get_optional_env_var("ROWAN_SOURCE_KEY", None)
    if result:
        return result
    if is_ropsten_testnet:
        # Ropsten requires that you manually set the ROWAN_SOURCE_KEY environment variable
        assert result
    else:
        return test_utilities.get_required_env_var("MONIKER")


@pytest.fixture
def sifnodecli_node():
    return test_utilities.get_optional_env_var("SIFNODE", None)


@pytest.fixture
def basedir():
    return test_utilities.get_required_env_var("BASEDIR")


@pytest.fixture
def ceth_fee():
    return max(test_utilities.burn_gas_cost, test_utilities.lock_gas_cost)


@pytest.fixture
def chain_id(is_ropsten_testnet):
    result = test_utilities.get_optional_env_var("CHAINNET", None)
    if result:
        return result
    else:
        id = "sandpit" if is_ropsten_testnet else "localnet"
        return id


@pytest.fixture
def ropsten_wait_time():
    return 30 * 60


@pytest.fixture
def is_ropsten_testnet(sifnodecli_node):
    """if sifnode_clinode is set, we're talking to ropsten/sandpit"""
    return sifnodecli_node


@pytest.fixture
def is_ganache(ethereum_network):
    """true if we're using ganache"""
    return not ethereum_network


@pytest.fixture
def sifchain_fees():
    return "100000rowan"


@pytest.fixture
def operator_account(smart_contracts_dir):
    return test_utilities.get_optional_env_var("OPERATOR_ACCOUNT", test_utilities.ganache_owner_account(smart_contracts_dir))


@pytest.fixture
def source_ethereum_address(is_ropsten_testnet, smart_contracts_dir):
    """account with some starting eth that can be transferred out"""
    addr = test_utilities.get_optional_env_var("ETHEREUM_ADDRESS", "")
    if addr:
        logging.debug("using ETHEREUM_ADDRESS provided for source_ethereum_address")
        return addr
    if is_ropsten_testnet:
        # Ropsten requires that you manually set the ETHEREUM_ADDRESS environment variable
        assert addr
    result = test_utilities.ganache_owner_account(smart_contracts_dir)
    logging.debug(
        f"Using source_ethereum_address {result} from ganache_owner_account.  (Set ETHEREUM_ADDRESS env var to set it manually)")
    assert result
    return result


@pytest.fixture(scope="function")
def ganache_timed_blocks(integration_dir):
    "restart ganache with timed blocks (keeps existing database)"
    logging.info("restart ganache with timed blocks (keeps existing database)")
    yield test_utilities.get_shell_output(f"{integration_dir}/ganache_start.sh 2")
    logging.info("restart ganache with instant mining (keeps existing database)")
    test_utilities.get_shell_output(f"{integration_dir}/ganache_start.sh")


@pytest.fixture(scope="function")
def no_whitelisted_validators(integration_dir):
    """restart sifchain with no whitelisted validators, execute test, then restart with validators"""
    yield test_utilities.get_shell_output(f"ADD_VALIDATOR_TO_WHITELIST= bash {integration_dir}/setup_sifchain.sh")
    test_utilities.get_shell_output(f". {integration_dir}/vagrantenv.sh; ADD_VALIDATOR_TO_WHITELIST=true bash {integration_dir}/setup_sifchain.sh")


@pytest.fixture(scope="function")
def ensure_relayer_restart(integration_dir, smart_contracts_dir):
    """restarts relayer after the test function completes.  Used by tests that need to stop the relayer."""
    yield None
    logging.info("restart ebrelayer after advancing wait blocks - avoids any interaction with replaying blocks")
    original_log_level = decrease_log_level(new_level=logging.WARNING)
    test_utilities.advance_n_ethereum_blocks(test_utilities.n_wait_blocks + 1, smart_contracts_dir)
    test_utilities.get_shell_output(f"{integration_dir}/sifchain_start_ebrelayer.sh")
    force_log_level(original_log_level)


@pytest.fixture(scope="function")
def basic_transfer_request(
        smart_contracts_dir,
        bridgebank_address,
        bridgetoken_address,
        ethereum_network,
        ceth_fee,
        sifnodecli_node,
        chain_id,
        sifchain_fees,
        solidity_json_path,
        is_ganache,
):
    """
    Creates a EthereumToSifchainTransferRequest with all the generic fields filled in.
    """
    return test_utilities.EthereumToSifchainTransferRequest(
        smart_contracts_dir=smart_contracts_dir,
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=bridgebank_address,
        bridgetoken_address=bridgetoken_address,
        ethereum_network=ethereum_network,
        ceth_amount=ceth_fee,
        sifnodecli_node=sifnodecli_node,
        manual_block_advance=is_ganache,
        chain_id=chain_id,
        sifchain_fees=sifchain_fees,
        solidity_json_path=solidity_json_path
    )


@pytest.fixture(scope="function")
def rowan_source_integrationtest_env_credentials(
        sifnodecli_homedir,
        validator_password,
        rowan_source_key,
        is_ganache,
):
    """
    Creates a SifchaincliCredentials with all the fields filled in
    to transfer rowan from an account that already has rowan.
    """
    return test_utilities.SifchaincliCredentials(
        keyring_backend="file" if is_ganache else "test",
        keyring_passphrase=validator_password,
        from_key=rowan_source_key,
        sifnodecli_homedir=sifnodecli_homedir
    )


@pytest.fixture(scope="function")
def rowan_source_integrationtest_env_transfer_request(
        basic_transfer_request,
        rowan_source
) -> test_utilities.EthereumToSifchainTransferRequest:
    """
    Creates a EthereumToSifchainTransferRequest with all the generic fields filled in
    for a transfer of rowan from an account that already has rowan.
    """
    result: test_utilities.EthereumToSifchainTransferRequest = copy.deepcopy(basic_transfer_request)
    result.sifchain_address = rowan_source
    result.sifchain_symbol = "rowan"
    return result
