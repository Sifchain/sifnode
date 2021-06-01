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
    sifchain_fees_entry = f"--gas auto"
    output = 'output.json'
    cmd = " ".join([
        "sifnodecli tx dispensation create",
        f"{claimType}",
        output,
        f"--from {signing_address}",
        f"--chain-id={chain_id}",
        f"{sifnodecli_node}",
        f"--fees 150000rowan",
        f"--generate-only", 
        f"--yes -o json"
        
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    return json_str

#CODE TO SIGN DISPENSATION BY A USER
def sign_txn(signingaddress, offlinetx):
    keyring_backend_entry = f"--keyring-backend test"
    cmd = " ".join([
        "sifnodecli tx sign",
        f"--from {signingaddress}",
        f"{offlinetx}",
        keyring_backend_entry,
        "--chain-id localnet",
        f"--yes -o json"
    ])
    json_str = get_shell_output_json(cmd)
    return json_str


#CODE TO BROADCAST SINGLE SIGNED TXN ON BLOCK
def broadcast_txn(signedtx):
    cmd = " ".join([
        "sifnodecli tx broadcast",
        f"{signedtx}",
        f"--yes -o json"
    ])
    json_str = get_shell_output_json(cmd)
    txn = json_str["txhash"]
    return txn

#CODE TO QUERY BLOCK FOR NEW DISPENSATION TXN
def query_block_claim(txnhash):
    cmd = " ".join([
        "sifnodecli q tx",
        f"{txnhash}",
        "--chain-id localnet",
        f"-o json"
    ])
    json_str = get_shell_output_json(cmd)
    return json_str

#TEST CODE TO ASSERT TAGS GENERATED ON A BLOCK WHEN A NEW UNSIGNED DISPENSATION IS CREATED
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_create_offline_singlekey_txn(claimType):
    sifchain_address, sifchain_name = create_new_sifaddr_and_key()
    logging.info(f"sifchain_address = {sifchain_address}, sifchain_name = {sifchain_name}")
    destaddress1, destname1 = create_new_sifaddr_and_key()
    logging.info(f"destaddress1 = {destaddress1}, destname1 = {destname1}")
    destaddress2, destname2 = create_new_sifaddr_and_key()
    logging.info(f"destaddress2 = {destaddress2}, destname2 = {destname2}")
    from_address = 'sifnodeadmin'
    dispensation_name = id_generator()
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '10000000rowan'
    fee = '50000'
    currency = 'rowan'
    sampleamount = '1000rowan'

    #THESE 3 TXNS ARE TO REGISTER NEW ACCOUNTS ON CHAIN
    send_sample_rowan(from_address,sifchain_address,amount,keyring_backend,chain_id)
    time.sleep(5)
    send_sample_rowan(from_address,destaddress1,sampleamount,keyring_backend,chain_id)
    time.sleep(5)
    send_sample_rowan(from_address,destaddress2,sampleamount,keyring_backend,chain_id)
    time.sleep(5)

    #CREATING TEST DATA HERE MIMICKING OUTPUT.JSON TO BE SUPPLIED BY NIKO'S API
    dict1 = {"denom": "rowan","amount": "5000"}
    dict2 = {"denom": "rowan","amount": "7000"}
    dict3 = {"address": destaddress1,"coins":[dict1]}
    dict4 = {"address": destaddress2,"coins":[dict2]}
    dict5 = {"Output":[dict3,dict4]}
    data = json.dumps(dict5)
    with open("output.json","w") as f:
        f.write(data)

    #READ OUTPUT.JSON WITH CLAIMING ADDRESSES AND AMOUNT
    with open("output.json","r") as f:
        data = f.read()
    d = json.loads(data)
    
  
    response = (create_offline_singlekey_txn(claimType,dispensation_name,sifchain_address,chain_id,sifnodecli_node))
    
    distributiontypetag = response['value']['msg'][0]['type']
    distributionvaluetags = response['value']['msg'][0]['value']
    actuallisttags = list(distributionvaluetags.keys())
    logging.info(f"dispensation create message= {distributiontypetag}")
    logging.info(f"dispensation message tags list= {actuallisttags}")
        
    assert str(distributiontypetag) == 'dispensation/create'
    assert actuallisttags[0] == 'distributor'
    assert actuallisttags[1] == 'distribution_name'
    assert actuallisttags[2] == 'distribution_type'
    assert actuallisttags[3] == 'output'
    try:
        os.remove('output.json')
    except OSError as e:
        print ("Error: %s - %s." % (e.filename, e.strerror))
        
#TEST CODE TO ASSERT TAGS GENERATED ON A BLOCK WHEN A NEW SIGNED DISPENSATION IS BROADCASTED on BLOCKCHAIN
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_broadcast_txn(claimType):
    sifchain_address, sifchain_name = create_new_sifaddr_and_key()
    logging.info(f"sifchain_address = {sifchain_address}, sifchain_name = {sifchain_name}")
    destaddress1, destname1 = create_new_sifaddr_and_key()
    destaddress2, destname2 = create_new_sifaddr_and_key()
    from_address = 'sifnodeadmin'
    dispensation_name = id_generator()
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '10000000rowan'
    fee='50000'
    currency = 'rowan'
    sampleamount = '1000rowan'

    #THESE 3 TXNS ARE TO REGISTER NEW ACCOUNTS ON CHAIN
    send_sample_rowan(from_address,sifchain_address,amount,keyring_backend,chain_id)
    time.sleep(5)
    send_sample_rowan(from_address,destaddress1,sampleamount,keyring_backend,chain_id)
    time.sleep(5)
    send_sample_rowan(from_address,destaddress2,sampleamount,keyring_backend,chain_id)
    time.sleep(5)

    #CREATING TEST DATA HERE MIMICKING OUTPUT.JSON TO BE SUPPLIED BY NIKO'S API
    dict1 = {"denom": "rowan","amount": "5000"}
    dict2 = {"denom": "rowan","amount": "7000"}
    dict3 = {"address": destaddress1,"coins":[dict1]}
    dict4 = {"address": destaddress2,"coins":[dict2]}
    dict5 = {"Output":[dict3,dict4]}
    data = json.dumps(dict5)
    with open("output.json","w") as f:
        f.write(data)

    #READ OUTPUT.JSON WITH CLAIMING ADDRESSES AND AMOUNT
    with open("output.json","r") as f:
        data = f.read()
    d = json.loads(data)
      
    response = (create_offline_singlekey_txn(claimType,dispensation_name,sifchain_address,chain_id,sifnodecli_node))
    with open("sample.json", "w") as outfile: 
        json.dump(response, outfile)
    
    sigresponse = sign_txn(sifchain_name, 'sample.json')
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
    assert broadcasttags[3] == 'output'
    assert list_of_values[0] == sifchain_address
    assert list_of_values[1] == dispensation_name
    try:
        os.remove('signed.json')
        os.remove('sample.json')
        os.remove('sample.json')
    except OSError as e:
            print ("Error: %s - %s." % (e.filename, e.strerror)) 
