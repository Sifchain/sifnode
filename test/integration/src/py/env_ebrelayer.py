import json
import os
import subprocess
import tempfile
from dataclasses import dataclass

import yaml

import env_utilities

ebrelayername = "ebrelayer"


@dataclass
class EbrelayerInput(env_utilities.SifchainCmdInput):
    ethereum_address: str
    ethereum_private_key: str
    web3_provider: str
    tendermint_node: str # something like tcp://0.0.0.0:26657
    bridge_registry_address: str
    moniker: str
    mnemonic: str
    chain_id: str
    home_dir: str
    gas: str
    gas_prices: str


def ebrelayer_cmd(args: EbrelayerInput):
    quote = "\""
    cmd = " ".join([
        "cd ~ &&",
        f"ETHEREUM_PRIVATE_KEY={args.ethereum_private_key}",
        "ebrelayer init",
        args.tendermint_node,
        args.web3_provider,
        args.bridge_registry_address,
        args.moniker,
        f"{quote}{args.mnemonic}{quote}",
        f"--chain-id {args.chain_id}",
        f"--home {args.home_dir}",
        f"--gas {args.gas}",
        f"--gas-prices {args.gas_prices}",
    ])
    return cmd


def run(args: EbrelayerInput):
    """runs ebrelayer"""
    cmd = ebrelayer_cmd(args)
    sout = open("/logs/ebrelayer.stdout.log", "w")
    serr = open("/logs/ebrelayer.stderr.log", "w")
    print(f"ebrelayercmd: \n{cmd}")
    subprocess.run(
        cmd,
        shell=True,
    )
    env_utilities.startup_complete(args, {})
