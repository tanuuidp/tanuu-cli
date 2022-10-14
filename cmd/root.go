package cmd

import (
	"fmt"
	"log"
	"os"

	cliconfig "github.com/k3d-io/k3d/v5/cmd/util/config"
	l "github.com/k3d-io/k3d/v5/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

var cfgFile string

func setDefaults() {
	viper.SetDefault("port", "8081")
	viper.SetDefault("localrepo", "/tmp/temp-manifests")
}

var rootCmd = &cobra.Command{
	Use:   "tanuu",
	Short: "Bootstrap a Tanuu installation",
	Long: `This app will start a local K3D based single container kubernetes 'cluster'
with the Tanuu setup tools installed and ready to boostrap a Tanuu management cluster.

First, setup your credentials for the provider you are deploying to. 

AWS:
echo "[default]
aws_access_key_id = <your access key ID>
aws_secret_access_key = <your secret key>
" > creds.conf.

Azure:
az ad sp create-for-rbac --name crossplane-demo --scopes /subscriptions/{SubID}/resourceGroups/{ResourceGroup} > creds.json


Then, tanuu start...`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	// cobra.OnInitialize(initConfig)
	err := initConfig()
	if err != nil {
		log.Fatal("Could not initialize")
	}
	setDefaults()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// Add my subcommand palette
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "tanuu.env", "config file (default is tanuu.env)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() error {

	viper.SetEnvPrefix("TANUU")
	viper.AutomaticEnv() // read in environment variables that match
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	// Viper for pre-processed config options
	ppViper.SetEnvPrefix("K3D")

	if l.Log().GetLevel() >= logrus.DebugLevel {

		c, _ := yaml.Marshal(ppViper.AllSettings())
		l.Log().Debugf("Additional CLI Configuration:\n%s", c)
	}

	return cliconfig.InitViperWithConfigFile(cfgViper, configFile)
}
