import logging
import os
import time
import json
import pytest
import string
import random
from dispensation_envutils import create_offline_singlekey_txn, create_new_sifaddr_and_key, send_sample_rowan, balance_check, \
     query_block_claim,sign_txn,broadcast_txn,broadcast_async_txn


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
    
  
    response = (create_offline_singlekey_txn(claimType,sifchain_address,chain_id,sifnodecli_node))
    
    distributiontypetag = response['value']['msg'][0]['type']
    distributionvaluetags = response['value']['msg'][0]['value']
    actuallisttags = list(distributionvaluetags.keys())
    logging.info(f"dispensation create message= {distributiontypetag}")
    logging.info(f"dispensation message tags list= {actuallisttags}")
        
    assert str(distributiontypetag) == 'dispensation/create'
    assert actuallisttags[0] == 'distributor'
    #assert actuallisttags[1] == 'distribution_name'
    assert actuallisttags[1] == 'distribution_type'
    assert actuallisttags[2] == 'output'
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
      
    response = (create_offline_singlekey_txn(claimType,sifchain_address,chain_id,sifnodecli_node))
    with open("sample.json", "w") as outfile: 
        json.dump(response, outfile)
    
    sigresponse = sign_txn(sifchain_name, 'sample.json')
    with open("signed.json", "w") as sigfile: 
        json.dump(sigresponse, sigfile)
      
    txhashbcast = broadcast_async_txn('signed.json')
    time.sleep(5)
    resp = query_block_claim(txhashbcast)
    distypebcast = (resp['tx']['value']['msg'][0]['type'])
    disvalsbcast = (resp['tx']['value']['msg'][0]['value'])
    list_of_values = [disvalsbcast[key] for key in disvalsbcast]
    broadcasttags = list(disvalsbcast.keys())
    assert str(distypebcast) == 'dispensation/create'
    assert broadcasttags[0] == 'distributor'
    assert broadcasttags[1] == 'distribution_type'
    assert broadcasttags[2] == 'output'
    assert list_of_values[0] == sifchain_address
    try:
        os.remove('signed.json')
        os.remove('sample.json')
        os.remove('output.json')
    except OSError as e:
            print ("Error: %s - %s." % (e.filename, e.strerror)) 


