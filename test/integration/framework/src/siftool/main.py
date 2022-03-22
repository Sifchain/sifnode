import sys
import time

from siftool import test_utils
from siftool.run_env import Integrator, UIStackEnvironment, Peggy2Environment, IBCEnvironment, IntegrationTestsEnvironment
from siftool.project import Project, killall, force_kill_processes
from siftool.common import *


def main(argv):
    # tmux usage:
    # tmux new-session -d -s env1
    # tmux main-pane-height -t env1 10
    # tmux split-window -h -t env1
    # tmux split-window -h -t env1
    # tmux select-layout -t env1 even-vertical
    # OR: tmux select-layout main-horizontal
    basic_logging_setup()
    what = argv[0] if argv else None
    cmd = Integrator()
    project = cmd.project
    if what == "project-init":
        project.init()
    elif what == "clean":
        project.clean()
    elif what == "build":
        project.build()
    elif what == "rebuild":
        project.rebuild()
    elif what == "project":
        return getattr(project, argv[1])(*argv[2:])
    elif what == "run-ui-env":
        e = UIStackEnvironment(cmd)
        e.stack_save_snapshot()
        e.stack_push()
    elif what == "run-env":
        if on_peggy2_branch:
            # Equivalent to future/devenv - hardhat, sifnoded, ebrelayer
            # I.e. cd smart-contracts; GOBIN=/home/anderson/go/bin npx hardhat run scripts/devenv.ts
            env = Peggy2Environment(cmd)
            processes = env.run()
        else:
            env = IntegrationTestsEnvironment(cmd)
            project.clean()
            # deploy/networks already included in run()
            processes = env.run()
            # TODO Cleanup:
            # - rm -rf test/integration/sifnoderelayerdb
            # - rm -rf networks/validators/localnet/$moniker/.sifnoded
            # - If you ran the execute_integration_test_*.sh you need to kill ganache-cli for proper cleanup
            #   as it might have been killed and started outside of our control
        if not in_github_ci:
            input("Press ENTER to exit...")
            killall(processes)
    elif what == "devenv":
        project.npx(["hardhat", "run", "scripts/devenv.ts"], cwd=project.smart_contracts_dir, pipe=False)
    elif what == "create_snapshot":
        # Snapshots are only supported in IntegrationTestEnvironment
        snapshot_name = argv[1]
        project.clean()
        env = IntegrationTestsEnvironment(cmd)
        processes = env.run()
        # Give processes some time to settle, for example relayerdb must init and create its "relayerdb"
        time.sleep(45)
        killall(processes)
        # processes1 = e.restart_processes()
        env.create_snapshot(snapshot_name)
    elif what == "restore_snapshot":
        # Snapshots are only supported in IntegrationTestEnvironment
        snapshot_name = argv[1]
        env = IntegrationTestsEnvironment(cmd)
        env.restore_snapshot(snapshot_name)
        processes = env.restart_processes()
        input("Press ENTER to exit...")
        killall(processes)
    elif what == "run-ibc-env":
        env = IBCEnvironment(cmd)
        processes = env.run()
    elif what == "run-integration-tests":
        # TODO After switching the branch,: cd smart-contracts; rm -rf node_modules; + cmd.install_smart_contract_dependencies() (yarn clean + yarn install)
        scripts = [
            "execute_integration_tests_against_test_chain_peg.sh",
            "execute_integration_tests_against_test_chain_clp.sh",
            "execute_integration_tests_against_any_chain.sh",
            "execute_integration_tests_with_snapshots.sh",
        ]
        for script in scripts:
            force_kill_processes(cmd)
            e = IntegrationTestsEnvironment(cmd)
            processes = e.run()
            cmd.execst(script, cwd=project.test_integration_dir)
            killall(processes)
            force_kill_processes(cmd)  # Some processes are restarted during integration tests so we don't own them
        log.info("Everything OK")
    elif what == "check-env":
        ctx = test_utils.get_env_ctx()
        ctx.sanity_check()
    elif what == "test-logging":
        ls_cmd = mkcmd(["ls", "-al", "."], cwd="/tmp")
        res = stdout_lines(cmd.execst(**ls_cmd))
        print(ls_cmd)
    elif what == "poc-geth":
        import geth
        g = geth.Geth(cmd)
        with open(cmd.mktempfile(), "w") as geth_log_file:
            datadir_for_running = cmd.mktempdir()
            datadir_for_keys = cmd.mktempdir()
            args = g.geth_cmd__test_integration_geth_branch(datadir=datadir_for_running)
            geth_proc = cmd.popen(args, log_file=geth_log_file)
            import hardhat
            for expected_addr, private_key in hardhat.default_accounts():
                addr = g.create_account("password", private_key, datadir=datadir_for_keys)
                assert addr == expected_addr
            input("Press ENTER to exit...")
            killall((geth_proc,))
    elif what == "inflate-tokens":
        import inflate_tokens
        inflate_tokens.run(*argv[1:])
    elif what == "recover-eth":
        test_utils.recover_eth_from_test_accounts()
    elif what == "run-peggy2-tests":
        cmd.execst(["yarn", "test"], cwd=project.smart_contracts_dir)
    elif what == "generate-python-protobuf-stubs":
        project.generate_python_protobuf_stubs()
    elif what == "localnet":
        import localnet
        localnet.run(cmd, argv[1:])
    elif what == "download-ibc-binaries":
        import localnet
        localnet.download_ibc_binaries(cmd, *argv[1:])
    else:
        raise Exception("Missing/unknown command")


if __name__ == "__main__":
    main(sys.argv[1:])
