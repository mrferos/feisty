package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"time"
)

func appsRestartCmdRun(args []string) error {
	ns := getNamespace()

	app, err := feistyClient.Applications(ns).Get(appName, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("could not load application %s\n%v\n", appName, err)
	}

	currentTime := time.Now()
	app.Spec.RestartTime = currentTime.Format("2006-01-02 15:04:05.000000000")

	if _, err := feistyClient.Applications(ns).Update(app); err != nil {
		return fmt.Errorf("there was an error restarting %s\n%v", app.Name, err)
	}

	fmt.Printf("%s in %s was restarted", app.Name, app.Namespace)

	return nil
}

var appsRestartCmd = &cobra.Command{
	Use:   "apps:restart",
	Short: "Restart application",
	Long: `Issue a restart of all running application pods. Example:

feisty apps:restart -a application-sample
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := appsRestartCmdRun(args); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	},
}

func init() {
	appsRestartCmd.Flags().StringVarP(&appName, "app name", "a", "", "target application")
	rootCmd.AddCommand(appsRestartCmd)
}
