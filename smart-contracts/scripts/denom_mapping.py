import subprocess
import json 
import os

# algorithm to get denom in peggy2.0
def getPeggy2Denom(network_descriptor: int, token_contract_address: str):
    assert token_contract_address.startswith("0x")
    assert network_descriptor >= 0
    assert network_descriptor <= 9999
    denom = f"sifBridge{network_descriptor:04d}{token_contract_address.lower()}"
    return denom

# network descriptor is 1 for ethereum
network_descriptor = 1
eth_contract_address = '0x0000000000000000000000000000000000000000'

def main():
    # get all entries from product
    result = subprocess.run(['sifnoded', 'query', 'tokenregistry', 'entries', '--node',  'tcp://rpc-archive.sifchain.finance:80'], stdout=subprocess.PIPE).stdout.decode('utf-8')
    denoms = json.loads(result)['entries']

    denom_address_map = {}
    data = json.load(open('../data/denom_contracts.json'))
    for entry in data:
        denom_address_map[entry['symbol']] = entry['token']

    # map denom in peggy 1.0 to peggy 2.0
    # if denom start with ibc, then denom is the same with peggy 2.0
    # if denom start with c, then remove c, call getPeggy2Denom
    # if denom start with x, will be the same as peggy 1.0

    missed_denom = []
    result = {}
    reverse_result = {}
    output_file_1 = open('../data/denom_mapping_peggy1_to_peggy2.json', 'w', encoding='utf-8')
    output_file_2 = open('../data/denom_mapping_peggy2_to_peggy1.json', 'w', encoding='utf-8')

    for item in denoms:
        denom = (item['denom'])
        if denom == 'rowan':
            result[denom] = denom
            reverse_result[denom] = denom
        elif denom == 'ceth':
            eth_denom = getPeggy2Denom(network_descriptor, eth_contract_address)
            result[denom] = eth_denom
            reverse_result[eth_denom] = denom
        elif denom.startswith('ibc/'):
            result[denom] = denom
            reverse_result[denom] = denom
        elif denom.startswith('x'):
            result[denom] = denom
            reverse_result[denom] = denom
        elif denom.startswith('c'):
            tmp_denom = denom[1:].upper()
            if tmp_denom in denom_address_map:
                peggy2_denom = getPeggy2Denom(network_descriptor, denom_address_map[tmp_denom])
                result[denom] = peggy2_denom
                reverse_result[peggy2_denom] = denom
            else:
                missed_denom.append(denom)
        else:
            missed_denom.append(denom)

    json.dump(result, output_file_1)
    json.dump(reverse_result, output_file_2)

    print("-------- items not found --------")
    # sort for better find out denom
    missed_denom.sort()
    for item in missed_denom:
        print(item, end=',')


if __name__ == "__main__":
    main()


