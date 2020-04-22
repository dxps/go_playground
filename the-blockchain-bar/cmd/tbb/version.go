package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	major = "0"
	minor = "1"
	fix   = "0"
	note  = "Txn Add && Balances List"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "It describes the version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s.%s.%s-beta %s\n", major, minor, fix, note)
	},
}
