import json
import os
import pathlib
import shutil
import subprocess
from dataclasses import dataclass
from typing import List, Tuple

import env_ethereum
import env_utilities

smartcontractrunner_name = "smartcontractrunner"


@dataclass
class SmartContractDeployInput(env_utilities.SifchainCmdInput):
    ethereum_config_file: str
    network_id: int
    ws_addr: str
    truffle_network: str
    operator_address: str
    operator_private_key: str
    ethereum_address: str  # an ethereum address with eth to hand out
    ethereum_private_key: str
    n_validators: int
    validator_ethereum_credentials: List[Tuple[str, str]]
    validator_powers: List[int]
    consensus_threshold: int
    deployment_dir: str


def smartcontractrunner_docker_compose(args: env_ethereum.EthereumInput):
    base = env_utilities.base_docker_compose(smartcontractrunner_name)
    network = "sifchaintest"
    return {
        smartcontractrunner_name: {
            **base,
            "networks": [network],
            "working_dir": "/sifnode/test/integration",
        }
    }


def smart_contract_dir(args: SmartContractDeployInput):
    return os.path.join(args.basedir, "smart-contracts")


def deploy_contracts_cmd(args: SmartContractDeployInput):
    print(f"argsare: {json.dumps(args.__dict__, indent=2)}")
    for i in ["OWNER", "PAUSER", "OPERATOR"]:
        os.environ[i] = args.operator_address
    os.environ["ETHEREUM_PRIVATE_KEY"] = args.operator_private_key
    os.environ["ETHEREUM_WEBSOCKET_ADDRESS"] = args.ws_addr
    os.environ["ETHEREUM_NETWORK_ID"] = str(args.network_id)
    os.environ["CONSENSUS_THRESHOLD"] = str(args.consensus_threshold)
    validator_addresses = ",".join(map(lambda x: x[0], args.validator_ethereum_credentials))
    valpowers = ",".join(map(lambda x: str(x), args.validator_powers))
    env_vars = " ".join([
        f"INITIAL_VALIDATOR_ADDRESSES={validator_addresses}",
        f"INITIAL_VALIDATOR_POWERS={valpowers}",
        f''
    ])
    return f"cd {smart_contract_dir(args)} && {env_vars} npx truffle deploy --network {args.truffle_network} --reset"


# def read_smart_contract_artifacts(args: SmartContractDeployInput):
#     contracts = ["BridgeBank", "BridgeRegistry", "BridgeToken"]
#     for c in contracts:
#         p = os.path.join(smart_contract_dir(args), "build/contracts", f"{c}.json")


def deploy_contracts(args: SmartContractDeployInput):
    cmd = deploy_contracts_cmd(args)
    print(f"deploy_contracts_cmd: {cmd}")
    output = subprocess.run(
        cmd,
        shell=True,
        text=True
    )
    build_artifacts_directory = os.path.join(args.basedir, "smart-contracts/build/contracts")
    build_artifacts = os.listdir(build_artifacts_directory)
    pathlib.Path(args.deployment_dir).mkdir(exist_ok=True)
    for jsfile in filter(lambda f: ".json" in f, build_artifacts):
        srcfile = os.path.join(build_artifacts_directory, jsfile)
        shutil.copy(srcfile, args.deployment_dir)
    env_utilities.startup_complete(
        args, {
            "stdout": output.stdout,
            "stderr": output.stderr,
        }
    )


def contract_address(dirname: str, contract_name: str, network: str):
    j = env_utilities.read_json_file(os.path.join(dirname, contract_name + ".json"))
    return j["networks"][str(network)]["address"]
