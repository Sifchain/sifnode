package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryAllDistributions:
			return queryAllDistributions(ctx, keeper)
		case types.QueryRecordsByDistrName:
			return queryDistributionRecordsForName(ctx, req, keeper)
		case types.QueryRecordsByRecipient:
			return queryDistributionRecordsForRecipient(ctx, req, keeper)
		case types.QueryClaimsByType:
			return queryClaimsByType(ctx, req, keeper)

		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown dispensation query endpoint")
		}
	}
}

func queryDistributionRecordsForName(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryRecordsByDistributionName

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	records := new(types.DistributionRecords)
	switch params.Status {
	case types.Pending:
		*records = keeper.GetRecordsForNamePending(ctx, params.DistributionName)
	case types.Completed:
		*records = keeper.GetRecordsForNameCompleted(ctx, params.DistributionName)
	default:
		*records = keeper.GetRecordsForNameAll(ctx, params.DistributionName)
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, records)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryDistributionRecordsForRecipient(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryRecordsByRecipientAddr

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	records := keeper.GetRecordsForRecipient(ctx, params.Address)
	res, err := codec.MarshalJSONIndent(keeper.cdc, records)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryAllDistributions(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	list := keeper.GetDistributions(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, list)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryClaimsByType(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryUserClaims
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	list := keeper.GetClaimsByType(ctx, params.UserClaimType)
	res, err := codec.MarshalJSONIndent(keeper.cdc, list)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
