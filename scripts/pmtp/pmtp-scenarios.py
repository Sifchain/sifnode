#!/usr/bin/env python
from decimal import *
from operator import truediv

import os
import json
import codecs
data = json.load(codecs.open(
    os.path.join(os.path.dirname(__file__), "pools.json"), 'r', 'utf-8-sig'))

policies = [
    {
        "params": {
            "min_create_pool_threshold": 100,
            "pmtp_period_governance_rate": 0.10,
            "pmtp_period_epoch_length": 1,
            "pmtp_period_start_block": 1,
            "pmtp_period_end_block": 40
        },
        "epoch": {
            "epoch_counter": 0,
            "block_counter": 0
        }
    },
    {
        "params": {
            "min_create_pool_threshold": 100,
            "pmtp_period_governance_rate": 0.02,
            "pmtp_period_epoch_length": 14400,
            "pmtp_period_start_block": 420001,
            "pmtp_period_end_block": 852000
        },
        "epoch": {
            "epoch_counter": 0,
            "block_counter": 0
        }
    },
        {
        "params": {
            "min_create_pool_threshold": 100,
            "pmtp_period_governance_rate": 0.0004,
            "pmtp_period_epoch_length": 14400,
            "pmtp_period_start_block": 340021,
            "pmtp_period_end_block": 397620
        },
        "epoch": {
            "epoch_counter": 0,
            "block_counter": 0
        }
    },
            {
        "params": {
            "min_create_pool_threshold": 100,
            "pmtp_period_governance_rate": 10.0000,
            "pmtp_period_epoch_length": 14400,
            "pmtp_period_start_block": 14401,
            "pmtp_period_end_block": 28800
        },
        "epoch": {
            "epoch_counter": 0,
            "block_counter": 0
        }
    },
]





scenarios = []

for policy in policies:
    startBlock = Decimal(policy["params"]["pmtp_period_start_block"])
    endBlock = Decimal(policy["params"]["pmtp_period_end_block"])
    epochLength = Decimal(policy["params"]["pmtp_period_epoch_length"])
    govRate = Decimal(policy["params"]["pmtp_period_governance_rate"])

    one = Decimal(1)

    numBlocksInPolicyPeriod = endBlock - startBlock + one
    numEpochsInPolicyPeriod = numBlocksInPolicyPeriod / epochLength
    blockRate = pow(one + govRate, (numEpochsInPolicyPeriod / numBlocksInPolicyPeriod)) - one

    scenario = {
        "params": policy["params"],
        "epoch": policy["epoch"],
        "expected_states": []
    }

    for pool in data['pools']:

        nativeBalance = Decimal(pool["native_asset_balance"])    
        externalBalance = Decimal(pool["external_asset_balance"])
        if pool['external_asset']['symbol'] != 'ceth':
            externalBalance = externalBalance * Decimal(1000000000000)

        inc = int(numBlocksInPolicyPeriod) / int(10)
        first = True
        for height in range(int(startBlock) - 1, int(endBlock) + int(inc), int(inc)):
            if first:
                height += 1
                first = False
            runningRate = pow(one + blockRate, Decimal(height) - startBlock + one) - one

            swapResultNative = externalBalance*(one-(nativeBalance/(one+nativeBalance))) * (one - (one / (one+nativeBalance)))
            swapPriceNative = swapResultNative * (one + runningRate)

            swapResultExternal = nativeBalance*(one-(externalBalance/(one+externalBalance))) * (one - (one / (one+externalBalance)))
            swapPriceExternal = swapResultExternal * (one + runningRate)

            scenario["expected_states"].append(
                {
                    "height": int(height),
                    "pool": {
                        "external_asset": pool["external_asset"],
                        "native_asset_balance": pool["native_asset_balance"],
                        "external_asset_balance": pool["external_asset_balance"],
                        "pool_units": pool["pool_units"]
                    },
                    "swap_price_native": "{:.18f}".format(round(swapPriceNative, 18)),
                    "swap_price_external": "{:.18f}".format(round(swapPriceExternal, 18)),
                    "pmtp_rate_params": {
                        "pmtp_period_block_rate": "{:.18f}".format(round(blockRate, 18)),
                        "pmtp_current_running_rate": "{:.18f}".format(round(runningRate, 18)),
                        "pmtp_inter_policy_rate": "{:.18f}".format(round(0, 18))
                    }
                }
            )

            # height = height + inc - 1
    
    scenarios.append(scenario)

with open(os.path.join(os.path.dirname(__file__), "scenarios.json"), 'w') as json_file:
    json.dump(scenarios, json_file)

        