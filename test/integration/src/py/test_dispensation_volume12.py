import logging
import os
import time
import json
import pytest
import string
import random
from dispensation_envutils import create_online_singlekey_txn, create_new_sifaddr_and_key, send_sample_rowan, balance_check, \
query_block_claim, create_online_singlekey_txn_with_runner, run_dispensation

# AUTOMATED TEST TO VALIDATE ONLINE TXN
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy'])
def test_run_online_morethan10distribution_txn(claimType):
    distributor_address, distributor_name = create_new_sifaddr_and_key()
    runner_address, runner_name = create_new_sifaddr_and_key()
    logging.info(f"distributor_address = {distributor_address}, distributor_name = {distributor_name}")
    logging.info(f"runner_address = {runner_address}, runner_name = {runner_name}")
    destaddress1, destname1 = create_new_sifaddr_and_key()
    destaddress2, destname2 = create_new_sifaddr_and_key()
    destaddress3, destname3 = create_new_sifaddr_and_key()
    destaddress4, destname4 = create_new_sifaddr_and_key()
    destaddress5, destname5 = create_new_sifaddr_and_key()
    destaddress6, destname6 = create_new_sifaddr_and_key()
    destaddress7, destname7 = create_new_sifaddr_and_key()
    destaddress8, destname8 = create_new_sifaddr_and_key()
    destaddress9, destname9 = create_new_sifaddr_and_key()
    destaddress10, destname10 = create_new_sifaddr_and_key()
    destaddress11, destname11 = create_new_sifaddr_and_key()
    destaddress12, destname12 = create_new_sifaddr_and_key()
    from_address = 'sifnodeadmin'
    keyring_backend = 'test'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    amount = '100000000rowan'
    fee = '150000'
    currency = 'rowan'
    sampleamount = '1000rowan'
    
    # THESE 4 TXNS ARE TO REGISTER NEW ACCOUNTS ON CHAIN
    send_sample_rowan(from_address, runner_address, amount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, distributor_address, amount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress1, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress2, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress3, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress4, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress5, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress6, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress7, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress8, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress9, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress10, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress11, sampleamount, keyring_backend, chain_id)
    time.sleep(5)
    send_sample_rowan(from_address, destaddress12, sampleamount, keyring_backend, chain_id)
    time.sleep(5)

    sorted_dest_address_list = sorted([destaddress1,destaddress2,destaddress3,destaddress4,destaddress5,destaddress6,destaddress7,destaddress8,destaddress9,destaddress10,destaddress11,destaddress12])
    logging.info(f"sorted_dest_address_list = {sorted_dest_address_list}")

    # CREATING TEST DATA HERE MIMICKING OUTPUT.JSON TO BE SUPPLIED BY NIKO'S API
    dict1 = {"denom": "rowan", "amount": "5000"}
    dict2 = {"denom": "rowan", "amount": "7000"}
    dict3 = {"denom": "rowan", "amount": "8000"}
    dict4 = {"denom": "rowan", "amount": "9000"}
    dict5 = {"denom": "rowan", "amount": "10000"}
    dict6 = {"denom": "rowan", "amount": "11000"}
    dict7 = {"denom": "rowan", "amount": "12000"}
    dict8 = {"denom": "rowan", "amount": "13000"}
    dict9 = {"denom": "rowan", "amount": "14000"}
    dict10 = {"denom": "rowan", "amount": "15000"}
    dict11 = {"denom": "rowan", "amount": "16000"}
    dict12 = {"denom": "rowan", "amount": "17000"}

    dict13 = {"address": destaddress1, "coins": [dict1]}
    dict14 = {"address": destaddress2, "coins": [dict2]}
    dict15 = {"address": destaddress3, "coins": [dict3]}
    dict16 = {"address": destaddress4, "coins": [dict4]}
    dict17 = {"address": destaddress5, "coins": [dict5]}
    dict18 = {"address": destaddress6, "coins": [dict6]}
    dict19 = {"address": destaddress7, "coins": [dict7]}
    dict20 = {"address": destaddress8, "coins": [dict8]}
    dict21 = {"address": destaddress9, "coins": [dict9]}
    dict22 = {"address": destaddress10, "coins": [dict10]}
    dict23 = {"address": destaddress11, "coins": [dict11]}
    dict24 = {"address": destaddress12, "coins": [dict12]}

    dict25 = {"Output": [dict13, dict14, dict15, dict16, dict17, dict18, dict19, dict20, dict21,dict22, dict23, dict24]}
    data = json.dumps(dict25)
    with open("output.json", "w") as f:
        f.write(data)

    # READ OUTPUT.JSON WITH CLAIMING ADDRESSES AND AMOUNT
    with open("output.json", "r") as f:
        data = f.read()
    d = json.loads(data)

    one_claiming_address = str(d['Output'][0]['address'])
    logging.info(f"one claiming address = {one_claiming_address}")

    # SENDER AND RECIPENT INITIAL BALANCE
    sender_initial_balance = int(balance_check(distributor_address, currency))
    claiming_address_initial_balance = int(balance_check(one_claiming_address, currency))
    logging.info(f"sender initial balance = {sender_initial_balance}")
    logging.info(f"one claiming address initial balance = {claiming_address_initial_balance}")

    # CREATE DISPENSATION TXN; GET TXN HASH
    txhash = str((create_online_singlekey_txn_with_runner(claimType, runner_address, distributor_name, chain_id, sifnodecli_node)))
    logging.info(f"txn hash for creatng a dispensation = {txhash}")
    time.sleep(5)

    # QUERY BLOCK USING TXN HASH
    resp = query_block_claim(txhash)
    
     
    distribution_name = resp['logs'][0]['events'][0]['attributes'][1]['value']
    distribution_type = resp['logs'][0]['events'][0]['attributes'][2]['value']
    logging.info(f"distribution_name = {distribution_name}, distribution_type = {distribution_type}")

    # RUN DISPENSATION TXN; GET TXN HASH    
    runtxnhash1 = run_dispensation(distribution_name, distribution_type, runner_address,chain_id,sifnodecli_node)
    logging.info(f"txn hash for running dispensation = {runtxnhash1}")
    time.sleep(5)
    runtxnhash2 = run_dispensation(distribution_name, distribution_type, runner_address,chain_id,sifnodecli_node)
    logging.info(f"txn hash for running dispensation = {runtxnhash2}")
    time.sleep(5)

    # QUERY BLOCK USING TXN HASH
    runresp1 = query_block_claim(runtxnhash1)
    logging.info(f"response from block for run dispensation = {runresp1}")

    runresp2 = query_block_claim(runtxnhash2)
    logging.info(f"response from block for run dispensation = {runresp2}")

    rundistributiontag1 = runresp1['logs'][0]['events'][0]['type']
    rundistname1 = runresp1['logs'][0]['events'][0]['attributes'][0]['value']
    runrunneraddress1 = runresp1['logs'][0]['events'][0]['attributes'][1]['value']
    # runtempdistreceiverlist1 = [runresp1['logs'][0]['events'][0]['attributes'][2]['value']]
    # rundistreceiverlist1 = [x for xs in runtempdistreceiverlist1 for x in xs.split(',')]
    # sortedrundistreceiverlist1 = sorted(rundistreceiverlist1)
    # logging.info(f"sortedrundistreceiverlist = {sortedrundistreceiverlist1}")
    # logging.info(f"sortedrundistreceiverlist first item = {sortedrundistreceiverlist1[0]}")
    # logging.info(f"sortedrundistreceiverlist second item  = {sortedrundistreceiverlist1[1]}")

    # RUN DISTRIBUTION TXN JSON TAGS ASSERTIONS
    # assert str(rundistributiontag1) == 'distribution_run'
    # assert str(rundistname1) == distribution_name
    # assert str(runrunneraddress1) == runner_address
    # assert sortedrundistreceiverlist1[0] == sorted_dest_address_list[0]
    # assert sortedrundistreceiverlist1[1] == sorted_dest_address_list[1]
    # assert sortedrundistreceiverlist1[2] == sorted_dest_address_list[2]
    # assert sortedrundistreceiverlist1[3] == sorted_dest_address_list[3]
    # assert sortedrundistreceiverlist1[4] == sorted_dest_address_list[4]
    # assert sortedrundistreceiverlist1[5] == sorted_dest_address_list[5]
    # assert sortedrundistreceiverlist1[6] == sorted_dest_address_list[6]
    # assert sortedrundistreceiverlist1[7] == sorted_dest_address_list[7]
    # assert sortedrundistreceiverlist1[8] == sorted_dest_address_list[8]
    # assert sortedrundistreceiverlist1[9] == sorted_dest_address_list[9]
    
    # READING TAGS FROM RUN DISPENSATION CMD   
    temprundistamount1 = runresp1['logs'][0]['events'][2]['attributes'][2]['value']
    logging.info(f"temp amount distributed 1 = {temprundistamount1}")
    temprundistamount2 = runresp1['logs'][0]['events'][2]['attributes'][5]['value']
    logging.info(f"temp amount distributed 2 = {temprundistamount2}")
    temprundistamount3 = runresp1['logs'][0]['events'][2]['attributes'][8]['value']
    logging.info(f"temp amount distributed 3 = {temprundistamount3}")
    temprundistamount4 = runresp1['logs'][0]['events'][2]['attributes'][11]['value']
    logging.info(f"temp amount distributed 4 = {temprundistamount4}")
    temprundistamount5 = runresp1['logs'][0]['events'][2]['attributes'][14]['value']
    logging.info(f"temp amount distributed 5 = {temprundistamount5}")
    temprundistamount6 = runresp1['logs'][0]['events'][2]['attributes'][17]['value']
    logging.info(f"temp amount distributed 6 = {temprundistamount6}")
    temprundistamount7 = runresp1['logs'][0]['events'][2]['attributes'][20]['value']
    logging.info(f"temp amount distributed 7 = {temprundistamount7}")
    temprundistamount8 = runresp1['logs'][0]['events'][2]['attributes'][23]['value']
    logging.info(f"temp amount distributed 8 = {temprundistamount8}")
    temprundistamount9 = runresp1['logs'][0]['events'][2]['attributes'][26]['value']
    logging.info(f"temp amount distributed 9 = {temprundistamount9}")
    temprundistamount10 = runresp1['logs'][0]['events'][2]['attributes'][29]['value']
    logging.info(f"temp amount distributed 10 = {temprundistamount10}")
    temprundistamount11 = runresp2['logs'][0]['events'][2]['attributes'][2]['value']
    logging.info(f"temp amount distributed 11 = {temprundistamount11}")
    temprundistamount12 = runresp2['logs'][0]['events'][2]['attributes'][5]['value']
    logging.info(f"temp amount distributed 12 = {temprundistamount12}")
    my_List = [temprundistamount1, temprundistamount2, temprundistamount3, temprundistamount4, temprundistamount5, temprundistamount6, temprundistamount7, temprundistamount8, temprundistamount9, temprundistamount10, temprundistamount11, temprundistamount12]
    logging.info(f"my list = {my_List}")
    rundistamount = [int(i[:-5]) for i in my_List]
    logging.info(f"temp amount distributed 2 = {rundistamount}")
    runrecipientaddress1 = runresp1['logs'][0]['events'][2]['attributes'][0]['value']
    runrecipientaddress2 = runresp1['logs'][0]['events'][2]['attributes'][3]['value']
    runrecipientaddress3 = runresp1['logs'][0]['events'][2]['attributes'][6]['value']
    runrecipientaddress4 = runresp1['logs'][0]['events'][2]['attributes'][9]['value']
    runrecipientaddress5 = runresp1['logs'][0]['events'][2]['attributes'][12]['value']
    runrecipientaddress6 = runresp1['logs'][0]['events'][2]['attributes'][15]['value']
    runrecipientaddress7 = runresp1['logs'][0]['events'][2]['attributes'][18]['value']
    runrecipientaddress8 = runresp1['logs'][0]['events'][2]['attributes'][21]['value']
    runrecipientaddress9 = runresp1['logs'][0]['events'][2]['attributes'][24]['value']
    runrecipientaddress10 = runresp1['logs'][0]['events'][2]['attributes'][27]['value']

    runrecipientaddress11 = runresp2['logs'][0]['events'][2]['attributes'][0]['value']
    runrecipientaddress12 = runresp2['logs'][0]['events'][2]['attributes'][3]['value']
    amount_distributed = [rundistamount[0], rundistamount[1], rundistamount[2], rundistamount[3], rundistamount[4], rundistamount[5], rundistamount[6], rundistamount[7], rundistamount[8], rundistamount[9], rundistamount[10], rundistamount[11]]
    recipient_dispensation_addresses = [runrecipientaddress1, runrecipientaddress2, runrecipientaddress3, runrecipientaddress4, runrecipientaddress5, runrecipientaddress6, runrecipientaddress7, runrecipientaddress8, runrecipientaddress9, runrecipientaddress10, runrecipientaddress11, runrecipientaddress12]
    logging.info(f"dispensation txn addresses = {recipient_dispensation_addresses}")
    logging.info(f"amount distributed = {amount_distributed}")

    total_amount_distributed = sum(int(i) for i in amount_distributed)
    recipient_with_respective_distributed_amount = dict(zip(recipient_dispensation_addresses, amount_distributed))

    logging.info(
        f"recipients and their respective distributed amounts = {recipient_with_respective_distributed_amount}")
    logging.info(f"total amount distributed = {total_amount_distributed}")

    sender_final_balance = int(balance_check(distributor_address, currency))
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