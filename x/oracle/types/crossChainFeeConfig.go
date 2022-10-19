package types

func (m *CrossChainFeeConfig) IsValid() bool {
	if len(m.FeeCurrency) == 0 ||
		m.FeeCurrencyGas.IsNegative() ||
		m.FirstBurnDoublePeggyCost.IsNegative() ||
		m.MinimumBurnCost.IsNegative() ||
		m.MinimumLockCost.IsNegative() {
		return false
	}
	return true
}
