# replay instruction
The role of ebrelayer is critical in Sifchain network. If some ebrelayers are offline, then prophecies on both Sifnode and Ethereum can't be finalized. To recover from the situation, the replay scripts are created to check if any missed cross chain transaction.

## Algorithm
We have two sub-commands 'ebrelayer replayCosmos' and 'ebrelayer replayEthereum'. For replayCosmos, we need input the scope of blocks the script to seach. So there are cosmosStartBlock, cosmosEndBlock, ethereumStartBlock and ethereumEndBlock as arguments. At first, script get all prophecies sent by this validator in Ethereum within the scope. Then get all lock/burn messages in Sifchain. If the lock/burn message already processed, via checking against the prophecies, then skip the message. Otherwise, script process the message and send the prophecy transaction to smart contract in Ethereum.

For replayEthereum, the algorithm is the same, just on the opposite direction.

## Command sample
ebrelayer replayCosmos tcp://localhost:26657 ws://localhost:7545/ 0xFB88dE099e13c3ED21F80a7a1E49f8CAEcF10df6 100 200  20 25 --chain-id=sifchain


ebrelayer replayEthereum tcp://localhost:26657 ws://localhost:7545/ 0xFB88dE099e13c3ED21F80a7a1E49f8CAEcF10df6 sif "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" 15 20 10 20 --chain-id=sifchain

## Recovery process
After noticed the ebrelayer down, can't get new message or can't send transaction to Ethereum/Sifchain. We need record the timestamp and estimate the block height, then restart the ebrealyer immediately. 

The block scope you input is not neccessary to be very accurate. But must include the blocks when ebrelayer not working. Otherwise we can't gunratee all messages will be replayed.

## Side effect
If the block scope set not correct, script may miss some messages or process the same messages twice. For duplicated message processing, there is no side effect on the network except extra gas fee. If missing some message, some prophecies may never be finalized. You can decide the recovery scope based on the working validators' power.
