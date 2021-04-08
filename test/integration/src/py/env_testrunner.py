import subprocess
import tempfile
from dataclasses import dataclass
from typing import List

import env_ethereum
import env_utilities
from env_utilities import wait_for_port

golangname = "golang"


@dataclass
class TestrunnerInput(env_utilities.SifchainCmdInput):
    ebrelayer_config_file: str
    base_dir: str
    sifnode: str
    deployment_name: str
    operator_address: str
    operator_private_key: str
    ethereum_address: str
    ethereum_private_key: str
    rowan_source: str
    ethereum_network: str
    ethereum_network_id: str
    ethereum_websocket_address: str
    infura_id: str


def testrunner_config_contents(args: TestrunnerInput):
    config = f"""
export BASEDIR={args.base_dir}
export SIFNODE={args.sifnode}
export DEPLOYMENT_NAME={args.deployment_name}
export OPERATOR_ADDRESS={args.operator_address}
export ETHEREUM_ADDRESS={args.ethereum_address}
export ROWAN_SOURCE={args.rowan_source}
export ETHEREUM_NETWORK={args.ethereum_network}
export ETHEREUM_NETWORK_ID={args.ethereum_network_id}
export ETHEREUM_WEBSOCKET_ADDRESS={args.ethereum_websocket_address}
export INFURA_PROJECT_ID={args.infura_id}
export ETHEREUM_PRIVATE_KEY={args.ethereum_private_key}
export OPERATOR_PRIVATE_KEY={args.operator_private_key}

. $BASEDIR/test/integration/environment_setup.sh
    """
    return f"cd {args.base_dir} && GOBIN={args.go_bin} PATH=$PATH:/usr/local/go/bin make install"



def golang_build(args: GolangInput):
    cmd = golang_build_cmd(args)
    subprocess.run(
        cmd,
        shell=True
    )
    env_utilities.startup_complete(args, {})
