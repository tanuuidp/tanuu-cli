/*
Copyright (c) 2022 tanuuidp
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	l "github.com/k3d-io/k3d/v5/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanuuidp/tanuu-cli/cmd/setup"
	"gopkg.in/src-d/go-git.v4"
)

var aws bool
var azure bool
var filename []byte
var configfile string
var clustername string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "(BETA) Start local Tanuu bootstrap server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		os.Setenv("KUBECONFIG", "kubeconfig.tmp")

		repo, _ := git.PlainClone(viper.GetString("localrepo"), false, &git.CloneOptions{
			URL: "https://github.com/tanuuidp/tanuu",
		})
		if repo != nil {
			return
		}
		if aws {
			filename, err := os.ReadFile(viper.GetString("creds"))
			if err != nil {
				l.Log().Fatal("Cannot read creds for the AWS credentials")
				panic(err.Error())
			}
			fmt.Println("using " + viper.GetString("creds"))
			l.Log().Trace(filename)
		} else if azure {
			filename, err := os.ReadFile(viper.GetString("creds"))
			if err != nil {
				l.Log().Fatal("Cannot read creds for the Azure credentials")
				panic(err.Error())
			}
			l.Log().Trace(filename)
		} else {
			l.Log().Warn("Deploying empty tanuu")
			// l.Log().Fatal("No valid deploy target specified.")
		}

		kubeconfig, err := NewCmdClusterCreate()
		if err != nil {
			l.Log().Fatal(err)
		}
		l.Log().Info("Waiting for Tanuu bootstrap local cluster to be ready, this might take 2-3 minutes, grab a cuppa!")
		l.Log().Debug(kubeconfig)

		if aws {
			filename, err := os.ReadFile(viper.GetString("creds"))
			if err != nil {
				l.Log().Fatal("Cannot read creds for the AWS credentials")
				panic(err.Error())
			}
			setup.SetCredentialSecrets(kubeconfig, string(filename), "crossplane-system", "aws-creds")
		} else if azure {
			filename, err := os.ReadFile(viper.GetString("creds"))
			if err != nil {
				l.Log().Fatal("Cannot read creds for the Azure credentials")
				panic(err.Error())
			}
			setup.SetCredentialSecrets(kubeconfig, string(filename), "crossplane-system", "azure-creds")
		} else {
			l.Log().Warn("Deployed empty tanuu")
			// l.Log().Fatal("No valid deploy target specified.")
			// panic(err.Error())
		}

		time.Sleep(120 * time.Second)
		setup.CheckSetup(kubeconfig)
		println("READY!")
		os.RemoveAll("/tmp/temp-manifests")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.PersistentFlags().StringVarP(&configfile, "creds", "c", "", "Credential file to be used. ex: creds.json for Azure, creds.conf for AWS.")
	rootCmd.PersistentFlags().BoolVar(&aws, "aws", false, "Use AWS")
	rootCmd.PersistentFlags().BoolVar(&azure, "azure", false, "Use Azure")
	rootCmd.MarkFlagsMutuallyExclusive("aws", "azure")
	viper.BindPFlag("creds", rootCmd.PersistentFlags().Lookup("creds"))
	viper.BindPFlag("aws", rootCmd.PersistentFlags().Lookup("aws"))
	viper.BindPFlag("azure", rootCmd.PersistentFlags().Lookup("azure"))

}
