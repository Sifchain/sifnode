import pathlib
import subprocess
import tempfile
from dataclasses import dataclass

import env_ethereum
import env_utilities
from env_utilities import wait_for_port

ganachename = "ganache"


@dataclass
class GanacheInput(env_ethereum.EthereumInput):
    block_delay: int
    mnemonic: str
    port: int
    db_dir: str


def ganache_cmd(args: GanacheInput, keysfile) -> str:
    # --db ${GANACHE_DB_DIR} --account_keys_path $GANACHE_KEYS_JSON > $GANACHE_LOG 2>&1"
    block_delay = f"-b {args.block_delay}" if args.block_delay and args.block_delay > 0 else ""
    mnemonic = args.mnemonic if args.mnemonic else "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat"
    cmd = " ".join([
        "ganache-cli",
        block_delay,
        "-h 0.0.0.0",
        f'-d --mnemonic "{mnemonic}"',
        f"--networkId {args.network_id}",
        f"--port {args.port}",
        f"--account_keys_path {keysfile}",
        f"--db {args.db_dir}",
        f"-e {args.starting_ether}"
    ])
    return cmd


def ganache_docker_compose(args: GanacheInput):
    base = env_utilities.base_docker_compose(ganachename)
    ports = [
        f"{args.port}:{args.port}",
    ]
    network = "sifchaintest"
    image = "sifdocker:latest"
    return {
        "ganache": {
            **base,
            "ports": ports,
            "networks": [network],
            "working_dir": "/sifnode/test/integration",
            "container_name": "ganache",
            "command": env_utilities.docker_compose_command("ganache")
        }
    }


def start_ganache(args: GanacheInput):
    """returns an object with a wait() method"""
    with tempfile.NamedTemporaryFile(mode="w", delete=False) as keysfile:
        cmd = ganache_cmd(args, keysfile.name)
        logparent = pathlib.Path(args.logfile).parent
        print(f"logparentis: {logparent}")
        pathlib.Path(args.logfile).parent.mkdir(exist_ok=True)
        with env_utilities.open_and_create_parent_dirs(args.logfile) as logfile:
            proc = subprocess.Popen(
                cmd,
                shell=True,
                # stdout=logfile,
                stdin=None,
                stderr=subprocess.STDOUT,
            )
            wait_for_port("localhost", args.port)
            keys = env_utilities.read_config_file(keysfile.name)
            env_utilities.startup_complete(args, keys)
            args.configoutputfile = "ethereum.json"
            env_utilities.startup_complete(args, keys)
            return proc
