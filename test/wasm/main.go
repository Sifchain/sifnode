package main

import (
	"encoding/json"
	"fmt"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/Sifchain/sifnode/app"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func main() {

	sifAddress, _ := sdk.AccAddressFromBech32("sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd")

	swapMsg := clptypes.NewMsgSwap(
		sifAddress,
		clptypes.NewAsset("rowan"),
		clptypes.NewAsset("ceth"),
		sdk.NewUint(1),
		sdk.NewUint(10),
	)

	rawSwapMessage, err := toReflectRawMsg(
		app.MakeTestEncodingConfig().Marshaler,
		&swapMsg,
	)
	if err != nil {
		panic(err)
	}

	reflectSwapMsg := ReflectHandleMsg{
		Reflect: &reflectPayload{
			Msgs: []wasmvmtypes.CosmosMsg{
				rawSwapMessage,
			},
		},
	}

	jsonReflectSwapMessage, err := json.Marshal(reflectSwapMsg)

	fmt.Printf("%s\n", jsonReflectSwapMessage)
}

// ReflectHandleMsg is used to encode handle messages
type ReflectHandleMsg struct {
	Reflect       *reflectPayload    `json:"reflect_msg,omitempty"`
	ReflectSubMsg *reflectSubPayload `json:"reflect_sub_msg,omitempty"`
	Change        *ownerPayload      `json:"change_owner,omitempty"`
}

type ownerPayload struct {
	Owner sdk.Address `json:"owner"`
}

type reflectPayload struct {
	Msgs []wasmvmtypes.CosmosMsg `json:"msgs"`
}

type reflectSubPayload struct {
	Msgs []wasmvmtypes.SubMsg `json:"msgs"`
}

type reflectCustomMsg struct {
	Debug string `json:"debug,omitempty"`
	Raw   []byte `json:"raw,omitempty"`
}

// toReflectRawMsg encodes an sdk msg using any type with json encoding.
// Then wraps it as an opaque message
func toReflectRawMsg(cdc codec.Codec, msg sdk.Msg) (wasmvmtypes.CosmosMsg, error) {
	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return wasmvmtypes.CosmosMsg{}, err
	}
	rawBz, err := cdc.MarshalJSON(any)
	if err != nil {
		return wasmvmtypes.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	customMsg, err := json.Marshal(reflectCustomMsg{
		Raw: rawBz,
	})
	res := wasmvmtypes.CosmosMsg{
		Custom: customMsg,
	}
	return res, nil
}
