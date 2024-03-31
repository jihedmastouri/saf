package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:       "saf",
	Short:     "",
	Long:      ``,
	Example:   "saf",
	ValidArgs: []string{
		// file or url
	},
	Version: "v0.0.0",
	PreRun: func(cmd *cobra.Command, args []string) {
		// logic to make sure that the flags are valid
		// and that we don't need flag for providing file and a url
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the server...")
	},
}

func main() {

	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
