import logging
import os

# to test against ropsten, define:
# ETHEREUM_ADDRESS
# ETHEREUM_PRIVATE_KEY
# ETHEREUM_NETWORK = ropsten
#
import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, get_shell_output


@pytest.mark.skip(reason="not implemented")
def test_transfer_erowan_to_rowan():
    assert False


@pytest.mark.skip(reason="not implemented")
def test_transfer_usdt_to_cusdt():
    assert False


@pytest.mark.skip(reason="not implemented")
def test_transfer_cusdt_to_usdt():
    assert False