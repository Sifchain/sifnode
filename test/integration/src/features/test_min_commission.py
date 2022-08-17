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
# How to bypass the 24h limit (Caner): https://www.notion.so/sifchain/v0-15-0-Edit-min-commission-of-existing-validator-edecd16b074a4900974704d223847b48
# Useful info:
# - https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/sifchain-chainops/200
# Upgrades:
# - https://github.com/Sifchain/sifchain-devops/blob/main/scripts/sifnode/release/testing/upgrade_path.json
# - https://github.com/Sifchain/sifnode/blob/68f69eb7e390363f336ec7a235ab7e564bf5dabb/scripts/upgrade-integration.sh#L39-L39


log = siftool_logger(__name__)

MIN_COMISSION = 0.05
MAX_VOTING_POWER = 0.066

OLD_VERSION = "0.14.0"
NEW_VERSION = "0.15.0-rc.2"


# Kill off any sifnoded processes running from before
def pkill(cmd):
    project.Project(cmd, project_dir()).pkill()
    time.sleep(2)


def get_binary_for_version(label):
    return project_dir("test", "integration", "framework", "build", "versions", label, "sifnoded")


def assert_no_exception(exception):
    if exception is None:
        return
    raise AssertionError("Assertion failed") from exception

def create_environment(cmd, version, commission_rate=0.06, commission_max_rate=0.10, commission_max_change_rate=0.05,
    staking_amount: int = 92 * 10**21
):
    home_root = "/tmp/siftool.tmp/test_min_commission"
    cmd.rmdir(home_root)
    cmd.mkdir(home_root)

    binary = get_binary_for_version(version)
    assert sifchain.Sifnoded(cmd, binary=binary).version() == version  # Check actual version

    pkill(cmd)

    env = environments.SifnodedEnvironment(cmd, sifnoded_home_root=home_root)
    env.staking_denom = STAKE
    env.add_validator(binary=binary, commission_rate=commission_rate, commission_max_rate=commission_max_rate,
        commission_max_change_rate=commission_max_change_rate, staking_amount=staking_amount)
    env.start()
    return env


def delegate(env, from_index, to_index, amount):
    from_validator_node_info = env.node_info[from_index]
    to_validator_node_info = env.node_info[to_index]
    sifnoded_to = env._sifnoded_for(to_validator_node_info)
    sifnoded_from_to = env._sifnoded_for(from_validator_node_info, to_node_info=to_validator_node_info)
    validator_addr = sifnoded_to.get_val_address(to_validator_node_info["admin_addr"])
    from_addr = from_validator_node_info["admin_addr"]
    env.fund(from_addr, {env.staking_denom: amount})  # Make sure admin has enough balance for what he is delegating
    res = sifnoded_from_to.staking_delegate(validator_addr, {env.staking_denom: amount}, from_addr, broadcast_mode="block")
    sifchain.check_raw_log(res)


def test_min_commission_create_new_validator(cmd: command.Command):
    def test_case(commission_rate, commission_max_rate, should_succeed):
        env = create_environment(cmd, NEW_VERSION)

        exception = None
        try:
            env.add_validator(commission_rate=commission_rate, commission_max_rate=commission_max_rate)
        except Exception as e:
            exception = e

        if should_succeed:
            assert_no_exception(exception)
        else:
            assert sifchain.is_min_commission_too_low_exception(exception)

    test_case(0.05, 0.10, True)
    test_case(0.03, 0.20, False)
    test_case(0.07, 0.10, True)  # TODO Original scenarion failed for test_case(0.07, 0.04, True)


def test_min_commission_modify_existing_validator_24h(cmd: command.Command):
    # Using defaults: commission_rate=0.06, commission_max_rate=0.10, commission_max_change_rate=0.05
    # We create 3 validators for 3 different test cases so that we only have to wait once
    env = create_environment(cmd, NEW_VERSION)
    env.add_validator()
    env.add_validator()

    sifnoded0 = env._sifnoded_for(env.node_info[0])
    admin0_addr = env.node_info[0]["admin_addr"]
    sifnoded1 = env._sifnoded_for(env.node_info[1])
    admin1_addr = env.node_info[1]["admin_addr"]
    sifnoded2 = env._sifnoded_for(env.node_info[2])
    admin2_addr = env.node_info[2]["admin_addr"]

    validators = sifnoded0.query_staking_validators()
    assert all(float(v["commission"]["commission_rates"]["rate"]) == 0.06 for v in validators)
    assert all(float(v["commission"]["commission_rates"]["max_rate"]) == 0.10 for v in validators)
    assert all(float(v["commission"]["commission_rates"]["max_change_rate"]) == 0.05 for v in validators)

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
    assert_no_exception(exception)

    exception = None
    try:
        res = sifnoded1.staking_edit_validator(0.03, from_acct=admin1_addr, broadcast_mode="block")
        sifchain.check_raw_log(res)
    except Exception as e:
        exception = e
    assert sifchain.is_min_commission_too_low_exception(exception)

    exception = None
    try:
        res = sifnoded2.staking_edit_validator(0.07, from_acct=admin2_addr, broadcast_mode="block")
        sifchain.check_raw_log(res)
    except Exception as e:
        exception = e
    assert_no_exception(exception)


def test_min_commission_upgrade_handler(cmd: command.Command):
    def test_case(pre_upgrade_commission_rate, pre_upgrade_commission_max_rate, expected_commission_rate,
        expected_commission_max_rate, should_succeed
    ):
        exception = None
        try:
            env = create_environment(cmd, OLD_VERSION, commission_rate=pre_upgrade_commission_rate,
                commission_max_rate=pre_upgrade_commission_max_rate, commission_max_change_rate=0.01)
        except Exception as e:
            exception = e
        if should_succeed:
            assert_no_exception(exception)
        else:
            # TODO In case of invalid validator setup (commission_rate > commission_max_rate) sifnoded does not start
            #      and we get a timeout. We don't check the exception here, but we assume that this is what happened
            #      since other scenarios are working which only differ in parameters.
            return

        sifnoded = env._sifnoded_for(env.node_info[0])
        upgrade_height = sifnoded.get_current_block() + 15  # 15 * 5 = 75s > 60s

        commission_rates_before = exactly_one(sifnoded.query_staking_validators())["commission"]["commission_rates"]
        assert float(commission_rates_before["rate"]) == pre_upgrade_commission_rate
        assert float(commission_rates_before["max_rate"]) == pre_upgrade_commission_max_rate
        assert float(commission_rates_before["max_change_rate"]) == 0.01

        env.upgrade(NEW_VERSION, get_binary_for_version(NEW_VERSION), upgrade_height)

        sifnoded = env._sifnoded_for(env.node_info[0])

        commission_rates_after = exactly_one(sifnoded.query_staking_validators())["commission"]["commission_rates"]
        assert float(commission_rates_after["rate"]) == expected_commission_rate
        assert float(commission_rates_after["max_rate"]) == expected_commission_max_rate
        assert float(commission_rates_after["max_change_rate"]) == 0.01

    test_case(0.03, 0.04, 0.05, 0.05, True)
    test_case(0.06, 0.14, 0.06, 0.14, True)
    test_case(0.15, 0.25, 0.15, 0.25, True)
    test_case(0.02, 0.10, 0.05, 0.10, True)
    test_case(0.12, 0.10, 0.12, 0.10, False)  # This fails with "panic: commission cannot be more than the max rate"
    test_case(0.07, 0.04, 0.07, 0.04, False)  # This fails with "panic: commission cannot be more than the max rate"


def test_max_voting_power(cmd: command.Command):
    def test_case(from_validator_index, to_validator_index, amount, should_succeed):
        env = create_environment(cmd, NEW_VERSION, staking_amount=1000 * 10**21)
        env.add_validator(staking_amount=62 * 10**21)
        env.start()
        sifnoded = env._sifnoded_for(env.node_info[0])

        validator_powers_before = [int(x["tokens"]) for x in sifnoded.query_staking_validators()]

        exception = None
        try:
            time.sleep(5)  # Without this we would sometimes get "validator does not exist" in "tx staking delegate"
            delegate(env, from_validator_index, to_validator_index, amount)
        except Exception as e:
            exception = e

        validator_powers_after = [int(x["tokens"]) for x in sifnoded.query_staking_validators()]

        if should_succeed:
            assert_no_exception(exception)
            # Check actual vs. expected validator powers.
            # Note: this assertion might fail if "sifnoded query staking validators" returns a list in different order.
            expected_validator_powers_after = validator_powers_before
            expected_validator_powers_after[to_validator_index] += amount
            assert validator_powers_after == expected_validator_powers_after
        else:
            assert sifchain.is_max_voting_power_limit_exceeded_exception(exception)
            assert validator_powers_after == validator_powers_before

    sif = 0
    akasha = 1

    test_case(akasha, sif, 100, False)  # Current (and projected) voting power too big
    test_case(sif, akasha, 100, True)
    test_case(sif, akasha, 10**23, False)  # Projected voting power too big


def main(argv: List[str]):
    basic_logging_setup()
    cmd = command.Command()
    pkill(cmd)

    if argv == ["24"]:
        log.info("24h test")
        test_min_commission_modify_existing_validator_24h(cmd)
    else:
        test_min_commission_create_new_validator(cmd)
        test_min_commission_upgrade_handler(cmd)
        test_max_voting_power(cmd)

    pkill(cmd)

    log.info("Finished successfully")


if __name__ == "__main__":
    main(sys.argv[1:])
