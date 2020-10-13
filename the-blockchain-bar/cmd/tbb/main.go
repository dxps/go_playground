package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	var cmd = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	cmd.AddCommand(versionCmd)
	cmd.AddCommand(balancesCmd())
	cmd.AddCommand(txCmd())

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
