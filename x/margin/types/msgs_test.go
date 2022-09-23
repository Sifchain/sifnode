package types_test

import (
	"fmt"
	"testing"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestTypes_ValidateAsset(t *testing.T) {
	validateTests := []struct {
		asset string
		valid bool
	}{
		{
			asset: "xxx",
			valid: true,
		},
		{
			asset: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			valid: false,
		},
	}
	for _, tt := range validateTests {
		tt := tt
		t.Run(fmt.Sprintf("asset: %v", tt.asset), func(t *testing.T) {
			got := types.Validate(tt.asset)
			require.Equal(t, got, tt.valid)
		})
	}
}

func TestTypes_MsgOpenValidateBasic(t *testing.T) {
	validateBasicTests := []struct {
		name      string
		msgOpen   types.MsgOpen
		err       error
		errString error
	}{
		{
			name:    "no signer",
			msgOpen: types.MsgOpen{},
			err:     sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ""),
		},
		{
			name: "collateral asset invalid",
			msgOpen: types.MsgOpen{
				Signer:          "xxx",
				CollateralAsset: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Position:        types.Position_LONG,
				Leverage:        sdk.NewDec(1),
			},
			err: sdkerrors.Wrap(clptypes.ErrInValidAsset, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		},
		{
			name: "borrow asset invalid",
			msgOpen: types.MsgOpen{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Position:        types.Position_LONG,
			},
			err: sdkerrors.Wrap(clptypes.ErrInValidAsset, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		},
		{
			name: "collateral amount is zero",
			msgOpen: types.MsgOpen{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(0),
				BorrowAsset:      "xxx",
				Position:         types.Position_LONG,
			},
			err: sdkerrors.Wrap(clptypes.ErrInValidAmount, sdk.NewUint(0).String()),
		},
		{
			name: "position invalid",
			msgOpen: types.MsgOpen{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(100),
				BorrowAsset:      "xxx",
				Position:         types.Position_UNSPECIFIED,
				Leverage:         sdk.NewDec(1),
			},
			err: sdkerrors.Wrap(types.ErrInvalidPosition, types.Position_UNSPECIFIED.String()),
		},
		{
			name: "all valid",
			msgOpen: types.MsgOpen{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(100),
				BorrowAsset:      "xxx",
				Position:         types.Position_LONG,
				Leverage:         sdk.NewDec(1),
			},
			err: nil,
		},
	}
	for _, tt := range validateBasicTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msgOpen.ValidateBasic()

			if tt.errString != nil {
				require.EqualError(t, got, tt.errString.Error())
			} else if tt.err == nil {
				require.NoError(t, got)
			} else {
				require.ErrorIs(t, got, tt.err)
			}
		})
	}
}

func TestTypes_MsgOpenGetSigners(t *testing.T) {
	getSignersTests := []struct {
		name      string
		msgOpen   types.MsgOpen
		errString string
	}{
		{
			name:      "no signer",
			msgOpen:   types.MsgOpen{},
			errString: "empty address string is not allowed",
		},
		{
			name: "wrong address",
			msgOpen: types.MsgOpen{
				Signer: "xxx",
			},
			errString: "decoding bech32 failed: invalid bech32 string length 3",
		},
		{
			name: "wrong prefix",
			msgOpen: types.MsgOpen{
				Signer: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			},
			errString: "invalid Bech32 prefix; expected cosmos, got sif",
		},
		{
			name: "returned address",
			msgOpen: types.MsgOpen{
				Signer: "cosmos1syavy2npfyt9tcncdtsdzf7kny9lh777pahuux",
			},
			errString: "",
		},
	}
	for _, tt := range getSignersTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.errString != "" {
				require.PanicsWithError(t, tt.errString, func() {
					tt.msgOpen.GetSigners()
				})
			} else {
				got := tt.msgOpen.GetSigners()
				require.Equal(t, got[0].String(), tt.msgOpen.Signer)
			}
		})
	}
}

func TestTypes_MsgCloseValidateBasic(t *testing.T) {
	validateBasicTests := []struct {
		name      string
		msgClose  types.MsgClose
		err       error
		errString error
	}{
		{
			name:     "no signer",
			msgClose: types.MsgClose{},
			err:      sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ""),
		},
		{
			name: "id invalid",
			msgClose: types.MsgClose{
				Signer: "xxx",
			},
			err: sdkerrors.Wrap(types.ErrMTPDoesNotExist, "no id specified"),
		},

		{
			name: "all valid",
			msgClose: types.MsgClose{
				Signer: "xxx",
				Id:     1,
			},
			err: nil,
		},
	}
	for _, tt := range validateBasicTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msgClose.ValidateBasic()

			if tt.errString != nil {
				require.EqualError(t, got, tt.errString.Error())
			} else if tt.err == nil {
				require.NoError(t, got)
			} else {
				require.ErrorIs(t, got, tt.err)
			}
		})
	}
}

func TestTypes_MsgCloseGetSigners(t *testing.T) {
	getSignersTests := []struct {
		name      string
		msgClose  types.MsgClose
		errString string
	}{
		{
			name:      "no signer",
			msgClose:  types.MsgClose{},
			errString: "empty address string is not allowed",
		},
		{
			name: "wrong address",
			msgClose: types.MsgClose{
				Signer: "xxx",
			},
			errString: "decoding bech32 failed: invalid bech32 string length 3",
		},
		{
			name: "wrong prefix",
			msgClose: types.MsgClose{
				Signer: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			},
			errString: "invalid Bech32 prefix; expected cosmos, got sif",
		},
		{
			name: "returned address",
			msgClose: types.MsgClose{
				Signer: "cosmos1syavy2npfyt9tcncdtsdzf7kny9lh777pahuux",
			},
			errString: "",
		},
	}
	for _, tt := range getSignersTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.errString != "" {
				require.PanicsWithError(t, tt.errString, func() {
					tt.msgClose.GetSigners()
				})
			} else {
				got := tt.msgClose.GetSigners()
				require.Equal(t, got[0].String(), tt.msgClose.Signer)
			}
		})
	}
}

func TestTypes_MsgForceCloseValidateBasic(t *testing.T) {
	validateBasicTests := []struct {
		name          string
		msgForceClose types.MsgForceClose
		err           error
		errString     error
	}{
		{
			name:          "no signer",
			msgForceClose: types.MsgForceClose{},
			err:           sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ""),
		},
		{
			name: "no mtp address",
			msgForceClose: types.MsgForceClose{
				Signer: "xxx",
			},
			err: sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ""),
		},
		{
			name: "all valid",
			msgForceClose: types.MsgForceClose{
				Signer:     "xxx",
				MtpAddress: "xxx",
				Id:         1,
			},
			err: nil,
		},
	}
	for _, tt := range validateBasicTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msgForceClose.ValidateBasic()

			if tt.errString != nil {
				require.EqualError(t, got, tt.errString.Error())
			} else if tt.err == nil {
				require.NoError(t, got)
			} else {
				require.ErrorIs(t, got, tt.err)
			}
		})
	}
}

func TestTypes_MsgForceCloseGetSigners(t *testing.T) {
	getSignersTests := []struct {
		name          string
		msgForceClose types.MsgForceClose
		errString     string
	}{
		{
			name:          "no signer",
			msgForceClose: types.MsgForceClose{},
			errString:     "empty address string is not allowed",
		},
		{
			name: "wrong address",
			msgForceClose: types.MsgForceClose{
				Signer: "xxx",
			},
			errString: "decoding bech32 failed: invalid bech32 string length 3",
		},
		{
			name: "wrong prefix",
			msgForceClose: types.MsgForceClose{
				Signer: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			},
			errString: "invalid Bech32 prefix; expected cosmos, got sif",
		},
		{
			name: "returned address",
			msgForceClose: types.MsgForceClose{
				Signer: "cosmos1syavy2npfyt9tcncdtsdzf7kny9lh777pahuux",
			},
			errString: "",
		},
	}
	for _, tt := range getSignersTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.errString != "" {
				require.PanicsWithError(t, tt.errString, func() {
					tt.msgForceClose.GetSigners()
				})
			} else {
				got := tt.msgForceClose.GetSigners()
				require.Equal(t, got[0].String(), tt.msgForceClose.Signer)
			}
		})
	}
}
