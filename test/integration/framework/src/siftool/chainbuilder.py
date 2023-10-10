import sys

from siftool import common, command, environments, project, sifchain


def __brutally_terminate_processes(cmd):
    prj = project.Project(cmd, common.project_dir())
    prj.pkill()


def install_testnet(cmd: command.Command, base_dir: str, chain_id: str):
    # mnemonics = {
    #     "sif": "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow",
    #     "akasha": "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard",
    #     "juniper": "clump genre baby drum canvas uncover firm liberty verb moment access draft erupt fog alter gadget elder elephant divide biology choice sentence oppose avoid",
    #     "ethereum_root": "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat",
    # }
    external_host = "147.135.105.196"
    env = environments.SifnodedEnvironment(cmd, chain_id=chain_id, sifnoded_home_root=base_dir)
    env.add_validator(moniker="node-0", external_host=external_host, pruning="default")
    env.add_validator(moniker="node-1", external_host=external_host, pruning="nothing")
    env.add_validator(moniker="node-2", external_host=external_host, pruning="everything")
    env.add_validator(moniker="node-3", external_host=external_host, pruning="everything")
    env.init(faucet_balance={sifchain.ROWAN: 10**30, sifchain.STAKE: 10**30})
    env.start()
    sifnoded = env._client_for()
    sifnoded.wait_for_block(sifnoded.get_current_block() + 10)
    __brutally_terminate_processes(cmd)


def main(*argv):
    cmd = command.Command()
    base_dir = argv[0]
    chain_id = argv[1]
    install_testnet(cmd, base_dir, chain_id)


if __name__ == "__main__":
    main(*sys.argv[1:])
