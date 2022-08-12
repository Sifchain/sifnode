from siftool.common import *
from siftool import command, environments

def test_environment():
    cmd = command.Command()
    env = environments.SifnodedEnvironment(cmd)
    # env.add_validator()