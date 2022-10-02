import argparse
import sys
import time

from siftool import test_utils, run_env, cosmos, diagnostics, sifchain, frontend, test_utils2
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
    log = siftool_logger(__name__)
    argparser = argparse.ArgumentParser()
    if what == "venv":
        log.info("Using Python {}.{}.{}".format(sys.version_info.major, sys.version_info.minor, sys.version_info.micro))
        log.info("sys.path={}".format(repr(sys.path)))
        log.info("Project root: {}".format(project.project_dir()))
        log.info("Project virtual environment location: {}".format(project.get_project_venv_dir()))
    elif what == "project-init":
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
        project.clean_run_env_state()
        argparser.add_argument("--type")
        args, remaining_args = argparser.parse_known_args(argv[1:])
        if args.type:
            class_name = args.type
        else:
            if on_peggy2_branch:
                class_name = "Peggy2Environment"
            else:
                class_name = "IntegrationTestsEnvironment"
        class_to_use = getattr(run_env, class_name)
        env = class_to_use(cmd)
        argparser = argparse.ArgumentParser()
        if class_to_use == Peggy2Environment:
            argparser.add_argument("--test-denom-count", type=int)
            argparser.add_argument("--geth", action="store_true", default=False)
            argparser.add_argument("--witness-count", type=int)
            argparser.add_argument("--ebrelayer-log-level")
            argparser.add_argument("--consensus-threshold", type=int)
            argparser.add_argument("--pkill", action="store_true", default=False)
            args = argparser.parse_args(remaining_args)
            if args.pkill:
                project.pkill()
            if args.witness_count is not None:
                env.witness_count = args.witness_count
            else:
                env.witness_count = int(os.getenv("WITNESS_COUNT", 2))
            if args.consensus_threshold is not None:
                env.consensus_threshold = args.consensus_threshold
            elif "CONSENSUS_THRESHOLD" in os.environ:
                env.consensus_threshold = int(os.getenv("CONSENSUS_THRESHOLD"))
            if args.ebrelayer_log_level:
                env.log_level_witness = env.log_level_relayer = args.ebrelayer_log_level
            env.use_geth_instead_of_hardhat = args.geth
            if args.test_denom_count:
                env.extra_balances_for_admin_account = {"test{}".format(i): 10**27 for i in range(args.test_denom_count)}
            hardhat_proc, sifnoded_proc, relayer0_proc, witness_procs = env.run()
            processes = [hardhat_proc, sifnoded_proc, relayer0_proc] + witness_procs
        elif class_to_use == IntegrationTestsEnvironment:
            project.clean()
            # deploy/networks already included in run()
            argparser.add_argument("--test-denom-count", type=int)
            args = argparser.parse_args(remaining_args)
            if args.test_denom_count:
                extra_balances = {"test{}".format(i): 10**27 for i in range(args.test_denom_count)}
                env.mint_amount = cosmos.balance_add(env.mint_amount, extra_balances)
            processes = env.run()
            # TODO Cleanup:
            # - rm -rf test/integration/sifnoderelayerdb
            # - rm -rf networks/validators/localnet/$moniker/.sifnoded
            # - If you ran the execute_integration_test_*.sh you need to kill ganache-cli for proper cleanup
            #   as it might have been killed and started outside of our control
        if not in_github_ci:
            wait_for_enter_key_pressed()
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
        wait_for_enter_key_pressed()
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
    elif what == "inflate-tokens":
        from siftool import inflate_tokens
        inflate_tokens.run(*argv[1:])
    elif what == "recover-eth":
        test_utils.recover_eth_from_test_accounts()
    elif what == "run-peggy2-tests":
        import glob
        from siftool.hardhat import Hardhat, default_accounts
        test_files = \
            glob.glob(os.path.join(project.smart_contracts_dir, "test", "*.js")) + \
            glob.glob(os.path.join(project.smart_contracts_dir, "test", "*.ts"))
        hardhat = Hardhat(cmd)
        # Running tests against geth is not working yet
        # hardhat_accounts = [private_key for _, private_key in default_accounts()]
        # script_runner = hardhat.script_runner("http://localhost:8545/", network="geth", accounts=hardhat_accounts)
        script_runner = hardhat.script_runner()
        script_runner.test(test_files)
    elif what == "generate-python-protobuf-stubs":
        project.generate_python_protobuf_stubs()
    elif what == "localnet":
        import localnet
        localnet.run(cmd, argv[1:])
    elif what == "download-ibc-binaries":
        import localnet
        localnet.download_ibc_binaries(cmd, *argv[1:])
    elif what == "geth":
        import siftool.geth, siftool.eth
        geth = siftool.geth.Geth(cmd)
        datadir = os.path.join(os.environ["HOME"], ".siftool-geth")
        datadir = None
        signer_addr, signer_private_key = siftool.eth.web3_create_account()
        ethereum_chain_id = 9999
        geth.init(ethereum_chain_id, [signer_addr], datadir)
    elif what == "dump-block-times":
        argparser.add_argument("--node", type=str, required=True)
        argparser.add_argument("--file", type=str, required=True)
        argparser.add_argument("--from-block", type=int)
        argparser.add_argument("--to-block", type=int)
        args = argparser.parse_args(argv[1:])
        sifnoded = sifchain.Sifnoded(cmd, node=args.node)
        from_block = args.from_block if args.from_block is not None else 1
        to_block = args.to_block if args.to_block is not None else sifnoded.get_current_block()
        block_times = diagnostics.get_block_times(sifnoded, from_block, to_block)
        block_times = [(block_times[i][0], (block_times[i][1] - block_times[i - 1][1]).total_seconds())
            for i in range(1, len(block_times))]
        lines = ["{},{:.3f}".format(t[0], t[1]) for t in block_times]
        with open(args.file, "w") as f:
            f.write(joinlines(lines))
    elif what == "create-wallets":
        argparser.add_argument("count", type=int)
        argparser.add_argument("--home", type=str)
        args = argparser.parse_args(argv[1:])
        test_utils2.PredefinedWallets.create(cmd, args.count, args.home)
    elif what == "run-ui":
        frontend.run_local_ui(cmd)
    else:
        raise Exception("Missing/unknown command")


if __name__ == "__main__":
    main(sys.argv[1:])
