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
from test_dispensation_onlinetxn import create_new_sifaddr_and_key, send_sample_rowan
from test_utilities import get_required_env_var, SifchaincliCredentials, get_optional_env_var, ganache_owner_account, \
    get_shell_output_json, get_shell_output, detect_errors_in_sifnodecli_output, get_transaction_result, amount_in_wei

#CODE TO GENERATE RANDOM STRING FOR DISPENSATION NAME
def id_generator(size=6, chars=string.ascii_uppercase + string.digits):
    return ''.join(random.choice(chars) for _ in range(size))

#CODE TO GENERATE OFFLINE DISPENSATION TXN
def create_offline_singlekey_txn(
        claimType,
        dispensation_name,
        signing_address,
        chain_id,
        sifnodecli_node
    ):
    logging.debug(f"create_unsigned_offline_dispensation_txn")
    sifchain_fees_entry = f"--gas 200064128"
    output = 'output.json'
    cmd = " ".join([
        "sifnodecli tx dispensation create",
        f"{dispensation_name}",
        f"{claimType}",
        output,
        sifchain_fees_entry,
        f"--from {signing_address}",
        f"--generate-only", 
        
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    return json_str

#CODE TO SIGN DISPENSATION BY A USER
def sign_txn(signingaddress, offlinetx):
    cmd = " ".join([
        "sifnodecli tx sign",
        f"--from {signingaddress}",
        f"{offlinetx}"
    ])
    json_str = get_shell_output_json(cmd)
    return json_str


#CODE TO BROADCAST SINGLE SIGNED TXN ON BLOCK
def broadcast_txn(signedtx):
    cmd = " ".join([
        "sifnodecli tx broadcast",
        f"{signedtx}"
    ])
    json_str = get_shell_output_json(cmd)
    txn = json_str["txhash"]
    return txn

#CODE TO QUERY BLOCK FOR NEW DISPENSATION TXN
def query_block_claim(txnhash):
    cmd = " ".join([
        "sifnodecli q tx",
        f"{txnhash}",
    ])
    json_str = get_shell_output_json(cmd)
    return json_str

#TEST CODE TO ASSERT TAGS GENERATED ON A BLOCK WHEN A NEW UNSIGNED DISPENSATION IS CREATED
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_create_offline_singlekey_txn(claimType):
    sifchain_address = str(create_new_sifaddr_and_key())
    from_address = 'sif'
    dispensation_name = id_generator()
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '10000000rowan'
    currency = 'rowan'
    send_sample_rowan(from_address,sifchain_address,amount,keyring_backend,chain_id)
    time.sleep(5)

    response = (create_offline_singlekey_txn(claimType,dispensation_name,sifchain_address,chain_id,sifnodecli_node))
    print(response)
    with open("sample.json", "w") as outfile: 
        json.dump(response, outfile)
    try:
        distype = response['value']['msg'][0]['type']
        imptags = response['value']['msg'][0]['value']
        actuallisttags = list(imptags.keys())
        logging.info(f"dispensation create message= {distype}")
        logging.info(f"dispensation message tags list= {actuallisttags}")
        
        assert str(distype) == 'dispensation/create'
        assert actuallisttags[0] == 'distributor'
        assert actuallisttags[1] == 'distribution_name'
        assert actuallisttags[2] == 'distribution_type'
        assert actuallisttags[3] == 'Output'
        try:
            os.remove('sample.json')
        except OSError as e:
            print ("Error: %s - %s." % (e.filename, e.strerror))

    except Exception as e:
            logging.error(f"error: {e}")   

#TEST CODE TO ASSERT TAGS GENERATED ON A BLOCK WHEN A NEW SIGNED DISPENSATION IS BROADCASTED on BLOCKCHAIN
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_broadcast_txn(claimType):
    sifchain_address = str(create_new_sifaddr_and_key())
    from_address = 'sif'
    dispensation_name = id_generator()
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '10000000rowan'
    currency = 'rowan'

    send_sample_rowan(from_address,sifchain_address,amount,keyring_backend,chain_id)
    time.sleep(5)
    
    response = (create_offline_singlekey_txn(claimType,dispensation_name,sifchain_address,chain_id,sifnodecli_node))
    with open("sample.json", "w") as outfile: 
        json.dump(response, outfile)
    try:
        sigresponse = sign_txn(sifchain_address, 'sample.json')
        with open("signed.json", "w") as sigfile: 
            json.dump(sigresponse, sigfile)
        
        txhashbcast = broadcast_txn('signed.json')
        time.sleep(5)
        resp = query_block_claim(txhashbcast)
        distypebcast = (resp['tx']['value']['msg'][0]['type'])
        disvalsbcast = (resp['tx']['value']['msg'][0]['value'])
        list_of_values = [disvalsbcast[key] for key in disvalsbcast]
        broadcasttags = list(disvalsbcast.keys())
        assert str(distypebcast) == 'dispensation/create'
        assert broadcasttags[0] == 'distributor'
        assert broadcasttags[1] == 'distribution_name'
        assert broadcasttags[2] == 'distribution_type'
        assert broadcasttags[3] == 'Output'
        assert list_of_values[0] == sifchain_address
        assert list_of_values[1] == dispensation_name
        try:
            os.remove('signed.json')
            os.remove('sample.json')
        except OSError as e:
            print ("Error: %s - %s." % (e.filename, e.strerror))

    except Exception as e:
            logging.error(f"error: {e}")
