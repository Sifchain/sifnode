import dataclasses
import json
import os
import sys
import time

import yaml

import env_ebrelayer
import env_ganache
import env_geth
import env_sifnoded
import env_smartcontractrunner
import env_utilities
from env_geth import start_geth
import os

configbase = "/configs"
logbase = "/logs/"

ganachename = "ganache"
gethname = "geth"
basedir = "/sifnode"
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
    http_port=7445,
    ethereum_addresses=4,
    configoutputfile=config_file_full_path(gethname)
)

ganache_ws_port = 7545
ganache_ws_addr = f"ws://ganache:{ganache_ws_port}"
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
    validator_ethereum_credentials=None,
    n_validators=n_validators,
    validator_powers=[100],
    consensus_threshold=100
)

sifnoded_input = env_sifnoded.SifnodedRunner(
    basedir=basedir,
    bin_prefix=os.path.join("/gobin"),
    logfile=log_file_full_path(ganachename),
    configoutputfile=config_file_full_path(env_sifnoded.sifnodename),
    rpc_port=26657,
    chain_id="localnet",
    network_config_file="/tmp/netconfig.yml",
    seed_ip_address="10.10.1.1",
    n_validators=n_validators
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
    }
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
elif component == "deploy_contracts":
    time.sleep(10)
    print(f"cffis: {smartcontractrunner_input.ethereum_config_file}")
    ethereum_config = env_utilities.read_config_file(smartcontractrunner_input.ethereum_config_file)
    print(f"ethereumconfig: {json.dumps(ethereum_config, indent=2)}")
    private_keys_stanza = ethereum_config["config"]["private_keys"]
    private_keys = list(private_keys_stanza.values())
    print(f"kis: {private_keys[0]}")
    i = dataclasses.replace(
        smartcontractrunner_input,
        operator_private_key=private_keys[0],
        validator_ethereum_credentials=list(private_keys_stanza.items())[1:smartcontractrunner_input.n_validators + 1]
    )
    env_smartcontractrunner.deploy_contracts(i)
elif component == "startsifnoded":
    f = config_file_full_path(env_smartcontractrunner.smartcontractrunner_name)
    print(f"about to read {f}")
    smart_contract_config = env_utilities.read_config_file(f)
    print(f"smart contract config: {json.dumps(smart_contract_config, indent=2)}")
    rslt = env_sifnoded.build_chain(sifnoded_input)
    print(f"build_chain result: \n{json.dumps(rslt)}")
    env_sifnoded.run(sifnoded_input, rslt)
    for v in rslt["validators"]:
        print(f"sifnode validator:\n{json.dumps(v, indent=2)}")
        x = env_ebrelayer.EbrelayerInput(
            basedir=basedir,
            logfile=log_file_full_path(env_ebrelayer.ebrelayername),
            configoutputfile=config_file_full_path(env_ebrelayer.ebrelayername),
            ethereum_private_key=smart_contract_config["input"]["validator_ethereum_credentials"][0][1],
            web3_provider="ws://ganache:7545",
            tendermint_node="tcp://0.0.0.0:26657",
            bridge_registry_address="0xf204a4Ef082f5c04bB89F7D5E6568B796096735a",
            moniker=v["moniker"],
            mnemonic=v["mnemonic"],
            chain_id="localnet",
            home_dir=v["sifnodeclipath"],
            gas="5000000000000",
            gas_prices="0.5rowan",
        )
        env_ebrelayer.run(x)
    time.sleep(10000)

# TODO
# start ganache
# start geth
# start ebrelayer
#