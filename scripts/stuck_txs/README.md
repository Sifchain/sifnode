# Get Stuck TXs

These scripts retrieve the stuck transactions between sifchain and terra (or vice versa) for a given IBC client on either chain.

## Prerequisites

`terrad` is required, follow the instructions here:

```
https://github.com/terra-money/core
```

## Running the scripts

The main script is `get_stuck_txs` which takes two arguments, `client` and `chain`. The `client` and `chain` combination identify the source of the stuck txs e.g.

```
./get_stuck_txs --client 07-tendermint-42 --chain sif
```

Writes a list of stuck txs sent from sifchain client 07-tendermint-42 to `processed/sif/connection-21/missing_txs.data`

NOTE: if you chose a client on sifchain which doesn't connect to terra (or vice versa) then the scripts will either error or return erroneous data

## How do the scripts work?

All the data needed is queried from public nodes, `https://rpc-archive.sifchain.finance:443` for sifchain and `http://public-node.terra.dev:26657` for terra. The steps are:

1. Query nodes to get the client, connection and channel information of both the send and receive chain - this can all be determined from the sending client id

2. Query the sending chain for `send_packet` data, filtering on send connection id

3. Query the receiving chain for `recv_packet`, filtering on receive connection id

4. Query the sending chain for `timeout_packet` data, filtering on send channel 

5. Extract the `packet_sequence` of each message in the `send_packet`, `recv_packet` and `timeout_packet` data

6. Find all `send_packet`s that do not have a corresponding `recv_packet` or `timeout packet` 

## TODOs

There is uncertainty around the uniqueness of channels and connections (for a given client and across clients) and the impact this would have on the results

## TXs of interest

The client of interest is '07-tendermint-42' on the `sif` chain (this was discovered by asking Brent):

```
./get_stuck_txs --client 07-tendermint-42 --chain sif
```
The list of stuck transactions, with transaction data is written to 'processed/sif/connection-21/missing_txs_full.csv'

The list of stuck transactions is written to `processed/sif/connection-21/missing_txs.data`:

```
389
390
391
392
393
394
395
396
397
398
399
400
401
402
403
404
405
406
407
408
409
410
411
412
413
414
415
416
420
423
424
425
426
427
428
429
430
435
471
472
473
474
479
480
551
552
553
554
555
556
557
558
559
560
561
562
563
564
565
566
567
568
569
570
571
572
573
574
575
576
577
578
579
580
581
582
584
585
586
587
588
589
590
591
592
593
594
595
596
597
598
599
600
601
602
603
604
605
606
607
608
```

And the reverse direction:

```
./get_stuck_txs --client 07-tendermint-19 --chain terra
```

The list of stuck transactions, with transaction data is written to 'processed/terra/connection-19/missing_txs_full.csv'

The list of stuck transactions is written to `processed/terra/connection-19/missing_txs.data`:

```
2183
```