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
	return func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var sifMsg SifchainMsg
		err := json.Unmarshal(msg, &sifMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		switch {
		case sifMsg.Swap != nil:
			return EncodeSwapMsg(sender, sifMsg.Swap)
		case sifMsg.AddLiquidity != nil:
			return EncodeAddLiquidityMsg(sender, sifMsg.AddLiquidity)
		case sifMsg.RemoveLiquidity != nil:
			return EncodeRemoveLiquidityMsg(sender, sifMsg.RemoveLiquidity)
		}

		return nil, fmt.Errorf("Unknown SifchainMsg type")
	}
}

func EncodeRemoveLiquidityMsg(sender sdk.AccAddress, msg *RemoveLiquidity) ([]sdk.Msg, error) {
	wBasisPoints, ok := sdk.NewIntFromString(msg.WBasisPoints)
	if !ok {
		return nil, fmt.Errorf("invalid w basis points %s", msg.WBasisPoints)
	}
	asymmetry, ok := sdk.NewIntFromString(msg.Asymmetry)
	if !ok {
		return nil, fmt.Errorf("invalid asymmetry %s", msg.Asymmetry)
	}
	removeLiquidityMsg := clptypes.NewMsgRemoveLiquidity(
		sender,
		clptypes.NewAsset(msg.ExternalAsset),
		sdk.Int(wBasisPoints),
		sdk.Int(asymmetry),
	)
	return []sdk.Msg{&removeLiquidityMsg}, nil
}

func EncodeAddLiquidityMsg(sender sdk.AccAddress, msg *AddLiquidity) ([]sdk.Msg, error) {

	nativeAssetAmount, ok := sdk.NewIntFromString(msg.NativeAssetAmount)
	if !ok {
		return nil, fmt.Errorf("invalid native asset amount %s", msg.NativeAssetAmount)
	}

	externalAssetAmount, ok := sdk.NewIntFromString(msg.ExternalAssetAmount)
	if !ok {
		return nil, fmt.Errorf("invalid external asset amount %s", msg.ExternalAssetAmount)
	}

	addLiquidityMsg := clptypes.NewMsgAddLiquidity(
		sender,
		clptypes.NewAsset(msg.ExternalAsset),
		sdk.Uint(nativeAssetAmount),
		sdk.Uint(externalAssetAmount),
	)

	return []sdk.Msg{&addLiquidityMsg}, nil

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
	swapMsg := clptypes.NewMsgSwap(
		sender,
		clptypes.NewAsset(msg.SentAsset),
		clptypes.NewAsset(msg.ReceivedAssed),
		sdk.Uint(sentAmount),
		sdk.Uint(minReceivedAmount),
	)
	return []sdk.Msg{&swapMsg}, nil
}
