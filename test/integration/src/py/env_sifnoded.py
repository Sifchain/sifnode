import json
import os
import subprocess
import tempfile
from dataclasses import dataclass
import time
import yaml

import env_utilities

sifnodename = "sifnoded"
key_password = "aaaaaaaa"


# Note that a validator doesn't have an intrinsic ethereum address and password; you assign
# it one when the validator starts
@dataclass
class SifnodedRunner(env_utilities.SifchainCmdInput):
    rpc_port: int
    n_validators: int
    chain_id: str
    network_config_file: str
    seed_ip_address: str
    bin_prefix: str
    go_build_config_path: str
    sifnode_host: str


# sequence
# rake genesis does:
# sifgen network create #{chainnet} #{validator_count} #{build_dir} #{seed_ip_address} #{network_config}")
def set_fields(args: SifnodedRunner):
    return {
        "networkdir": tempfile.mkdtemp(),
        "network_config_file": tempfile.NamedTemporaryFile(delete=False, suffix=".yml").name
    }


def sifgen_network_create_cmd(args: SifnodedRunner, fields) -> str:
    ndir = fields["networkdir"]
    nf = fields["network_config_file"]
    return f"{args.bin_prefix}/sifgen network create {args.chain_id} {args.n_validators} {ndir} {args.seed_ip_address} {nf}"


def build_chain(args: SifnodedRunner):
    # we don't care about the actual go build configuration,
    # we just need to know the build happened
    env_utilities.read_config_file(args.go_build_config_path)
    fields = set_fields(args)
    network_create_cmd = sifgen_network_create_cmd(args, fields)
    ox = subprocess.run(
        network_create_cmd,
        capture_output=True,
        shell=True
    )
    if ox.returncode != 0:
        raise Exception(f"failed to execute {network_create_cmd}: {ox}")
    with open(fields["network_config_file"]) as f:
        validators = yaml.safe_load(f)
        print(f"ncc: \n{network_create_cmd}")
    for v in validators:
        p = v["password"]
        nd = fields["networkdir"]
        base_path = os.path.join(nd, "validators", args.chain_id, v["moniker"])
        sifnodeclipath = os.path.join(base_path, ".sifnodecli")
        sifnodedpath = os.path.join(base_path, ".sifnoded")
        v["sifnodeclipath"] = sifnodeclipath
        v["sifnodedpath"] = sifnodedpath
        m = v["moniker"]
        o = subprocess.run(
            f"yes {p} | {args.bin_prefix}/sifnodecli keys show -a --bech val {m} --home {sifnodeclipath}",
            shell=True,
            text=True,
            capture_output=True
        )
        valoper = o.stdout.strip()
        v["sifvaloper"] = valoper
        subprocess.run(
            f"{args.bin_prefix}/sifnoded add-genesis-validators {valoper} --home {sifnodedpath}",
            shell=True,
            check=True
        )
    print(f"validators: \n{json.dumps(validators)}")

    # need a new account to be the administrator
    adminusercmd = f"yes | {args.bin_prefix}/sifnodecli keys add sifnodeadmin --keyring-backend test -o json"
    adminuseroutput = subprocess.run(
        adminusercmd,
        capture_output=True,
        shell=True,
        text=True,
        check=True
    )
    adminuser = json.loads(adminuseroutput.stderr)
    adminuseraddr = adminuser["address"]
    subprocess.run(
        f"{args.bin_prefix}/sifnoded add-genesis-account {adminuseraddr} 100000000000000000000rowan --home {sifnodedpath}",
        check=True,
        shell=True,
    )
    subprocess.run(
        f"{args.bin_prefix}/sifnoded set-genesis-oracle-admin {adminuseraddr} --home {sifnodedpath}",
        check=True,
        shell=True,
    )

    return {
        **fields,
        "validators": validators,
        "adminuser": adminuser,
    }


def sifnoded_docker_compose(args: SifnodedRunner):
    base = env_utilities.base_docker_compose(sifnodename)
    ports = [
        f"{args.rpc_port}:{args.rpc_port}",
    ]
    network = {
        "sifchaintest": {
            "ipv4_address": "10.0.0.30"
        }
    }

    return {
        sifnodename: {
            **base,
            "ports": ports,
            "networks": network,
            "command": env_utilities.docker_compose_command("startsifnoded"),
        }
    }


def run(args: SifnodedRunner, sifnoded_chain_data):
    p = args.rpc_port
    for v in sifnoded_chain_data["validators"]:
        sndp = v["sifnodedpath"]
        addr = v["address"]
        cmd = f'{args.bin_prefix}/sifnoded start --minimum-gas-prices 0.5rowan --rpc.laddr tcp://0.0.0.0:{p} --home {sndp}'
        sifnoded_process = subprocess.Popen(cmd, shell=True)

        # It takes a while before the validator account is available, so need to wait for that
        subprocess.run(
            f"python3 src/py/wait_for_sif_account.py 0 {addr}",
            shell=True
        )

        env_utilities.startup_complete(args, sifnoded_chain_data)

        sifnoded_process.wait()


def export_key(
        key_name: str,
        backend: str,
        password: str
):
    cmd = f'yes {password} | sifnodecli keys export --keyring-backend {backend} {key_name}'
    output = subprocess.run(
        cmd,
        shell=True,
        text=True,
        capture_output=True,
        check=True
    )
    key_contents = output.stderr
    print(f"outputis: {json.dumps(output.__dict__, indent=2)}")
    output = import_key("mykey", "test", key_password, key_contents)
    print(f"outputis: {json.dumps(output.__dict__, indent=2)}")
    return output.stdout


def import_key(
        key_name: str,
        backend: str,
        password: str,
        key_contents: str
):
    # note that you can also get keys with
    # sifnodecli keys add --keyring-backend test --recover fnord
    # using a mnemonic
    with tempfile.NamedTemporaryFile(mode="w", delete=False) as keyfile:
        keyfile.write(key_contents)
        cmd = f'yes {password} | sifnodecli keys import {key_name} --keyring-backend {backend} {keyfile.name}'
        output = subprocess.run(
            cmd,
            shell=True,
            text=True,
            capture_output=True,
            check=True
        )
        print(f"outputis: {json.dumps(output.__dict__, indent=2)}")
        return output.stdout


def recover_key(
        key_name: str,
        backend: str,
        mnemonic: str,
):
    cmd = f'echo "{mnemonic}" | sifnodecli keys add --keyring-backend {backend} {key_name} --recover'
    subprocess.run(
        cmd,
        shell=True,
        check=True
    )