# Burn bug

Get message for burn: `/ethbridge/burn`

```ts
async function burn(params: {
  fromAddress: string;
  ethereumRecipient: string;
  assetAmount: AssetAmount;
}) {
  const web3 = await ensureWeb3();
  const ethereumChainId = await web3.eth.net.getId();
  const tokenAddress =
    (params.assetAmount.asset as Token).address ?? ETH_ADDRESS;

  return await sifUnsignedClient.burn({
    ethereum_receiver: params.ethereumRecipient,
    base_req: {
      chain_id: sifChainId,
      from: params.fromAddress,
    },
    amount: params.assetAmount.amount.toString(),
    symbol: params.assetAmount.asset.symbol,
    cosmos_sender: params.fromAddress,
    ethereum_chain_id: `${ethereumChainId}`,
    token_contract_address: tokenAddress,
  });
},

```

```json
{
  "ethereum_receiver": "0x627306090abaB3A6e1400e9345bC60c78a8BEf57",
  "base_req": {
    "chain_id": "sifchain",
    "from": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"
  },
  "amount": "1",
  "symbol": "ceth",
  "cosmos_sender": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
  "ethereum_chain_id": "5777",
  "token_contract_address": "0x0000000000000000000000000000000000000000"
}
```

```
curl 'http://127.0.0.1:1317/ethbridge/burn' \
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
  --data-binary '{"ethereum_receiver":"0x627306090abaB3A6e1400e9345bC60c78a8BEf57","base_req":{"chain_id":"sifchain","from":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"},"amount":"1","symbol":"ceth","cosmos_sender":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5","ethereum_chain_id":"5777","token_contract_address":"0x0000000000000000000000000000000000000000"}' \
  --compressed
```

Returns:

```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "ethbridge/MsgBurn",
        "value": {
          "cosmos_sender": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
          "amount": "1",
          "symbol": "ceth",
          "ethereum_chain_id": "5777",
          "ethereum_receiver": "0x627306090abaB3A6e1400e9345bC60c78a8BEf57"
        }
      }
    ],
    "fee": { "amount": [], "gas": "200000" },
    "signatures": null,
    "memo": ""
  }
}
```

Now we have our message we sign and broadcast it.

```ts
const result = await api.EthbridgeService.burn({
  assetAmount,
  ethereumRecipient: store.wallet.eth.address,
  fromAddress: store.wallet.sif.address,
});

return await api.SifService.signAndBroadcast(result.value.msg);
```

So at this stage msg (`result.value.msg`) is the following array:

```json
[
  {
    "type": "ethbridge/MsgBurn",
    "value": {
      "cosmos_sender": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
      "amount": "1",
      "symbol": "ceth",
      "ethereum_chain_id": "5777",
      "ethereum_receiver": "0x627306090abaB3A6e1400e9345bC60c78a8BEf57"
    }
  }
]
```

`client` here is an instance of SigningCosmosClient which has been initialized with a wallet mnemonic of akasha:

'hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard'

Which becomes the following address:

sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5

```ts

 async function signAndBroadcast(msg: Msg | Msg[], memo?: string) {
      if (!client) throw "No client. Please sign in.";
      try {
        const fee = {
          amount: coins(0, "rowan"),
          gas: "200000", // need gas fee for tx to work - see genesis file
        };

        const msgArr = Array.isArray(msg) ? msg : [msg];

        console.log("signAndBroadcast:", JSON.stringify({ msgArr, fee, memo }));
        const txHash = await client.signAndBroadcast(msgArr, fee, memo);

        if (isBroadcastTxFailure(txHash)) {
          console.log(txHash);
          console.log(txHash.rawLog);
          throw new Error(txHash.rawLog);
        }

        triggerUpdate();

        return txHash;
      } catch (err) {
        console.error(err);
      }
    },
```

The signing process makes the following rest call:

```json
{
  "tx": {
    "msg": [
      {
        "type": "ethbridge/MsgBurn",
        "value": {
          "cosmos_sender": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
          "amount": "1",
          "symbol": "ceth",
          "ethereum_chain_id": "5777",
          "ethereum_receiver": "0x627306090abaB3A6e1400e9345bC60c78a8BEf57"
        }
      }
    ],
    "fee": { "amount": [{ "amount": "0", "denom": "rowan" }], "gas": "200000" },
    "memo": "",
    "signatures": [
      {
        "pub_key": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq"
        },
        "signature": "KVryJDDdYV/zDRrY9K6QdF7mgI094JdMdcWegExxzYAPDGsfoVg/fNwhr5LfAYCwDigOMl0YMxGhDa1UN3DBTA=="
      }
    ]
  },
  "mode": "block"
}
```

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
 --data-binary '{"tx":{"msg":[{"type":"ethbridge/MsgBurn","value":{"cosmos_sender":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5","amount":"1","symbol":"ceth","ethereum_chain_id":"5777","ethereum_receiver":"0x627306090abaB3A6e1400e9345bC60c78a8BEf57"}}],"fee":{"amount":[{"amount":"0","denom":"rowan"}],"gas":"200000"},"memo":"","signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq"},"signature":"KVryJDDdYV/zDRrY9K6QdF7mgI094JdMdcWegExxzYAPDGsfoVg/fNwhr5LfAYCwDigOMl0YMxGhDa1UN3DBTA=="}]},"mode":"block"}' \
 --compressed
```

Which returns the following value

```json
{
  "height": "0",
  "txhash": "39C66705D5A676A4AEBF9A15F2C02366CD90F1CF4E131D39B3664C586A34ED8C",
  "codespace": "sdk",
  "code": 4,
  "raw_log": "unauthorized: signature verification failed; verify correct account sequence and chain-id",
  "gas_wanted": "200000",
  "gas_used": "29340"
}
```

---

# Analysis of swap signing rest calls

Call the swap REST endpoint with the following data in the payload:

```json
{
  "base_req": {
    "chain_id": "sifchain",
    "from": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"
  },
  "received_asset": {
    "source_chain": "sifchain",
    "symbol": "cbtk",
    "ticker": "cbtk"
  },
  "sent_amount": "1000",
  "sent_asset": {
    "source_chain": "sifchain",
    "symbol": "catk",
    "ticker": "catk"
  },
  "signer": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"
}
```

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
  --data-binary '{"base_req":{"chain_id":"sifchain","from":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"},"received_asset":{"source_chain":"sifchain","symbol":"cbtk","ticker":"cbtk"},"sent_amount":"1000","sent_asset":{"source_chain":"sifchain","symbol":"catk","ticker":"catk"},"signer":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"}' \
  --compressed
```

Response

```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "clp/Swap",
        "value": {
          "Signer": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
          "SentAsset": { "symbol": "catk" },
          "ReceivedAsset": { "symbol": "cbtk" },
          "SentAmount": "1000"
        }
      }
    ],
    "fee": { "amount": [], "gas": "200000" },
    "signatures": null,
    "memo": ""
  }
}
```

---

We then sign it using `signAndBroadcast` to deliver the REST call

```ts
const tx = await api.ClpService.swap({
  fromAddress: state.address,
  sentAmount,
  receivedAsset,
});

return await api.SifService.signAndBroadcast(tx.value.msg);
```

Request

```json
{
  "tx": {
    "msg": [
      {
        "type": "clp/Swap",
        "value": {
          "Signer": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
          "SentAsset": { "symbol": "catk" },
          "ReceivedAsset": { "symbol": "cbtk" },
          "SentAmount": "1000"
        }
      }
    ],
    "fee": { "amount": [{ "amount": "0", "denom": "rowan" }], "gas": "200000" },
    "memo": "",
    "signatures": [
      {
        "pub_key": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq"
        },
        "signature": "xavoNyYUCiqgh+u1mafHgpmlrhJUENuCuF5VhO3nBkwyCB0T4oER5FZWiaWWc2BkNGBr6QOwwiPpk7ORbpuzng=="
      }
    ]
  },
  "mode": "block"
}
```

```
curl 'http://127.0.0.1:1317/txs' \
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
  --data-binary '{"tx":{"msg":[{"type":"clp/Swap","value":{"Signer":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5","SentAsset":{"symbol":"catk"},"ReceivedAsset":{"symbol":"cbtk"},"SentAmount":"1000"}}],"fee":{"amount":[{"amount":"0","denom":"rowan"}],"gas":"200000"},"memo":"","signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq"},"signature":"xavoNyYUCiqgh+u1mafHgpmlrhJUENuCuF5VhO3nBkwyCB0T4oER5FZWiaWWc2BkNGBr6QOwwiPpk7ORbpuzng=="}]},"mode":"block"}' \
  --compressed
```

```json
{
  "height": "0",
  "txhash": "6E6B994C7B7E993F6B3FBF16816D40B06C924827BC2F76B4A17CAD381E2707CE",
  "code": 19
}
```
