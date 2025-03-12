package main

import (
	"fmt"

	"github.com/FreekingDean/redfish_exporter/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Redfish Exporter",
		Long:  `Print the version number of Redfish Exporter`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Redfish Exporter version: %s\n", version.Version())
		},
	}
}
