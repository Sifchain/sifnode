package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

var (
	_ sdk.Msg = &MsgRemoveLiquidity{}
	_ sdk.Msg = &MsgRemoveLiquidityUnits{}
	_ sdk.Msg = &MsgCreatePool{}
	_ sdk.Msg = &MsgAddLiquidity{}
	_ sdk.Msg = &MsgSwap{}
	_ sdk.Msg = &MsgDecommissionPool{}
	_ sdk.Msg = &MsgUnlockLiquidityRequest{}
	_ sdk.Msg = &MsgUpdateRewardsParamsRequest{}
	_ sdk.Msg = &MsgAddRewardPeriodRequest{}
	_ sdk.Msg = &MsgModifyPmtpRates{}
	_ sdk.Msg = &MsgUpdatePmtpParams{}
	_ sdk.Msg = &MsgUpdateStakingRewardParams{}
	_ sdk.Msg = &MsgSetSymmetryThreshold{}
	_ sdk.Msg = &MsgCancelUnlock{}
	_ sdk.Msg = &MsgUpdateLiquidityProtectionParams{}
	_ sdk.Msg = &MsgModifyLiquidityProtectionRates{}
	_ sdk.Msg = &MsgAddProviderDistributionPeriodRequest{}
	_ sdk.Msg = &MsgUpdateSwapFeeParamsRequest{}

	_ legacytx.LegacyMsg = &MsgRemoveLiquidity{}
	_ legacytx.LegacyMsg = &MsgRemoveLiquidityUnits{}
	_ legacytx.LegacyMsg = &MsgCreatePool{}
	_ legacytx.LegacyMsg = &MsgAddLiquidity{}
	_ legacytx.LegacyMsg = &MsgSwap{}
	_ legacytx.LegacyMsg = &MsgDecommissionPool{}
	_ legacytx.LegacyMsg = &MsgUnlockLiquidityRequest{}
	_ legacytx.LegacyMsg = &MsgUpdateRewardsParamsRequest{}
	_ legacytx.LegacyMsg = &MsgAddRewardPeriodRequest{}
	_ legacytx.LegacyMsg = &MsgModifyPmtpRates{}
	_ legacytx.LegacyMsg = &MsgUpdatePmtpParams{}
	_ legacytx.LegacyMsg = &MsgUpdateStakingRewardParams{}
	_ legacytx.LegacyMsg = &MsgSetSymmetryThreshold{}
	_ legacytx.LegacyMsg = &MsgCancelUnlock{}
	_ legacytx.LegacyMsg = &MsgAddProviderDistributionPeriodRequest{}
	_ legacytx.LegacyMsg = &MsgUpdateSwapFeeParamsRequest{}
)

func (m MsgCancelUnlock) Route() string {
	return RouterKey
}

func (m MsgCancelUnlock) Type() string {
	return "cancel_unlock"
}

func (m MsgCancelUnlock) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ExternalAsset.Symbol)
	}
	return nil
}

func (m MsgCancelUnlock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCancelUnlock) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgUpdateStakingRewardParams) Route() string {
	return RouterKey
}

func (m MsgUpdateStakingRewardParams) Type() string {
	return "update_staking_reward"
}

func (m MsgUpdateStakingRewardParams) ValidateBasic() error {
	return m.Params.Validate()
}

func (m MsgUpdateStakingRewardParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgUpdateStakingRewardParams) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgAddRewardPeriodRequest) Route() string {
	return RouterKey
}

func (m MsgAddRewardPeriodRequest) Type() string {
	return "add_reward_period"
}

func (m MsgAddRewardPeriodRequest) ValidateBasic() error {
	for _, period := range m.RewardPeriods {
		if period.RewardPeriodId == "" {
			return fmt.Errorf("reward period id must be non-empty: %d", period.RewardPeriodStartBlock)
		}
		if period.RewardPeriodEndBlock < period.RewardPeriodStartBlock {
			return fmt.Errorf("reward period start block must be before end block: %d %d", period.RewardPeriodStartBlock, period.RewardPeriodEndBlock)
		}
		for _, multiplier := range period.RewardPeriodPoolMultipliers {
			if multiplier.Multiplier.LT(sdk.ZeroDec()) {
				return fmt.Errorf("pool multiplier should be less than 0 | pool : %s , multiplier : %s", multiplier.PoolMultiplierAsset, multiplier.Multiplier.String())
			}
			if multiplier.Multiplier.GT(sdk.MustNewDecFromStr("10.00")) {
				return fmt.Errorf("pool multiplier should be greater than 10 | pool : %s , multiplier : %s", multiplier.PoolMultiplierAsset, multiplier.Multiplier.String())
			}
		}
		if period.RewardPeriodDefaultMultiplier.LT(sdk.ZeroDec()) {
			return fmt.Errorf("default should be less than 0 |multiplier : %s", period.RewardPeriodDefaultMultiplier.String())
		}
		if period.RewardPeriodDefaultMultiplier.GT(sdk.MustNewDecFromStr("10.00")) {
			return fmt.Errorf("default multiplier should be greater than 10 | multiplier : %s", period.RewardPeriodDefaultMultiplier.String())
		}
	}
	return nil
}

func (m MsgAddRewardPeriodRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgAddRewardPeriodRequest) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgUpdateRewardsParamsRequest) Route() string {
	return RouterKey
}

func (m MsgUpdateRewardsParamsRequest) Type() string {
	return "update_reward_params"
}

func (m MsgUpdateRewardsParamsRequest) ValidateBasic() error {
	return nil
}

func (m MsgUpdateRewardsParamsRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgUpdateRewardsParamsRequest) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdatePmtpParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return err
	}
	if m.PmtpPeriodEpochLength <= 0 {
		return fmt.Errorf("pmtp epoch length must be greated than zero: %d", m.PmtpPeriodEpochLength)
	}
	if m.PmtpPeriodStartBlock < 0 {
		return fmt.Errorf("pmtp start block cannot be negative: %d", m.PmtpPeriodStartBlock)
	}
	// End block must be at-least 1
	if m.PmtpPeriodEndBlock <= 0 {
		return fmt.Errorf("pmtp end block cannot be negative: %d", m.PmtpPeriodStartBlock)
	}
	if m.PmtpPeriodEndBlock < m.PmtpPeriodStartBlock {
		return fmt.Errorf(
			"end block (%d) must be after begin block (%d)",
			m.PmtpPeriodEndBlock, m.PmtpPeriodStartBlock,
		)
	}

	if (m.PmtpPeriodEndBlock-m.PmtpPeriodStartBlock+1)%m.PmtpPeriodEpochLength != 0 {
		return fmt.Errorf("all epochs must have equal number of blocks : %d", m.PmtpPeriodEpochLength)
	}

	return nil
}

func (m *MsgUpdatePmtpParams) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgUpdatePmtpParams) Route() string {
	return RouterKey
}

func (m MsgUpdatePmtpParams) Type() string {
	return "update_pmtp_params"
}

func (m MsgUpdatePmtpParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m *MsgModifyPmtpRates) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return err
	}
	return nil
}

func (m *MsgModifyPmtpRates) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgModifyPmtpRates) Route() string {
	return RouterKey
}

func (m MsgModifyPmtpRates) Type() string {
	return "modify_pmtp_rates"
}

func (m MsgModifyPmtpRates) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func NewMsgDecommissionPool(signer sdk.AccAddress, symbol string) MsgDecommissionPool {
	return MsgDecommissionPool{Signer: signer.String(), Symbol: symbol}
}

func (m MsgDecommissionPool) Route() string {
	return RouterKey
}

func (m MsgDecommissionPool) Type() string {
	return "decommission_pool"
}

func (m MsgDecommissionPool) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !VerifyRange(len(strings.TrimSpace(m.Symbol)), 0, MaxSymbolLength) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, m.Symbol)
	}
	return nil
}

func (m MsgDecommissionPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgDecommissionPool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgSwap(signer sdk.AccAddress, sentAsset Asset, receivedAsset Asset, sentAmount sdk.Uint, minReceivingAmount sdk.Uint) MsgSwap {
	return MsgSwap{Signer: signer.String(), SentAsset: &sentAsset, ReceivedAsset: &receivedAsset, SentAmount: sentAmount, MinReceivingAmount: minReceivingAmount}
}

func (m MsgSwap) Route() string {
	return RouterKey
}

func (m MsgSwap) Type() string {
	return "swap"
}

func (m MsgSwap) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.SentAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.SentAsset.Symbol)
	}
	if !m.ReceivedAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ReceivedAsset.Symbol)
	}
	if m.SentAsset.Equals(*m.ReceivedAsset) {
		return sdkerrors.Wrap(ErrInValidAsset, "Sent And Received asset cannot be the same")
	}
	if m.SentAmount.IsZero() {
		return sdkerrors.Wrap(ErrInValidAmount, m.SentAmount.String())
	}
	return nil
}

func (m MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgSwap) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgRemoveLiquidity(signer sdk.AccAddress, externalAsset Asset, wBasisPoints sdk.Int, asymmetry sdk.Int) MsgRemoveLiquidity {
	return MsgRemoveLiquidity{Signer: signer.String(), ExternalAsset: &externalAsset, WBasisPoints: wBasisPoints, Asymmetry: asymmetry}
}

func (m MsgRemoveLiquidity) Route() string {
	return RouterKey
}

func (m MsgRemoveLiquidity) Type() string {
	return "remove_liquidity"
}

func (m MsgRemoveLiquidity) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ExternalAsset.Symbol)
	}
	if !(m.WBasisPoints.IsPositive()) || m.WBasisPoints.GT(sdk.NewInt(MaxWbasis)) {
		return sdkerrors.Wrap(ErrInvalidWBasis, m.WBasisPoints.String())
	}
	if m.Asymmetry.GT(sdk.NewInt(10000)) || m.Asymmetry.LT(sdk.NewInt(-10000)) {
		return sdkerrors.Wrap(ErrInvalidAsymmetry, m.Asymmetry.String())
	}
	return nil
}

func (m MsgRemoveLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgRemoveLiquidityUnits(signer sdk.AccAddress, externalAsset Asset, withdrawUnits sdk.Uint) MsgRemoveLiquidityUnits {
	return MsgRemoveLiquidityUnits{Signer: signer.String(), ExternalAsset: &externalAsset, WithdrawUnits: withdrawUnits}
}

func (m MsgRemoveLiquidityUnits) Route() string {
	return RouterKey
}

func (m MsgRemoveLiquidityUnits) Type() string {
	return "remove_liquidity"
}

func (m MsgRemoveLiquidityUnits) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ExternalAsset.Symbol)
	}
	if !m.WithdrawUnits.GT(sdk.ZeroUint()) {
		return sdkerrors.Wrap(ErrInValidAmount, fmt.Sprintf("Units must be greater than 0 : %s", m.WithdrawUnits.String()))
	}
	return nil
}

func (m MsgRemoveLiquidityUnits) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgRemoveLiquidityUnits) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgAddLiquidity(signer sdk.AccAddress, externalAsset Asset, nativeAssetAmount sdk.Uint, externalAssetAmount sdk.Uint) MsgAddLiquidity {
	return MsgAddLiquidity{Signer: signer.String(), ExternalAsset: &externalAsset, NativeAssetAmount: nativeAssetAmount, ExternalAssetAmount: externalAssetAmount}
}

func (m MsgAddLiquidity) Route() string {
	return RouterKey
}

func (m MsgAddLiquidity) Type() string {
	return "add_liquidity"
}

func (m MsgAddLiquidity) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ExternalAsset.Symbol)
	}
	if m.ExternalAsset.Equals(GetSettlementAsset()) {
		return sdkerrors.Wrap(ErrInValidAsset, "External asset cannot be rowan")
	}
	if (!m.NativeAssetAmount.GT(sdk.ZeroUint())) && (!m.ExternalAssetAmount.GT(sdk.ZeroUint())) {
		return sdkerrors.Wrap(ErrInValidAmount, fmt.Sprintf("Both asset ammounts cannot be 0 %s / %s", m.NativeAssetAmount.String(), m.ExternalAssetAmount.String()))
	}

	return nil
}

func (m MsgAddLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgCreatePool(signer sdk.AccAddress, externalAsset Asset, nativeAssetAmount sdk.Uint, externalAssetAmount sdk.Uint) MsgCreatePool {
	return MsgCreatePool{Signer: signer.String(), ExternalAsset: &externalAsset, NativeAssetAmount: nativeAssetAmount, ExternalAssetAmount: externalAssetAmount}
}

func (m MsgCreatePool) Route() string {
	return RouterKey
}

func (m MsgCreatePool) Type() string {
	return "create_pool"
}

func (m MsgCreatePool) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ExternalAsset.Symbol)
	}
	if m.ExternalAsset.Equals(GetSettlementAsset()) {
		return sdkerrors.Wrap(ErrInValidAsset, "External Asset cannot be rowan")
	}
	if !(m.NativeAssetAmount.GT(sdk.ZeroUint())) {
		return sdkerrors.Wrap(ErrInValidAmount, m.NativeAssetAmount.String())
	}
	if !(m.ExternalAssetAmount.GT(sdk.ZeroUint())) {
		return sdkerrors.Wrap(ErrInValidAmount, m.NativeAssetAmount.String())
	}
	return nil
}

func (m MsgCreatePool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCreatePool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgUnlockLiquidityRequest) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ExternalAsset.Symbol)
	}
	return nil
}

func (m MsgUnlockLiquidityRequest) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgUnlockLiquidityRequest) Route() string {
	return RouterKey
}

func (m MsgUnlockLiquidityRequest) Type() string {
	return "unlock_liquidity"
}

func (m MsgUnlockLiquidityRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgSetSymmetryThreshold) Route() string {
	return RouterKey
}

func (m MsgSetSymmetryThreshold) Type() string {
	return "set_symmetry_threshold"
}

func (m MsgSetSymmetryThreshold) ValidateBasic() error {
	if m.Signer == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	return nil
}

func (m MsgSetSymmetryThreshold) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgSetSymmetryThreshold) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateLiquidityProtectionParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return err
	}
	if m.EpochLength <= 0 {
		return fmt.Errorf("liquidity protection epoch length must be greated than zero: %d", m.EpochLength)
	}
	return nil
}

func (m *MsgUpdateLiquidityProtectionParams) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgUpdateLiquidityProtectionParams) Route() string {
	return RouterKey
}

func (m MsgUpdateLiquidityProtectionParams) Type() string {
	return "update_liquidity_protection_params"
}

func (m MsgUpdateLiquidityProtectionParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m *MsgModifyLiquidityProtectionRates) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return err
	}
	return nil
}

func (m *MsgModifyLiquidityProtectionRates) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgModifyLiquidityProtectionRates) Route() string {
	return RouterKey
}

func (m MsgModifyLiquidityProtectionRates) Type() string {
	return "modify_liquidity_protection_rates"
}

func (m MsgModifyLiquidityProtectionRates) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgAddProviderDistributionPeriodRequest) Route() string {
	return RouterKey
}

func (m MsgAddProviderDistributionPeriodRequest) Type() string {
	return "add_provider_distribution_period"
}

func (m MsgAddProviderDistributionPeriodRequest) ValidateBasic() error {
	if m.Signer == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	for _, period := range m.DistributionPeriods {
		if period.DistributionPeriodStartBlock > period.DistributionPeriodEndBlock {
			return fmt.Errorf("provider distribution period start block must be < end block: %d %d", period.DistributionPeriodStartBlock, period.DistributionPeriodEndBlock)
		}

		if period.DistributionPeriodBlockRate.LT(sdk.NewDec(0)) ||
			period.DistributionPeriodBlockRate.GT(sdk.NewDec(1)) {
			return fmt.Errorf("provider distribution period block rate must be >= 0 and <= 1 but is: %s", period.DistributionPeriodBlockRate.String())
		}

		if period.DistributionPeriodMod == 0 {
			return fmt.Errorf("provider distribution period modulo must be > 0")
		}
	}

	return nil
}

func (m MsgAddProviderDistributionPeriodRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgAddProviderDistributionPeriodRequest) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m MsgUpdateSwapFeeParamsRequest) Route() string {
	return RouterKey
}

func (m MsgUpdateSwapFeeParamsRequest) Type() string {
	return "update_swap_fee_rate"
}

func (m MsgUpdateSwapFeeParamsRequest) ValidateBasic() error {
	if m.Signer == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	if m.DefaultSwapFeeRate.LT(sdk.ZeroDec()) {
		return fmt.Errorf("swap rate fee must be greater than or equal to zero")
	}

	if m.DefaultSwapFeeRate.GT(sdk.OneDec()) {
		return fmt.Errorf("swap rate fee must be less than or equal to one")
	}

	for _, p := range m.TokenParams {
		if p.SwapFeeRate.LT(sdk.ZeroDec()) {
			return fmt.Errorf("swap rate fee must be greater than or equal to zero")
		}

		if p.SwapFeeRate.GT(sdk.OneDec()) {
			return fmt.Errorf("swap rate fee must be less than or equal to one")
		}
	}

	return nil
}

func (m MsgUpdateSwapFeeParamsRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgUpdateSwapFeeParamsRequest) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
