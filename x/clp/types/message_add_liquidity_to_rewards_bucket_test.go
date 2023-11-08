package types

import (
	"testing"

	"github.com/Sifchain/sifnode/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAddLiquidityToRewardsBucketRequest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddLiquidityToRewardsBucketRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddLiquidityToRewardsBucketRequest{
				Signer: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgAddLiquidityToRewardsBucketRequest{
				Signer: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
