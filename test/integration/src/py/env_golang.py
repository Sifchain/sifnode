import subprocess
import tempfile
from dataclasses import dataclass
from typing import List

import env_ethereum
import env_utilities
from env_utilities import wait_for_port

golangname = "golang"


@dataclass
class GolangInput(env_utilities.SifchainCmdInput):
    go_bin: str
    base_dir: str


def golang_build_cmd(args: GolangInput):
    return f"cd {args.base_dir} && GOBIN={args.go_bin} PATH=$PATH:/usr/local/go/bin make install"


def golang_build(args: GolangInput):
    cmd = golang_build_cmd(args)
    subprocess.run(
        cmd,
        shell=True
    )
    env_utilities.startup_complete(args, {})
