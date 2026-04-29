

package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"



type ClpProposalHandler struct {
	keeper Keeper
}

func NewClpProposalHandler(keeper Keeper) ClpProposalHandler {
	return ClpProposalHandler{
		keeper: keeper,
	}
}

func (h ClpProposalHandler) HandleProposal(ctx sdk.Context, proposal govtypes.Proposal) error {
	switch proposal := proposal.(type) {
	case *types.UpdatePmtpParamsProposal:
		return HandleUpdatePmtpParamsProposal(ctx, h.keeper, proposal)
	default:
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized clp proposal type: %T", proposal)
	}
}


)

// HandleUpdatePmtpParamsProposal is a handler for updating pmtp params
func HandleUpdatePmtpParamsProposal(ctx sdk.Context, k Keeper, proposal *types.UpdatePmtpParamsProposal) error {
	proposer, err := sdk.AccAddressFromBech32(proposal.RunAs)
	if err != nil {
		return err
	}
	if !k.adminKeeper.IsAdminAccount(ctx, types.AdminType_CLP, proposer) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "proposer is not a recognized admin account")
	}

	pmtpParams := proposal.PmtpParams
	if pmtpParams == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pmtp params cannot be nil")
	}

	k.SetPmtpParams(ctx, pmtpParams)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdatePmtpParams,
			sdk.NewAttribute(types.AttributeKeyPmtpPeriodStartBlock, fmt.Sprintf("%d", pmtpParams.PmtpPeriodStartBlock)),
			sdk.NewAttribute(types.AttributeKeyPmtpPeriodEndBlock, fmt.Sprintf("%d", pmtpParams.PmtpPeriodEndBlock)),
			sdk.NewAttribute(types.AttributeKeyPmtpBlockRate, pmtpParams.PmtpBlockRate.String()),
			sdk.NewAttribute(types.AttributeKeyPmtpCurrentRunningRate, pmtpParams.PmtpCurrentRunningRate.String()),
			sdk.NewAttribute(types.AttributeKeyPmtpInterimAsset, pmtpParams.PmtpInterimAsset.String()),
			sdk.NewAttribute(types.AttributeKeyPmtpCurrentDepth, pmtpParams.PmtpCurrentDepth.String()),
			sdk.NewAttribute(types.AttributeKeyPmtpFarmingFeeRate, pmtpParams.PmtpFarmingFeeRate.String()),
		),
	)

	return nil
}

