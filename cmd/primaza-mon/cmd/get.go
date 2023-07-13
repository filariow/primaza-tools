/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "use this command to get information from a primaza tenant",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
