import sys
import argparse
import json
import time
import siftool_path
from siftool.common import *
from siftool.sifchain import ROWAN, STAKE
from siftool import command, project, environments, sifchain


log = siftool_logger(__name__)


# Min commission / max voting power
# Design document: https://github.com/Sifchain/sifnode/blob/feature/min-commission/docs/tutorials/commission.md
# Useful info:
# - https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/sifchain-chainops/200




def should_not_add_validator_with_commission_less_than_5_percent(cmd: command.Command, prj: project.Project,
    commission_rate: float
):
    prj.pkill()
    time.sleep(2)

    tmpdir = cmd.tmpdir("siftool.tmp", "test_max_voting_power")
    cmd.rmdir(tmpdir)
    cmd.mkdir(tmpdir)

    chain_id = "localnet"

    sifnoded_akasha = sifchain.Sifnoded(cmd, chain_id=chain_id)
    cmd.rmdir(sifnoded_akasha.get_effective_home())
    akasha_name = "akasha"
    akasha = sifnoded_akasha.keys_add(akasha_name)
    akasha_addr = akasha["address"]
    sifnoded_akasha.init(akasha_name)
    akasha_pubkey = sifnoded_akasha.tendermint_show_validator()

    stake = {STAKE: 92 * 10**21}

    env = environments.SifnodedEnvironment(cmd)
    env.chain_id = chain_id
    env.sifnoded_home_root = tmpdir
    env.node_external_ip_address = LOCALHOST
    env.staking_denom = STAKE
    env.admin0_stake = stake
    env.validator_account_balance = {ROWAN: 10**30, STAKE: 92* 10**21}
    env.genesis_balances[akasha_addr] = {ROWAN: 10**30, STAKE: 92* 10**21}
    env.init()

    sifnoded = env.sifnoded[0]
    validators_before = sifnoded.query_staking_validators()

    assert len(validators_before) == 1

    sifnoded_akasha = sifchain.Sifnoded(cmd, home=sifnoded_akasha.home, chain_id=sifnoded.chain_id, node=sifnoded.node)
    sifnoded_akasha.staking_create_validator(stake, akasha_pubkey, akasha_name,
        commission_rate,
        0.20,  # commision_max_rate
        0.10,  # comission_max_change_rate
        1000000,  # min_self_delegation
        akasha_addr
    )
    sifnoded.wait_for_last_transaction_to_be_mined()

    validators_after = sifnoded.query_staking_validators()
    assert len(validators_after) == 2

    pass


def main(argv: List[str]):
    basic_logging_setup()
    parser = argparse.ArgumentParser()
    args = parser.parse_args(argv[1:])

    cmd = command.Command()
    prj = project.Project(cmd, project_dir())
    # should_not_add_validator_with_commission_less_than_5_percent(cmd, prj, 0.10)  # Should work
    should_not_add_validator_with_commission_less_than_5_percent(cmd, prj, 0.03)  # Should fail


if __name__ == "__main__":
    main(sys.argv)
