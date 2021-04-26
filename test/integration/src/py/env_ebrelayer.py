import subprocess
from dataclasses import dataclass
import re
import env_utilities

ebrelayername = "ebrelayer"


@dataclass
class EbrelayerInput(env_utilities.SifchainCmdInput):
    ethereum_address: str
    ethereum_private_key: str
    web3_provider: str
    tendermint_node: str  # something like tcp://0.0.0.0:26657
    rpc_url: str
    bridge_registry_address: str
    moniker: str
    mnemonic: str
    chain_id: str
    home_dir: str
    gas: str
    gas_prices: str


def ebrelayer_cmd(args: EbrelayerInput):
    quote = "\""
    # ebrelayer wants private keys without a leading 0x
    pk = re.sub(r"0x", "", args.ethereum_private_key)
    cmd = " ".join([
        "cd ~ &&",
        f"ETHEREUM_PRIVATE_KEY={pk}",
        "ebrelayer init",
        args.tendermint_node,
        args.web3_provider,
        args.bridge_registry_address,
        args.moniker,
        f"{quote}{args.mnemonic}{quote}",
        f"--rpc-url {args.rpc_url}",
        f"--chain-id {args.chain_id}",
        f"--home {args.home_dir}",
        f"--gas {args.gas}",
        f"--gas-prices {args.gas_prices}",
    ])
    return cmd


def relayer_docker_compose(i: int):
    name = f"{ebrelayername}{i}"
    base = env_utilities.base_docker_compose(name)
    network = "sifchaintest"
    return {
        name: {
            **base,
            "networks": [network],
        }
    }


def run(args: EbrelayerInput):
    """runs ebrelayer"""
    cmd = ebrelayer_cmd(args)
    sout = open("/logs/ebrelayer.stdout.log", "w")
    serr = open("/logs/ebrelayer.stderr.log", "w")
    print(f"ebrelayercmd: \n{cmd}")
    process = subprocess.Popen(
        cmd,
        shell=True,
    )
    env_utilities.startup_complete(args, {})
    return process
