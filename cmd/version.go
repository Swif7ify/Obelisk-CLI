package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Obelisk",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Obelisk CLI\n")
		fmt.Printf("  Version:    %s\n", Version)
		fmt.Printf("  Commit:     %s\n", Commit)
		fmt.Printf("  Built:      %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
