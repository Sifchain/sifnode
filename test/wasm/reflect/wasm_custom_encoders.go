package reflect

import (
	"encoding/json"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// reflectEncoders needs to be registered in to handle custom message callbacks
func ReflectEncoders(cdc codec.Codec) *wasmkeeper.MessageEncoders {
	return &wasmkeeper.MessageEncoders{
		Custom: EncodeSifchainMessage(cdc),
	}
}

// EncodeSifchainMessage decodes msg.Data to an sdk.Msg using proto Any and json
// encoding. This needs to be registered on the Encoders
func EncodeSifchainMessage(cdc codec.Codec) wasmkeeper.CustomEncoder {
	return func(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var sifMsg SifchainMsg
		err := json.Unmarshal(msg, &sifMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		fmt.Printf("@@@ SifchainMsg: %#v\n", sifMsg)

		if sifMsg.Swap != nil {
			return EncodeSwapMsg(_sender, sifMsg.Swap)
		}

		return nil, fmt.Errorf("@@@ Unknown SifchainMsg type")
	}
}

func EncodeSwapMsg(sender sdk.AccAddress, msg *Swap) ([]sdk.Msg, error) {
	// ATTENTION
	// cosmwasm tends to always user sender as signer
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	sentAmount, ok := sdk.NewIntFromString(msg.SentAmount)
	if !ok {
		return nil, fmt.Errorf("Invalid sent amount %s", msg.SentAmount)
	}

	minReceivedAmount, ok := sdk.NewIntFromString(msg.MinReceivedAmount)
	if !ok {
		return nil, fmt.Errorf("Invalid min received amount %s", msg.MinReceivedAmount)
	}

	swapMsg := clptypes.NewMsgSwap(
		signer,
		clptypes.NewAsset(msg.SentAsset),
		clptypes.NewAsset(msg.ReceivedAssed),
		sdk.Uint(sentAmount),
		sdk.Uint(minReceivedAmount),
	)

	fmt.Printf("@@@ sent amount: %v\n", swapMsg.SentAmount)
	fmt.Printf("@@@ min received amount: %v\n", swapMsg.MinReceivingAmount)

	return []sdk.Msg{&swapMsg}, nil
}
