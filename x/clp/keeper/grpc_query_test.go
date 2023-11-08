package keeper_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	keepertest "github.com/Sifchain/sifnode/testutil/keeper"
	"github.com/Sifchain/sifnode/testutil/nullify"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQuerier_GetPool(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetPool(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetPools(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetPools(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetPools_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.PoolsReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	_, err := querier.GetPools(ctx, req)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProvider(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetLiquidityProvider(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviderData(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetLiquidityProviderData(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetAssetList(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetAssetList(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetAssetList_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.AssetListReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	_, err := querier.GetAssetList(ctx, req)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviderList(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetLiquidityProviderList(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviderList_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.LiquidityProviderListReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	_, err := querier.GetLiquidityProviderList(ctx, req)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviders(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetLiquidityProviders(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviders_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.LiquidityProvidersReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	_, err := querier.GetLiquidityProviders(ctx, req)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetPoolShareEstimate(t *testing.T) {
	testcases := []struct {
		name                        string
		createBalance               bool
		createPool                  bool
		poolAsset                   string
		address                     string
		poolNativeAssetBalance      sdk.Uint
		poolExternalAssetBalance    sdk.Uint
		poolNativeLiabilities       sdk.Uint
		poolExternalLiabilities     sdk.Uint
		poolUnits                   sdk.Uint
		poolAssetPermissions        []tokenregistrytypes.Permission
		nativeAssetPermissions      []tokenregistrytypes.Permission
		RequestNativeAssetAmount    sdk.Uint
		RequestExternalAssetAmount  sdk.Uint
		expectedExternalAssetAmount sdk.Uint
		expectedNativeAssetAmount   sdk.Uint
		expectedPercentage          sdk.Dec
		expectedSwapStatus          types.SwapStatus
		expectedSwapFee             sdk.Uint
		expectedSwapFeeRate         sdk.Dec
		expectedSwapAmount          sdk.Uint
		expectedSwapResult          sdk.Uint
		err                         error
		errString                   error
	}{
		{
			name:                        "symmetric",
			createBalance:               true,
			createPool:                  true,
			poolAsset:                   "eth",
			address:                     "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			poolNativeAssetBalance:      sdk.NewUint(1000),
			poolExternalAssetBalance:    sdk.NewUint(1000),
			poolUnits:                   sdk.NewUint(1000),
			poolAssetPermissions:        []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			RequestNativeAssetAmount:    sdk.NewUint(200),
			RequestExternalAssetAmount:  sdk.NewUint(200),
			expectedExternalAssetAmount: sdk.NewUint(200),
			expectedNativeAssetAmount:   sdk.NewUint(200),
			expectedPercentage:          sdk.MustNewDecFromStr("0.166666666666666667"),
			expectedSwapStatus:          types.SwapStatus_NO_SWAP,
			expectedSwapFee:             sdk.ZeroUint(),
			expectedSwapFeeRate:         sdk.ZeroDec(),
			expectedSwapAmount:          sdk.ZeroUint(),
			expectedSwapResult:          sdk.ZeroUint(),
		},
		{
			name:                        "asymmetric",
			createBalance:               true,
			createPool:                  true,
			poolAsset:                   "eth",
			address:                     "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			poolNativeAssetBalance:      sdk.NewUint(1000),
			poolExternalAssetBalance:    sdk.NewUint(1000),
			poolUnits:                   sdk.NewUint(1000),
			poolAssetPermissions:        []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			RequestNativeAssetAmount:    sdk.NewUint(200),
			RequestExternalAssetAmount:  sdk.ZeroUint(),
			expectedExternalAssetAmount: sdk.NewUint(115),
			expectedNativeAssetAmount:   sdk.NewUint(138),
			expectedPercentage:          sdk.MustNewDecFromStr("0.115826702033598585"),
			expectedSwapStatus:          types.SwapStatus_SELL_NATIVE,
			expectedSwapFee:             sdk.ZeroUint(),
			expectedSwapFeeRate:         sdk.MustNewDecFromStr("0.003"),
			expectedSwapAmount:          sdk.NewUint(61),
			expectedSwapResult:          sdk.NewUint(114),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				trGs := &tokenregistrytypes.GenesisState{
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: tc.poolAsset, BaseDenom: tc.poolAsset, Decimals: 18, Permissions: tc.poolAssetPermissions},
							{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: tc.nativeAssetPermissions},
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(trGs)
				genesisState["tokenregistry"] = bz

				if tc.createPool {
					pools := []*types.Pool{
						{
							ExternalAsset:        &types.Asset{Symbol: tc.poolAsset},
							NativeAssetBalance:   tc.poolNativeAssetBalance,
							ExternalAssetBalance: tc.poolExternalAssetBalance,
							PoolUnits:            tc.poolUnits,
							NativeLiabilities:    tc.poolNativeLiabilities,
							ExternalLiabilities:  tc.poolExternalLiabilities,
						},
					}
					clpGs := types.DefaultGenesisState()
					clpGs.Params = types.Params{
						MinCreatePoolThreshold: 100,
					}
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			querier := clpkeeper.Querier{app.ClpKeeper}

			req := &types.PoolShareEstimateReq{
				ExternalAsset:       &types.Asset{Symbol: tc.poolAsset},
				NativeAssetAmount:   tc.RequestNativeAssetAmount,
				ExternalAssetAmount: tc.RequestExternalAssetAmount,
			}
			res, err := querier.GetPoolShareEstimate(sdk.WrapSDKContext(ctx), req)

			require.NoError(t, err)

			require.Equal(t, tc.expectedExternalAssetAmount.String(), res.ExternalAssetAmount.String())
			require.Equal(t, tc.expectedNativeAssetAmount.String(), res.NativeAssetAmount.String())
			require.Equal(t, tc.expectedPercentage.String(), res.Percentage.String())
			require.Equal(t, tc.expectedSwapStatus, res.SwapInfo.Status)
			require.Equal(t, tc.expectedSwapFee.String(), res.SwapInfo.Fee.String())
			require.Equal(t, tc.expectedSwapFeeRate.String(), res.SwapInfo.FeeRate.String())
			require.Equal(t, tc.expectedSwapAmount.String(), res.SwapInfo.Amount.String())
			require.Equal(t, tc.expectedSwapResult.String(), res.SwapInfo.Result.String())
		})
	}
}

func TestRewardsBucketQuerySingle(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	querier := clpkeeper.Querier{Keeper: *keeper}
	msgs := createNRewardsBucket(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.RewardsBucketReq
		response *types.RewardsBucketRes
		err      error
	}{
		{
			desc: "First",
			request: &types.RewardsBucketReq{
				Denom: msgs[0].Denom,
			},
			response: &types.RewardsBucketRes{RewardsBucket: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.RewardsBucketReq{
				Denom: msgs[1].Denom,
			},
			response: &types.RewardsBucketRes{RewardsBucket: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.RewardsBucketReq{
				Denom: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := querier.GetRewardsBucket(wctx, tc.request)
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

func TestRewardsBucketQueryPaginated(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	querier := clpkeeper.Querier{Keeper: *keeper}
	msgs := createNRewardsBucket(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.AllRewardsBucketReq {
		return &types.AllRewardsBucketReq{
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
			resp, err := querier.GetRewardsBucketAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RewardsBucket), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.RewardsBucket),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := querier.GetRewardsBucketAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RewardsBucket), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.RewardsBucket),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := querier.GetRewardsBucketAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.RewardsBucket),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := querier.GetRewardsBucketAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
