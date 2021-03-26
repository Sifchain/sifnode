package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagMultiSigKey = "address"
	FlagInputList   = "input"
	FlagOutputList  = "output"
)

// common flagsets to add to various functions
var (
	FsMultiSigKey = flag.NewFlagSet("", flag.ContinueOnError)
	FsInputList   = flag.NewFlagSet("", flag.ContinueOnError)
	FsOutputListt = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {

	FsMultiSigKey.String(FlagMultiSigKey, "", "Multisig Key For transfer")
	FsInputList.String(FlagInputList, "", "Input List")
	FsOutputListt.String(FlagOutputList, "", "Output List")

}
