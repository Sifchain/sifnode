import logging

import pytest

import burn_lock_functions
import integration_env_credentials
import test_utilities
from burn_lock_functions import EthereumToSifchainTransferRequest
from pytest_utilities import generate_test_account
from test_utilities import SifchaincliCredentials


def test_update_whitelist_validator_command(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        validator_address,
        sifchain_admin_account_credentials,
):
    admin_account = test_utilities.get_required_env_var("SIFCHAIN_ADMIN_ACCOUNT")
    basic_transfer_request.sifchain_address = validator_address
    admin_user_credentials = sifchain_admin_account_credentials
    test_utilities.update_whitelist_validator(
        admin_account=admin_account,
        validator_account="sifvaloper18qcnjcy3hrp6svzmxaegh3vz96vwwn4a42c9z4",
        transfer_request=basic_transfer_request,
        credentials=admin_user_credentials
    )
