package keeper

import (
	"errors"
	"fmt"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

// TODO: move to x/oracle

// NewQuerier is the module level router for state queries
func NewQuerier(keeper types.OracleKeeper, cdc *codec.Codec, bridgerKeeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryEthProphecy:
			return queryEthProphecy(ctx, cdc, req, keeper)
		case types.QueryGasPrice:
			return queryGasPrice(ctx, bridgerKeeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown ethbridge query endpoint")
		}
	}
}

func queryEthProphecy(
	ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, keeper types.OracleKeeper,
) ([]byte, error) {
	var params types.QueryEthProphecyParams

	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(types.ErrJSONMarshalling, fmt.Sprintf("failed to parse params: %s", err.Error()))
	}

	id := strconv.Itoa(params.EthereumChainID) + strconv.Itoa(params.Nonce) + params.EthereumSender.String()
	prophecy, found := keeper.GetProphecy(ctx, id)
	if !found {
		return nil, sdkerrors.Wrap(oracletypes.ErrProphecyNotFound, id)
	}

	bridgeClaims, err := types.MapOracleClaimsToEthBridgeClaims(
		params.EthereumChainID, params.BridgeContractAddress, params.Nonce, params.Symbol, params.TokenContractAddress,
		params.EthereumSender, prophecy.ValidatorClaims, types.CreateEthClaimFromOracleString)
	if err != nil {
		return nil, err
	}

	response := types.NewQueryEthProphecyResponse(prophecy.ID, prophecy.Status, bridgeClaims)

	return cdc.MarshalJSONIndent(response, "", "  ")
}

func queryGasPrice(
	ctx sdk.Context, keeper Keeper,
) ([]byte, error) {
	gasPrice := keeper.GetEthGasPrice(ctx)
	gasMultiPlier := keeper.GetGasMultiplier(ctx)

	if gasPrice == nil || gasMultiPlier == nil {
		return nil, errors.New("not get the gas price from keeper")
	}

	*gasPrice = gasPrice.Mul(*gasMultiPlier)

	return []byte(gasPrice.String()), nil
}
