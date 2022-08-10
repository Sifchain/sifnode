import sys
import re
from typing import Tuple

import siftool_path
from siftool.common import *
from siftool.sifchain import ROWAN, ROWAN_DECIMALS, STAKE
from siftool import command, cosmos, project, environments, sifchain


# How to use:
# Install Python 3.8-3.10
# Run test/integration/framework/siftool venv
# Prepare sifnoded binaries
# - Compile sifnoded binary for v0.14.0 and put them into test/integration/framework/build/versions/0.14.0/sifnoded
# - Compile sifnoded binary for v0.15.0-rc.1 and put them into test/integration/framework/build/versions/0.15.0-rc.1/sifnoded
# - For exact versions se OLD_VERSION and NEW_VERSION above
# Run test/integration/framework/venv/bin/python3 test/integration/src/features/test_min_commission.py
# To watch live logs: tail -F /tmp/siftool.tmp/test_min_commission/sifnoded-0/sifnoded.log
#
# More information about min commission / max voting power:
# Test scenarios (Kevin): https://github.com/Sifchain/sifnode/blob/feature/min-commission/docs/tutorials/commission.md
# Test scenarios (James): https://www.notion.so/sifchain/Minimum-Commissions-Max-Voting-Power-Test-Scenarios-Draft-729620045e2d41f8b18f3a5df28b623b
# Useful info:
# - https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/sifchain-chainops/200
# Upgrades:
# - https://github.com/Sifchain/sifchain-devops/blob/main/scripts/sifnode/release/testing/upgrade_path.json
# - https://github.com/Sifchain/sifnode/blob/68f69eb7e390363f336ec7a235ab7e564bf5dabb/scripts/upgrade-integration.sh#L39-L39


log = siftool_logger(__name__)

MIN_COMISSION = 0.05
MAX_VOTING_POWER = 0.066

OLD_VERSION = "0.14.0"
NEW_VERSION = "0.15.0-rc.1"


def get_binary_for_version(label):
    return project_dir("test", "integration", "framework", "build", "versions", label, "sifnoded")


def create_environment(cmd, version, commission_rate=0.06, commission_max_rate=0.10, commission_max_change_rate=0.05,
    default_staking_amount: int = 92 * 10**21
):
    home_root = "/tmp/siftool.tmp/test_min_commission"
    cmd.rmdir(home_root)
    cmd.mkdir(home_root)
    env = environments.SifnodedEnvironment(cmd)
    env.staking_denom = STAKE
    env.default_validator_balance = {ROWAN: 10**25, STAKE: 10**25}
    env.default_binary = get_binary_for_version(version)
    env.default_commission_rate = commission_rate
    env.default_commission_max_rate = commission_max_rate
    env.default_commission_max_change_rate = commission_max_change_rate
    env.default_staking_amount = default_staking_amount
    env.sifnoded_home_root = home_root
    env.init()
    env.start()
    return env


def upgrade(env, new_version):
    sifnoded = env.sifnoded[0]
    admin_addr = env.node_info[0]["admin_addr"]

    # Whoever makes the proposal has to put in  deposit.
    # Deposit must be >= genesis::app_state.gov.deposit_params.min_deposit
    deposit = {env.staking_denom: env.default_staking_amount}
    env.fund(admin_addr, deposit)

    upgrade_info = "{\"binaries\":{\"linux/amd64\":\"url_with_checksum\"}}"
    upgrade_height = env.sifnoded[0].get_current_block() + 15  # Note: must be > 60s (as per app config)

    proposals_before = sifnoded.query_gov_proposals()
    res = sifnoded.gov_submit_software_upgrade(NEW_VERSION, admin_addr, deposit, upgrade_height, upgrade_info,
        "test_release", "Test Release", broadcast_mode="block"
    )
    sifchain.check_raw_log(res)
    sifnoded.wait_for_last_transaction_to_be_mined()
    proposals_after = sifnoded.query_gov_proposals()
    new_proposal_ids = {p["proposal_id"] for p in proposals_after}.difference({p["proposal_id"] for p in proposals_before})
    active_proposal = exactly_one([p for p in proposals_after if p["proposal_id"] in new_proposal_ids])
    proposal_id = int(active_proposal["proposal_id"])

    res = sifnoded.gov_vote(1, True, admin_addr, broadcast_mode="block")
    sifchain.check_raw_log(res)

    sifnoded.wait_for_block(upgrade_height)
    time.sleep(5)
    for p in env.running_processes:
        p.terminate()
        p.wait()
    for f in env.open_log_files:
        f.close()
    time.sleep(5)
    sifnoded.binary = get_binary_for_version(new_version)
    env._sifnoded_start(0)
    assert sifnoded.version() == new_version


def delegate(env, from_index, to_index, amount):
    from_validator_node_info = env.node_info[from_index]
    to_validator_node_info = env.node_info[to_index]
    sifnoded_tmp = env.sifnoded_from_to(from_validator_node_info, to_validator_node_info)
    validator_addr = env.sifnoded[to_index].get_val_address(to_validator_node_info["admin_addr"])
    from_addr = from_validator_node_info["admin_addr"]
    res = sifnoded_tmp.staking_delegate(validator_addr, {env.staking_denom: amount}, from_addr, broadcast_mode="block")
    sifchain.check_raw_log(res)


def should_not_add_validator_with_commission_less_than_5_percent(cmd: command.Command, prj: project.Project):
    # Min commission - blocking MsgCreateValidator messages
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
    assert sifchain.is_min_commission_too_low_exception(exception)

    assert len(sifnoded.query_staking_validators()) == 2  # Cross check

    # Min commission - blocking MsgEditValidator messages

    # Try to change the first validator to 3%. Since this is less than allowed 3%, it should fail
    res = akasha_sifnoded.staking_edit_validator(0.30, akasha_admin_addr, broadcast_mode="block")
    sifchain.check_raw_log(res)


def test_min_commission_create_new_validator(cmd: command.Command, prj: project.Project):
    env = create_environment(cmd, NEW_VERSION)

    validator1 = env.add_validator(commission_rate=0.05)

    exception = None
    try:
        validator2 = env.add_validator(commission_rate=0.03)
    except Exception as e:
        exception = e
    assert sifchain.is_min_commission_too_low_exception(exception)

    validator3 = env.add_validator(commission_rate=0.07)


def test_min_commission_modify_existing_validator(cmd: command.Command, prj: project.Project):
    # Using defaults: commission_rate=0.06, commission_max_rate=0.10, commission_max_change_rate=0.05
    # We create 3 validators for 3 different test cases so that we only have to wait once
    env = create_environment(cmd, NEW_VERSION)
    env.add_validator()
    env.add_validator()

    sifnoded0 = env.sifnoded[0]
    admin0_addr = env.node_info[0]["admin_addr"]
    sifnoded1 = env.sifnoded[1]
    admin1_addr = env.node_info[1]["admin_addr"]
    sifnoded2 = env.sifnoded[2]
    admin2_addr = env.node_info[2]["admin_addr"]

    # Commission cannot be changed more than once in 24h.
    log.info("Sleeping for 24h...")
    time.sleep(24 * 3600 + 5 * 60)  # 1 day + 5 minutes
    log.info("Sleep over")

    exception = None
    try:
        res = sifnoded0.staking_edit_validator(0.05, from_acct=admin0_addr, broadcast_mode="block")
        sifchain.check_raw_log(res)
    except Exception as e:
        exception = e
    assert exception is None

    exception = None
    try:
        res = sifnoded1.staking_edit_validator(0.03, from_acct=admin1_addr, broadcast_mode="block")
        sifchain.check_raw_log(res)
    except Exception as e:
        exception = e
    assert sifchain.is_min_commission_too_low_exception(exception)

    try:
        res = sifnoded2.staking_edit_validator(0.07, from_acct=admin2_addr, broadcast_mode="block")
        sifchain.check_raw_log(res)
    except Exception as e:
        exception = e
    assert exception is None


def test_min_commission_upgrade_handler(cmd: command.Command, prj: project.Project):
    env = create_environment(cmd, OLD_VERSION, commission_rate=0.03, commission_max_rate=0.04, commission_max_change_rate=0.01)

    commission_rates_before = exactly_one(env.sifnoded[0].query_staking_validators())["commission"]["commission_rates"]
    assert float(commission_rates_before["rate"]) == 0.03
    assert float(commission_rates_before["max_rate"]) == 0.04
    assert float(commission_rates_before["max_change_rate"]) == 0.01

    upgrade(env, NEW_VERSION)

    commission_rates_after = exactly_one(env.sifnoded[0].query_staking_validators())["commission"]["commission_rates"]
    assert float(commission_rates_after["rate"]) == 0.05
    assert float(commission_rates_after["max_rate"]) == 0.05
    assert float(commission_rates_after["max_change_rate"]) == 0.01


def test_max_voting_power(cmd: command.Command, prj: project.Project):
    stake0 = 1000 * 10**21
    stake1 = 62 * 10**21
    stake2 = 100 * 10**21
    stake3 = 3 * 10**21

    env = create_environment(cmd, NEW_VERSION, default_staking_amount=stake0)
    env.default_staking_amount = 10**24

    sifnoded = env.sifnoded[0]

    sif_val = env.sifnoded[0]
    akasha_val = env.add_validator(staking_amount=stake1)

    sif_validator_index = 0
    akasha_validator_index = 1

    validator_powers_0 = sifnoded.query_staking_validators()

    exception = None
    try:
        delegate(env, akasha_validator_index, sif_validator_index, stake2)
    except Exception as e:
        exception = e
    assert sifchain.is_max_voting_power_limit_exceeded_exception(exception)

    validator_powers_1 = sifnoded.query_staking_validators()

    # First time it's OK
    delegate(env, sif_validator_index, akasha_validator_index, stake3)

    validator_powers_2 = sifnoded.query_staking_validators()

    # Second time we're already over 6.6%
    exception = None
    try:
        delegate(env, akasha_validator_index, sif_validator_index, stake3)
    except Exception as e:
        exception = e
    assert sifchain.is_max_voting_power_limit_exceeded_exception(exception)

    validator_powers_3 = sifnoded.query_staking_validators()

    return


def main(argv: List[str]):
    basic_logging_setup()
    cmd = command.Command()

    # Check versions
    for version in [OLD_VERSION, NEW_VERSION]:
        reported_version = sifchain.Sifnoded(cmd, binary=get_binary_for_version(version)).version()
        assert reported_version == version

    prj = project.Project(cmd, project_dir())
    # Kill off any sifnoded processes running from before
    prj.pkill()
    time.sleep(2)

    if argv == ["24"]:
        log.info("24h test")
        test_min_commission_modify_existing_validator(cmd, prj)
    else:
        test_min_commission_create_new_validator(cmd, prj)

        prj.pkill()
        time.sleep(2)

        test_min_commission_upgrade_handler(cmd, prj)

        prj.pkill()
        time.sleep(2)

        test_max_voting_power(cmd, prj)

    prj.pkill()
    time.sleep(2)

    log.info("Finished successfully")


if __name__ == "__main__":
    main(sys.argv[1:])
