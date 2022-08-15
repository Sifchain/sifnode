from siftool.common import *
from siftool import command, environments, project, sifchain


def get_validators(env):
    sifnoded = env._sifnoded_for(env.node_info[0])
    return {v["description"]["moniker"]: v for v in sifnoded.query_staking_validators()}


def assert_validators_working(env, expected_monikers):
    assert set(get_validators(env)) == set("sifnoded-{}".format(i) for i in range(3))
    for i in range(len(env.node_info)):
        node_info = env.node_info[i]
        sifnoded = env._sifnoded_for(node_info)
        sifnoded.send_and_check(node_info["admin_addr"], sifnoded.create_addr(), {sifchain.ROWAN: 10**sifchain.ROWAN_DECIMALS})

def test_environment_1():
    cmd = command.Command()
    sifnoded_home_root = cmd.tmpdir("siftool.tmp")
    cmd.rmdir(sifnoded_home_root)
    cmd.mkdir(sifnoded_home_root)
    prj = project.Project(cmd, project_dir())
    prj.pkill()
    env = environments.SifnodedEnvironment(cmd, sifnoded_home_root=sifnoded_home_root)
    env.add_validator()
    env.add_validator()
    env.init()
    env.start()
    env.add_validator()
    assert_validators_working(env, set("sifnoded-{}".format(i) for i in range(3)))

def test_environment_2():
    cmd = command.Command()
    sifnoded_home_root = cmd.tmpdir("siftool.tmp")
    cmd.rmdir(sifnoded_home_root)
    cmd.mkdir(sifnoded_home_root)
    prj = project.Project(cmd, project_dir())
    prj.pkill()
    env = environments.SifnodedEnvironment(cmd, sifnoded_home_root=sifnoded_home_root)
    env.add_validator()
    env.start()
    env.add_validator()
    assert_validators_working(env, set("sifnoded-{}".format(i) for i in range(3)))
