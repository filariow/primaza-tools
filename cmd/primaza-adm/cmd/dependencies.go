/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/primaza/primaza-tools/pkg/dependencies"
	"github.com/primaza/primaza-tools/pkg/mermaid"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// dependenciesCmd represents the dependencies command
var dependenciesCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := client.New(config.GetConfigOrDie(), client.Options{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating Kubernetes API Client: %s", err)
			return nil
		}

		tenant := args[0]
		cr := dependencies.NewServiceDependeciesCrawler(cli)
		sdd, err := cr.CrawlServiceDependencies(cmd.Context(), tenant)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error crawling service dependencies: %s", err)
		}

		d, err := json.Marshal(sdd)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(d))
		return nil
	},
}

func init() {
	listCmd.AddCommand(dependenciesCmd)

	listCmd.Flags().BoolP("tenant", "t", false, "Help message for toggle")
}

func listDependencies() {
	m := mermaid.Graph{
		Name: "tenant-name",
		Adjacencies: []mermaid.Adjancency{
			{Start: "A", End: "B"},
			{Start: "A", End: "C"},
			{Start: "B", End: "D"},
		},
	}
	fmt.Println(m.String())
}
