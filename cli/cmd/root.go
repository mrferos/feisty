package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var appName string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "A Feisty CLI",
	Long:  `The Feisty CLI used to deploy and configure apps on your cluster`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&appName, "app name", "a", "", "target application")
}
