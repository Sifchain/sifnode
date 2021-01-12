package txs

// DONTCOVER

import (
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/Sifchain/sifnode/x/ethbridge"
	bridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
)

// RelayToCosmos applies validator's signature to an EthBridgeClaim message containing
// information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayToCosmos(cosmosContext *types.CosmosContext, claim *bridgetypes.EthBridgeClaim) error {
	// Packages the claim as a Tendermint message
	msg := ethbridge.NewMsgCreateEthBridgeClaim(*claim)
	return SendMsgToCosmos(cosmosContext, msg)
}

// SendMsgToCosmos send message to Sifchain
func SendMsgToCosmos(cosmosContext *types.CosmosContext, msg sdk.Msg) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	// Prepare tx
	txBldr, err := utils.PrepareTxBuilder(cosmosContext.TxBldr, cosmosContext.CliCtx)
	if err != nil {
		return err
	}

	// Build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(cosmosContext.ValidatorName, cosmosContext.TempPassword, []sdk.Msg{msg})
	if err != nil {
		return err
	}

	// Broadcast to a Tendermint node
	res, err := cosmosContext.CliCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return err
	}

	if err = cosmosContext.CliCtx.PrintOutput(res); err != nil {
		return err
	}
	return nil
}

// SendOutMessage send cosmos message to Sifchain
func SendOutMessage(cosmosContext *types.CosmosContext, message *types.CosmosMsg) error {
	if message.ClaimType == types.MsgLock {
		msg := bridgetypes.NewMsgLock(message.EthereumChainID, sdk.AccAddress(message.CosmosSender), bridgetypes.EthereumAddress(message.EthereumReceiver),
			message.Amount, message.Symbol, message.CethAmount, message.MessageType)
		msg.SetSequence(uint64(message.CosmosSenderSequence.Int64()))
		return SendMsgToCosmos(cosmosContext, msg)
	}

	if message.ClaimType == types.MsgBurn {
		msg := bridgetypes.NewMsgBurn(message.EthereumChainID, sdk.AccAddress(message.CosmosSender), bridgetypes.EthereumAddress(message.EthereumReceiver),
			message.Amount, message.Symbol, message.CethAmount, message.MessageType)
		msg.SetSequence(uint64(message.CosmosSenderSequence.Int64()))
		return SendMsgToCosmos(cosmosContext, msg)
	}
	return nil
}
