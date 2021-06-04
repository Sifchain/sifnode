import logging
import os
import time
import json
import pytest
import string
import random
from dispensation_envutils import create_online_singlekey_txn, create_new_sifaddr_and_key, send_sample_rowan, balance_check, query_block_claim

# AUTOMATED TEST TO VALIDATE ONLINE TXN
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_create_online_singlekey_txn(claimType):
    sifchain_address, sifchain_name = create_new_sifaddr_and_key()
    logging.info(f"sifchain_address = {sifchain_address}, sifchain_name = {sifchain_name}")
    destaddress1, destname1 = create_new_sifaddr_and_key()
    destaddress2, destname2 = create_new_sifaddr_and_key()
    from_address = 'sifnodeadmin'
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '10000000rowan'
    fee = '50000'
    currency = 'rowan'
    sampleamount = '1000rowan'
    
    # THESE 3 TXNS ARE TO REGISTER NEW ACCOUNTS ON CHAIN
    send_sample_rowan(from_address, sifchain_address, amount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress1, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress2, sampleamount, keyring_backend, chain_id)
    time.sleep(5)

    # CREATING TEST DATA HERE MIMICKING OUTPUT.JSON TO BE SUPPLIED BY NIKO'S API
    dict1 = {"denom": "rowan", "amount": "5000"}
    dict2 = {"denom": "rowan", "amount": "7000"}
    dict3 = {"address": destaddress1, "coins": [dict1]}
    dict4 = {"address": destaddress2, "coins": [dict2]}
    dict5 = {"Output": [dict3, dict4]}
    data = json.dumps(dict5)
    with open("output.json", "w") as f:
        f.write(data)

    # READ OUTPUT.JSON WITH CLAIMING ADDRESSES AND AMOUNT
    with open("output.json", "r") as f:
        data = f.read()
    d = json.loads(data)

    one_claiming_address = str(d['Output'][0]['address'])
    logging.info(f"one claiming address = {one_claiming_address}")

    # SENDER AND RECIPENT INITIAL BALANCE
    sender_initial_balance = int(balance_check(sifchain_address, currency))
    claiming_address_initial_balance = int(balance_check(one_claiming_address, currency))
    logging.info(f"sender initial balance = {sender_initial_balance}")
    logging.info(f"one claiming address initial balance = {claiming_address_initial_balance}")

    # ACTUAL DISPENSATION TXN; GET TXN HASH
    txhash = str((create_online_singlekey_txn(claimType, sifchain_name, chain_id, sifnodecli_node)))
    logging.info(f"txn hash = {txhash}")
    time.sleep(5)

    # QUERY BLOCK USING TXN HASH
    resp = query_block_claim(txhash)
    logging.info(f"valid hash response = {txhash}")

     # READ DISPENSATION TXN JSON TAGS
    distributionstartedtag = resp['logs'][0]['events'][0]['type']
    distributionattributesttags = resp['logs'][0]['events'][0]['attributes'][0]
    distributiontypetag = resp['tx']['value']['msg'][0]['type']
    distributionvaluetags = resp['tx']['value']['msg'][0]['value']
    account_key = str((distributionattributesttags['key']))
    chaintags = list(distributionvaluetags.keys())
    list_of_values = [distributionvaluetags[key] for key in distributionvaluetags]
    print(chaintags)
    print(list_of_values)
    # DISTRIBUTION TXN JSON TAGS ASSERTIONS
    assert str(distributionstartedtag) == 'distribution_started'
    assert str(account_key) == 'module_account'
    assert str(distributiontypetag) == 'dispensation/create'
    assert chaintags[0] == 'distributor'
    assert chaintags[1] == 'runner'
    assert chaintags[2] == 'distribution_type'
    assert list_of_values[0] == sifchain_address
    
    txn_signer_sender_address = resp['tx']['value']['msg'][0]['value']['distributor']
    distributionaddresslist = resp['tx']['value']['msg'][0]['value']['output']
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

    logging.info(
        f"recipients and their respective distributed amounts = {recipient_with_respective_distributed_amount}")
    logging.info(f"total amount distributed = {total_amount_distributed}")

    sender_final_balance = int(balance_check(sifchain_address, currency))
    recipient_address_final_balance = int(balance_check(one_claiming_address, currency))

    logging.info(f"sender initial balance = {sender_initial_balance}")
    logging.info(f"sender final balance = {sender_final_balance}")

    claimed_amount_single_recipient = int(recipient_with_respective_distributed_amount[one_claiming_address])

    # BALANCES ASSERTIONS
    assert int(total_amount_distributed) == int((sender_initial_balance - sender_final_balance) - int(fee))
    assert int(claimed_amount_single_recipient) == (recipient_address_final_balance - claiming_address_initial_balance)
    logging.info(
        f"balance transferred including fee from sender's address  = {(sender_initial_balance - sender_final_balance)}")
    logging.info(f"total amount distributed  = {total_amount_distributed}")

    logging.info(f"amount claimed by one recipient  = {claimed_amount_single_recipient}")
    logging.info(
        f"balance transferred in one recipient address  = {(recipient_address_final_balance - claiming_address_initial_balance)}")


# AUTOMTED TEST TO VALIDATE IF FUNDING ADDRESS DOESN'T HAVE ENOUGH BALANCE
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy', 'LiquidityMining'])
def test_insufficient_funds_dispensation_txn(claimType):
    sifchain_address, sifchain_name = create_new_sifaddr_and_key()
    destaddress1, destname1 = create_new_sifaddr_and_key()
    destaddress2, destname2 = create_new_sifaddr_and_key()
    from_address = 'sifnodeadmin'
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '70000rowan'
    fee = '50000'
    currency = 'rowan'
    sampleamount = '1000rowan'

    # THESE 3 TXNS ARE TO REGISTER NEW ACCOUNTS ON CHAIN
    send_sample_rowan(from_address, sifchain_address, amount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress1, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress2, sampleamount, keyring_backend, chain_id)
    time.sleep(5)

    # CREATING TEST DATA HERE MIMICKING OUTPUT.JSON TO BE SUPPLIED BY NIKO'S API
    dict1 = {"denom": "rowan", "amount": "5000"}
    dict2 = {"denom": "rowan", "amount": "19000"}
    dict3 = {"address": destaddress1, "coins": [dict1]}
    dict4 = {"address": destaddress2, "coins": [dict2]}
    dict5 = {"Output": [dict3, dict4]}
    data = json.dumps(dict5)
    with open("output.json", "w") as f:
        f.write(data)

    # READ OUTPUT.JSON WITH CLAIMING ADDRESSES AND AMOUNT
    with open("output.json", "r") as f:
        data = f.read()
    d = json.loads(data)

    one_claiming_address = str(d['Output'][0]['address'])
    logging.info(f"one claiming address = {one_claiming_address}")

    # SENDER AND RECIPENT INITIAL BALANCE
    sender_initial_balance = int(balance_check(sifchain_address, currency))
    claiming_address_initial_balance = int(balance_check(one_claiming_address, currency))
    logging.info(f"sender initial balance = {sender_initial_balance}")
    logging.info(f"one claiming address initial balance = {claiming_address_initial_balance}")

    # ACTUAL DISPENSATION TXN; TXN RAISES AN EXCEPTION ABOUT INSUFFICIENT FUNDS, CAPTURED HERE AND TEST IS MARKED PASS
    with pytest.raises(Exception) as execinfo:
        txhash = str(
            (create_online_singlekey_txn(claimType, sifchain_address, chain_id, sifnodecli_node)))
        assert str(
            execinfo.value) == f"for address  : {sifchain_address}: Failed in collecting funds for airdrop: failed to execute message; message index: 0: failed to simulate tx"
        logging.info(f"Insufficient Funds Message = {resp['raw_log']}")
