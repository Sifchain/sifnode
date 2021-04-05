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

### replay sifchain step by step
1. If notice the ebrelayer is offline, have some errors to send prophecy to smart contract in Ethereum or cross-chain transfer from Sifchain to Ethereum doesn't work, we need check the timestamp and estimate the block scope in sifchain side.
2. We choose those failed ebrelayer nodes, check their voting power and compute how many nodes we need run replay scripts on those nodes.
3. We also need check the blocks in Ethereum side to find out all prophecies transaction already sent by the ebrelayer, avoid send the same transaction again. It will not mint or burn the token twice since we have the unique prophecy ID. but it will waste of some gas fee. 
4. After confirm the from/to block number in both Ethereum and Sifchain, we can run ebrelayer replayEthereum. the command usage could be found in Command sample segment.
5. Check the output of replay script, it will print out all the cross-chain token transfer messages and tell you if the ebrelayer processed it or not.
6. Need run the replay script in several nodes until the prophecy completed in Ethereum side. It depends on each ebrelayer's voting power.

## Side effect
If the block scope set not correct, script may miss some messages or process the same messages twice. For duplicated message processing, there is no side effect on the network except extra gas fee. If missing some message, some prophecies may never be finalized. You can decide the recovery scope based on the working validators' power.
