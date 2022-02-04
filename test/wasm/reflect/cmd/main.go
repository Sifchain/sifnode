package main

import (
	"encoding/json"
	"fmt"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/test/wasm/reflect"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {

	app.SetConfig(true)

	contractAddress, err := sdk.AccAddressFromBech32("sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6")
	if err != nil {
		panic(err)
	}

	swapMsg := clptypes.NewMsgSwap(
		contractAddress,
		clptypes.NewAsset("rowan"),
		clptypes.NewAsset("ceth"),
		sdk.NewUint(20000),
		sdk.NewUint(0),
	)

	rawSwapMessage, err := reflect.ToReflectRawMsg(
		app.MakeTestEncodingConfig().Marshaler,
		&swapMsg,
	)
	if err != nil {
		panic(err)
	}

	reflectSwapMsg := reflect.ReflectHandleMsg{
		Reflect: &reflect.ReflectPayload{
			Msgs: []wasmvmtypes.CosmosMsg{
				rawSwapMessage,
			},
		},
	}

	jsonReflectSwapMessage, err := json.Marshal(reflectSwapMsg)

	fmt.Printf("%s\n", jsonReflectSwapMessage)
}
