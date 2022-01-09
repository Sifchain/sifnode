import os
import sys

# Temporary workaround to include integration framework
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), *([os.path.pardir] * 4)))
base_dir = os.path.join(project_root, "test", "integration", "framework")
enabled = False
for p in sys.path:
    enabled = enabled or os.path.realpath(p) == os.path.realpath(base_dir)
if not enabled:
    sys.path = sys.path + [base_dir]

import command
import cosmos
import eth
import main
import common
import project
import geth
import hardhat
import truffle
import test_utils
import inflate_tokens
import sifchain
