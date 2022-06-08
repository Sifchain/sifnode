import argparse
import sys

import siftool_path
from siftool import test_utils, cosmos, command
from siftool.common import *
from load_testing import *


log = logging.getLogger(__name__)


def get_ctx():
    return test_utils.get_env_ctx()


def get_balances(ctx: test_utils.EnvCtx, address: cosmos.Address):
    balances = ctx.get_sifchain_balance_large(address)
    log.info("Number of balances: {}".format(len(balances)))


if __name__ == "__main__":
    basic_logging_setup()
    parser = argparse.ArgumentParser()
    parser.add_argument("address")
    args = parser.parse_args(sys.argv[1:])
    ctx = test_utils.get_env_ctx()
    get_balances(ctx, args.address)
