import subprocess
import json 
import os

# network descriptor is 1 for ethereum
network_descriptor = 1

bridge_bank_address = '0xB5F54ac4466f5ce7E0d8A5cB9FE7b8c0F35B7Ba8'

# get all entries from product
result = subprocess.run(['sifnoded', 'query', 'tokenregistry', 'entries', '--node',  'tcp://rpc-archive.sifchain.finance:80'], stdout=subprocess.PIPE).stdout.decode('utf-8')
denoms = json.loads(result)['entries']

# algorithm to get denom in peggy2.0
def getPeggy2Denom(network_descriptor, token_contract_address):
    assert token_contract_address.startswith("0x")
    assert network_descriptor > 0
    assert network_descriptor <= 9999
    denom = f"sifBridge{network_descriptor:04d}{token_contract_address.lower()}"
    return denom

# define the format of output file
def composeMapping(old_denom, new_denom):
    return {
        'peggy1': old_denom,
        'peggy2': new_denom
    }

# get all token address from mainnet by filter logs in bridge bank 
def get_token_address():
    result = subprocess.run(['yarn', 'integrationtest:whitelistedTokens', 
    '--json_path', '../deployments/sifchain-1',
    '--ethereum_network', 'mainnet',
    '--bridgebank_address', bridge_bank_address,
    '--network', 'mainnet'])

folder = './'

directory = os.fsencode(folder)
denom_address_map = {}   
for file in os.listdir(directory):
     filename = os.fsdecode(file)
     # all deployed contract in the files with prefix whitelist
     if filename.startswith('whitelist'):
        print(filename)
        file = open(folder + '/' + filename)
        data = json.load(file)

        entries = data['array']
        for entry in entries:
            if 'symbol' in entry and 'address' in entry:
                denom_address_map[entry['symbol']] = entry['address']


# map denom in peggy 1.0 to peggy 2.0
# if denom start with ibc, then denom is the same with peggy 2.0
# if denom start with c, then remove c, call getPeggy2Denom
# if denom start with x, then remove x, call getPeggy2Denom

missed_denom = []
result = []
output_file = open('denom_mapping.json', 'w', encoding='utf-8')

for item in denoms:
    denom = (item['denom'])

    if denom == 'rowan':
        result.append(composeMapping(denom, denom))
    elif denom.startswith('ibc/'):
        result.append(composeMapping(denom, denom))
    elif denom.startswith('c'):
        tmp_denom = denom[1:].upper()
        if tmp_denom in denom_address_map:
            result.append(composeMapping(denom, getPeggy2Denom(network_descriptor, denom_address_map[tmp_denom])))
        else:
            missed_denom.append(denom)
    else:
        missed_denom.append(denom)

json.dump({'array': result}, output_file)

print("-------- items not found --------")
# sort for better find out denom
missed_denom.sort()
for item in missed_denom:
    print(item, end=',')

