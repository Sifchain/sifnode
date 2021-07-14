import copy
import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import create_new_sifaddr_and_credentials
from test_utilities import get_shell_output, amount_in_wei, \
    create_new_currency


def build_request_for_new_sifchain_address(basic_transfer_request, source_ethereum_address, new_currency, amount):
    sifaddress, _ = create_new_sifaddr_and_credentials()
    request = copy.deepcopy(basic_transfer_request)
    request.ethereum_symbol = new_currency["newtoken_address"]
    request.ethereum_address = source_ethereum_address
    request.sifchain_symbol = "c" + new_currency["newtoken_symbol"]
    request.sifchain_address = sifaddress
    request.amount = amount
    return request


@pytest.mark.parametrize("token_length", [7, 12])
@pytest.mark.usefixtures("operator_private_key")
def test_can_create_a_new_token_and_peg_it(
        token_length: int,
        basic_transfer_request: EthereumToSifchainTransferRequest,
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
        operator_address,
        ethereum_network,
        source_ethereum_address,
):
    new_account_key = ("a" + get_shell_output("uuidgen").replace("-", ""))[:token_length]
    token_name = new_account_key
    amount = amount_in_wei(9)
    new_currency = create_new_currency(
        amount=amount,
        symbol=new_account_key,
        token_name=token_name,
        decimals=18,
        smart_contracts_dir=smart_contracts_dir,
        bridgebank_address=bridgebank_address,
        solidity_json_path=solidity_json_path,
        operator_address=operator_address,
        ethereum_network=ethereum_network
    )
    request = build_request_for_new_sifchain_address(
        basic_transfer_request,
        source_ethereum_address,
        new_currency,
        amount / 10
    )
    burn_lock_functions.transfer_ethereum_to_sifchain(request)
