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

#CODE TO CREATE A NEW MULTI-KEY UNSIGNED TXN
#INPUT.JSON CONTAINS FUNDING ADDRESSES; OUTPUT.JSON CONTAINS CLAIM RECIPIENT ADDRESSES
def create_unsigned_multikey_txn(
        claimType,
        dispensation_name,
        chain_id,
        sifnodecli_node
    ):
    logging.debug(f"create_unsigned_dispensation")
    sifchain_fees_entry = f"--gas 200064128"
    input = 'input.json'
    output = 'output.json'
    cmd = " ".join([
        "sifnodecli tx dispensation create mkey",
        f"{dispensation_name}",
        f"{claimType}",
        input,
        output,
        sifchain_fees_entry,
        f"--generate-only", 
        
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    return json_str

#CODE TO SIGN DISPENSATION BY USER1
def sig1_txn(address1, offlinetx):
    cmd = " ".join([
        "sifnodecli tx sign",
        f"--multisig",
        f"$(sifnodecli keys show mkey -a)",
        f"--from $(sifnodecli keys show {address1} -a)",
        f"{offlinetx}"
    ])
    json_str = get_shell_output_json(cmd)
    return json_str

#CODE TO SIGN DISPENSATION BY USER2
def sig2_txn(address2, offlinetx):
    cmd = " ".join([
        "sifnodecli tx sign",
        f"--multisig",
        f"$(sifnodecli keys show mkey -a)",
        f"--from $(sifnodecli keys show {address2} -a)",
        f"{offlinetx}"
    ])
    json_str = get_shell_output_json(cmd)
    return json_str

#CODE TO GENERATE MULTI_SIGN TX
def multisign_txn(offlinetx, sig1, sig2):
    cmd = " ".join([
        "sifnodecli tx multisign",
        f"{offlinetx}",
        f"mkey",
        f"{sig1}",
        f"{sig2}"
    ])
    json_str = get_shell_output_json(cmd)
    return json_str

#CODE TO BROADCAST MULTISIGNED TXN ON BLOCK
def broadcast_multisign_txn(signedtx):
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

#CODE TO QUERY A NEW DISPENSATION
def query_created_dispensation(dispName):
    cmd = " ".join([
        "sifnodecli q dispensation records-by-name-all",
        f"{dispName}",
    ])
    json_str = get_shell_output_json(cmd)
    return json_str

#TEST CODE TO ASSERT TAGS GENERATED ON A BLOCK WHEN A NEW UNSIGNED DISPENSATION IS CREATED
def test_create_unsigned_multikey_txn():
    claimType = 'ValidatorSubsidy'
    dispensation_name = 'test_cd'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    address1 = 'sif'
    address2 = 'akasha'
    response = (create_unsigned_multikey_txn(claimType,dispensation_name,chain_id,sifnodecli_node))
    with open("sample.json", "w") as outfile: 
        json.dump(response, outfile)
    try:
        distype = (response['value']['msg'][0]['type'])
        imptags = response['value']['msg'][0]['value']
        actuallisttags = list(imptags.keys())
        print(actuallisttags)
        print(distype)
        assert str(distype) == 'dispensation/create'
        assert actuallisttags[0] == 'Signer'
        assert actuallisttags[1] == 'distribution_name'
        assert actuallisttags[2] == 'distribution_type'
        assert actuallisttags[3] == 'Input'
        assert actuallisttags[4] == 'Output'
        try:
            os.remove('sample.json')
        except OSError as e:
            print ("Error: %s - %s." % (e.filename, e.strerror))

    except KeyError:
        with pytest.raises(Exception, match='User trying to create a duplicate claim'):
            raise Exception

#TEST CODE TO ASSERT TAGS WNEN INDIVIDUAL SIGNED TXNS ARE CREATED
def test_sig1_sig2_txn():
    claimType = 'ValidatorSubsidy'
    dispensation_name = 'test_cd'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    address1 = 'sif'
    address2 = 'akasha'
    response = (create_unsigned_multikey_txn(claimType,dispensation_name,chain_id,sifnodecli_node))
    with open("sample.json", "w") as outfile: 
        json.dump(response, outfile)
    try:
        sig1response = sig1_txn(address1, 'sample.json')
        sig2response = sig2_txn(address2, 'sample.json')
        print(sig1response)
        sig1tags = list(sig1response.keys())
        sig2tags = list(sig2response.keys())
        print(sig1tags)
        assert sig1tags[0] == 'pub_key'
        assert sig1tags[1] == 'signature'
        assert sig2tags[0] == 'pub_key'
        assert sig2tags[1] == 'signature'
        try:
            os.remove('sample.json')
        except OSError as e:
            print ("Error: %s - %s." % (e.filename, e.strerror))
    except KeyError:
        with pytest.raises(Exception, match='User trying to create a duplicate claim'):
            raise Exception

#TEST CODE TO ASSERT TAGS WNEN A MULT_SIGN TXN IS CREATED USING INDIVIDUAL SIGNED TXNS
def test_multisign_txn():
    claimType = 'ValidatorSubsidy'
    dispensation_name = 'test_cd'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    address1 = 'sif'
    address2 = 'akasha'
    response = (create_unsigned_multikey_txn(claimType,dispensation_name,chain_id,sifnodecli_node))
    with open("sample.json", "w") as outfile: 
        json.dump(response, outfile)
    try:
        sig1response = sig1_txn(address1, 'sample.json')
        sig2response = sig2_txn(address2, 'sample.json')
        with open("sig1.json", "w") as sig1file: 
            json.dump(sig1response, sig1file)
        
        with open("sig2.json", "w") as sig2file: 
            json.dump(sig2response, sig2file)

        multisigresponse = multisign_txn('sample.json', 'sig1.json', 'sig2.json')
        distype = (multisigresponse['value']['msg'][0]['type'])
        imptags = multisigresponse['value']['msg'][0]['value']
        actuallisttags = list(imptags.keys())
        print(actuallisttags)
        print(distype)
        assert str(distype) == 'dispensation/create'
        assert actuallisttags[0] == 'Signer'
        assert actuallisttags[1] == 'distribution_name'
        assert actuallisttags[2] == 'distribution_type'
        assert actuallisttags[3] == 'Input'
        assert actuallisttags[4] == 'Output'
        try:
            os.remove('sig1.json')
            os.remove('sig2.json')
            os.remove('sample.json')
        except OSError as e:
            print ("Error: %s - %s." % (e.filename, e.strerror))

    except KeyError:
        with pytest.raises(Exception, match='User trying to create a duplicate claim'):
            raise Exception

#TEST CODE TO ASSERT TAGS GENERATED ON A BLOCK WHEN A NEW SIGNED DISPENSATION IS BROADCASTED on BLOCKCHAIN
def test_broadcast_multisign_txn():
    claimType = 'ValidatorSubsidy'
    dispensation_name = 'test_cd'
    chain_id = 'localnet'
    sifnodecli_node = 'tcp://127.0.0.1:1317'
    address1 = 'sif'
    address2 = 'akasha'
    response = (create_unsigned_multikey_txn(claimType,dispensation_name,chain_id,sifnodecli_node))
    with open("sample.json", "w") as outfile: 
        json.dump(response, outfile)
    try:
        sig1response = sig1_txn(address1, 'sample.json')
        sig2response = sig1_txn(address2, 'sample.json')
        with open("sig1.json", "w") as sig1file: 
            json.dump(sig1response, sig1file)
        
        with open("sig2.json", "w") as sig2file: 
            json.dump(sig2response, sig2file)

        multisigresponse = multisign_txn('sample.json', 'sig1.json', 'sig2.json')
        with open("multisig.json", "w") as multisigfile: 
            json.dump(multisigresponse, multisigfile)

        txhash = broadcast_multisign_txn('multisig.json')
        time.sleep(5)
        resp = query_block_claim(txhash)
        distype = (resp['tx']['value']['msg'][0]['type'])
        disvals = (resp['tx']['value']['msg'][0]['value'])
        bcasttags = list(disvals.keys())
        assert bcasttags[0] == 'Signer'
        assert bcasttags[1] == 'distribution_name'
        assert bcasttags[2] == 'distribution_type'
        assert bcasttags[3] == 'Input'
        assert bcasttags[4] == 'Output'
        try:
            os.remove('multisig.json')
            os.remove('sig1.json')
            os.remove('sig2.json')
            os.remove('sample.json')
        except OSError as e:
            print ("Error: %s - %s." % (e.filename, e.strerror))

    except KeyError:
        with pytest.raises(Exception, match='User trying to create a duplicate claim'):
            raise Exception

