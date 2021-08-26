package ibctransfer_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ibctransfer"
	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func TestOnAcknowledgementMaybeConvert(t *testing.T) {
	type args struct {
		ctx               sdk.Context
		sdkTransferKeeper sctransfertypes.SDKTransferKeeper
		whitelistKeeper   tokenregistrytypes.Keeper
		bankKeeper        bankkeeper.Keeper
		packet            channeltypes.Packet
		acknowledgement   []byte
	}
	tests := []struct {
		name   string
		args   args
		events sdk.Events
		result *sdk.Result
		err    error
	}{
		// TODO: Add test cases.
		{name: "Ack err causes refund - success"},
		{name: "Ack err causes refund - failed to "},
		{name: "Ack err when sender is source - "},
		{name: "Ack err when sender is sink (not source) - "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := ibctransfer.OnAcknowledgementMaybeConvert(tt.args.ctx, tt.args.sdkTransferKeeper, tt.args.whitelistKeeper, tt.args.bankKeeper, tt.args.packet, tt.args.acknowledgement)
			require.ErrorIs(t, err, tt.err)
			// Assert events have recorded what happened.
		})
	}
}
