# Current Peggy Architecture

We have moved to a new architecture, where instead of having any queue where we store data, instead, we just query ethereum blocks that are 50 blocks in the past and get the log data from those blocks. That way, we know that data is not going to be moved around and we don’t have to do validation on our client side and keep blockchain data in memory

Anytime a new block comes in, we look 50 blocks in the past and see if there were any logs going to the bridgebank in that block. If there were, then we pick that up, package into a cosmos tx and send it sifchain. If there is no data 50 blocks back, we just leave it. Then, we store that block we just looked at in a leveldb database on disk so that if the relayer crashes, it can pick back up where it left off
