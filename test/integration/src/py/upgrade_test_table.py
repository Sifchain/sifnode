from decimal import *
import json
import codecs
getcontext().prec = 60
poolsAfterData = json.load(codecs.open('pools-after.json', 'r', 'utf-8-sig'))
poolsBeforeData = json.load(codecs.open('pools-before.json', 'r', 'utf-8-sig'))
lpAfterData = json.load(codecs.open('lp-after.json', 'r', 'utf-8-sig'))
lpBeforeData = json.load(codecs.open('lp-before.json', 'r', 'utf-8-sig'))
    
lpUnitsExpected = []
poolUnitsExpected = []
poolcount = 0
lptabledata = {  '_info': {
    'desc': "Inputs and expected results for calculating pool units.",
    'symbol': "external asset symbol",
    'r': "native asset added",
    'a': "external asset added",
    'R': "native Balance (before)",
    'A': "external Balance (before)",
    'P': "existing Pool Units"
  }, "PoolUnits": []}
for pool in poolsBeforeData['Pools']:
  symbol = pool['external_asset']['symbol']
  poolUnits = 0
  externalPool = 0
  nativePool = 0
  lpcount = 0

  for lp in lpBeforeData['LiquidityProviders']:
    if lp['asset']['symbol'] == symbol:
      share = Decimal(lp['liquidity_provider_units']) / Decimal(pool['pool_units'])
      externalBalance = share * Decimal(pool['external_asset_balance'])
      nativeBalance =  share * Decimal(pool['native_asset_balance'])
      P = poolUnits
      A = externalPool
      R = nativePool
      a = externalBalance
      r = nativeBalance
      if P == 0:
        answer = r
      else:
        answer = ((P*(R*a+r*A))/(2*R*A))*(1-abs((R*a-r*A)/((r+R)*(a+A))))
        
      poolUnits += round(answer, 0)
      externalPool += round(a, 0)
      nativePool +=  round(r, 0)
      output = "{:.0f}".format(round(answer, 0))
      lptabledata['PoolUnits'].append({'symbol': symbol, 'r': "{:.0f}".format(round(r, 0)), 'a': "{:.0f}".format(round(a, 0)), 'R': "{:.0f}".format(round(R, 0)), 'A': "{:.0f}".format(round(A, 0)), 'P': "{:.0f}".format(round(P, 0)), 'expected': output})

with open('real_pool_units.json', 'w') as json_file:
  json.dump(lptabledata, json_file)  