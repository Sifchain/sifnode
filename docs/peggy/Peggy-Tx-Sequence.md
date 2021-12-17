<!-- MermaidJS doesnt support sequence diagram title sadface -->

<!-- Initiating a burn flow on cosmos -->
:::mermaid
sequenceDiagram
    %% ethbridge.keeper.msg_server::Burn
    ethbridge_module ->> keeper.globalSequence: getSequence(networkDescriptor)

    keeper.globalSequence ->> ethbridge_module: sequence
    Note over ethbridge_module, keeper.globalSequence: 1 if not found
    %% TODO: Consider using box notation to denote func location?
    %% ethbridge.keeper.globalNonce::UpdateGlobalSequence
    ethbridge_module ->> keeper.globalSequence: store [networkDescriptor -> globalSequence+1]
    ethbridge_module ->> keeper.globalSequenceToBlock: store [networkDescriptor, globalSequence -> blockHeight]
    ethbridge_module ->> cosmosEventBus: burnMsg {netwrokDescriptor, prophecyId, globalSequence}
:::

<!-- Ebrelayer cosmos witness -->
:::mermaid
sequenceDiagram
    %% cmd.ebrelayer.relayer.cosmos::CheckSequenceAndProcess
    %% TODO: Denote this is looped forever
    ebrelayer ->> grpcClient: witnessLockBurnSequenceRequest(networkDescriptor, relayerAddress)
    %% TODO: Should i add detail here? This is fulfilled by oracle.keeper
    grpcClient ->> ebrelayer: sequence
    Note over grpcClient, ebrelayer: 0 if not found
    ebrelayer ->> grpcClient: globalSequenceBlockNumberRequest(networkDescriptor, sequence+1)
    grpcClient ->> ebrelayer: blockHeight
    note over grpcClient, ebrelayer: 0 if not found

    %% ebrelayer ->> ebrelayer: process from blockHeight to currentHeight or max
    loop blockHeight -> endBlock
        loop tx
            ebrelayer ->> ethbridge_module: Broadcast MsgSignProphecy{cosmosSender, networkDescriptor, prophecyId, EthereumAddr, Signature}
        end
    end
:::

<!-- TODO: Might not need to separate them, coz wanna show nonce in keeper -->
<!-- Cosmos ethbrdige module SignProphecy msg handling-->
:::mermaid
sequenceDiagram
    %% ebrelayer.relayer.cosmos::witnessSignProphecyId
    ebrelayer ->> ethbridge_module: MsgSignProphecy
    %% oracle.keeper.keeper.go
    ethbridge_module ->> keeper: witnessLockBurnSequence(networkDescriptor, validatorAddress)
    keeper ->> ethbridge_module: sequence
    note over keeper, ethbridge_module: 0 if not found
    note over ethbridge_module: exit if sequence != 0 AND not 1+sequence in msg

    ethbridge_module ->> ethbridge_module: AppendValidator, AppendSignature, UpdateProphecyStatus
    ethbridge_module ->> keeper: SetProphecy(Prophecy)
    ethbridge_module ->> keeper: SetWitnessLockBurnNonce(networkDescriptor, validatorAddress, prophecyId.GlobalSequence)


    alt newPropehcyStatus == Success
        ethbridge_module ->> cosmosEventBus: ProphecyComplete
    end
:::

:::mermaid
sequenceDiagram
    loop interval
        ebrelayer ->> cosmosBridgeEvmContract: GetLastNonceSubmitted()
        cosmosBridgeEvmContract ->> ebrelayer: nonce
        note over ebrelayer: Nonce here === globalSequence the contract has seen
        ebrelayer ->> keeper: ProphciesCompletedRequest(networkDescriptor, nonce+1)
        loop while prophecy is valid
            keeper ->> keeper: GetProphecyId(networkDescriptor, globalSequence)
        end
        keeper ->> ebrelayer: []ProphecyInfo
        ebrelayer ->> cosmosBridgeEvmContract: BatchSubmitProphecyClaimAggregatedSigs(batchClaimData, batchSig)
    end

:::