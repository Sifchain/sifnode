import logging
import os
import time
import json
import pytest
import string
import random

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
import test_utilities
from pytest_utilities import generate_test_account
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, SifchaincliCredentials, get_optional_env_var, ganache_owner_account, \
    get_shell_output_json, get_shell_output, detect_errors_in_sifnodecli_output, get_transaction_result, amount_in_wei

#CODE TO CREATE A NEW CLAIM
def create_claim(
        sifchain_address,
        claimType,
        keyring_backend,
        chain_id,
        sifnodecli_node
    ):
    logging.debug(f"create_claim")
    keyring_backend_entry = f"--keyring-backend {keyring_backend}"     
    sifchain_fees_entry = f"--fees 100000rowan"
    cmd = " ".join([
        "sifnodecli tx dispensation claim",
        f"{claimType}",
        f"--from {sifchain_address}",
        sifchain_fees_entry,
        keyring_backend_entry,
        f"--yes -o json" 
        
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    txn = json_str["txhash"]
    return txn

#CODE TO QUERY BLOCK FOR NEW CLAIM TXN
def query_block_claim(txnhash):
    cmd = " ".join([
        "sifnodecli q tx",
        f"{txnhash}",
    ])
    json_str = get_shell_output_json(cmd)
    return json_str

#CODE TO QUERY A NEW CLAIM 
def query_created_claim(claimType):
    cmd = " ".join([
        "sifnodecli q dispensation claims-by-type",
        f"{claimType}",
    ])
    json_str = get_shell_output_json(cmd)
    return json_str
#CODE TO GENERATE NEW ADDRESS    
def create_new_sifaddr_and_key():
    new_account_key = test_utilities.get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    return new_addr["address"]

#CODE TO SEND SOME SAMPLE TOKEN TO NEW ADDRESS
def send_sample_rowan(from_address,to_address,amount,keyring_backend,chain_id):
    logging.debug(f"transfer_rowan")
    sifchain_fees_entry = f"--fees 150000rowan"
    keyring_backend_entry = f"--keyring-backend {keyring_backend}"     
    output = 'output.json'
    cmd = " ".join([
        "sifnodecli tx send",
        f"{from_address}",
        f"{to_address}",
        f"{amount}",
        keyring_backend_entry,
        sifchain_fees_entry,
        f"--chain-id={chain_id}",
        f"--yes -o json"
        
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    return json_str

#TEST CODE TO ASSERT TAGS GENERATED ON A BLOCK WHEN A NEW CLAIM IS CREATED
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_create_new_claim(claimType):
    sifchain_address = str(create_new_sifaddr_and_key())
    keyring_backend = 'test'
    chain_id = 'localnet'
    from_address = 'sifnodeadmin'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '10000000rowan'
    send_sample_rowan(from_address,sifchain_address,amount,keyring_backend,chain_id)
    time.sleep(5)
    txnhash = (create_claim(sifchain_address,claimType,keyring_backend,chain_id,sifnodecli_node))
    time.sleep(5)
    response = (query_block_claim(str(txnhash)))
    try:
        data = (response['logs'][0]['events'][1]['attributes'])
        expectedOutputTagsList = []
        for value in data:
            expectedOutputTagsList.append(value['key'])
            expectedOutputTagsList.append(value['value'])
        print(txnhash)
        assert response['txhash'] == txnhash
        assert expectedOutputTagsList[0] == 'userClaim_creator'
        assert expectedOutputTagsList[2] == 'userClaim_type'
        assert expectedOutputTagsList[3] == claimType
        assert expectedOutputTagsList[4] == 'userClaim_creationTime'
    except KeyError:
        with pytest.raises(Exception, match='User trying to create a duplicate claim'):
            raise Exception

#TEST CODE TO ASSERT TAGS RETURNED BY A CLAIM QUERY COMMAND
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_query_created_claim(claimType):
    queryresponse = query_created_claim(claimType)
    queryresponse = query_created_claim(claimType)
    querydata = (queryresponse['claims'][0])
    queryexpectedtags = list(querydata.keys())
    assert queryexpectedtags[0] == 'user_address'
    assert queryexpectedtags[1] == 'user_claim_type'
    assert queryexpectedtags[2] == 'user_claim_time'
    
