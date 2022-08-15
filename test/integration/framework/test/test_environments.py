from siftool.common import *
from siftool import command, environments, project


def test_environment():
    cmd = command.Command()
    sifnoded_home_root = cmd.tmpdir("siftool.tmp")
    cmd.rmdir(sifnoded_home_root)
    cmd.mkdir(sifnoded_home_root)
    prj = project.Project(cmd, project_dir())
    prj.pkill()
    env = environments.SifnodedEnvironment2(cmd, sifnoded_home_root=sifnoded_home_root)
    env.add_validator()
    env.add_validator()
    env.init()
    env.start()
    env.add_validator()
    env.sif
    return

