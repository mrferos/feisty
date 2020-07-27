package cmd

import (
	"fmt"
	"github.com/mrferos/feisty/cli/output"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func configsListCmdRun(args []string) error {
	ns := getNamespace()

	appConfig, err := feistyClient.ApplicationConfigs(ns).Get(appName, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("There was an error getting configurations for %s\n%v\n", appName, err)
	}

	configKeyVals := appConfig.Spec.KeyValuePairs
	headers := []string{"KEY", "VALUE"}
	data := [][]string{{}}
	for k, v := range configKeyVals {
		data = append(data, []string{k, v})
	}

	output.OutputTable(headers, data)

	return nil
}

var configsListCmd = &cobra.Command{
	Use:   "configs:list",
	Short: "List configurations for application",
	Long: `List configurations for application. Example:

feisty configs:list -a application-sample


`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := configsListCmdRun(args); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	},
}

func init() {
	configsListCmd.Flags().StringVarP(&appName, "app name", "a", "", "target application")
	rootCmd.AddCommand(configsListCmd)
}
