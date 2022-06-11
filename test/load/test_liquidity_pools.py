import sys
import time
import logging
import argparse

import siftool_path
from siftool import cosmos, eth, sifchain, test_utils, command
from siftool.cosmos import balance_add, balance_sub, balance_equal, balance_mul
from siftool.common import *
from load_testing import *


log = logging.getLogger(__name__)


# CREATE-POOL
# =============
# sifnoded tx clp add-liquidity \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol cusdc \
#   --nativeAmount 25990373000000000000000000 \
#   --externalAmount 123010000000 \
#   --fees 100000000000000000rowan \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# ADD-LIQUIDITY
# =================
# sifnoded tx clp add-liquidity \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol cusdc \
#   --nativeAmount 25990373000000000000000000 \
#   --externalAmount 123010000000 \
#   --fees 100000000000000000rowan \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# UNBOND-LIQUIDITY
# =================
# sifnoded tx clp unbond-liquidity \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol cusdc \
#   --units 77971121144444445014456057 \
#   --fees 100000000000000000rowan \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# CANCEL-UNBOND-LIQUIDITY
# ===========================
# sifnoded tx clp cancel-unbond\
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol cusdc \
#   --units 77971121144444445014456057 \
#   --fees 100000000000000000rowan \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# REMOVE-LIQUIDITY
# ===========================
# sifnoded tx clp remove-liquidity-units \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol cusdc \
#   --withdrawUnits 77971121144444445014456057 \
#   --fees 100000000000000000rowan \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# SWAP-IN
# ============================
# sifnoded tx clp swap \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --sentSymbol cusdc \
#   --receivedSymbol rowan \
#   --sentAmount 184515000000 \
#   --minReceivingAmount 0 \
#   --fees 100000000000000000rowan \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# SWAP-OUT
# ============================
# sifnoded tx clp swap \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --sentSymbol rowan \
#   --receivedSymbol cusdc \
#   --sentAmount 77971121000000000000000000 \
#   --minReceivingAmount 0 \
#   --fees 100000000000000000rowan \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# QUERY-POOLS
# =============
# sifnoded q clp pools \
#   --chain-id $SIFNODE_CHAIN_ID

def test_liquidity_pools(ctx: test_utils.EnvCtx, example_param: int):
    # Note: start the environment in a separate window with "siftool --test-denom-count 10" to automatically create 10
    # test denoms in sifnoded genesis file. They will be named test0, test1, ..., test9 and credited to admin
    # address. When you create a new account via "create_sifchain_addr()" below, you can then use "fund_amounts"
    # parameter to automatically transfer some tokens from admin acount to the newly created account.

    # Test setup
    symbol = "test1"
    test_sif_addr = ctx.create_sifchain_addr(fund_amounts={ROWAN: 10**20, symbol: 10**18})
    assert ctx.get_sifchain_balance(test_sif_addr)[symbol] == 10**18
    ctx.sifnode_client.tx_clp_create_pool(test_sif_addr, symbol, 25990373000000000000000000, 123010000000)
    pools = ctx.sifnode_client.query_pools()

    # Test execution
    for i in range(example_param):
        log.info("Running loop {} of {}...".format(i, example_param))
        assert True


if __name__ == "__main__":
    basic_logging_setup()
    ctx = test_utils.get_env_ctx()
    parser = argparse.ArgumentParser()
    parser.add_argument("--example-param", type=int, default=10)
    args = parser.parse_args(sys.argv[1:])
    example_param = args.example_param
    test_liquidity_pools(ctx, example_param)
    log.info("Finished successfully")
