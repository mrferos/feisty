package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strconv"
)

func appSetCmdRun(args []string) error {
	ns := getNamespace()

	app, err := feistyClient.Applications(ns).Get(appName, v1.GetOptions{});
	if err != nil {
		return fmt.Errorf("could not load application %s\n%v\n", appName, err)

	}

	parsedArgs, err := parseArgs(args)
	if err != nil {
		return fmt.Errorf("could not parse arguments")
	}

	for key, val := range parsedArgs {
		switch key {
		case "image":
			app.Spec.Image = val
		case "replicas":
			replicas, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("could not parse replicas: %s", val)
			} else {
				app.Spec.Replicas = replicas
			}
		case "port":
			port, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("could not parse port: %s", val)
			} else {
				app.Spec.Port = port
			}
		case "routingEnabled":
			if val == "true" {
				app.Spec.RoutingEnabled = true
			} else if val == "false" {
				app.Spec.RoutingEnabled = false
			} else {
				return fmt.Errorf("could not parse routing enabled; %s", val)
			}
		}
	}

	if _, err := feistyClient.Applications(ns).Update(app); err != nil {
		return fmt.Errorf("there was an error updating %s\n%v", app.Name, err)
	}

	return nil
}

var appsSetCmd = &cobra.Command{
	Use:   "apps:set",
	Short: "Set options for applications",
	Long: `Set application options. Example:

feisty apps:set image=nginx/helloworld -a application-sample

Supported options:
	* image - the image to deploy
	* replicas - how many instances of the application should be running
	* port - the application's exposed port
	* routingEnabled - true/false value to manage an ingress for the application
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("an application name is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := appSetCmdRun(args); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	},
}

func init() {
	appsSetCmd.Flags().StringVarP(&appName, "app name", "a", "", "target application")
	rootCmd.AddCommand(appsSetCmd)
}
