package bank

// Transfer funds.
func (b Bank) Transfer(fromKeyPassword, fromKeyAddress, toKeyAddress string, coins string) error {
	_, err := b.CLI.TransferFunds(fromKeyPassword, fromKeyAddress, toKeyAddress, coins)
	if err != nil {
		return err
	}

	return nil
}
