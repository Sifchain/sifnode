package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/Sifchain/sifnode/testutil/keeper"
	"github.com/Sifchain/sifnode/testutil/nullify"
	"github.com/Sifchain/sifnode/x/epochs/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestEpochInfoQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.EpochsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createEpochInfos(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryCurrentEpochRequest
		response *types.QueryCurrentEpochResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryCurrentEpochRequest{
				Identifier: msgs[0].Identifier,
			},
			response: &types.QueryCurrentEpochResponse{CurrentEpoch: msgs[0].CurrentEpoch},
		},
		{
			desc: "Second",
			request: &types.QueryCurrentEpochRequest{
				Identifier: msgs[1].Identifier,
			},
			response: &types.QueryCurrentEpochResponse{CurrentEpoch: msgs[1].CurrentEpoch},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryCurrentEpochRequest{
				Identifier: strconv.Itoa(100000),
			},
			err: status.Errorf(codes.NotFound, "epoch info not found: %s", strconv.Itoa(100000)),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.CurrentEpoch(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestEntryQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.EpochsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createEpochInfos(keeper, ctx)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryEpochsInfoRequest {
		return &types.QueryEpochsInfoRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.EpochInfos(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Epochs), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Epochs),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.EpochInfos(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Epochs), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Epochs),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.EpochInfos(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Epochs),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.EpochInfos(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "empty request"))
	})
}
