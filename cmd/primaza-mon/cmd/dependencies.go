/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/primaza/primaza-tools/pkg/connections"
	"github.com/primaza/primaza-tools/pkg/console"
	"github.com/primaza/primaza-tools/pkg/primaza"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// connectionsCmd represents the connections command
var connectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "list the connections between workloads and services",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		printer := console.NewPrinterOrDie(console.Format(outputFlag))

		cli, err := primaza.NewClient(config.GetConfigOrDie())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating Kubernetes API Client: %s", err)
			return nil
		}

		tenant := args[0]
		cr := connections.NewServiceDependeciesCrawler(cli)
		sdd, err := cr.CrawlServiceConnections(cmd.Context(), tenant)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error crawling service connections: %s", err)
		}

		errs := []error{}
		for _, sd := range sdd {
			sd := sd

			if err := printer.Println(&sd); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	},
}

func init() {
	getCmd.AddCommand(connectionsCmd)
}
