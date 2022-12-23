/*
Copyright (c) 2022 tanuuidp
*/
package cmd

import (
	"os"
	"time"

	l "github.com/k3d-io/k3d/v5/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanuuidp/tanuu-cli/cmd/setup"
	"gopkg.in/src-d/go-git.v4"
)

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
		os.RemoveAll("kubeconfig.tmp")

		repo, _ := git.PlainClone(viper.GetString("localrepo"), false, &git.CloneOptions{
			URL: "https://github.com/tanuuidp/tanuu",
		})
		if repo != nil {
			return
		}
		l.Log().Warn("Deploying tanuu")
		// l.Log().Fatal("No valid deploy target specified.")

		kubeconfig, err := NewCmdClusterCreate()
		if err != nil {
			l.Log().Fatal(err)
		}
		l.Log().Info("Waiting for Tanuu bootstrap local cluster to be ready, this might take 2-3 minutes, grab a cuppa!")
		l.Log().Debug(kubeconfig)

		l.Log().Warn("Deployed tanuu")

		time.Sleep(120 * time.Second)
		setup.CheckSetup(kubeconfig)
		println("READY!")
		os.RemoveAll("/tmp/temp-manifests")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
