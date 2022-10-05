package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewLegacyQuerier is the module level router for state queries
func NewLegacyQuerier(keeper Keeper) sdk.Querier {
	querier := NewQuerier(keeper)

	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryAllDistributions:
			return queryAllDistributions(ctx, querier)
		case types.QueryRecordsByDistrName:
			return queryDistributionRecordsForName(ctx, req, querier)
		case types.QueryRecordsByRecipient:
			return queryDistributionRecordsForRecipient(ctx, req, querier)
		case types.QueryClaimsByType:
			return queryClaimsByType(ctx, req, querier)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown dispensation query endpoint")
		}
	}
}

func queryClaimsByType(ctx sdk.Context, req abci.RequestQuery, querier Querier) ([]byte, error) {
	var params types.QueryClaimsByTypeRequest
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	queryResponse, _ := querier.ClaimsByType(sdk.WrapSDKContext(ctx), &params)
	res, err := types.ModuleCdc.MarshalJSON(queryResponse)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryDistributionRecordsForName(ctx sdk.Context, req abci.RequestQuery, querier Querier) ([]byte, error) {
	var params types.QueryRecordsByDistributionNameRequest

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// RecordsByDistributionName cannot return error
	res, _ := querier.RecordsByDistributionName(sdk.WrapSDKContext(ctx), &params)
	bz, err := types.ModuleCdc.MarshalJSON(res.DistributionRecords)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryDistributionRecordsForRecipient(ctx sdk.Context, req abci.RequestQuery, querier Querier) ([]byte, error) {
	var params types.QueryRecordsByRecipientAddrRequest

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// RecordsByRecipient cannot return error
	res, _ := querier.RecordsByRecipient(sdk.WrapSDKContext(ctx), &params)

	bz, err := types.ModuleCdc.MarshalJSON(res.DistributionRecords)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryAllDistributions(ctx sdk.Context, querier Querier) ([]byte, error) {
	res, err := querier.AllDistributions(sdk.WrapSDKContext(ctx), &types.QueryAllDistributionsRequest{})
	if err != nil {
		return nil, err
	}

	return types.ModuleCdc.MarshalJSON(res)
}
