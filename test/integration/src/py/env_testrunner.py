import subprocess
import tempfile
from dataclasses import dataclass
from typing import List

import env_ethereum
import env_utilities
from env_utilities import wait_for_port

testrunnername="testrunner"

@dataclass
class TestrunnerInput(env_utilities.SifchainCmdInput):
    ebrelayer_config_file: str
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
    sifchain_admin_address: str


# Builds a collection of all the values needed to run tests
def build_testrunner_input(
        basedir: str,
        logfile: str,
        configoutputfile: str,
        ebrelayer_config_file: str,
        sifnode_config_file: str,
        deployment_name: str,
        ethereum_config_file: str,
        smart_contract_config_file: str,
):
    sifnode_config = env_utilities.read_config_file(sifnode_config_file)
    ethereum_config = env_utilities.read_config_file(ethereum_config_file)
    smart_contract_config = env_utilities.read_config_file(smart_contract_config_file)
    return TestrunnerInput(
        basedir=basedir,
        logfile=logfile,
        configoutputfile=configoutputfile,
        ebrelayer_config_file=ebrelayer_config_file,
        sifnode=f'tcp://{sifnode_config["input"]["sifnode_host"]}:{sifnode_config["input"]["rpc_port"]}',
        deployment_name=deployment_name,
        operator_address=smart_contract_config["input"]["operator_address"],
        operator_private_key=smart_contract_config["input"]["operator_private_key"],
        ethereum_address=smart_contract_config["input"]["ethereum_address"],
        ethereum_private_key=smart_contract_config["input"]["ethereum_private_key"],
        rowan_source=sifnode_config["config"]["adminuser"]["address"],
        ethereum_network=smart_contract_config["input"]["truffle_network"],
        ethereum_network_id=smart_contract_config["input"]["network_id"],
        ethereum_websocket_address=smart_contract_config["input"]["ws_addr"],
        infura_id=8,
        sifchain_admin_address=sifnode_config["config"]["adminuser"]["address"]
    )
    pass


# Writes out test values to an environment file for use with the
# current python integration tests
def testrunner_config_contents(args: TestrunnerInput):
    config = f"""
export BASEDIR={args.basedir}
export TEST_INTEGRATION_DIR={args.basedir}/test/integration
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
export SIFCHAIN_ADMIN_ACCOUNT={args.sifchain_admin_address}

. $BASEDIR/test/integration/environment_setup.sh
    """
    return config
