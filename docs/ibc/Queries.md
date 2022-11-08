# Common queries 

## Channel Related

- Query all channels for a chain
```shell
sifnoded q ibc channel channels --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
```
```json
- channel_id: channel-0
  connection_hops:
  - connection-0
  counterparty:
    channel_id: channel-82
    port_id: transfer
  ordering: ORDER_UNORDERED
  port_id: transfer
  state: STATE_OPEN
  version: ics20-1

```

- Query all channels with a connection ID 
```shell
sifnoded q ibc channel connections connection-110 --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
```
```json
channels:
- channel_id: channel-101
  connection_hops:
  - connection-110
  counterparty:
    channel_id: channel-3
    port_id: transfer
  ordering: ORDER_UNORDERED
  port_id: transfer
  state: STATE_OPEN
  version: ics20-1
height:
  revision_height: "1846741"
  revision_number: "1"
```
- Query the client-state with channel-id

```shell
sifnoded q ibc channel client-state transfer channel-101  --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
````
```json
client_id: 07-tendermint-173
client_state:
  '@type': /ibc.lightclients.tendermint.v1.ClientState
  allow_update_after_expiry: false
  allow_update_after_misbehaviour: false
  chain_id: cosmoshub-testnet
  ... Other values are not that important for debugging

```

## Packet Related 
The transfer command emits the packet sequence in the events 
Sequence can be considered similar to nonce , but it is specific to a channel
```json
     {
                     "key":"packet_sequence",  
                     "value":"31"
     }
```

The packet_sequence can be used to query the state of the packet (in the order mentioned) 
```shell
1 sifnoded q ibc channel packet-commitment transfer channel-101 31  --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
2 sifnoded q ibc channel packet-receipt transfer channel-101 31 --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
3 sifnoded q ibc channel packet-ack transfer channel-101 31 --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
```


Query if packet is not received in the recipient chain
```shell
sifnoded q ibc channel unreceived-packets transfer channel-101 --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
```


Query if the ack for packet receipt is not present in the sending chain
```shell
sifnoded q ibc channel  unreceived-acks transfer channel-101 --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
```

## Other Helpful queries
- Get the denom trace from the hash
```shell
sifnoded q ibc-transfer denom-trace C782C1DE5F380BC8A5B7D490684894B439D31847A004B271D7B7BA07751E582A --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
```
```json
denom_trace:
  base_denom: uphoton
  path: transfer/channel-101

```

- Get Escrow address for a channel and port combination
Tokens send from Chain-1 to Chain-2 ,are escrowed in Chain-1,instead of Burning .
When a token comes back to the source chain, it gets released from the escrow address to the user instead of being Minted again.  
```shell
sifnoded q ibc-transfer escrow-address transfer channel-101 --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
```
```json
sif1j3mmq2dsfws0pv5fut3ce2252w0ere8g2alrvd
```

## Sifchain Related Queries 
- Query sifchain tokenregistry 
```shell
 sifnoded q ibc-transfer escrow-address transfer channel-101 --node=https://rpc-devnet.sifchain.finance:443 --chain-id=sifchain-devnet-1
```
```json
{
            "address": "",
            "base_denom": "uphoton",
            "decimals": "6",
            "denom": "ibc/C782C1DE5F380BC8A5B7D490684894B439D31847A004B271D7B7BA07751E582A",
            "display_name": "",
            "display_symbol": "",
            "external_symbol": "",
            "ibc_channel_id": "channel-101",
            "ibc_counterparty_chain_id": "cosmoshub-testnet",
            "ibc_counterparty_channel_id": "channel-3",
            "ibc_counterparty_denom": "",
            "network": "",
            "path": "transfer/channel-101",
            "permissions": [
                "CLP",
                "IBCEXPORT",
                "IBCIMPORT"
            ],
            "transfer_limit": "",
            "unit_denom": ""
        }
```
## Stuck Transfers

Use `sifnoded ibc-diag stuck-txs` to get a list of stuck IBC transfers.

For example, to get the list of stuck transfer from Sifchain to Terra:

```
sifnoded ibc-diag stuck-txs \
--src-node http://rpc.sifchain.finance:80 \
--dst-node http://public-node.terra.dev:26657 \
--src-channel channel-18 \
--dst-channel channel-7
```

and from Terra to Sifchain:

```
sifnoded ibc-diag stuck-txs \
--src-node http://public-node.terra.dev:26657 \
--dst-node http://rpc.sifchain.finance:80 \
--src-channel channel-7 \
--dst-channel channel-18
```

Use the regular IBC commands described above to find which channel ids to use.