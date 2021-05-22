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

#CODE TO GENERATE RANDOM STRING FOR DISPENSATION NAME AS DISPENSATION NAME IS A UNIQUE KEY
def id_generator(size=6, chars=string.ascii_uppercase + string.digits):
    return ''.join(random.choice(chars) for _ in range(size))

#CODE TO CREATE A NEW SINGLE-KEY ONLINE TXN
#FROM ADDRESS IS THE SIGNING/FUNDING ADDRESS; OUTPUT.JSON CONTAINS CLAIM RECIPIENT ADDRESSES
def create_online_singlekey_txn(
        claimType,
        dispensation_name,
        signing_address,
        chain_id,
        sifnodecli_node
    ):
    logging.debug(f"create_online_dispensation")
    sifchain_fees_entry = f"--gas 200064128"
    output = 'output.json'
    cmd = " ".join([
        "sifnodecli tx dispensation create",
        f"{dispensation_name}",
        f"{claimType}",
        output,
        sifchain_fees_entry,
        f"--from {signing_address}", 
        f"{chain_id}",
        f"{sifnodecli_node}",
        f"--yes"
        
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    txn = json_str["txhash"]
    return txn
 
#CODE TO GENERATE NEW ADDRESS    
def create_new_sifaddr_and_key():
    new_account_key = test_utilities.get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    return new_addr["address"]

#CODE TO SEND SOME SAMPLE TOKEN TO NEW ADDRESS
def send_sample_rowan(from_address,to_address,amount,keyring_backend,chain_id):
    logging.debug(f"transfer_rowan")
    sifchain_fees_entry = f"--fees 10000rowan"
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
        f"--yes"
        
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    return json_str

#CODE TO QUERY BLOCK FOR NEW DISPENSATION TXN
def query_block_claim(txnhash):
    cmd = " ".join([
        "sifnodecli q tx",
        f"{txnhash}",
    ])
    json_str = get_shell_output_json(cmd)
    return json_str

#CODE TO CHECK ACCOUNT BALANCE
def balance_check(address,currency):
    logging.debug(f"check_balance")
    cmd = " ".join([
        "sifnodecli query account",
        f"{address}",
      
    ])
    json_str = get_shell_output_json(cmd)
    amountbalance = json_str['value']['coins']
    for i in amountbalance:
        if i['denom'] == currency:
            balance = i['amount']
    return (balance)

#AUTOMATED TEST TO VALIDATE ONLINE TXN
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_create_online_singlekey_txn(claimType):
    sifchain_address = str(create_new_sifaddr_and_key())
    from_address = 'sif'
    #claimType = 'ValidatorSubsidy'
    dispensation_name = id_generator()
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '10000000rowan'
    currency = 'rowan'

    #THIS TXN IS TO REGISTER A NEW ACCOUNT ON CHAIN
    send_sample_rowan(from_address,sifchain_address,amount,keyring_backend,chain_id)
    time.sleep(5)

    #READ OUTPUT.JSON WITH CLAIMING ADDRESSES AND AMOUNT
    with open('output.json') as f:
        opjson = json.load(f)
    one_claiming_address = str(opjson['Output'][0]['address'])
    logging.info(f"one claiming address = {one_claiming_address}")
    
     #SENDER AND RECIPENT INITIAL BALANCE 
    sender_initial_balance = int(balance_check(sifchain_address,currency))
    claiming_address_initial_balance = int(balance_check(one_claiming_address,currency))
    logging.info(f"sender initial balance = {sender_initial_balance}")
    logging.info(f"one claiming address initial balance = {claiming_address_initial_balance}")
    
    try:
        #ACTUAL DISPENSATION TXN; GET TXN HASH
        txhash = str((create_online_singlekey_txn(claimType,dispensation_name,sifchain_address,chain_id,sifnodecli_node)))
        logging.info(f"txn hash = {txhash}")
        time.sleep(5)

        #QUERY BLOCK USING TXN HASH
        resp = query_block_claim(txhash)

        #READ SPECIFIC DISPENSATION TXN JSON TAGS
        distbstart = resp['logs'][0]['events'][0]['type']
        dispattb = resp['logs'][0]['events'][0]['attributes'][0]
        distype = resp['tx']['value']['msg'][0]['type']
        disvals = resp['tx']['value']['msg'][0]['value']
        account_key = str((dispattb['key']))
        bcasttags = list(disvals.keys())
        list_of_values = [disvals[key] for key in disvals]

        #DISTRIBUTION TXN JSON TAGS ASSERTIONS
        assert str(distbstart) == 'distribution_started'  
        assert str(account_key) == 'module_account'
        assert str(distype) == 'dispensation/create'
        assert bcasttags[0] == 'distributor'
        assert bcasttags[1] == 'distribution_name'
        assert bcasttags[2] == 'distribution_type'
        assert bcasttags[3] == 'Output'
        assert list_of_values[0] == sifchain_address
        assert list_of_values[1] == dispensation_name

        txn_signer_sender_address = resp['tx']['value']['msg'][0]['value']['distributor']
        distributionaddresslist = resp['tx']['value']['msg'][0]['value']['Output']
        recipient_dispensation_addresses = [] 
        amount_distributed = []
        for dic in distributionaddresslist:
            recipient_dispensation_addresses.append(dic['address']) 
            for val in dic['coins']:
                amount_distributed.append(val['amount'])

        logging.info(f"dispensation txn addresses = {recipient_dispensation_addresses}")  
        logging.info(f"amount distributed = {amount_distributed}") 
    
        total_amount_distributed = sum(int(i) for i in amount_distributed)
        recipient_with_respective_distributed_amount = dict(zip(recipient_dispensation_addresses, amount_distributed))

        logging.info(f"recipients and their respective distributed amounts = {recipient_with_respective_distributed_amount}") 
        logging.info(f"total amount distributed = {total_amount_distributed}") 

        sender_final_balance = int(balance_check(sifchain_address,currency))
        recipient_address_final_balance = int(balance_check(one_claiming_address,currency))  
    
        logging.info(f"sender initial balance = {sender_initial_balance}") 
        logging.info(f"sender final balance = {sender_final_balance}") 
    
        claimed_amount_single_recipient = int(recipient_with_respective_distributed_amount[one_claiming_address]) 
        
        #BALANCES ASSERTIONS
        assert int(total_amount_distributed) == int(sender_initial_balance - sender_final_balance)
        assert int(claimed_amount_single_recipient) == (recipient_address_final_balance - claiming_address_initial_balance) 
        logging.info(f"balance transferred from sender's address  = {(sender_initial_balance - sender_final_balance)}")  
        logging.info(f"total amount distributed  = {total_amount_distributed}")
        
        logging.info(f"amount claimed by one recipient  = {claimed_amount_single_recipient}")
        logging.info(f"balance transferred in one recipient address  = {(recipient_address_final_balance - claiming_address_initial_balance)}")

    except Exception as e:
            logging.error(f"error: {e}")   

#AUTOMTED TEST TO VALIDATE IF FUNDING ADDRESS DOESN'T HAVE ENOUGH BALANCE
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_insufficient_funds_dispensation_txn(claimType):
    sifchain_address = str(create_new_sifaddr_and_key())
    from_address = 'sif'
    claimType = 'ValidatorSubsidy'
    dispensation_name = id_generator()
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '100rowan'
    currency = 'rowan'

    #THIS TXN IS TO REGISTER A NEW ACCOUNT ON CHAIN
    send_sample_rowan(from_address,sifchain_address,amount,keyring_backend,chain_id)
    time.sleep(5)

    #READ OUTPUT.JSON WITH CLAIMING ADDRESSES AND AMOUNT
    with open('output.json') as f:
        opjson = json.load(f)
    one_claiming_address = str(opjson['Output'][0]['address'])
    logging.info(f"one claiming address = {one_claiming_address}")
    
     #SENDER AND RECIPENT INITIAL BALANCE 
    sender_initial_balance = int(balance_check(sifchain_address,currency))
    claiming_address_initial_balance = int(balance_check(one_claiming_address,currency))
    logging.info(f"sender initial balance = {sender_initial_balance}")
    logging.info(f"one claiming address initial balance = {claiming_address_initial_balance}")
    
    try:
        #ACTUAL DISPENSATION TXN; GET TXN HASH
        txhash = str((create_online_singlekey_txn(claimType,dispensation_name,sifchain_address,chain_id,sifnodecli_node)))
        logging.info(f"txn hash = {txhash}")
        time.sleep(5)

        #QUERY BLOCK USING TXN HASH
        resp = query_block_claim(txhash)
        assert resp['raw_log'] == f"for address  : {sifchain_address}: Failed in collecting funds for airdrop: failed to execute message; message index: 0"
        logging.info(f"Insufficient Funds Message = {resp['raw_log']}")
    except Exception as e:
            logging.error(f"error: {e}")

    
    
