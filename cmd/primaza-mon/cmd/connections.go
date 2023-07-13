/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/primaza/tools/pkg/connections"
	"github.com/primaza/tools/pkg/console"
	"github.com/primaza/tools/pkg/primaza"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// connectionsCmd represents the connections command
var connectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "get the connections between workloads and services",
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
		sdd, ww, err := cr.CrawlServiceConnections(cmd.Context(), tenant)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error crawling service connections: %s", err)
			return nil
		}

		if verboseFlag && len(ww) > 0 {
			for _, w := range ww {
				log.Println(w)
			}
		}

		errs := []error{}
		for _, sd := range sdd {
			sd := sd

			if err := printer.Println(&sd); err != nil {
				errs = append(errs, err)
			}
		}

		if err := errors.Join(errs...); err != nil {
			fmt.Fprintf(os.Stderr, "error printing output: %v", err)
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(connectionsCmd)
}
