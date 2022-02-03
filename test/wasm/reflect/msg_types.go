package reflect

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
