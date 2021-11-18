#IBC (Testing and Debugging)

General Flow 
1. User requests to transfer Tokens from an address in Chain-1 to an address in Chain-2
2. The Sending Chain
    - if is is the token source chain (origin chain), then the exported tokens are locked up in an escrow address
    - if the token did **not** orginate in the sending chain, then the tokens are burned on the sending chain, so that they can be unlocked on the other chain
3. If the sending chain is sifchain, the token is checked against the tokenregistry for the following:
    - Check if we need to modify decimal precision
    - Check if the token has permission for IBCEXPORT
4. The transfer packet then goes through the IBC Send,Receive,Ack flow. More details on events https://github.com/cosmos/ibc-go/v2/blob/main/modules/core/spec/06_events.md
5. The Receiving Chain
    - if is is the token source chain (origin chain), then the tokens are unlocked from the escrow address
    - if the token did **not** originate in the sending chain, then the tokens are minted with a new denom. The new denom is created by appending the port and the channel to the existing denom to create the denom trace.
6. If the receiving chain is sifchain, the token is checked against the token registry 
    - Check if we need to modify the decimal precision
    - Check if the token has IBCIMPORT permission
    - Check if the token is Whitelisted
    
##CLI

Transfer Funds from Chain-1 to Chain-2
```shell
sifnoded tx ibc-transfer transfer 
transfer  -> Port
channel-101  -> Channel
cosmos1syavy2npfyt9tcncdtsdzf7kny9lh777pahuux  -> Receiver( in Chain-2)
1ibc/C782C1DE5F380BC8A5B7D490684894B439D31847A004B271D7B7BA07751E582A (Tokens sent)
--from=sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd (Sender in Chain-1)
--node=https://rpc-devnet.sifchain.finance:443 (Broadcasting Node for Chain-1)
--chain-id=sifchain-devnet-1  (chainID for Chain-1)
--gas-prices=1rowan (Gas Prices to pay in Chain-1)
--gas=5000000 (Max gas the Sender is willing to pay)
--y  (Auto-Confirm)
--packet-timeout-timestamp=600000000000 (The transfer operation will timieout after this number + concensus.CurrentTimestamp)
--keyring-backend=test (Keyring Backend To use)
```

Sample Response
```json
{
   "height":"1846220",
   "txhash":"7FB00D756A19D08BE28926192AE0C4EBB8D2C6F7AD18C7293CFE5F4191C293ED",
   "codespace":"",
   "code":0,
   "data":"0A0A0A087472616E73666572",
   "raw_log":"[{\"events\":[{\"type\":\"ibc_transfer\",\"attributes\":[{\"key\":\"sender\",\"value\":\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\"},{\"key\":\"receiver\",\"value\":\"cosmos1syavy2npfyt9tcncdtsdzf7kny9lh777pahuux\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"transfer\"},{\"key\":\"sender\",\"value\":\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\"},{\"key\":\"module\",\"value\":\"ibc_channel\"},{\"key\":\"module\",\"value\":\"transfer\"}]},{\"type\":\"send_packet\",\"attributes\":[{\"key\":\"packet_data\",\"value\":\"{\\\"amount\\\":\\\"1\\\",\\\"denom\\\":\\\"transfer/channel-101/uphoton\\\",\\\"receiver\\\":\\\"cosmos1syavy2npfyt9tcncdtsdzf7kny9lh777pahuux\\\",\\\"sender\\\":\\\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\\\"}\"},{\"key\":\"packet_timeout_height\",\"value\":\"0-104721\"},{\"key\":\"packet_timeout_timestamp\",\"value\":\"1629903792589870129\"},{\"key\":\"packet_sequence\",\"value\":\"31\"},{\"key\":\"packet_src_port\",\"value\":\"transfer\"},{\"key\":\"packet_src_channel\",\"value\":\"channel-101\"},{\"key\":\"packet_dst_port\",\"value\":\"transfer\"},{\"key\":\"packet_dst_channel\",\"value\":\"channel-3\"},{\"key\":\"packet_channel_ordering\",\"value\":\"ORDER_UNORDERED\"},{\"key\":\"packet_connection\",\"value\":\"connection-110\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"sif1yl6hdjhmkf37639730gffanpzndzdpmhtzelcg\"},{\"key\":\"sender\",\"value\":\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\"},{\"key\":\"amount\",\"value\":\"1ibc/C782C1DE5F380BC8A5B7D490684894B439D31847A004B271D7B7BA07751E582A\"}]}]}]",
   "logs":[
      {
         "msg_index":0,
         "log":"",
         "events":[
            {
               "type":"ibc_transfer",
               "attributes":[
                  {
                     "key":"sender",
                     "value":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
                  },
                  {
                     "key":"receiver",
                     "value":"cosmos1syavy2npfyt9tcncdtsdzf7kny9lh777pahuux"
                  }
               ]
            },
            {
               "type":"message",
               "attributes":[
                  {
                     "key":"action",
                     "value":"transfer"
                  },
                  {
                     "key":"sender",
                     "value":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
                  },
                  {
                     "key":"module",
                     "value":"ibc_channel"
                  },
                  {
                     "key":"module",
                     "value":"transfer"
                  }
               ]
            },
            {
               "type":"send_packet",
               "attributes":[
                  {
                     "key":"packet_data",
                     "value":"{\"amount\":\"1\",\"denom\":\"transfer/channel-101/uphoton\",\"receiver\":\"cosmos1syavy2npfyt9tcncdtsdzf7kny9lh777pahuux\",\"sender\":\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\"}"
                  },
                  {
                     "key":"packet_timeout_height",
                     "value":"0-104721"
                  },
                  {
                     "key":"packet_timeout_timestamp",
                     "value":"1629903792589870129"
                  },
                  {
                     "key":"packet_sequence",   // Can be used to query status of the packet
                     "value":"31"
                  },
                  {
                     "key":"packet_src_port",
                     "value":"transfer"
                  },
                  {
                     "key":"packet_src_channel",
                     "value":"channel-101"
                  },
                  {
                     "key":"packet_dst_port",
                     "value":"transfer"
                  },
                  {
                     "key":"packet_dst_channel",
                     "value":"channel-3"
                  },
                  {
                     "key":"packet_channel_ordering",
                     "value":"ORDER_UNORDERED"
                  },
                  {
                     "key":"packet_connection",
                     "value":"connection-110"
                  }
               ]
            },
            {
               "type":"transfer",
               "attributes":[
                  {
                     "key":"recipient",
                     "value":"sif1yl6hdjhmkf37639730gffanpzndzdpmhtzelcg"
                  },
                  {
                     "key":"sender",
                     "value":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
                  },
                  {
                     "key":"amount",
                     "value":"1ibc/C782C1DE5F380BC8A5B7D490684894B439D31847A004B271D7B7BA07751E582A"
                  }
               ]
            }
         ]
      }
   ],
   "info":"",
   "gas_wanted":"5000000",
   "gas_used":"213491",
   "tx":null,
   "timestamp":""
}
```


