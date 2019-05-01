package cmd

import (
	"fmt"
	"os"

	"github.com/arjunrn/dodns/pkg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "dodns is a modern dynamic dns",
	RunE:  synchronizeCmd,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func synchronizeCmd(cmd *cobra.Command, args []string) error {
	return pkg.Run(args)
}
