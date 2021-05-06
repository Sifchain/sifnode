import logging
import os
import time
import json
import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
import test_utilities
from pytest_utilities import generate_test_account
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
        f"--yes", 
        
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    txn = json_str["txhash"]
    #logging.debug(f"resulting tx: {tx}")
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

#TEST CODE TO ASSERT TAGS GENERATED ON A BLOCK WHEN A NEW CLAIM IS CREATED
def test_create_new_claim():
    sifchain_address = 'akasha'
    claimType = 'ValidatorSubsidy'
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    txnhash = (create_claim(sifchain_address,claimType,keyring_backend,chain_id,sifnodecli_node))
    time.sleep(5)
    response = (query_block_claim(str(txnhash)))
    print(response['raw_log'])
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
def test_query_created_claim():
    claimType = 'ValidatorSubsidy'
    queryresponse = query_created_claim(claimType)
    queryresponse = query_created_claim(claimType)
    querydata = (queryresponse['claims'][0])
    queryexpectedtags = list(querydata.keys())
    assert queryexpectedtags[0] == 'user_address'
    assert queryexpectedtags[1] == 'user_claim_type'
    assert queryexpectedtags[2] == 'user_claim_time'
    assert queryexpectedtags[3] == 'locked'



    

    

