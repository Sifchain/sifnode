import json
import codecs
poolsAfterData = json.load(codecs.open('pools-after.json', 'r', 'utf-8-sig'))
poolsBeforeData = json.load(codecs.open('pools-before.json', 'r', 'utf-8-sig'))
lpAfterData = json.load(codecs.open('lp-after.json', 'r', 'utf-8-sig'))
lpBeforeData = json.load(codecs.open('lp-before.json', 'r', 'utf-8-sig'))

lpsAfter = []
for pool in poolsAfterData['Pools']:
  symbol = pool['external_asset']['symbol']
  for lp in lpAfterData['LiquidityProviders']:
    if lp['asset']['symbol'] == symbol:
      share = float(lp['liquidity_provider_units']) / float(pool['pool_units'])
      lpsAfter.append({'lp': lp['liquidity_provider_address'], 'share': share})
    
lpsBefore = []
for pool in poolsBeforeData['Pools']:
  symbol = pool['external_asset']['symbol']
  for lp in lpBeforeData['LiquidityProviders']:
    if lp['asset']['symbol'] == symbol:
      share = float(lp['liquidity_provider_units']) / float(pool['pool_units'])
      lpsBefore.append({'lp': lp['liquidity_provider_address'], 'share': share})


diff = []
# print(lpsAfter)
# print(lpsBefore)
count = 0
for lp in lpsBefore:
  if lpsBefore[count]['share'] != lpsAfter[count]['share']:
    diff.append({'lpaddress': lpsBefore[count]['lp'], 'shareBefore': lpsBefore[count]['share'], 'shareAfter': lpsAfter[count]['share'], 'diff': lpsAfter[count]['share'] - lpsBefore[count]['share']})
    print(lpsBefore[count])
    print(lpsAfter[count])
  count += 1

with open('afterShares.json', 'w') as json_file:
    json.dump(lpsAfter, json_file)

with open('beforeShares.json', 'w') as json_file:
  json.dump(lpsBefore, json_file)  

with open('diff.json', 'w') as json_file:
  json.dump(diff, json_file)  