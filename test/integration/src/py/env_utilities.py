import argparse
import json
import os
import pathlib
import socket
import tempfile
import time
from dataclasses import dataclass
from pathlib import Path


@dataclass
class SifchainCmdParameters():
    def as_json(self):
        return json.dumps(self.__dict__)


@dataclass
class SifchainCmdInput(SifchainCmdParameters):
    basedir: str
    logfile: str
    configoutputfile: str


@dataclass
class SifchainCmdOutput(SifchainCmdParameters):
    pass


def base_docker_compose(name: str):
    volumes = [
        "../..:/sifnode",
        "./configs:/configs",
        "./logs:/logs",
        "./gobin:/gobin",
        "./gocache:/gocache",
    ]
    image = "sifdocker:latest"
    return {
        "image": image,
        "volumes": volumes,
        "working_dir": "/sifnode/test/integration",
        "container_name": name,
        "command": docker_compose_command(name),
    }


def wait_for_port(host, port) -> bool:
    while True:
        try:
            with socket.create_connection((host, port)):
                return True
        except (ConnectionRefusedError, OSError) as e:
            time.sleep(1)
    return False


def atomic_write(s: str, filename: str):
    output = Path(filename)
    output.unlink(missing_ok=True)
    with tempfile.NamedTemporaryFile(mode="w", dir=os.path.dirname(output)) as temp:
        temp.write(s)
        temp.flush()
        os.link(temp.name, output)


def default_cmdline_parser():
    return argparse.ArgumentParser(
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )


def sifchain_cmd_input_parser(parser=None) -> argparse.ArgumentParser:
    """Turn command line arguments into EthereumInput"""
    if parser is None:
        parser = default_cmdline_parser()
    parser.add_argument('--logfile', required=True)
    parser.add_argument('--configoutputfile', required=True)
    return parser


def startup_complete(args, config):
    Path(args.configoutputfile).write_text(json.dumps({
        "input": args.__dict__,
        "config": config
    }, indent=2))


def docker_compose_command(component: str) -> str:
    return f"""bash -c "PATH=${{PATH}}:/gobin python3 -u src/py/env_framework.py {component}" """


def open_and_create_parent_dirs(f: str):
    pathlib.Path(f).parent.mkdir(exist_ok=True)
    return open(f, "w")


def read_json_file(json_filename):
    with open(json_filename, mode="r") as json_file:
        contents = json_file.read()
        return json.loads(contents)


def wait_for_file(f: str):
    while True:
        try:
            if os.path.getsize(f) > 0:
                return f
        except Exception as e:
            time.sleep(1)
            pass


def read_config_file(f: str):
    print(f"reading config file {f}")
    wait_for_file(f)
    return read_json_file(f)
