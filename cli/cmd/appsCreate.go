package cmd

import (
	"errors"
	"fmt"
	"github.com/mrferos/feisty/api/v1alpha1"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

var appsCreateCmd = &cobra.Command{
	Use:   "apps:create",
	Short: "Create application",
	Long: `Create a new application. Example:

feisty apps:create cool-app
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("an application name is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ns := getNamespace()
		app := v1alpha1.Application{
			ObjectMeta: v1.ObjectMeta{
				Name: args[0],
				Namespace: ns,
			},
		}

		if _, err := feistyClient.Applications(ns).Create(&app); err != nil {
			fmt.Printf("There was an error creating the application: \n")
			fmt.Printf("%v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("The application %s in namespace %s was created!\n", app.Name, app.Namespace)
		}
	},
}

func init() {
	rootCmd.AddCommand(appsCreateCmd)
}
