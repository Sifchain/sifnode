import dataclasses
import json
import os
import sys
import time
from uuid import uuid4

import yaml

import env_ebrelayer
import env_ganache
import env_geth
import env_golang
import env_sifnoded
import env_smartcontractrunner
import env_testrunner
import env_utilities
from env_geth import start_geth

configbase = "/configs"
logbase = "/logs/"

ganachename = "ganache"
gethname = "geth"
basedir = "/sifnode"
deployment_name = "localnet"
n_validators = 1


def config_file_full_path(name: str):
    return os.path.join(configbase, f"{name}.json")


def log_file_full_path(name: str):
    return os.path.join(logbase, f"{name}.log")


geth_input = env_geth.GethInput(
    basedir=basedir,
    logfile=log_file_full_path(gethname),
    chain_id=3,
    network_id=3,
    starting_ether=123,
    ws_port=8646,
    http_port=7890,
    ethereum_addresses=4,
    configoutputfile=config_file_full_path(gethname)
)

golang_input = env_golang.GolangInput(
    basedir=basedir,
    logfile=log_file_full_path(gethname),
    configoutputfile=config_file_full_path(env_golang.golangname),
    go_bin="/gobin",
    base_dir="/sifnode"
)

ganache_ws_port = 7545
ganache_ws_addr = f"http://ganache:{ganache_ws_port}"
ganache_network_id = 5777

ganache_input = env_ganache.GanacheInput(
    basedir=basedir,
    logfile=log_file_full_path(ganachename),
    network_id=ganache_network_id,
    chain_id=3,
    starting_ether=123,
    port=7545,
    block_delay=1,
    mnemonic=None,
    db_dir="/tmp/ganachedb",
    configoutputfile=config_file_full_path(ganachename),
)

smartcontractrunner_input = env_smartcontractrunner.SmartContractDeployInput(
    basedir=basedir,
    network_id=ganache_network_id,
    ethereum_config_file=config_file_full_path(ganachename),
    logfile=log_file_full_path(env_smartcontractrunner.smartcontractrunner_name + "deploy"),
    configoutputfile=config_file_full_path(env_smartcontractrunner.smartcontractrunner_name),
    ws_addr=ganache_ws_addr,
    truffle_network="dynamic",
    operator_private_key=None,
    operator_address=None,
    ethereum_address=None,
    ethereum_private_key=None,
    validator_ethereum_credentials=None,
    n_validators=n_validators,
    validator_powers=[100],
    consensus_threshold=100,
    deployment_dir=os.path.join(basedir, f"smart-contracts/deployments/{deployment_name}")
)

sifnoded_input = env_sifnoded.SifnodedRunner(
    basedir=basedir,
    bin_prefix=os.path.join("/gobin"),
    logfile=log_file_full_path(ganachename),
    configoutputfile=config_file_full_path(env_sifnoded.sifnodename),
    rpc_port=26657,
    chain_id=deployment_name,
    network_config_file="/tmp/netconfig.yml",
    seed_ip_address="10.10.1.1",
    n_validators=n_validators,
    go_build_config_path=config_file_full_path(env_golang.golangname),
    sifnode_host=env_sifnoded.sifnodename
)

geth_docker = env_geth.geth_docker_compose(geth_input)
ganache_docker = env_ganache.ganache_docker_compose(ganache_input)
smartcontractrunner_docker = env_smartcontractrunner.smartcontractrunner_docker_compose(ganache_input)
sifnodedrunner = env_sifnoded.sifnoded_docker_compose(sifnoded_input)

shared_docker = {
    "version": "3.9",
    "networks": {
        "sifchaintest": {
            "ipam": {
                "driver": "default",
                "config": [
                    {"subnet": "10.0.0.0/24"}
                ]
            },
        }
    },
}

component = sys.argv[1] if len(sys.argv) > 1 else "dockerconfig"

if component == "dockerconfig":
    print(yaml.dump({
        **shared_docker,
        "services": {
            **ganache_docker,
            **geth_docker,
            **smartcontractrunner_docker,
            **sifnodedrunner,
        }
    }))
elif component == "geth":
    print(f"starting geth, configuration is {geth_input}")
    start_geth(geth_input).wait()
elif component == "ganache":
    print(f"starting ganache, configuration is {yaml.dump(ganache_input)}")
    env_ganache.start_ganache(ganache_input).wait()
elif component == "smartcontractrunner":
    time.sleep(100000000)
elif component == "golang_build":
    env_golang.golang_build(golang_input)
elif component == "deploy_contracts":
    ethereum_config = env_utilities.read_config_file(smartcontractrunner_input.ethereum_config_file)
    private_keys_stanza = ethereum_config["config"]["private_keys"]
    private_keys = list(private_keys_stanza.values())
    ethereum_addresses = list(private_keys_stanza.keys())
    i = dataclasses.replace(
        smartcontractrunner_input,
        operator_address=ethereum_addresses[0],
        operator_private_key=private_keys[0],
        validator_ethereum_credentials=list(private_keys_stanza.items())[1:smartcontractrunner_input.n_validators + 1],
        ethereum_address=ethereum_addresses[-1],
        ethereum_private_key=private_keys[-1]
    )
    env_smartcontractrunner.deploy_contracts(i)
elif component == "startsifnoded":
    sifnoded_chain_data = env_sifnoded.build_chain(sifnoded_input)
    print(f"build_chain result: \n{json.dumps(sifnoded_chain_data)}")
    env_sifnoded.run(sifnoded_input, sifnoded_chain_data)
elif "relayer" in component:
    f = config_file_full_path(env_smartcontractrunner.smartcontractrunner_name)
    smart_contract_config = env_utilities.read_config_file(f)
    print(f"smart contract config: {json.dumps(smart_contract_config, indent=2)}")
    sifnodedconfig = env_utilities.read_config_file(config_file_full_path(env_sifnoded.sifnodename))["config"]
    bridge_registry_address = env_smartcontractrunner.contract_address(
        smart_contract_config["input"]["deployment_dir"],
        "BridgeRegistry",
        smart_contract_config["input"]["network_id"]
    )
    for v in sifnodedconfig["validators"]:
        print(f"sifnode validator:\n{json.dumps(v, indent=2)}")
        x = env_ebrelayer.EbrelayerInput(
            basedir=basedir,
            logfile=log_file_full_path(env_ebrelayer.ebrelayername),
            configoutputfile=config_file_full_path(env_ebrelayer.ebrelayername),
            ethereum_address=smart_contract_config["input"]["validator_ethereum_credentials"][0][0],
            ethereum_private_key=smart_contract_config["input"]["validator_ethereum_credentials"][0][1],
            web3_provider="ws://ganache:7545",
            tendermint_node="tcp://0.0.0.0:26657",
            bridge_registry_address=bridge_registry_address,
            moniker=v["moniker"],
            mnemonic=v["mnemonic"],
            chain_id=deployment_name,
            home_dir=v["sifnodeclipath"],
            gas="5000000000000",
            gas_prices="0.5rowan",
        )
        env_ebrelayer.run(x)
    time.sleep(10000)
elif component == "tr":
    i = env_testrunner.build_testrunner_input(
        basedir=basedir,
        logfile=log_file_full_path(env_testrunner.testrunnername),
        configoutputfile=config_file_full_path(env_testrunner.testrunnername),
        ebrelayer_config_file=config_file_full_path(env_ebrelayer.ebrelayername),
        sifnode_config_file=config_file_full_path(env_sifnoded.sifnodename),
        deployment_name=deployment_name,
        ethereum_config_file=config_file_full_path(env_ganache.ganachename),
        smart_contract_config_file=config_file_full_path(env_smartcontractrunner.smartcontractrunner_name),
    )
    j = env_testrunner.testrunner_config_contents(i)
    print(j)
elif component == "sifnodekeys":
    sifnodedconfig = env_utilities.read_config_file(config_file_full_path(env_sifnoded.sifnodename))["config"]
    env_sifnoded.recover_key(uuid4().hex, "test", sifnodedconfig["adminuser"]["mnemonic"])
    for v in sifnodedconfig["validators"]:
        env_sifnoded.recover_key(uuid4().hex, "test", v["mnemonic"])

# TODO
# start ganache
# start geth
# start ebrelayer
#
