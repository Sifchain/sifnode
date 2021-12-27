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

func TestTypes_MsgOpenLongValidateBasic(t *testing.T) {
	validateBasicTests := []struct {
		name        string
		msgOpenLong types.MsgOpenLong
		err         error
		errString   error
	}{
		{
			name:        "no signer",
			msgOpenLong: types.MsgOpenLong{},
			err:         sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ""),
		},
		{
			name: "collateral asset invalid",
			msgOpenLong: types.MsgOpenLong{
				Signer:          "xxx",
				CollateralAsset: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			err: sdkerrors.Wrap(clptypes.ErrInValidAsset, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		},
		{
			name: "borrow asset invalid",
			msgOpenLong: types.MsgOpenLong{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			err: sdkerrors.Wrap(clptypes.ErrInValidAsset, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		},
		{
			name: "collateral amount is zero",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(0),
				BorrowAsset:      "xxx",
			},
			err: sdkerrors.Wrap(clptypes.ErrInValidAmount, sdk.NewUint(0).String()),
		},
		{
			name: "all valid",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(100),
				BorrowAsset:      "xxx",
			},
			err: nil,
		},
	}
	for _, tt := range validateBasicTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msgOpenLong.ValidateBasic()

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

func TestTypes_MsgOpenLongGetSigners(t *testing.T) {
	getSignersTests := []struct {
		name        string
		msgOpenLong types.MsgOpenLong
		errString   string
	}{
		{
			name:        "no signer",
			msgOpenLong: types.MsgOpenLong{},
			errString:   "empty address string is not allowed",
		},
		{
			name: "wrong address",
			msgOpenLong: types.MsgOpenLong{
				Signer: "xxx",
			},
			errString: "decoding bech32 failed: invalid bech32 string length 3",
		},
		{
			name: "wrong prefix",
			msgOpenLong: types.MsgOpenLong{
				Signer: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			},
			errString: "invalid Bech32 prefix; expected cosmos, got sif",
		},
		{
			name: "returned address",
			msgOpenLong: types.MsgOpenLong{
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
					tt.msgOpenLong.GetSigners()
				})
			} else {
				got := tt.msgOpenLong.GetSigners()
				require.Equal(t, got[0].String(), tt.msgOpenLong.Signer)
			}
		})
	}
}

func TestTypes_MsgCloseLongValidateBasic(t *testing.T) {
	validateBasicTests := []struct {
		name         string
		msgCloseLong types.MsgCloseLong
		err          error
		errString    error
	}{
		{
			name:         "no signer",
			msgCloseLong: types.MsgCloseLong{},
			err:          sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ""),
		},
		{
			name: "collateral asset invalid",
			msgCloseLong: types.MsgCloseLong{
				Signer:          "xxx",
				CollateralAsset: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			err: sdkerrors.Wrap(clptypes.ErrInValidAsset, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		},
		{
			name: "borrow asset invalid",
			msgCloseLong: types.MsgCloseLong{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			err: sdkerrors.Wrap(clptypes.ErrInValidAsset, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		},
		{
			name: "all valid",
			msgCloseLong: types.MsgCloseLong{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
			},
			err: nil,
		},
	}
	for _, tt := range validateBasicTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msgCloseLong.ValidateBasic()

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

func TestTypes_MsgCloseLongGetSigners(t *testing.T) {
	getSignersTests := []struct {
		name         string
		msgCloseLong types.MsgCloseLong
		errString    string
	}{
		{
			name:         "no signer",
			msgCloseLong: types.MsgCloseLong{},
			errString:    "empty address string is not allowed",
		},
		{
			name: "wrong address",
			msgCloseLong: types.MsgCloseLong{
				Signer: "xxx",
			},
			errString: "decoding bech32 failed: invalid bech32 string length 3",
		},
		{
			name: "wrong prefix",
			msgCloseLong: types.MsgCloseLong{
				Signer: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			},
			errString: "invalid Bech32 prefix; expected cosmos, got sif",
		},
		{
			name: "returned address",
			msgCloseLong: types.MsgCloseLong{
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
					tt.msgCloseLong.GetSigners()
				})
			} else {
				got := tt.msgCloseLong.GetSigners()
				require.Equal(t, got[0].String(), tt.msgCloseLong.Signer)
			}
		})
	}
}
