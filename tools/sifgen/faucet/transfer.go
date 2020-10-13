package faucet

// Transfer funds.
func (f Faucet) Transfer(fromKeyPassword, fromKeyAddress, toKeyAddress, coins string) error {
	_, err := f.CLI.TransferFunds(fromKeyPassword, fromKeyAddress, toKeyAddress, coins)
	if err != nil {
		return err
	}

	return nil
}
