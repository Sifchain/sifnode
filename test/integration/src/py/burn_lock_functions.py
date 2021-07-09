import argparse
import json
import logging
import os
import tempfile
import textwrap
from typing import List

import sys
import time

from test_utilities import get_sifchain_addr_balance, advance_n_ethereum_blocks, \
    n_wait_blocks, print_error_message, wait_for_sifchain_addr_balance, send_from_ethereum_to_sifchain, \
    get_eth_balance, send_from_sifchain_to_ethereum, wait_for_eth_balance, \
    wait_for_ethereum_block_number, send_from_sifchain_to_sifchain, wait_for_sif_account, \
    get_shell_output_json, EthereumToSifchainTransferRequest, SifchaincliCredentials, RequestAndCredentials, \
    sifnoded_binary

default_timeout_for_ganache = 160


def decrease_log_level(new_level=logging.WARNING):
    logger = logging.getLogger()
    existing_level = logger.level
    if new_level > existing_level:
        logger.setLevel(new_level)
    return existing_level


def force_log_level(new_level):
    logger = logging.getLogger()
    existing_level = logger.level
    logger.setLevel(new_level)
    return existing_level


def transfer_ethereum_to_sifchain(transfer_request: EthereumToSifchainTransferRequest,
                                  max_seconds: int = default_timeout_for_ganache):
    logging.debug(f"transfer_ethereum_to_sifchain {transfer_request.as_json()}")
    assert transfer_request.ethereum_address
    assert transfer_request.sifchain_address

    # it's possible that this is the first transfer to the address, so there's
    # no balance to retrieve.  Catch that exception.

    original_log_level = decrease_log_level()

    try:
        sifchain_starting_balance = get_sifchain_addr_balance(
            transfer_request.sifchain_address,
            transfer_request.sifnodecli_node,
            transfer_request.sifchain_symbol
        )
    except:
        logging.debug(f"transfer_ethereum_to_sifchain failed to get starting balance, this is probably a new account")
        sifchain_starting_balance = 0

    status = {
        "action": "transfer_ethereum_to_sifchain",
        "sifchain_starting_balance": sifchain_starting_balance,
        "transfer_request": transfer_request.__dict__,
    }
    logging.debug(f"transfer_ethereum_to_sifchain_json: {json.dumps(status)}", )

    force_log_level(original_log_level)
    starting_block = send_from_ethereum_to_sifchain(transfer_request)
    original_log_level = decrease_log_level()
    logging.debug(f"send_from_ethereum_to_sifchain ethereum block number: {starting_block}")

    half_n_wait_blocks = n_wait_blocks / 2
    logging.debug("wait half the blocks, transfer should not complete")
    if transfer_request.manual_block_advance:
        advance_n_ethereum_blocks(half_n_wait_blocks, transfer_request.smart_contracts_dir)
        # we really want to wait for ebrelayer to catch up, but that's not possible yet
        time.sleep(5)
    else:
        wait_for_ethereum_block_number(
            block_number=starting_block + half_n_wait_blocks,
            transfer_request=transfer_request
        )

    # we still may not have an account
    try:
        sifchain_balance_before_required_elapsed_blocks = get_sifchain_addr_balance(
            transfer_request.sifchain_address,
            transfer_request.sifnodecli_node,
            transfer_request.sifchain_symbol
        )
    except:
        sifchain_balance_before_required_elapsed_blocks = 0

    # need to be able to turn off checking the balance after waiting half the blocks
    # because we want to be able to run some tests in parallel.  If parallel tests
    # are manually advancing blocks, you can't be sure where you are.
    if transfer_request.check_wait_blocks and sifchain_balance_before_required_elapsed_blocks != sifchain_starting_balance:
        print_error_message(
            f"balance should not have changed yet.  Starting balance {sifchain_starting_balance},"
            f" current balance {sifchain_balance_before_required_elapsed_blocks}"
        )

    if transfer_request.manual_block_advance:
        advance_n_ethereum_blocks(half_n_wait_blocks, transfer_request.smart_contracts_dir)
    else:
        wait_for_ethereum_block_number(
            block_number=starting_block + n_wait_blocks,
            transfer_request=transfer_request
        )

    target_balance = sifchain_starting_balance + transfer_request.amount

    # You can't get the balance of an account that doesn't exist yet,
    # so wait for the account to be there before asking for the balance
    logging.debug(f"wait for account {transfer_request.sifchain_address}")
    wait_for_sif_account(
        sif_addr=transfer_request.sifchain_address,
        sifchaincli_node=transfer_request.sifnodecli_node,
        max_seconds=max_seconds
    )

    wait_for_sifchain_addr_balance(
        sifchain_address=transfer_request.sifchain_address,
        symbol=transfer_request.sifchain_symbol,
        sifchaincli_node=transfer_request.sifnodecli_node,
        target_balance=target_balance,
        max_seconds=max_seconds,
        debug_prefix=f"transfer_ethereum_to_sifchain waiting for balance {transfer_request}"
    )

    force_log_level(original_log_level)

    result = {
        **status,
        "sifchain_ending_balance": target_balance,
    }
    logging.debug(f"transfer_ethereum_to_sifchain completed {json.dumps(result)}")
    return result


def transfer_sifchain_to_ethereum(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials,
        max_seconds: int = 90
):
    logging.debug(f"transfer_sifchain_to_ethereum_json: {transfer_request.as_json()}")

    original_log_level = decrease_log_level()
    ethereum_starting_balance = get_eth_balance(transfer_request)

    sifchain_starting_balance = get_sifchain_addr_balance(
        transfer_request.sifchain_address,
        transfer_request.sifnodecli_node,
        transfer_request.sifchain_symbol
    )

    status = {
        "action": "transfer_sifchain_to_ethereum",
        "ethereum_starting_balance": ethereum_starting_balance,
        "sifchain_starting_balance": sifchain_starting_balance,
    }
    logging.debug(status)

    force_log_level(original_log_level)
    raw_output = send_from_sifchain_to_ethereum(transfer_request, credentials)
    original_log_level = decrease_log_level()

    target_balance = ethereum_starting_balance + transfer_request.amount

    wait_for_eth_balance(
        transfer_request=transfer_request,
        target_balance=ethereum_starting_balance + transfer_request.amount,
        max_seconds=max_seconds
    )

    sifchain_ending_balance = get_sifchain_addr_balance(
        transfer_request.sifchain_address,
        transfer_request.sifnodecli_node,
        transfer_request.sifchain_symbol
    )

    result = {
        **status,
        "sifchain_ending_balance": sifchain_ending_balance,
        "ethereum_ending_balance": target_balance,
    }
    logging.debug(f"transfer_sifchain_to_ethereum_complete_json: {json.dumps(result)}")
    force_log_level(original_log_level)
    return result


def transfer_sifchain_to_sifchain(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials,
        max_seconds: int = 30
):
    logging.debug(f"transfer_sifchain_to_sifchain: {transfer_request.as_json()}")

    try:
        sifchain_starting_balance = get_sifchain_addr_balance(
            transfer_request.sifchain_destination_address,
            transfer_request.sifnodecli_node,
            transfer_request.sifchain_symbol
        )
    except Exception as e:
        # this is a new account, so the balance is 0
        sifchain_starting_balance = 0

    status = {
        "action": "transfer_sifchain_to_sifchain",
        "sifchain_starting_balance": sifchain_starting_balance,
    }
    logging.info(status)

    send_from_sifchain_to_sifchain(
        transfer_request,
        credentials
    )
    target_balance = transfer_request.amount + sifchain_starting_balance
    wait_for_sif_account(
        sif_addr=transfer_request.sifchain_destination_address,
        sifchaincli_node=transfer_request.sifnodecli_node,
        max_seconds=max_seconds
    )
    wait_for_sifchain_addr_balance(
        sifchain_address=transfer_request.sifchain_destination_address,
        symbol=transfer_request.sifchain_symbol,
        target_balance=target_balance,
        sifchaincli_node=transfer_request.sifnodecli_node,
        max_seconds=max_seconds,
        debug_prefix=f"transfer_sifchain_to_sifchain {transfer_request}"
    )

    return {
        **status,
        "sifchain_ending_balance": target_balance,
    }


def transfer_argument_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        formatter_class=argparse.RawDescriptionHelpFormatter,
        description=textwrap.dedent("""
    Transfer from Ethereum to Sifchain
    """))
    parser.add_argument(
        '--sifchain_address',
        type=str,
        nargs=1,
        required=True,
        help="A SifChain address like sif132tc0acwt8klntn53xatchqztl3ajfxxxsawn8"
    )
    parser.add_argument(
        '--sifchain_destination_address',
        type=str,
        nargs=1,
        required=False,
        default=[""],
        help="A SifChain address like sif132tc0acwt8klntn53xatchqztl3ajfxxxsawn8, used for transferring between sifchain addresses"
    )
    parser.add_argument(
        '--ethereum_address',
        type=str,
        nargs=1,
        required=True,
        help="An ethereum address like X"
    )
    parser.add_argument(
        '--ethereum_symbol',
        type=str,
        nargs=1,
        required=True,
        help="An ethereum symbol like eth"
    )
    parser.add_argument(
        '--sifchain_symbol',
        type=str,
        nargs=1,
        required=True,
        help="A symbol like ceth"
    )
    parser.add_argument(
        '--amount',
        type=str,
        nargs=1,
        required=True,
        help="An amount like 1000000000000000000"
    )
    parser.add_argument(
        '--smart_contracts_dir',
        type=str,
        nargs=1,
        required=True,
        help="The smart_contracts directory"
    )
    parser.add_argument(
        '--ethereum_chain_id',
        type=str,
        nargs=1,
        required=True,
    )
    parser.add_argument(
        '--logfile',
        type=str,
        nargs=1,
        default=["/dev/null"],
        help="A filename for logging (use - for stdout)"
    )
    parser.add_argument(
        '--loglevel',
        type=str,
        nargs=1,
        default=["debug"],
    )
    parser.add_argument(
        '--n_wait_blocks',
        type=str,
        nargs=1,
        default=[50],
    )
    parser.add_argument(
        '--chain_id',
        type=str,
        nargs=1,
        required=True
    )
    parser.add_argument(
        '--bridgebank_address',
        type=str,
        nargs=1,
        required=True
    )
    parser.add_argument(
        '--bridgetoken_address',
        type=str,
        nargs=1,
        required=True
    )
    parser.add_argument(
        '--sifnodecli_node',
        type=str,
        nargs=1,
        default="tcp://localhost:26657",
    )
    parser.add_argument('--manual_block_advance', action='store_true')
    return parser


def add_credentials_arguments(parser: argparse.ArgumentParser) -> argparse.ArgumentParser:
    parser.add_argument(
        '--keyring_backend',
        type=str,
        nargs=1,
        required=True,
        help="file,test,os"
    )
    parser.add_argument(
        '--keyring_passphrase_env_var',
        type=str,
        nargs=1,
        default=[""],
        help="The name of an environment variable holding the password"
    )
    parser.add_argument(
        '--from_key',
        type=str,
        nargs=1,
        default=[""],
        help="--from argument for sifnoded"
    )
    parser.add_argument(
        '--sifnodecli_homedir',
        type=str,
        nargs=1,
        required=True,
        help="The smart_contracts directory"
    )
    return parser


def configure_logging(args):
    logfile = args.logfile[0]

    if logfile == "-":
        handlers = [logging.StreamHandler(sys.stdout)]
    elif logfile:
        handlers = [logging.FileHandler(args.logfile[0])]
    else:
        tf = tempfile.NamedTemporaryFile(delete=False)
        args.logfile = [tf.name]
        handlers = [logging.FileHandler(tf)]

    logging.basicConfig(
        level=str.upper(args.loglevel[0]),
        format="%(asctime)s [%(levelname)s] %(message)s",
        handlers=handlers
    )


def process_args(cmdline: List[str]) -> RequestAndCredentials:
    arg_parser = transfer_argument_parser()
    args = add_credentials_arguments(arg_parser).parse_args(args=cmdline)
    configure_logging(args)

    logging.debug(f"command line arguments: {sys.argv} {args}")

    transfer_request = EthereumToSifchainTransferRequest.from_args(args)

    credentials = SifchaincliCredentials(
        keyring_passphrase=os.environ.get(args.keyring_passphrase_env_var[0]),
        from_key=args.from_key[0],
        keyring_backend=args.keyring_backend[0],
        sifnodecli_homedir=args.sifnodecli_homedir[0],
    )

    return RequestAndCredentials(transfer_request, credentials, args)


def create_new_sifaddr(
        credentials: SifchaincliCredentials,
        keyname
):
    """returns something like {"name":"9cbf3bd4-f15c-4128-bae6-a534fc8d6877","type":"local","address":"sif19u4xtckuvy2zk9r2l4063g93s3r8qc4vw0a20t","pubkey":"sifpub1addwnpepqw88ns6dmy3xwjqh4mkvuda6ezn056nxy8ldrtpkrfuvuamexv9hxyzhxm7","mnemonic":"surprise fire cupboard orange scatter boat cruel ability oven gap accident purity delay"}"""
    keyring_passphrase = credentials.keyring_passphrase
    yes_subcmd = f"yes {keyring_passphrase} |" if keyring_passphrase else ""
    keyring_backend_subcmd = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    # Note that keys-add prints to stderr
    cmd = f"{yes_subcmd} {sifnoded_binary} keys add {keyname} --home {credentials.sifnodecli_homedir} {keyring_backend_subcmd} --output json 2>&1"
    return get_shell_output_json(cmd)
