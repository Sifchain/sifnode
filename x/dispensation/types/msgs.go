package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

var (
	_ sdk.Msg = &MsgAirdrop{}
)

type MsgAirdrop struct {
	Signer      sdk.AccAddress `json:"Signer"`
	AirdropName string         `json:"airdrop_name"`
	Input       []bank.Input   `json:"Input"`
	Output      []bank.Output  `json:"Output"`
}

func NewMsgAirdrop(signer sdk.AccAddress, name string, input []bank.Input, output []bank.Output) MsgAirdrop {
	return MsgAirdrop{Signer: signer, AirdropName: name, Input: input, Output: output}
}

func (m MsgAirdrop) Route() string {
	return RouterKey
}

func (m MsgAirdrop) Type() string {
	return "airdrop"
}

func (m MsgAirdrop) ValidateBasic() error {
	if m.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer.String())
	}
	if m.AirdropName == "" {
		return sdkerrors.Wrap(ErrInvalid, "Name cannot be empty")
	}
	err := bank.ValidateInputsOutputs(m.Input, m.Output)
	if err != nil {
		return err
	}
	return nil
}

func (m MsgAirdrop) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgAirdrop) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}
