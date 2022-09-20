package types_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestTypes_NewMTP(t *testing.T) {
	got := types.NewMTP("signer", "collateralAsset", "borrowAsset", types.Position_LONG, sdk.OneDec())

	require.Equal(t, got.Address, "signer")
	require.Equal(t, got.CollateralAsset, "collateralAsset")
	require.Equal(t, got.CustodyAsset, "borrowAsset")
	require.Equal(t, got.Position, types.Position_LONG)
	require.Equal(t, got.Leverage, sdk.OneDec())
}

func TestTypes_MtpValidate(t *testing.T) {
	validateTests := []struct {
		name      string
		mtp       types.MTP
		err       error
		errString error
	}{
		{
			name:      "missing asset",
			mtp:       types.MTP{},
			errString: sdkerrors.Wrap(types.ErrMTPInvalid, "no asset specified"),
		},
		{
			name: "missing address",
			mtp: types.MTP{
				CollateralAsset: "xxx",
			},
			errString: sdkerrors.Wrap(types.ErrMTPInvalid, "no address specified"),
		},
		{
			name: "all filled",
			mtp: types.MTP{
				CollateralAsset: "xxx",
				Address:         "xxx",
				Id:              1,
				Position:        types.Position_LONG,
			},
			errString: nil,
		},
	}
	for _, tt := range validateTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mtp.Validate()

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

func TestTypes_GetSettlementAsset(t *testing.T) {
	got := types.GetSettlementAsset()

	require.Equal(t, got, "rowan")
}
