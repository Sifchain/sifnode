package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagAssetSourceChain         = "sourceChain"
	FlagAssetSymbol              = "symbol"
	FlagAssetTicker              = "ticker"
	FlagSentAssetSourceChain     = "sentSourceChain"
	FlagSentAssetSymbol          = "sentSymbol"
	FlagSentAssetTicker          = "sentTicker"
	FlagReceivedAssetSourceChain = "receivedSourceChain"
	FlagReceivedAssetSymbol      = "receivedSymbol"
	FlagReceivedAssetTicker      = "receivedTicker"
	FlagNativeAssetAmount        = "nativeAmount"
	FlagExternalAssetAmount      = "externalAmount"
	FlagWBasisPoints             = "wBasis"
	FlagAsymmetry                = "asymmetry"
	FlagAmount                   = "sentAmount"
)

// common flagsets to add to various functions
var (
	FsAssetSourceChain         = flag.NewFlagSet("", flag.ContinueOnError)
	FsAssetSymbol              = flag.NewFlagSet("", flag.ContinueOnError)
	FsAssetTicker              = flag.NewFlagSet("", flag.ContinueOnError)
	FsNativeAssetAmount        = flag.NewFlagSet("", flag.ContinueOnError)
	FsExternalAssetAmount      = flag.NewFlagSet("", flag.ContinueOnError)
	FsWBasisPoints             = flag.NewFlagSet("", flag.ContinueOnError)
	FsAsymmetry                = flag.NewFlagSet("", flag.ContinueOnError)
	FsSentAssetSourceChain     = flag.NewFlagSet("", flag.ContinueOnError)
	FsSentAssetSymbol          = flag.NewFlagSet("", flag.ContinueOnError)
	FsSentAssetTicker          = flag.NewFlagSet("", flag.ContinueOnError)
	FsReceivedAssetSourceChain = flag.NewFlagSet("", flag.ContinueOnError)
	FsReceivedAssetSymbol      = flag.NewFlagSet("", flag.ContinueOnError)
	FsReceivedAssetTicker      = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount                   = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {

	FsAssetSourceChain.String(FlagAssetSourceChain, "", "SoureChain for Asset")
	FsAssetSymbol.String(FlagAssetSymbol, "", "Symbol for Asset")
	FsAssetTicker.String(FlagAssetTicker, "", "Ticker for Asset")
	FsNativeAssetAmount.String(FlagNativeAssetAmount, "", "Native Asset Amount")
	FsExternalAssetAmount.String(FlagExternalAssetAmount, "", "External Asset Amount")
	FsWBasisPoints.String(FlagWBasisPoints, "", "WBasis Points ")
	FsAsymmetry.String(FlagAsymmetry, "", "Asymmetry")
	FsSentAssetSourceChain.String(FlagSentAssetSourceChain, "", "SoureChain for Sent Asset")
	FsSentAssetSymbol.String(FlagSentAssetSymbol, "", "Symbol for Sent Asset")
	FsSentAssetTicker.String(FlagSentAssetTicker, "", "Ticker for Sent Asset")
	FsReceivedAssetSourceChain.String(FlagReceivedAssetSourceChain, "", "SoureChain for Received Asset")
	FsReceivedAssetSymbol.String(FlagReceivedAssetSymbol, "", "Symbol for Received Asset")
	FsReceivedAssetTicker.String(FlagReceivedAssetTicker, "", "Ticker for Received Asset")
	FsAmount.String(FlagAmount, "", "Sent amount")

}
