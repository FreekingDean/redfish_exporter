package main

import (
	"log"

	"github.com/FreekingDean/redfish_exporter/internal/cmds"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "redfish_exporter",
		Short: "Redfish Exporter for Prometheus",
	}

	rootCmd.AddCommand(
		cmds.NewServeCmd(),
		cmds.NewVersionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
