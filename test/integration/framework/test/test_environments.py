from siftool.common import *
from siftool import command, environments, project, sifchain


def get_validators(env):
    sifnoded = env._sifnoded_for(env.node_info[0])
    return {v["description"]["moniker"]: v for v in sifnoded.query_staking_validators()}


def test_transfer(env, i):
    node_info = env.node_info[i]
    sifnoded = env._sifnoded_for(node_info)
    sifnoded.send_and_check(node_info["admin_addr"], sifnoded.create_addr(), {sifchain.ROWAN: 10 ** sifchain.ROWAN_DECIMALS})


def assert_validators_working(env, expected_monikers):
    assert set(get_validators(env)) == expected_monikers
    for i in range(len(env.node_info)):
        test_transfer(env, 0)


class Test:
    def setup_method(self):
        self.cmd = command.Command()
        self.sifnoded_home_root = self.cmd.tmpdir("siftool.tmp")
        self.cmd.rmdir(self.sifnoded_home_root)
        self.cmd.mkdir(self.sifnoded_home_root)
        prj = project.Project(self.cmd, project_dir())
        prj.pkill()

    def teardown_method(self):
        prj = project.Project(self.cmd, project_dir())
        prj.pkill()

    def test_environment_setup_basic(self):
        env = environments.SifnodedEnvironment(self.cmd, sifnoded_home_root=self.sifnoded_home_root)
        env.add_validator()
        env.start()
        assert_validators_working(env, set("sifnoded-{}".format(i) for i in range(1)))

    def test_environment_mixed(self):
        env = environments.SifnodedEnvironment(self.cmd, sifnoded_home_root=self.sifnoded_home_root)
        env.add_validator()
        env.add_validator()
        env.init()
        env.start()
        env.add_validator()
        assert_validators_working(env, set("sifnoded-{}".format(i) for i in range(3)))

    def test_environment_fails_to_start_if_commission_rate_is_over_max3(self):
        env = environments.SifnodedEnvironment(self.cmd, sifnoded_home_root=self.sifnoded_home_root)
        env.add_validator(commission_rate=0.10, commission_max_rate=0.05)
        exception = None
        try:
            env.start()
        except Exception as e:
            exception = e
        assert type(exception) == sifchain.SifnodedException

    def test_need_2_out_of_3_validators_running_for_consensus(self):
        env = environments.SifnodedEnvironment(self.cmd, sifnoded_home_root=self.sifnoded_home_root)
        env.add_validator()
        env.add_validator()
        env.add_validator()
        env.add_validator()
        env.start()
        assert len(env.running_processes) == 4
        test_transfer(env, 0)  # 4 out of 4 => should work
        env.running_processes[-1].terminate()
        env.running_processes[-1].wait()
        env.open_log_files[-1].close()
        env.running_processes.pop()
        env.open_log_files.pop()
        test_transfer(env, 0)  # 3 out of 4 => should work
        env.running_processes[-1].terminate()
        env.running_processes[-1].wait()
        env.open_log_files[-1].close()
        env.running_processes.pop()
        env.open_log_files.pop()
        exception = None
        try:
            test_transfer(env, 0)  # 2 out of 4 => should fail
        except Exception as e:
            exception = e
        assert type(exception) == sifchain.SifnodedException
