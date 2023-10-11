import sys

from siftool import common, command, environments, project, sifchain


def __brutally_terminate_processes(cmd):
    prj = project.Project(cmd, common.project_dir())
    prj.pkill()


def install_testnet(cmd: command.Command, base_dir: str, chain_id: str):
    faucet_mnemonic = "fiction cousin fragile allow fruit slogan useless sting exile virus scale dress fatigue eight clay sort tape between cargo flag civil rude umbrella sing".split()
    node0_admin_mnemonic = "frog skin business valve fish fat glory syrup chicken skin slow ensure sun luggage wild click into paper swamp car ecology infant thought squeeze".split()
    node1_admin_mnemonic = "system faculty master promote among arrive dose zone cream fame barrel warm slice please creek puzzle boat excess rain lonely cupboard flame punch shed".split()
    node2_admin_mnemonic = "box fix inmate zoo night model inject gesture inquiry slice treat curve reopen portion absent adjust toilet lyrics resist same goddess dad damage sort".split()
    node3_admin_mnemonic = "hundred usual invite burger chat final collect acquire magnet repair upon venue initial ride street other tail vanish bicycle soap icon camp tragic material".split()
    external_host = "147.135.105.196"
    env = environments.SifnodedEnvironment(cmd, chain_id=chain_id, sifnoded_home_root=base_dir)
    env.add_validator(moniker="node-0", admin_mnemonic=node0_admin_mnemonic, external_host=external_host, pruning="default")
    env.add_validator(moniker="node-1", admin_mnemonic=node1_admin_mnemonic, external_host=external_host, pruning="nothing")
    env.add_validator(moniker="node-2", admin_mnemonic=node2_admin_mnemonic, external_host=external_host, pruning="everything")
    env.add_validator(moniker="node-3", admin_mnemonic=node3_admin_mnemonic, external_host=external_host, pruning="everything")
    env.init(faucet_balance={sifchain.ROWAN: 10**30, sifchain.STAKE: 10**30}, faucet_mnemonic=faucet_mnemonic)
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
