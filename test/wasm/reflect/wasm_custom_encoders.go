package reflect

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// reflectEncoders needs to be registered in to handle custom message callbacks
func ReflectEncoders(cdc codec.Codec) *wasmkeeper.MessageEncoders {
	return &wasmkeeper.MessageEncoders{
		Custom: FromReflectRawMsg(cdc),
	}
}

// FromReflectRawMsg decodes msg.Data to an sdk.Msg using proto Any and json
// encoding. This needs to be registered on the Encoders
func FromReflectRawMsg(cdc codec.Codec) wasmkeeper.CustomEncoder {
	return func(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var custom ReflectCustomMsg
		err := json.Unmarshal(msg, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		if custom.Raw != nil {
			var any codectypes.Any
			if err := cdc.UnmarshalJSON(custom.Raw, &any); err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			var msg sdk.Msg
			if err := cdc.UnpackAny(&any, &msg); err != nil {
				return nil, err
			}
			return []sdk.Msg{msg}, nil
		}
		if custom.Debug != "" {
			return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Custom Debug: %s", custom.Debug)
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom message variant")
	}
}

// ToReflectRawMsg encodes an sdk msg using any type with json encoding.
// Then wraps it as an opaque message
func ToReflectRawMsg(cdc codec.Codec, msg sdk.Msg) (wasmvmtypes.CosmosMsg, error) {
	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return wasmvmtypes.CosmosMsg{}, err
	}
	rawBz, err := cdc.MarshalJSON(any)
	if err != nil {
		return wasmvmtypes.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	customMsg, err := json.Marshal(ReflectCustomMsg{
		Raw: rawBz,
	})
	res := wasmvmtypes.CosmosMsg{
		Custom: customMsg,
	}
	return res, nil
}
