package main

func filterAccounts(accounts []Account, filterAddresses []string) []Account {
	filterMap := make(map[string]struct{})
	for _, addr := range filterAddresses {
		filterMap[addr] = struct{}{}
	}

	newAccounts := []Account{}
	for _, account := range accounts {
		if shouldFilterAccount(account, filterMap) {
			continue
		}
		newAccounts = append(newAccounts, account)
	}
	return newAccounts
}
