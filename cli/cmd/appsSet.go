package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// appsSetCmd represents the appsSet command
var appsSetCmd = &cobra.Command{
	Use:   "apps:set",
	Short: "Set a value on the Application object",
	Long: `Can be used to set values on the application object. For example:

feisty apps:set image=nginxdemos/hello:plain-text -a application-sample
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("appsSet called")
	},
}

func init() {
	rootCmd.AddCommand(appsSetCmd)
}
