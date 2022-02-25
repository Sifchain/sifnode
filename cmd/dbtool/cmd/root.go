package cmd

import (
	"fmt"
	"os"

	"github.com/Sifchain/sifnode/cmd/dbtool/utils"
	"github.com/spf13/cobra"
)

var (
	datadir = fmt.Sprintf("%s/.sifnoded/data", homeDir())
	verbose = false
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dbtool",
		Short: "A tool to query the sifnode database",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				utils.SetVerbose()
			}
			err := utils.OpenDB(datadir)
			if err != nil {
				panic(err)
			}
		},
	}
	addPersistentFlags(rootCmd)
	addCommands(rootCmd)
	return rootCmd
}

func addPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&datadir, "data", "d", datadir, "Data directory")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "Verbose output")
}

func addCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewSearchCmd(),
		NewIBCCmd(),
	)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}

func homeDir() string {
	hd, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return hd
}
