import argparse
import textwrap
from dataclasses import dataclass
from typing import List

from env_utilities import SifchainCmdInput, SifchainCmdOutput, sifchain_cmd_input_parser


@dataclass
class EthereumInput(SifchainCmdInput):
    logfile: str
    chain_id: str
    network_id: str
    starting_ether: int


@dataclass
class EthereumOutput(SifchainCmdOutput):
    """geth has no special output that we need to use"""
    pass


def ethereum_args_parser(parser = None) -> argparse.ArgumentParser:
    if parser is None:
        parser = sifchain_cmd_input_parser()
    parser.add_argument('--chain_id', required=True)
    parser.add_argument('--network_id', required=True)
    parser.add_argument('--ws_port', required=True)
    parser.add_argument('--http_port', required=True)
    parser.add_argument('--ethereum_addresses', required=True)
    parser.add_argument('--starting_ethereum', type=int, default=100)
    return parser


def parsed_args_to_ethereum_input(args: argparse.Namespace) -> EthereumInput:
    return EthereumInput(
        logfile=args.logfile,
        chain_id=args.chain_id,
        network_id=args.network_id,
        http_port=args.http_port,
        ws_port=args.ws_port,
        ethereum_addresses=args.ethereum_addresses.split(','),
        configoutputfile=args.configoutputfile,
        starting_ether=args.starting_ether,
    )
