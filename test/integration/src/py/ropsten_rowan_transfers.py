import logging
from copy import deepcopy

import burn_lock_functions
import test_utilities
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test, create_new_sifaddr_and_credentials
from test_utilities import get_required_env_var, get_shell_output, SifchaincliCredentials, get_optional_env_var

ethereum_network = get_optional_env_var("ETHEREUM_NETWORK", "")
amount = int(get_optional_env_var("AMOUNT", "20000"))
bridgebank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")


def build_request(
        sifnodecli_node,
        source_ethereum_address,
        chain_id,
        smart_contracts_dir
) -> (EthereumToSifchainTransferRequest, SifchaincliCredentials):
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    credentials.from_key = new_addr["name"]
    ceth_fee = 2 * (10 ** 16)
    request = EthereumToSifchainTransferRequest(
        sifchain_address=new_addr["address"],
        smart_contracts_dir=smart_contracts_dir,
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=bridgebank_address,
        ethereum_network=ethereum_network,
        amount=amount + ceth_fee,
        ceth_amount=ceth_fee,
        sifnodecli_node=sifnodecli_node,
        manual_block_advance=False,
        chain_id=chain_id,
        sifchain_fees="100000rowan",
    )
    return request, credentials


def test_sample(basic_transfer_request):
    a, b = create_new_sifaddr_and_credentials()
    logging.info(f"xis: {a}")
    logging.info(f"xis: {b}")
    logging.info(f"basic trransfer: {basic_transfer_request}")


def test_transfer_rowan_to_erowan_and_back(
        ropsten_wait_time,
        rowan_source,
        sifnodecli_node,
        source_ethereum_address,
        chain_id,
        smart_contracts_dir,
        basic_transfer_request: EthereumToSifchainTransferRequest,
        ceth_fee,
        bridgetoken_address,
):
    rq: EthereumToSifchainTransferRequest = deepcopy(basic_transfer_request)
    rq.ethereum_address = "0xa584E4Fd44425937649A52846bF95A783564fCda"
    rq.ethereum_symbol = bridgetoken_address
    bx = test_utilities.get_eth_balance(rq)
    logging.info(f"bx is {bx}")
    raise Exception("stop test")
    logging.info(f"transfer rowan from {rowan_source} to a newly created account")
    sifaddr, credentials = create_new_sifaddr_and_credentials()
    rowan_transfer_from_source = deepcopy(basic_transfer_request)
    rowan_transfer_from_source.sifchain_address = rowan_source
    rowan_transfer_from_source.sifchain_destination_address = sifaddr
    amt = 20000000
    rowan_transfer_from_source.amount = amt
    rowan_transfer_from_source.sifchain_symbol = "rowan"
    test_utilities.send_from_sifchain_to_sifchain(rowan_transfer_from_source, credentials)

    logging.info(f"add ceth to new sif account to pay lock fees")

    eth_transfer: EthereumToSifchainTransferRequest = deepcopy(basic_transfer_request)
    eth_transfer.ethereum_address = source_ethereum_address
    eth_transfer.sifchain_address = sifaddr
    eth_transfer.amount = ceth_fee * 2

    logging.info("get balances just to have those commands in the history")
    try:
        test_utilities.get_sifchain_addr_balance(sifaddr, basic_transfer_request.sifnodecli_node, "ceth")
    except Exception as e:
        logging.info(f"got exception while checking balance: {e}")
    test_utilities.get_eth_balance(eth_transfer)

    logging.info("execute transfer of eth => ceth to enable fee payment")
    burn_lock_functions.transfer_ethereum_to_sifchain(eth_transfer, ropsten_wait_time)

    ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    logging.info(f"lock rowan from {rowan_transfer_from_source.sifchain_destination_address} to {ethereum_address}")

    rowan_lock: EthereumToSifchainTransferRequest = deepcopy(rowan_transfer_from_source)
    rowan_lock.sifchain_address = sifaddr
    rowan_lock.ethereum_address = ethereum_address
    burn_lock_functions.transfer_sifchain_to_ethereum(rowan_lock, credentials, ropsten_wait_time)

    logging.info(f"send erowan back to {sifaddr} from ethereum {ethereum_address}")
    return_request = deepcopy(rowan_lock)
    return_request.amount = amt / 2
    burn_lock_functions.transfer_sifchain_to_ethereum(return_request, credentials, ropsten_wait_time)
