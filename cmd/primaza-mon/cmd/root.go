package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/primaza/tools/pkg/console"
	"github.com/primaza/tools/pkg/flags"
	"github.com/spf13/cobra"
)

var (
	outputFlag  flags.OutputFlag = flags.OutputFlag(console.FormatJson)
	verboseFlag bool             = false
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "primaza-mon",
	Short: "primaza-mon helps monitoring a Primaza Tenant",
	Long: `use primaza-mon to list the connections between workloads and services
in your Primaza Tenants`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&verboseFlag,
		"verbose",
		"v",
		false,
		fmt.Sprintf("Verbose output"))

	rootCmd.PersistentFlags().VarP(
		&outputFlag,
		"output",
		"o",
		fmt.Sprintf("Set the output format. Allowed format: %s",
			strings.Join(console.AllowedFormats(), ", ")))
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
