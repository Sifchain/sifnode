package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
)

var (
	_ sdk.Msg = &MsgDistribution{}
)

// Basic message type to create a new distribution
// TODO modify this struct to keep adding more fields to identify different types of distributions
type MsgDistribution struct {
	Distributor      multisig.PubKeyMultisigThreshold `json:"distributor"`
	DistributionName string                           `json:"distribution_name"`
	DistributionType DistributionType                 `json:"distribution_type"`
	Input            []bank.Input                     `json:"Input"`
	Output           []bank.Output                    `json:"Output"`
}

func NewMsgDistribution(signer multisig.PubKeyMultisigThreshold, DistributionName string, DistributionType DistributionType, input []bank.Input, output []bank.Output) MsgDistribution {
	return MsgDistribution{Distributor: signer, DistributionName: DistributionName, DistributionType: DistributionType, Input: input, Output: output}
}

func (m MsgDistribution) Route() string {
	return RouterKey
}

//TODO Replace with constant defined in keys.go with value CreateDispensation
func (m MsgDistribution) Type() string {
	return "airdrop"
}

func (m MsgDistribution) ValidateBasic() error {
	if m.DistributionName == "" {
		return sdkerrors.Wrap(ErrInvalid, "Name cannot be empty")
	}
	err := bank.ValidateInputsOutputs(m.Input, m.Output)
	if err != nil {
		return err
	}
	err = VerifyInputList(m.Input, m.Distributor.PubKeys)
	if err != nil {
		return err
	}
	return nil
}

func (m MsgDistribution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgDistribution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.Distributor.Address())}
}

// Create a user claim
type MsgCreateClaim struct {
	Signer        sdk.AccAddress   `json:"signer"`
	UserClaimType DistributionType `json:"user_claim_type"`
}

func NewMsgCreateClaim(Signer sdk.AccAddress, userClaimType DistributionType) MsgCreateClaim {
	return MsgCreateClaim{Signer: Signer, UserClaimType: userClaimType}
}

func (m MsgCreateClaim) Route() string {
	return RouterKey
}

func (m MsgCreateClaim) Type() string {
	return "createClaim"
}

func (m MsgCreateClaim) ValidateBasic() error {
	if m.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer.String())
	}
	return nil
}

func (m MsgCreateClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCreateClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

//more comments required
//Possible refactor : https://github.com/deckarep/golang-set
func VerifyInputList(inputList []bank.Input, pubKeys []crypto.PubKey) error {

	addressCount := len(pubKeys)
	for _, i := range inputList {
		addressFound := false
		for _, signPubKeys := range pubKeys {
			if bytes.Equal(signPubKeys.Address().Bytes(), i.Address.Bytes()) {
				addressFound = true
				continue
			}
		}
		if !addressFound {
			return errors.Wrap(ErrKeyInvalid, i.Address.String())
		}
		addressCount = addressCount - 1
	}
	if addressCount != 0 {
		return errors.Wrap(ErrKeyInvalid, "Input list and MultiSig Key have a different address count")
	}
	return nil
}
