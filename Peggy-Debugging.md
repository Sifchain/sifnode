# Current Peggy Issues
## Issue 1
Currently peggy faces an issue where transactions are getting stuck in the queue as they are written to a block number in a mapping that the relayer forgets about as those transactions are not added to the array of that block number’s key in time. To summarize, there are too many transactions that are being written to a single index in the mapping, and the mapping gets read and cleared at that index by the relayer. Then, additional events are added to that index of the mapping, but the relayer has already moved on so those transactions are effectively forgotten.

There are two channels in the relayer that actually matter. 1, the new blocks that come in, and 2, the one that reads in the individual events from the blocks. The 2nd channel reads in transactions one by one by listening to the bridgebank contract. In doing this, it can take many blocks for a previous block’s log data to be completely added to that block. This causes an issue as you could be in a state where you see a new ethereum block, look back x number of blocks in your queue, read data out of the queue, clear the queue of that block, and the 2nd channel is still writing logs to that block. At this point, the 1st channel has moved on and will no longer look for blocks that far in the past. This results in transactions from that block that were added to the mapping after this happened to be dropped and never sent to sifchain.

## Issue 2
One other issue we have encountered is a transaction about to be broadcast to sifchain, but it never does and returns an error. Before, we have not logged this error out. Now, we have added error handling so that next time we see this issue, we can catch it. This issue occured once after 2000 tx's from ethereum so seems rarer.


# How to reproduce
1. run ```setup.sh``` in the root folder of sifnode. This will start up your environment and get sifnode, the relayer and ethereum going.

2. In a new terminal, cd into ```test/integration/vagrant/data/logs```.

3. In your old terminal, run this command to ensure your sifchain account has no balance:
```sifnodecli q account sif19yrthf758lvl5nrfrkhg8ndllve7ngpjanaxgn```

3.5 (Optional) in the ```test/integration/vagrant/data/logs``` folder, run this command to see the output and what is happening with the relayer and sifnode live ```tail -F sifnoded.log ebrelayer.log```.

4. In another terminal, cd into the ```smart-contracts``` directory and run ```COUNT=150 BRIDGEBANK_ADDRESS=0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4 truffle exec scripts/sendLockTx.js  sif19yrthf758lvl5nrfrkhg8ndllve7ngpjanaxgn eth 1``` This should send 150 transactions across the bridge in quick succession.

5. Flip back to the ```test/integration/vagrant/data/logs``` folder and run ```cat ebrelayer.log | grep "Add event into buffer" | wc -l ; cat ebrelayer.log | grep "Witnessed tx" | wc -l``` to see how many transactions have made their way into the relayer. The end state of running this command after the relayer processes all of the transactions should output 450 on two separate lines.

6. After the script in the smart contracts directory finishes, keep running the script mentioned in step 5 to see how many transactions the relayer has observed. It may take a few minutes for all of the transactions to be observed.

7. Run step 3 again after all transactions clear in the relayer and make sure that this sifnode user has 450 ceth.

8. While running the following commands, it is very likely that the relayer enters into an invalid state where transactions will be left in the queue. You can find the explanation for why this happens at the top of this document, but here how you can find that. 1, if your balance is not 450 ceth by the time this is all over, then transactions were dropped. In order to see this, pull up the window where you are displaying the spew dump of all of the block data. If that block data consistently has data in it from an older block and doesn't clear, then you know you have hit the first issue.

Debugging the other issue is outside of the scope of this document and I can write more on that later. For now, we are trying to replicate this issue where tx's get stuck in the queue indefinitely.