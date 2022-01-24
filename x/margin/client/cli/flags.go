package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagBaseUrl = "baseUrl"
)

// common flagsets to add to various functions
var (
	FsBaseUrl = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {

	FsBaseUrl.String(FlagBaseUrl, "", "BaseUrl to query")

}
