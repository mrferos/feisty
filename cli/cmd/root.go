package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/mrferos/feisty/cli/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string
var appNamespace string
var appName string
var feistyClient *client.FeistyV1Alpha1Client

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

func getNamespace() string {
	if appNamespace == "" {
		if defaultNamespace := viper.GetString("defaultNamespace"); defaultNamespace != "" {
			return defaultNamespace
		}

		return appName
	}

	return appNamespace
}

func init() {
	cobra.OnInitialize(initConfig)
	//rootCmd.PersistentFlags().StringVarP(&appName, "app name", "a", "", "target application")
	rootCmd.PersistentFlags().StringVarP(&appNamespace, "namespace", "n", "", "target namespace (defaults to app name)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.feisty.yml")

	var err error
	feistyClient, err = client.GetFeistyClient()
	if err != nil {
		// TODO: make this cleaner
		panic(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".feisty")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Could not load config file: ", viper.ConfigFileUsed())
	}
}
