import os
import subprocess
import re
import enum
import time
import json
import logging

json_data_file = os.environ.get('JSON_DATA_FILE')
sifchain_network = os.environ.get('SIFCHAIN_NETWORK')
logging_flag = os.environ.get('LOG')
get_gas_info_flag = os.environ.get('GAS_INFO')

formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s')
general_logger = ''
gas_logger = ''


def setup_logger(name, log_file, level=logging.DEBUG):
    handler = logging.FileHandler(log_file)
    handler.setFormatter(formatter)

    logger = logging.getLogger(name)
    logger.setLevel(level)
    logger.addHandler(handler)

    return logger


if logging_flag:
    general_logger = setup_logger('general logger', 'decimal_precision.log')
    gas_logger = setup_logger('gas logger', 'ibc_transfers_gas.log')
# if we don't output to a file and set default logging level to WARNING, basically nothing extra (INFO, DEBUG) will be logged
else:
    general_logger = setup_logger(
        'general logger', 'decimal_precision.log', logging.WARNING)
    gas_logger = setup_logger(
        'gas logger', 'ibc_transfers_gas.log', logging.WARNING)

with open(json_data_file) as json_file:
    data = json.load(json_file)


class Chain(enum.Enum):
    SIFCHAIN = 'sifchain'
    COSMOS = 'cosmos'
    AKASH = 'akash'


class TxType(enum.Enum):
    INCREASE = 'increase'
    DEDUCT = 'deduct'


def get_tokenregistry_entries():
    cmd = cmd_sif_q_tokenregistry_entries
    gas_logger.info(cmd)
    result = subprocess.run(cmd.split(' '),
                            capture_output=True, text=True)
    if result.returncode != 0:
        print(result.stderr)
        gas_logger.error(result.stderr)
        exit(1)

    gas_logger.info(result.stdout)
    return json.loads(result.stdout)


def query_balance(asset, chain=Chain.SIFCHAIN):
    if (chain == Chain.SIFCHAIN):
        cmd = cmd_sif_q_balance
    else:
        cmd = cmd_external_q_balance

    general_logger.info(cmd)
    result = subprocess.run(cmd.split(' '),
                            capture_output=True, text=True)
    if result.returncode != 0:
        print(result.stderr)
        general_logger.error(result.stderr)
        exit(1)
    general_logger.info(result.stdout)
    balances = result.stdout.split('\n')
    denom_index = [index for index, value in enumerate(
        balances) if value.find(f'denom: {asset}') != -1]
    if (len(denom_index) > 0):
        # index of amount is always less by 1, update : it's is NOT!! for iris (it's +1)
        if (chain != 'iris'):
            balance = re.sub(
                r'^.*amount: ', '', balances[denom_index[0] - 1])
        else:
            balance = re.sub(
                r'^.*amount: ', '', balances[denom_index[0] + 1])
    else:
        raise Exception(f'Denom balance for {asset} not found.')

    general_logger.info(f'{chain}:{asset} balance = {balance}')
    # remove surrounding '"' chars and possible '.'
    return balance.replace('"', '').replace('.', '')


# def transfer_tx(denom, tx_amount, source_chain=Chain.SIFCHAIN.value, dest_chain=Chain.AKASH.value):
def transfer_tx(source_chain=Chain.SIFCHAIN.value, dest_chain=Chain.AKASH.value):
    if source_chain == Chain.SIFCHAIN.value:
        cmd = cmd_tx_sif_to_external
    elif dest_chain == Chain.SIFCHAIN.value:
        cmd = cmd_tx_external_to_sif
    else:
        raise Exception(
            f'Transaction from {source_chain} to {dest_chain} is not supported.')

    # print(cmd)  # delete me
    general_logger.info(cmd)
    result = subprocess.run(cmd.split(' '),
                            capture_output=True, text=True)
    if result.returncode != 0:
        print(result.stderr)
        general_logger.error(result.stderr)
        exit(1)

    general_logger.info(result.stdout)
    json_data = json.loads(result.stdout)

    tx_info = dict()
    tx_info['txhash'] = json_data["txhash"]
    tx_info['gas_used'] = json_data["gas_used"]
    return tx_info


def query_tx_hash(cmd):
    gas_logger.info(cmd)
    result = subprocess.run(cmd.split(' '),
                            capture_output=True, text=True)
    if result.returncode != 0:
        print(result.stderr)
        gas_logger.error(result.stderr)
        exit(1)

    gas_logger.info(result.stdout)
    json_data = json.loads(result.stdout)

    tx_info = dict()
    tx_info['gas_amount'] = json_data["tx"]["auth_info"]["fee"]["amount"][0]
    return tx_info


# method used to truncate 18dp transferred amount, i.e.
# 123456789012345678 -> 123456789000000000
# 123456789 -> 100000000
# 500 -> 000
# it digits_to_truncate is negative, then we are adding instead of truncating
# 123456 -> 123456000000000000


def convert_amount(value, digits_to_truncate=8):
    # create a string from value and then manipulate it
    value = str(value)
    if digits_to_truncate > 0:
        if len(value) <= digits_to_truncate:
            result = '0'.ljust(digits_to_truncate, '0')
        else:
            result = value[0:len(value)-digits_to_truncate]
    else:
        result = value.ljust(len(value)-digits_to_truncate,
                             '0')  # it's in fact adding here

    general_logger.info(f'Converted amount = {result}')
    return result


def calculate_expected_value(input_value, tx_amount, type=TxType.INCREASE):
    expected_value = 0
    if isinstance(input_value, str):
        input_value = int(input_value)
    if isinstance(tx_amount, str):
        tx_amount = int(tx_amount)

    if type == TxType.INCREASE:
        expected_value = input_value + tx_amount
    else:
        expected_value = input_value - tx_amount
    return str(expected_value)


def assert_expected_value(expected, asset, chain=Chain.SIFCHAIN):
    timeout_start = time.time()
    actual = 0

    while time.time() < timeout_start + assertion_timeout:
        actual = query_balance(asset, chain)
        print(f'\t\t{actual} {chain}')
        if expected == actual:
            print(f'\t\tAssertion passed for {chain}')
            break
        time.sleep(5)
        print(f'\t\tAssertion retry: expected {expected}, actual {actual}')

    assert expected == actual, f'\t\tAssertion failed for {chain}, {asset}: expected {expected}, got {actual}.'


# 1000 and dp=6 => 0.001000
# 1000 and dp=18 => 0.001000000000000000

def get_dot_notation(amount, dp_value):
    if isinstance(amount, int):
        amount = str(amount)
    dot_notation = '0'
    if len(amount) < dp_value:
        missing_dp = dp_value - len(amount)
        dot_notation = '0.'.ljust(missing_dp+2, '0') + amount
    elif len(amount) == dp_value:
        dot_notation = '0.' + amount
    else:
        dot_notation = amount[0] + '.' + amount[1:]

    print(f'dot notation = {dot_notation}')
    return dot_notation


#     precision_diff = abs(int(sif_asset_dp_value) - int(ibc_denom_dp_value))
#     normalized = 0
#     # new code
#     if len(amount) < precision_diff:
#         normalized = int(amount.ljust(len(amount) + precision_diff, '0'))
#     else:
#         normalized = int(amount)

#     general_logger.info(f'Normalized amount = {normalized}')
#     return normalized


def save_tx_info(denom, tx_amount, source_chain, tx_info):
    if not tx_amount in tx_info_table[source_chain][denom]:
        tx_info_table[source_chain][denom][tx_amount] = {}
    tx_info_table[source_chain][denom][tx_amount].update(tx_info)
    gas_logger.info(tx_info_table)


def save_tx_info_new_denom(denom):
    tx_info_table['sifchain'][denom] = {}
    tx_info_table[ibc_network_name][denom] = {}


def get_ibc_counterparty_chains():
    ibc_chains = []
    entries = get_tokenregistry_entries()["entries"]
    for chain in entries:
        if chain["ibc_counterparty_chain_id"]:
            ibc_chains.append(chain)
    return ibc_chains


sif_wallet = data["wallet"]["sif"]

sif_chain_config = [x for x in data["chain"]
                    if x['name'] == sifchain_network][0]
sif_node = sif_chain_config["node"]
sif_chain_id = sif_chain_config["chain_id"]

cmd_sif_q_balance = f'sifnoded query bank balances {sif_wallet} --node {sif_node} --chain-id {sif_chain_id}'
cmd_sif_q_tokenregistry_entries = f'sifnoded query tokenregistry entries --node {sif_node} --output json'


ibc_chains = get_ibc_counterparty_chains()
assertion_timeout = 180  # seconds

general_logger.info('++++++++++ New run started ++++++++++')

for ibc_chain in ibc_chains:
    ibc_network_name = ibc_chain["display_name"]
    if ibc_network_name != "akash":  # tmp
        continue

    json_external_chain_config = [x for x in data["chain"]  # todo: once tokenregistry entries store network name, it can be deleted
                                  if x['name'] == ibc_network_name][0]

    tx_info_table = {"sifchain": {}, ibc_network_name: {}}
    external_cli_tool = json_external_chain_config["cli_tool"]
    # todo: once tokenregistry entries store network node, it can be retrieved from external_chain_config["address"]
    external_node = json_external_chain_config["node"]
    external_gas_price = f'{json_external_chain_config["gas_price"]}{ibc_chain["external_symbol"]}'
    external_wallet = data["wallet"][ibc_network_name]
    channel = ibc_chain['ibc_channel_id']
    counterparty_channel = ibc_chain['ibc_counterparty_channel_id']
    external_chain_id = ibc_chain["ibc_counterparty_chain_id"]

    cmd_external_q_balance = f'{external_cli_tool} query bank balances {external_wallet} --node {external_node} --chain-id {external_chain_id}'

    for tx_data in data["tx"]:
        sif_asset = ibc_chain["denom"]
        ibc_base_denom = ibc_chain["base_denom"]
        sif_asset_dp_value = int(ibc_chain["decimals"])
        ibc_denom_dp_value = 18  # todo: verify if this can be hardcoded value
        save_tx_info_new_denom(sif_asset)
        for tx_amount in tx_data["amount"]:
            # it might be a negative value, then we should add instead of truncate
            digits_to_truncate = sif_asset_dp_value - ibc_denom_dp_value

            tx_amount_converted = int(
                convert_amount(tx_amount, digits_to_truncate))
            tx_amount_dot_notation = get_dot_notation(
                tx_amount, abs(sif_asset_dp_value))
            cmd_tx_sif_to_external = f'sifnoded tx ibc-transfer transfer transfer {channel} {external_wallet} {tx_amount}{sif_asset} --from={sif_wallet} --keyring-backend=test --node={sif_node} --chain-id={sif_chain_id} -y --packet-timeout-timestamp=6000000000000 --fees=100000000000000000rowan --broadcast-mode=block --output=json'
            # think there is a bug on iris cli. When it transfers from iris then it expects value with 6dp (instead of 18. And 18 was used for tx sif->iris, 18 is returned by iris q balance)
            cmd_tx_external_to_sif = f'{external_cli_tool} tx ibc-transfer transfer transfer {counterparty_channel} {sif_wallet} {tx_amount_dot_notation}{ibc_base_denom} --from={external_wallet} --keyring-backend=test --chain-id={external_chain_id} --node={external_node} -y --gas-prices={external_gas_price} --gas=500000 --packet-timeout-timestamp=600000000000 --broadcast-mode=block --output=json'
            # cmd_tx_external_to_sif = f'{external_cli_tool} tx ibc-transfer transfer transfer {counterparty_channel} {sif_wallet} {tx_amount_converted}{ibc_denom} --from={external_wallet} --keyring-backend=test --chain-id={external_chain_id} --node={external_node} -y --gas-prices={external_gas_price} --gas=500000 --packet-timeout-timestamp=600000000000 --broadcast-mode=block --output=json'

            print(
                f'++++ {sif_asset} ==== sif->{ibc_network_name} (tx {tx_amount}) and {ibc_network_name}->sif (tx {tx_amount_dot_notation}) ====')
            # f'++++ {tx_data["sif_asset"]} ==== sif->{external_network} (tx {tx_amount}) and {external_network}->sif (tx {tx_amount_converted}) ====')
            print(f'\tTransferring sif->{ibc_network_name}')

            sif_asset_balance = query_balance(sif_asset, Chain.SIFCHAIN)
            print(f'\t{sif_asset_balance}')
            external_asset_balance = query_balance(
                ibc_base_denom, ibc_network_name)
            print(f'\t{external_asset_balance}')

            tx_info = transfer_tx(dest_chain=ibc_network_name)
            save_tx_info(sif_asset, tx_amount, Chain.SIFCHAIN.value, tx_info)
            assert_expected_value(calculate_expected_value(
                sif_asset_balance, tx_amount, TxType.DEDUCT), sif_asset, Chain.SIFCHAIN)
            assert_expected_value(calculate_expected_value(
                external_asset_balance, tx_amount_converted, TxType.INCREASE), ibc_base_denom, ibc_network_name)

            if get_gas_info_flag:
                cmd_tx_hash_info = f'sifnoded q tx {tx_info["txhash"]} --node={sif_node} --output=json'
                tx_info = query_tx_hash(cmd_tx_hash_info)
                save_tx_info(sif_asset, tx_amount,
                             Chain.SIFCHAIN.value, tx_info)

            time.sleep(2)
            print(f'\tTransferring {ibc_network_name}->sif')
            if tx_amount_converted > 0:
                sif_asset_balance = query_balance(sif_asset, Chain.SIFCHAIN)
                print(f'\t{sif_asset_balance}')
                external_asset_balance = query_balance(
                    ibc_base_denom, ibc_network_name)
                print(f'\t{external_asset_balance}')

                tx_info = transfer_tx(
                    source_chain=ibc_network_name, dest_chain=Chain.SIFCHAIN.value)
                save_tx_info(sif_asset, tx_amount_converted,
                             ibc_network_name, tx_info)
                assert_expected_value(calculate_expected_value(
                    external_asset_balance, tx_amount_converted, TxType.DEDUCT), ibc_base_denom, ibc_network_name)
                assert_expected_value(calculate_expected_value(
                    sif_asset_balance, tx_amount, TxType.INCREASE), sif_asset, Chain.SIFCHAIN)
                if get_gas_info_flag:
                    cmd_tx_hash_info = f'sifnoded q tx {tx_info["txhash"]} --node={external_node} --output=json'
                    tx_info = query_tx_hash(cmd_tx_hash_info)
                    save_tx_info(sif_asset, tx_amount_converted,
                                 ibc_network_name, tx_info)
            else:
                print("\t\tSkipping: tx amount = 0")

    if get_gas_info_flag:
        gas_logger.info(json.dumps(tx_info_table, indent=4))
