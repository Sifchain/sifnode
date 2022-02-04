package reflect

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ReflectHandleMsg is used to encode handle messages
type ReflectHandleMsg struct {
	Reflect       *ReflectPayload    `json:"reflect_msg,omitempty"`
	ReflectSubMsg *ReflectSubPayload `json:"reflect_sub_msg,omitempty"`
	Change        *OwnerPayload      `json:"change_owner,omitempty"`
}

type OwnerPayload struct {
	Owner sdk.Address `json:"owner"`
}

type ReflectPayload struct {
	Msgs []wasmvmtypes.CosmosMsg `json:"msgs"`
}

type ReflectSubPayload struct {
	Msgs []wasmvmtypes.SubMsg `json:"msgs"`
}

type ReflectCustomMsg struct {
	Debug string `json:"debug,omitempty"`
	Raw   []byte `json:"raw,omitempty"`
}
