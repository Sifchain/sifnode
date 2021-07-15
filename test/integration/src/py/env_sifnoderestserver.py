from typing import List
import yaml
import subprocess
from dataclasses import dataclass

import env_ethereum
import env_utilities
from env_utilities import wait_for_port
import env_ethereum


@dataclass
class SifnodeRestServer(env_utilities.SifchainCmdInput):
    port: int


def sifnoderestserver_docker_compose(args: env_ethereum.EthereumInput):
    ports = [
        f"{args.port}:{args.port}",
    ]
    network = "sifchaintest"
    volumes = [
        "../..:/sifnode"
    ]
    name = "sifnoderestserver"
    image = "sifdocker:latest"
    return {
        name: {
            "image": image,
            "ports": ports,
            "networks": [network],
            "volumes": volumes,
            "working_dir": "/sifnode/test/integration",
            "container_name": name,
            "command": env_utilities.docker_compose_command(name)
        }
    }


def start_sifnoderestserver(conf: SifnodeRestServer):
    cmd = sifnoderestserver_cmd(conf)
    subprocess.run(

    )

def sifnoderestserver_cmd(conf: SifnodeRestServer):
    return f"sifnodecli rest-server --laddr tcp://0.0.0.0:{conf.port}"


def deploy_contracts_cmd(args: SmartContractDeployInput):
    return f"cd {args.basedir}/smart-contracts && npx truffle deploy --network {args.truffle_network} --reset"


def deploy_contracts(args: SmartContractDeployInput):
    cmd = deploy_contracts_cmd(args)
    result = subprocess.run(
        cmd,
        shell=True,
        stdout=args.logfile,
        stdin=None,
        stderr=subprocess.STDOUT,
    )

