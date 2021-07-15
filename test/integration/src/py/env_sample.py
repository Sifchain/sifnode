import pathlib
import subprocess
import tempfile
from dataclasses import dataclass
from typing import List

import time

import env_ethereum
import env_utilities
from env_utilities import wait_for_port
from env_utilities import SifchainCmdInput, SifchainCmdOutput, sifchain_cmd_input_parser

samplename = "sample"


@dataclass
class SampleInput(SifchainCmdInput):
    shark: str


def sample_docker_compose(args: env_ethereum.EthereumInput):
    base = env_utilities.base_docker_compose(samplename)
    return {
        samplename: {
            **base,
        }
    }


def start_sample(args: SampleInput):
    """returns an object with a wait() method"""
    cmd = f"sleep infinity"
    print(f"cmd for {samplename} is {cmd}")
    # t = time.localtime()
    # print(f"time is {t}")
    proc = subprocess.Popen(
        cmd,
        shell=True,
        text=True,
        stderr = subprocess.STDOUT,
    )
    proc.wait()
    return proc