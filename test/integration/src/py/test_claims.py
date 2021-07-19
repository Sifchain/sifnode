import logging
import os
import time
import json
import pytest
import string
import random
from dispensation_envutils import create_new_sifaddr_and_key, send_sample_rowan, balance_check, \
query_block_claim, create_claim, query_created_claim


# TEST CODE TO ASSERT TAGS GENERATED ON A BLOCK WHEN A NEW CLAIM IS CREATED
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_create_new_claim(claimType):
    sifchain_address, sifchain_name = create_new_sifaddr_and_key()
    keyring_backend = 'test'
    chain_id = 'localnet'
    from_address = 'sifnodeadmin'
    amount = '10000000rowan'
    send_sample_rowan(from_address, sifchain_address, amount, keyring_backend, chain_id,"")
    time.sleep(5)
    txnhash = (create_claim(sifchain_address, claimType, keyring_backend, chain_id))
    logging.info(f"txnhash= ${txnhash}")
    time.sleep(5)
    response = (query_block_claim(str(txnhash)))
    logging.info(f"response= ${response}")
    try:
        data = (response['logs'][0]['events'][1]['attributes'])
        expectedOutputTagsList = []
        for value in data:
            expectedOutputTagsList.append(value['key'])
            expectedOutputTagsList.append(value['value'])
        assert response['txhash'] == txnhash
        assert expectedOutputTagsList[0] == 'userClaim_creator'
        assert expectedOutputTagsList[2] == 'userClaim_type'
        assert expectedOutputTagsList[4] == 'userClaim_creationTime'
        assert str(expectedOutputTagsList[1]) == sifchain_address
        assert str(expectedOutputTagsList[3]) in ['DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY', 'DISTRIBUTION_TYPE_LIQUIDITY_MINING']
    except KeyError:
        with pytest.raises(Exception, match='User trying to create a duplicate claim'):
            raise Exception


# TEST CODE TO ASSERT TAGS RETURNED BY A CLAIM QUERY COMMAND
@pytest.mark.parametrize("claimType", ['ValidatorSubsidy','LiquidityMining'])
def test_query_created_claim(claimType):
    queryresponse = query_created_claim(claimType)
    queryresponse = query_created_claim(claimType)
    querydata = (queryresponse['claims'][0])
    query_expected_tags = list(querydata.keys())
    assert query_expected_tags[0] == 'user_address'
    assert query_expected_tags[1] == 'user_claim_type'
    assert query_expected_tags[2] == 'user_claim_time'
    query_expected_values = list(querydata.values())
    assert query_expected_values[1] in ['DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY', 'DISTRIBUTION_TYPE_LIQUIDITY_MINING']