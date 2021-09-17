import os
import subprocess
import re
import enum
import time

sif_wallet = os.environ.get('SIF_WALLET')
cosmos_wallet = os.environ.get('COSMOS_WALLET')
sif_asset = os.environ.get('SIF_ASSET')
ibc_denom = os.environ.get('IBC_DENOM')
tx_amount = os.environ.get('TX_AMOUNT')
sif_asset_dp_value = os.environ.get('SIF_ASSET_DP_VALUE')
ibc_denom_dp_value = os.environ.get('IBC_DENOM_DP_VALUE')


def remove_last_8_digits(input_value):
    # create a string from value and then manipulate it
    input_value = str(input_value)
    if len(input_value) - 8 > 0:
        return input_value[0:len(input_value)-8]
    else:
        return input_value


tx_amount_10dp = remove_last_8_digits(tx_amount)
cmd_tx_sif_to_cosmos = f'sifnoded tx ibc-transfer transfer transfer channel-114 {cosmos_wallet} {tx_amount}{sif_asset} --from={sif_wallet} --keyring-backend=test --node=tcp://rpc-devnet.sifchain.finance:80 --chain-id=sifchain-devnet-1 -y --packet-timeout-timestamp=6000000000000 --gas=5000000 --gas-prices=0.5rowan'
cmd_tx_cosmos_to_sif = f'gaiad tx ibc-transfer transfer transfer channel-26 {sif_wallet} {tx_amount_10dp}{ibc_denom} --from={cosmos_wallet} --keyring-backend=test --chain-id=cosmoshub-testnet --node=https://rpc.testnet.cosmos.network:443 -y --gas-prices=50.0uphoton --gas=500000 --packet-timeout-timestamp 600000000000'
cmd_sif_q_balance = f'sifnoded query bank balances {sif_wallet} --node tcp://rpc-devnet.sifchain.finance:80 --chain-id sifchain-devnet-1'
cmd_cosmos_q_balance = f'gaiad query bank balances {cosmos_wallet} --node https://rpc.testnet.cosmos.network:443 --chain-id cosmoshub-testnet'

assertion_timeout = 180  # seconds


class Chain(enum.Enum):
    SIFCHAIN = 'sifchain'
    COSMOS = 'cosmos'


class TxType(enum.Enum):
    INCREASE = 'increase'
    DEDUCT = 'deduct'


def query_balance(asset, chain=Chain.SIFCHAIN):
    if (chain == Chain.SIFCHAIN):
        cmd = cmd_sif_q_balance
    else:
        cmd = cmd_cosmos_q_balance

    result = subprocess.run(cmd.split(' '),
                            capture_output=True, text=True)
    if result.returncode != 0:
        print(result.stderr)
        exit(1)
    balances = result.stdout.split('\n')
    denom_index = [index for index, value in enumerate(
        balances) if value.find(f'denom: {asset}') != -1]
    if (len(denom_index) > 0):
        # index of amount is always less by 1
        sif_initial_balance = re.sub(
            r'^.*amount: ', '', balances[denom_index[0] - 1])
    return sif_initial_balance.replace('"', '')  # remove surrounding '"' chars


def transfer_tx(sourceChain=Chain.SIFCHAIN, destChain=Chain.COSMOS):
    if sourceChain == Chain.SIFCHAIN and destChain == Chain.COSMOS:
        cmd = cmd_tx_sif_to_cosmos
    elif sourceChain == Chain.COSMOS and destChain == Chain.SIFCHAIN:
        cmd = cmd_tx_cosmos_to_sif
    else:
        raise Exception(
            f'Transaction from {sourceChain} to {destChain} is not supported.')

    result = subprocess.run(cmd.split(' '),
                            capture_output=True, text=True)
    if result.returncode != 0:
        print(result.stderr)
        exit(1)

# method used to truncate 18dp transferred amount, i.e.
# 123456789012345678 -> 123456789000000000
# 123456789 -> 100000000
# 500 -> 000


def truncate_18dp_amount(value):
    # create a string from value and then manipulate it
    value = str(value)
    if len(value) <= 8:
        result = '00000000'
    else:
        result = value[0:len(value)-8].ljust(len(value), '0')

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
        print(actual, chain)
        if expected == actual:
            print(f'\t\tAssertion passed for {chain}')
            break
        time.sleep(5)
        print(f'\t\tAssertion retry: expected {expected}, actual {actual}')

    assert expected == actual, f'\t\tAssertion failed for {chain}, {asset}: expected {expected}, got {actual}.'


def normalize_ibc_amount_to_sif_dp(amount):
    if isinstance(amount, int):
        amount = str(amount)

    precision_diff = int(sif_asset_dp_value) - int(ibc_denom_dp_value)
    if precision_diff > 0:  # i.e. 18 for ceth/rowan, 6 for cusdc, 8 for ccro
      return int(amount.ljust(len(amount) + precision_diff, '0'))
    else:
      return int(amount)


print(
    f'==== sif->cosmos (tx {tx_amount}) and cosmos->sif (tx {tx_amount_10dp}) ====')
print(f'\tTransferring sif->cosmos')

sif_asset_balance = query_balance(sif_asset, Chain.SIFCHAIN)
print(f'\t{sif_asset_balance}')
cosmos_asset_balance = query_balance(ibc_denom, Chain.COSMOS)
print(f'\t{cosmos_asset_balance}')

transfer_tx()
truncated_18dp_amount = int(truncate_18dp_amount(tx_amount))
assert_expected_value(calculate_expected_value(
    sif_asset_balance, truncated_18dp_amount, TxType.DEDUCT), sif_asset, Chain.SIFCHAIN)
assert_expected_value(calculate_expected_value(cosmos_asset_balance, remove_last_8_digits(
    truncated_18dp_amount), TxType.INCREASE), ibc_denom, Chain.COSMOS)

time.sleep(5)
print(f'\tTransferring cosmos->sif')

sif_asset_balance = query_balance(sif_asset, Chain.SIFCHAIN)
print(f'\t{sif_asset_balance}')
cosmos_asset_balance = query_balance(ibc_denom, Chain.COSMOS)
print(f'\t{cosmos_asset_balance}')

transfer_tx(sourceChain=Chain.COSMOS, destChain=Chain.SIFCHAIN)
assert_expected_value(calculate_expected_value(
    cosmos_asset_balance, tx_amount_10dp, TxType.DEDUCT), ibc_denom, Chain.COSMOS)
assert_expected_value(calculate_expected_value(
    sif_asset_balance, normalize_ibc_amount_to_sif_dp(tx_amount_10dp), TxType.INCREASE), sif_asset, Chain.SIFCHAIN)
