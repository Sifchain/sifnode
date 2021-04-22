import subprocess
import tempfile
from dataclasses import dataclass
from typing import List

import env_ethereum
import env_utilities
from env_utilities import wait_for_port

gethname = "geth"


@dataclass
class GethInput(env_ethereum.EthereumInput):
    http_port: int
    ws_port: int
    ethereum_addresses: int


def geth_cmd(args: env_ethereum.EthereumInput) -> str:
    apis = "personal,eth,net,web3,debug"
    cmd = " ".join([
        "geth",
        "--datadir /tmp/gethdata",
        f"--networkid {args.network_id}",
        f"--ws --ws.addr 0.0.0.0 --ws.port {args.ws_port} --ws.api {apis}",
        f"--http --http.addr 0.0.0.0 --http.port {args.http_port} --http.api {apis}",
        "--dev --dev.period 1",
        "--mine --miner.threads=1",
    ])
    return cmd


def create_initial_accounts(n: int):
    return list(map(lambda _: create_account(), range(n)))


def fund_initial_accounts(addresses: List[str], starting_amount: int):
    quote = '"'
    print(f"addresses: {addresses}")
    for addr in addresses:
        quotedaddr = f"\\{quote}{addr}\\{quote}"
        cmd = f'geth attach /tmp/gethdata/geth.ipc --exec "eth.sendTransaction({{from:eth.coinbase, to:{quotedaddr}, value:{starting_amount * 10 ** 18}}})"'
        subprocess.run(cmd, shell=True, check=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for addr in addresses:
        quotedaddr = f"\\{quote}{addr}\\{quote}"
        while True:
            cmd = f'geth attach /tmp/gethdata/geth.ipc --exec "eth.getBalance({quotedaddr})"'
            balance_result = subprocess.run(
                cmd,
                check=True,
                text=True,
                shell=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                timeout=10,
            )
            print(f"getbal: {balance_result}")
            balance = int(float(balance_result.stdout))
            if balance >= starting_amount:
                break;


def geth_docker_compose(args: env_ethereum.EthereumInput):
    base = env_utilities.base_docker_compose(gethname)
    ports = [
        f"{args.ws_port}:{args.ws_port}",
        f"{args.http_port}:{args.http_port}",
    ]
    network = "sifchaintest"
    return {
        gethname: {
            **base,
            "ports": ports,
            "networks": [network],
            "working_dir": "/sifnode/test/integration",
        }
    }


def format_new_accounts(accts):
    result = {}
    for a in accts:
        public_address, private_key = a
        result[public_address] = private_key
    return {
        "private_keys": result
    }


def start_geth(args: GethInput):
    """returns an object with a wait() method"""
    print("starting geth")
    cmd = geth_cmd(args)
    logfile = open(args.logfile, "w")
    proc = subprocess.Popen(
        cmd,
        shell=True,
        stdout=logfile,
        stdin=None,
        stderr=subprocess.STDOUT,
    )
    wait_for_port("localhost", args.ws_port)
    wait_for_port("localhost", args.http_port)
    new_accounts = create_initial_accounts(args.ethereum_addresses)
    fund_initial_accounts(map(lambda a: a[0], new_accounts), args.starting_ether)
    env_utilities.startup_complete(args, format_new_accounts(new_accounts))
    return proc


def create_account() -> (str, str):
    """returns a pair of public_address, private_key"""
    bad_password = "notasecret"
    with tempfile.NamedTemporaryFile(mode="w", delete=False) as tf:
        print(bad_password, file=tf, flush=True)
        cmd = f"geth account new --password {tf.name}"
        output = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            check=True,
            shell=True
        )
        s = output.stdout
        lines = s.split("\n")
        for l in lines:
            if "Public address of the key: " in l:
                _, public = l.split(": ")
            if "Path of the secret key file: " in l:
                _, keyfilepath = l.split(": ")
        print(f"rslt: {output.stdout} | {output.stderr}")

        cmd = f"web3 account extract --keyfile {keyfilepath} --password {bad_password}"
        output = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            shell=True
        )
        print(f"outputis: {output}")
        s = output.stdout
        lines = s.split("\n")
        for l in lines:
            if "Private key: " in l:
                _, private_key = l.split(": ")
            if "Public address: " in l:
                _, public_address = l.split(": ")
        return public_address, private_key
