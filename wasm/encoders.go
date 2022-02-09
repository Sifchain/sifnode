package wasm

import (
	"encoding/json"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Encoders needs to be registered in order to handle custom sifchain messages
func Encoders(cdc codec.Codec) *wasmkeeper.MessageEncoders {
	return &wasmkeeper.MessageEncoders{
		Custom: EncodeSifchainMessage(cdc),
	}
}

// EncodeSifchainMessage encodes the contents of a SifchainMsg into an SDK msg
// destined to a sifchain-specific module
func EncodeSifchainMessage(cdc codec.Codec) wasmkeeper.CustomEncoder {
	return func(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var sifMsg SifchainMsg
		err := json.Unmarshal(msg, &sifMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		if sifMsg.Swap != nil {
			return EncodeSwapMsg(_sender, sifMsg.Swap)
		}
		return nil, fmt.Errorf("Unknown SifchainMsg type")
	}
}

func EncodeSwapMsg(sender sdk.AccAddress, msg *Swap) ([]sdk.Msg, error) {
	sentAmount, ok := sdk.NewIntFromString(msg.SentAmount)
	if !ok {
		return nil, fmt.Errorf("Invalid sent amount %s", msg.SentAmount)
	}
	minReceivedAmount, ok := sdk.NewIntFromString(msg.MinReceivedAmount)
	if !ok {
		return nil, fmt.Errorf("Invalid min received amount %s", msg.MinReceivedAmount)
	}
	// XXX note that the swap signer is the contract
	swapMsg := clptypes.NewMsgSwap(
		sender,
		clptypes.NewAsset(msg.SentAsset),
		clptypes.NewAsset(msg.ReceivedAssed),
		sdk.Uint(sentAmount),
		sdk.Uint(minReceivedAmount),
	)
	return []sdk.Msg{&swapMsg}, nil
}
