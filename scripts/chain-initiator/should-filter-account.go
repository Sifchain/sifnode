package main

func shouldFilterAccount(account Account, filterAddresses map[string]struct{}) bool {
	if account.BaseAccount != nil {
		if _, exists := filterAddresses[account.BaseAccount.Address]; exists {
			return true
		}
	}
	if account.ModuleAccount != nil {
		if _, exists := filterAddresses[account.ModuleAccount.BaseAccount.Address]; exists {
			return true
		}
	}
	return false
}
