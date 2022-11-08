# Draft - not used yet
#
# Design document: https://docs.google.com/document/d/1yxxQ3RtftCvCJp_vDlSR5MiXvtOEJxlGXg3GLvf9dls/edit?skip_itp2_check=true
# Test environment for testing the new Sifchain public SDK: https://docs.google.com/document/d/1MAlg-I0xMnUvbavAZdAN---WuqbyuRyKw-6Lfgfe130/edit
# IBC localnet tool: https://github.com/Sifchain/sifchain-deploy/commit/0424d2ec81f6f739a3a2cbc77feba1a6f297430d
# scripts/init-multichain.sh

import json
import sys

from siftool.command import Command

chains = {
    "akash": {"binary": "akash", "relayer": "ibc"},
    "iris": {"binary": "iris", "relayer": "hermes"},
    "sentinel": {"binary": "sentinelhub", "relayer": "ibc"},
    "persistence": {"binary": "persistenceCore", "relayer": "hermes"},
    "sifchain": {"binary": "sifnoded", "relayer": "ibc"},
}

# TODO Define format and usage
configs = {
    "local": {
        "akash": "http://127.0.0.1:26656",
        "sifchain": "http://127.0.0.1:26657",
        "chainId": "akash-testnet-6",
        "fees": 5000,
        "gas": 1000,
        "denom": "uakt",
        "start_chain_locally": True,
    },
    "ci": {
        "start_chain_locally": False,
    },
}

# Runs the command synchronously, checks that the exit code is 0 and returns standard output and error.
def run_command(args, stdin=None, cwd=None, env=None, check_exit=False):
    return Command().execst(args, cwd=cwd, env=env, check_exit=check_exit, stdin=stdin)

# Starts a process and returns a Popen object for it
def start_process(args):
    return None  # TODO

def get_binary_for_chain(chain_name):
    return chains[chain_name]["binary"]

def get_config(config_name):
    return None  # TODO

# Generates a sifnoded key and stores it into test keyring. Returns the mnemonic that can be used to
# recreate it.
def add_new_key_to_keyring(chain, key_name):
    binary = get_binary_for_chain(chain)
    res = run_command([binary, "keys", "add", key_name, "--keyring-backend", "test", "--output", "json"], stdin=["y"])
    return json.loads(res.stdout)["mnemonic"]

def add_existing_key_to_keyring(chain, key_name, mnemonic, overwrite=True):
    binary = get_binary_for_chain(chain)
    if overwrite:
        run_command([binary, "keys", "delete", key_name, "--keyring-backend", "test", "-y"], check_exit=False)
    run_command([binary, "keys", "add", key_name, "-i", "--recover", "--keyring-backend", "test"],
        stdin=[mnemonic, ""])

def start_chain(chain):
    pass

# Can be initialized either manually or from genesis file
def init_chain(chain):
    pass

# Both hermes and ts-relayer are IBC-compliant relayers.
# (1) hermes is written in Rust (Cephalopod, Informal systems), (2) ts-relayer is written in Typescript, (3) Go
# implementation. We use ts-relayer for most of the cases, because hermes was working unreliable.
# Their core functionality should be similar, but for some reason we might choose one over the other in specific cases.
# ts-relayer has the ability to act as multiple relayers in one process.

def start_hermes_relayer():
    pass

def start_ts_relayer():
    pass

def start_relayer(chain_a, chain_b, config, channel_id, counterchannel_id):
    # TODO Determine which relayer to use
    relayer_binary = None
    relayer_args = []
    relayer_process = start_process([relayer_binary] + relayer_args)
    return relayer_process

def send_transaction(chain, channel, amount, denom, src_addr, dst_addr, sequence, chain_id, node, broadcast_mode,
    fees, gas, account_number, dry_run=False
):
    if not broadcast_mode in ["async", "block"]:
        raise ValueError("Invalid broadcast_mode '{}'".format(broadcast_mode))
    args = [get_binary_for_chain(chain), "tx", "ibc-transfer", "transfer", "transfer", f"channel-{channel}",
        dst_addr, f"{amount}{denom}", "--from", src_addr, "--keyring-backend", "test", "--chain-id", chain_id,
        "--node", node, "--sequence", str(sequence), "--account-number", account_number] + \
        (["--fees", f"{fees}{denom}"] if fees else []) + \
        (["--gas", gas] if gas else []) + \
        (["--broadcast-mode", broadcast_mode if broadcast_mode else "async"]) + \
        (["--offline"] if broadcast_mode == "async" else []) + \
        (["--dry-run"] if dry_run else [])
    run_command(args)

def query_bank_balance(chain, addr, denom, config):
    node = None  # TODO
    chain_id = None  # TODO
    result = json.loads(run_command([get_binary_for_chain(chain), "q", "bank", "balances", addr, "--node", node,
        "--chain-id", chain_id, "--output", "json"]).stdout)
    return int(result[denom])

def get_initial_account_and_sequence_number(config, chain, src_addr, node, chain_id):
    res = run_command([get_binary_for_chain(chain), "q", "auth", "account", src_addr, "--node", node, "--chain-id",
        chain_id, "--output", "json"]).stdout
    account_number, sequence = res["account_number"], res["sequence"]
    return account_number, sequence

def run_tests_for_one_chain_in_one_direction(config, other_chain, direction_flag, number_of_iterations):
    # Assume we're always sending transactions between other_chain and sifchain.
    # from_chain is chain sending the assets, to_chain is the receiving chain.
    from_chain = "sifchain" if direction_flag else other_chain
    to_chain = other_chain if direction_flag else "sifchain"
    broadcast_mode = "block"  # TODO
    chain_id = int(config["chain_id"])
    denom = config["denom"]
    channel_id = int(config["channel_id" if direction_flag else "counterchannel_id"])
    from_account = config["from_account"]
    to_account = config["to_account"]
    amount = int(config["amount"])
    node = config[chain_id]["node"]
    fees = config["fees"]
    gas = config["gas"]
    sifchain_proc = start_chain("sifchain")
    other_chain_proc = start_chain(other_chain)
    if config["init_chain"]:
        init_chain(from_chain)
        init_chain(to_chain)
        relayer_proc = start_relayer(from_chain, to_chain, channel_id, counterchannel_id)
    mnemonic = add_new_key_to_keyring("sifchain", from_account)
    add_existing_key_to_keyring(other_chain, to_account, mnemonic)
    sequence, account_number = get_initial_account_and_sequence_number(config, from_chain, from_account, node, chain_id)

    from_balance_before = query_bank_balance(from_chain, from_account, denom)
    to_balance_before = query_bank_balance(to_chain, to_account, denom)

    # TODO Check that from_balance_before >= number_of_iterations * amount + fees + gas
    # We can know the exact gas in block mode

    for i in range(number_of_iterations):
        send_transaction(from_chain, channel_id, amount, denom, from_account,
            to_account, sequence + i, chain_id, node, broadcast_mode, fees, gas, account_number)

    if broadcast_mode == "async":
        # TODO Wait for transactions to complete before querying the balances
        pass

    from_balance_after = query_bank_balance(from_chain, from_account, denom)
    to_balance_after = query_bank_balance(to_chain, to_account, denom)
    relayer_proc.stop()
    sifchain_proc.stop()
    other_chain_proc.stop()
    assert from_balance_after == from_balance_before - number_of_iterations * amount # TODO Account for fees and gas
    assert to_balance_after == to_balance_before + number_of_iterations * amount

def run_tests_for_all_chains_in_both_directions(config, number_of_iterations):
    for chain in chains:
        run_tests_for_one_chain_in_one_direction(config, chain, True, number_of_iterations)
        run_tests_for_one_chain_in_one_direction(config, chain, False, number_of_iterations)

# This is called from GitHub CI/CD (i.e. .github/workflows)
def run_from_ci(args):
    config = get_config("ci")
    run_tests_for_all_chains_in_both_directions(config, 1000)

def run_locally(args):
    config = get_config("local")
    other_chain = args[0]
    direction_flag = args[1] == "receiver"
    number_of_iterations = int(args[2])
    run_tests_for_one_chain_in_one_direction(config, other_chain, direction_flag, number_of_iterations)

def main(argv):
    action = argv[0]
    action_args = argv[1:]
    if action == "ci":
        run_from_ci(action_args)
    elif action == "local":
        run_locally(action_args)

if __name__ == "__main__":
    main(sys.argv)
