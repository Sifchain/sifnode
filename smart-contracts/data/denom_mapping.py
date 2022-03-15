import subprocess
import json 
import os

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

# network descriptor is 1 for ethereum
network_descriptor = 1

def main():
    # get all entries from product
    result = subprocess.run(['sifnoded', 'query', 'tokenregistry', 'entries', '--node',  'tcp://rpc-archive.sifchain.finance:80'], stdout=subprocess.PIPE).stdout.decode('utf-8')
    denoms = json.loads(result)['entries']

    denom_address_map = {}
    data = json.load(open('denom_contracts.json'))
    for entry in data:
        denom_address_map[entry['symbol']] = entry['token']

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


if __name__ == "__main__":
    main()


