package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/primaza/primaza-tools/pkg/console"
	"github.com/primaza/primaza-tools/pkg/flags"
	"github.com/spf13/cobra"
)

var (
	outputFlag flags.OutputFlag = flags.OutputFlag(console.FormatJson)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "primaza-adm",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
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
