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
		Custom: FromReflectCustomMsg(cdc),
	}
}

// FromReflectCustomMsg decodes msg.Data to an sdk.Msg using proto Any and json
// encoding. This needs to be registered on the Encoders
func FromReflectCustomMsg(cdc codec.Codec) wasmkeeper.CustomEncoder {
	return func(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var sifMsg SifchainMsg
		err := json.Unmarshal(msg, &sifMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		fmt.Printf("@@@ SifchainMsg: %#v\n", sifMsg)

		if sifMsg.Swap != nil {
			fmt.Printf("@@@ swap: %#v\n", sifMsg.Swap)

			// TODO populate swap message with data in ReflectCustomMsg
			contractAddress, err := sdk.AccAddressFromBech32("sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6")
			if err != nil {
				return nil, err
			}
			swapMsg := clptypes.NewMsgSwap(
				contractAddress,
				clptypes.NewAsset("rowan"),
				clptypes.NewAsset("ceth"),
				sdk.NewUint(uint64(sifMsg.Swap.Amount)),
				sdk.NewUint(0),
			)
			return []sdk.Msg{&swapMsg}, nil
		}

		return nil, fmt.Errorf("@@@ Empty SifchainMsg")
	}
}
