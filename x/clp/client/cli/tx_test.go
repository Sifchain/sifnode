package cli_test

/* TODO: convert to integration tests

func SetupViper() {
	viper.Set(flags.FlagKeyringBackend, flags.DefaultKeyringBackend)
	viper.Set(flags.FlagGenerateOnly, true)
	viper.Set(flags.FlagChainID, "sifchainTest")

}

func TestGetCmdCreatePool(t *testing.T) {
	clpcmd := cli.GetCmdCreatePool()
	SetupViper()
	viper.Set(cli.FlagExternalAssetAmount, "100")
	viper.Set(cli.FlagNativeAssetAmount, "100")
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--externalAmount", "100",
		"--nativeAmount", "100",
		"--symbol", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdAddLiquidity(t *testing.T) {

	clpcmd := cli.GetCmdAddLiquidity()
	SetupViper()
	viper.Set(cli.FlagExternalAssetAmount, "5000000000000000000")
	viper.Set(cli.FlagNativeAssetAmount, "5000000000000000000000")
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--externalAmount", "5000000000000000000",
		"--nativeAmount", "5000000000000000000000",
		"--symbol", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdRemoveLiquidity(t *testing.T) {

	clpcmd := cli.GetCmdRemoveLiquidity()
	SetupViper()
	viper.Set(cli.FlagWBasisPoints, "100")
	viper.Set(cli.FlagAsymmetry, "1000")
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--symbol", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err := clpcmd.Execute()
	assert.NoError(t, err)

	viper.Set(cli.FlagWBasisPoints, "%%")
	viper.Set(cli.FlagAsymmetry, "1000")
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--symbol", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err = clpcmd.Execute()
	assert.Error(t, err)

	viper.Set(cli.FlagWBasisPoints, "100")
	viper.Set(cli.FlagAsymmetry, "asdef")
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--symbol", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err = clpcmd.Execute()
	assert.Error(t, err)
}

func TestGetCmdSwap(t *testing.T) {

	clpcmd := cli.GetCmdSwap()
	SetupViper()

	viper.Set(cli.FlagSentAssetSymbol, "eth")
	viper.Set(cli.FlagReceivedAssetSymbol, "dash")
	viper.Set(cli.FlagAmount, "100")
	viper.Set(cli.FlagMinimumReceivingAmount, "90")
	clpcmd.SetArgs([]string{
		"--sentAmount", "100",
		"--minReceivingAmount", "90",
		"--receivedSymbol", "dash",
		"--sentSymbol", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdDecommissionPool(t *testing.T) {

	clpcmd := cli.GetCmdDecommissionPool()
	SetupViper()
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--symbol", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

*/