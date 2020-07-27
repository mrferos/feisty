package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func configsSetCmdRun(args []string) error {
	ns := getNamespace()

	appConfig, err := feistyClient.ApplicationConfigs(ns).Get(appName, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("could not load config %s\n%v\n", appName, err)

	}

	parsedArgs, err := parseArgs(args)
	if err != nil {
		return fmt.Errorf("could not parse configs")
	}

	for key, val := range parsedArgs {
		appConfig.Spec.KeyValuePairs[key] = val
	}

	if _, err := feistyClient.ApplicationConfigs(ns).Update(appConfig); err != nil {
		return fmt.Errorf("there was an error updating %s\n%v", appConfig.Name, err)
	}

	fmt.Printf("%s in %s was updated", appConfig.Name, appConfig.Namespace)

	return nil
}

var configsSetCmd = &cobra.Command{
	Use:   "configs:set",
	Short: "Set application configs",
	Long: `Set application configs. Example:

feisty configs:set ONE=value TWO=more POSSIBLE=values -a application-sample

`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("application configs are required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := configsSetCmdRun(args); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	},
}

func init() {
	configsSetCmd.Flags().StringVarP(&appName, "app name", "a", "", "target application")
	rootCmd.AddCommand(configsSetCmd)
}
