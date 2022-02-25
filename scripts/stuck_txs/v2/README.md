# Finding Stuck IBC Transfers

1. Get list of stuck packets from Sichain to Terra:

```
./get_stuck_txs.sh > ~/stuck_packets.txt
```

2. Use the Go dbtool to get packet data for each packet:

```
go run ../../cmd/dbtool/main.go ibc get-transfers ~/stuck_packets.txt channel-18 http://rpc.sifchain.finance:80 > ~/stuck_transfers.csv
```

This outputs a csv file with all the stuck transfers data