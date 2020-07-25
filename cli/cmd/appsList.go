package cmd

import (
	"fmt"
	"github.com/mrferos/feisty/cli/output"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var appsListCmd = &cobra.Command{
	Use:   "apps:list",
	Short: "List applications and the namespace they belong in",
	Long: `Used to list applications. Example:

feisty apps:set image=nginxdemos/hello:plain-text -a application-sample
`,
	Run: func(cmd *cobra.Command, args []string) {
		ns := getNamespace()
		apps, err := feistyClient.Applications(ns).List(v1.ListOptions{})
		if err != nil {
			fmt.Printf("Could not list applications, err: %v", err)
		}

		tableHeaders := []string{"NAMESPACE", "NAME"}
		tableData := [][]string{{}}
		for _, app := range apps.Items {
			tableData = append(tableData, []string{
				app.Namespace,
				app.Name,
			})
		}

		output.OutputTable(tableHeaders, tableData)
	},
}

func init() {
	appsListCmd.Flags().Bool("all", false, "return applications across all namespaces you have access to")
	rootCmd.AddCommand(appsListCmd)
}
