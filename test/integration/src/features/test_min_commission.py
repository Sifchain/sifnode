import sys

import siftool_path
from siftool.common import *
from siftool.sifchain import ROWAN, ROWAN_DECIMALS, STAKE
from siftool import command, cosmos, project, environments, sifchain


log = siftool_logger(__name__)


MIN_COMISSION = 0.05
MAX_VOTING_POWER = 0.066


# Min commission / max voting power
# Design document: https://github.com/Sifchain/sifnode/blob/feature/min-commission/docs/tutorials/commission.md
# Also: https://www.notion.so/sifchain/Minimum-Commissions-Max-Voting-Power-Test-Scenarios-Draft-729620045e2d41f8b18f3a5df28b623b
# Useful info:
# - https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/sifchain-chainops/200
# Upgrades:
# - https://github.com/Sifchain/sifchain-devops/blob/main/scripts/sifnode/release/testing/upgrade_path.json
# - https://github.com/Sifchain/sifnode/blob/68f69eb7e390363f336ec7a235ab7e564bf5dabb/scripts/upgrade-integration.sh#L39-L39

def should_not_add_validator_with_commission_less_than_5_percent(cmd: command.Command, prj: project.Project):
    # Setup local environment (by default with a single validator)
    env = environments.SifnodedEnvironment(cmd)
    env.staking_denom = STAKE
    env.validator_account_balance = {ROWAN: 10**30, STAKE: 10**30}
    env.init()
    env.start()

    sifnoded = exactly_one(env.sifnoded)  # Use the initial validator (only one at this point)
    validators_before = sifnoded.query_staking_validators()
    assert len(validators_before) == 1

    # This should succeed since the commission rate is higher than minimal (5%)
    akasha_index = env.add_validator(moniker="akasha", extra_funds={ROWAN: 10**25}, commission_rate=0.10)
    akasha_sifnoded = env.sifnoded[akasha_index]
    akasha_info = env.node_info[akasha_index]
    akasha_admin_addr = akasha_info["admin_addr"]

    validators_after = sifnoded.query_staking_validators()
    assert len(validators_after) == 2
    assert "akasha" in {v["description"]["moniker"] for v in validators_after}

    # This should fail since the commission rate is higher than minimal (5%)
    exception = None
    try:
        env.add_validator(moniker="juno", extra_funds={ROWAN: 10**25}, commission_rate=0.03)
    except Exception as e:
        exception = e
    assert type(exception) == sifchain.SifnodedException
    assert exception.message == 'validator commission 0.030000000000000000 cannot be lower than minimum of 0.050000000000000000: invalid request'

    assert len(sifnoded.query_staking_validators()) == 2  # Cross check

    # Try to change the first validator to 3%. Since this is less than allowed 3%, it should fail

    res = akasha_sifnoded.staking_edit_validator(0.30, akasha_admin_addr, broadcast_mode="block")
    sifchain.check_raw_log(res)


def main(argv: List[str]):
    basic_logging_setup()
    cmd = command.Command()
    prj = project.Project(cmd, project_dir())
    # Kill off any sifnoded processes running from before
    prj.pkill()
    time.sleep(2)
    should_not_add_validator_with_commission_less_than_5_percent(cmd, prj)


if __name__ == "__main__":
    main(sys.argv)
