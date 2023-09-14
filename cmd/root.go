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
var demo bool

func setDefaults() {
	viper.SetDefault("port1", "8081")
	viper.SetDefault("port2", "8082")
	viper.SetDefault("port3", "7007")
	viper.SetDefault("port4", "8084")
	viper.SetDefault("port5", "8085")
	viper.SetDefault("localrepo", "/tmp/temp-manifests")
	viper.SetDefault("clustername", "tanuu")
}

var rootCmd = &cobra.Command{
	Use:   "tanuu",
	Short: "(BETA) Bootstrap a Tanuu installation",
	Long: `This app will start a local K3D based single container kubernetes 'cluster'
with the Tanuu setup tools installed and ready to boostrap a Tanuu management cluster.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if ll := os.Getenv("LOG_LEVEL"); ll != "" {
		level, err := logrus.ParseLevel(ll)
		if err == nil {
			l.Log().SetLevel(level)
		}
		if level == logrus.DebugLevel || level == logrus.TraceLevel {
			// If LOG_LEVEL is debug or trace, we assume
			// that the user wants to see the line numbers
			// and function name.
			l.Log().SetReportCaller(true)
		}
		l.Log().Info("Log level is ", level)
	}
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
