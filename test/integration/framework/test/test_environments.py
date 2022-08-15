from siftool.common import *
from siftool import command, environments, project


def get_validators(env):
    sifnoded = env._sifnoded_for(env.node_info[0])
    return {v["description"]["moniker"]: v for v in sifnoded.query_staking_validators()}


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
    sifnoded = env._sifnoded_for(env.node_info[0])
    assert set(get_validators(env)) == set("sifnoded-{}".format(i) for i in range(3))
    return

def test_environment_2():
    cmd = command.Command()
    sifnoded_home_root = cmd.tmpdir("siftool.tmp")
    cmd.rmdir(sifnoded_home_root)
    cmd.mkdir(sifnoded_home_root)
    prj = project.Project(cmd, project_dir())
    prj.pkill()
    env = environments.SifnodedEnvironment(cmd, sifnoded_home_root=sifnoded_home_root)
    env.add_validator()
    env.init()
    env.add_validator()
    env.start()
    env.add_validator()
    sifnoded = env._sifnoded_for(env.node_info[0])
    assert set(get_validators(env)) == set("sifnoded-{}".format(i) for i in range(3))
    return
