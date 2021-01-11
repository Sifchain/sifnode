```
curl 'http://127.0.0.1:1317/clp/getPools' \
  -H 'Connection: keep-alive' \
  -H 'Pragma: no-cache' \
  -H 'Cache-Control: no-cache' \
  -H 'sec-ch-ua: "Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"' \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'DNT: 1' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36' \
  -H 'Origin: http://localhost:8080' \
  -H 'Sec-Fetch-Site: cross-site' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Referer: http://localhost:8080/' \
  -H 'Accept-Language: en-US,en;q=0.9,la;q=0.8' \
  --compressed
```

RES

```json
{
  "height": "84",
  "result": {
    "Pools": [
      {
        "external_asset": {
          "symbol": "catk"
        },
        "native_asset_balance": "10000000000000000000000000",
        "external_asset_balance": "10000000000000000000000000",
        "pool_units": "10000000000000000000000000"
      },
      {
        "external_asset": {
          "symbol": "cbtk"
        },
        "native_asset_balance": "10000000000000000000000000",
        "external_asset_balance": "10000000000000000000000000",
        "pool_units": "10000000000000000000000000"
      },
      {
        "external_asset": {
          "symbol": "ceth"
        },
        "native_asset_balance": "10000000000000000000000000",
        "external_asset_balance": "8300000000000000000000",
        "pool_units": "10000000000000000000000000"
      },
      {
        "external_asset": {
          "symbol": "clink"
        },
        "native_asset_balance": "10000000000000000000000000",
        "external_asset_balance": "588235000000000000000000",
        "pool_units": "10000000000000000000000000"
      },
      {
        "external_asset": {
          "symbol": "cusdc"
        },
        "native_asset_balance": "10000000000000000000000000",
        "external_asset_balance": "10000000000000000000000000",
        "pool_units": "10000000000000000000000000"
      }
    ],
    "clp_module_address": "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85",
    "height": "84"
  }
}
```

```
sifnodecli query clp pools
```

---

```
curl 'http://127.0.0.1:1317/clp/swap' \
  -H 'Connection: keep-alive' \
  -H 'Pragma: no-cache' \
  -H 'Cache-Control: no-cache' \
  -H 'sec-ch-ua: "Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"' \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'DNT: 1' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36' \
  -H 'Content-Type: application/json;charset=UTF-8' \
  -H 'Origin: http://localhost:8080' \
  -H 'Sec-Fetch-Site: cross-site' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Referer: http://localhost:8080/' \
  -H 'Accept-Language: en-US,en;q=0.9,la;q=0.8' \
  --data-binary '{"base_req":{"chain_id":"sifchain","from":"sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl"},"received_asset":{"source_chain":"sifchain","symbol":"clink","ticker":"clink"},"sent_amount":"100000000000000000000","sent_asset":{"source_chain":"sifchain","symbol":"cusdc","ticker":"cusdc"},"signer":"sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl"}' \
  --compressed
```

REQ

```json
{
  "base_req": {
    "chain_id": "sifchain",
    "from": "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl"
  },
  "received_asset": {
    "source_chain": "sifchain",
    "symbol": "clink",
    "ticker": "clink"
  },
  "sent_amount": "100000000000000000000",
  "sent_asset": {
    "source_chain": "sifchain",
    "symbol": "cusdc",
    "ticker": "cusdc"
  },
  "signer": "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl"
}
```

RES

```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "clp/Swap",
        "value": {
          "Signer": "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl",
          "SentAsset": { "symbol": "cusdc" },
          "ReceivedAsset": { "symbol": "clink" },
          "SentAmount": "100000000000000000000"
        }
      }
    ],
    "fee": { "amount": [], "gas": "200000" },
    "signatures": null,
    "memo": ""
  }
}
```

```bash
sifnodecli query clp pools

# Check balance
sifnodecli q auth account $(sifnodecli keys show juniper -a)

# Swap
sifnodecli tx clp swap --from juniper --sentSymbol cusdc --receivedSymbol clink --sentAmount 100000000000000000000

sifnodecli q auth account $(sifnodecli keys show juniper -a)
```

---

```
curl 'http://127.0.0.1:1317/txs' \
 -H 'Connection: keep-alive' \
 -H 'Pragma: no-cache' \
 -H 'Cache-Control: no-cache' \
 -H 'sec-ch-ua: "Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"' \
 -H 'Accept: application/json, text/plain, _/_' \
 -H 'DNT: 1' \
 -H 'sec-ch-ua-mobile: ?0' \
 -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36' \
 -H 'Content-Type: application/json;charset=UTF-8' \
 -H 'Origin: http://localhost:8080' \
 -H 'Sec-Fetch-Site: cross-site' \
 -H 'Sec-Fetch-Mode: cors' \
 -H 'Sec-Fetch-Dest: empty' \
 -H 'Referer: http://localhost:8080/' \
 -H 'Accept-Language: en-US,en;q=0.9,la;q=0.8' \
 --data-binary '{"tx":{"msg":[{"type":"clp/Swap","value":{"Signer":"sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl","SentAsset":{"symbol":"cusdc"},"ReceivedAsset":{"symbol":"clink"},"SentAmount":"100000000000000000000"}}],"fee":{"amount":[{"amount":"0","denom":"rowan"}],"gas":"200000"},"memo":"","signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"AniBCyImSLLldpjRDyunZz+aerZhkWYWVmYbD96BDo5g"},"signature":"RPmVcIPggAfw60Zwn8WcWsBZHTpW6+GRvKifTcmyVhdOVcW3fBXW8lxvaw+uKRdh2d3yowxDfU6zy7SgxxVncA=="}]},"mode":"block"}' \
 --compressed
```

REQ

```json
{
  "tx": {
    "msg": [
      {
        "type": "clp/Swap",
        "value": {
          "Signer": "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl",
          "SentAsset": { "symbol": "cusdc" },
          "ReceivedAsset": { "symbol": "clink" },
          "SentAmount": "100000000000000000000"
        }
      }
    ],
    "fee": { "amount": [{ "amount": "0", "denom": "rowan" }], "gas": "200000" },
    "memo": "",
    "signatures": [
      {
        "pub_key": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AniBCyImSLLldpjRDyunZz+aerZhkWYWVmYbD96BDo5g"
        },
        "signature": "RPmVcIPggAfw60Zwn8WcWsBZHTpW6+GRvKifTcmyVhdOVcW3fBXW8lxvaw+uKRdh2d3yowxDfU6zy7SgxxVncA=="
      }
    ]
  },
  "mode": "block"
}
```

RES

```json
{
  "height": "85",
  "txhash": "94C8DB6043A96ACE525B8B44B9B3D1A343DB4E3E11D1851E81AB65242D474473",
  "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"swap\"},{\"key\":\"sender\",\"value\":\"sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl\"},{\"key\":\"sender\",\"value\":\"sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85\"},{\"key\":\"module\",\"value\":\"clp\"},{\"key\":\"sender\",\"value\":\"sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl\"}]},{\"type\":\"swap\",\"attributes\":[{\"key\":\"swap_amount\",\"value\":\"5882114714235019420\"},{\"key\":\"liquidity_fee\",\"value\":\"1058799971037049\"},{\"key\":\"trade_slip\",\"value\":\"0\"},{\"key\":\"height\",\"value\":\"85\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85\"},{\"key\":\"sender\",\"value\":\"sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl\"},{\"key\":\"amount\",\"value\":\"100000000000000000000cusdc\"},{\"key\":\"recipient\",\"value\":\"sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl\"},{\"key\":\"sender\",\"value\":\"sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85\"},{\"key\":\"amount\",\"value\":\"99998000029999600004clink\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "swap"
            },
            {
              "key": "sender",
              "value": "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl"
            },
            {
              "key": "sender",
              "value": "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"
            },
            {
              "key": "module",
              "value": "clp"
            },
            {
              "key": "sender",
              "value": "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl"
            }
          ]
        },
        {
          "type": "swap",
          "attributes": [
            {
              "key": "swap_amount",
              "value": "5882114714235019420"
            },
            {
              "key": "liquidity_fee",
              "value": "1058799971037049"
            },
            {
              "key": "trade_slip",
              "value": "0"
            },
            {
              "key": "height",
              "value": "85"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"
            },
            {
              "key": "sender",
              "value": "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl"
            },
            {
              "key": "amount",
              "value": "100000000000000000000cusdc"
            },
            {
              "key": "recipient",
              "value": "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl"
            },
            {
              "key": "sender",
              "value": "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"
            },
            {
              "key": "amount",
              "value": "99998000029999600004clink"
            }
          ]
        }
      ]
    }
  ],
  "gas_wanted": "200000",
  "gas_used": "109427"
}
```
