package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "redfish_exporter",
		Short: "Redfish Exporter for Prometheus",
	}

	rootCmd.AddCommand(
		newServeCmd(),
		newVersionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
