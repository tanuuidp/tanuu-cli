/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Tanuu to existing k8s cluster (non-functional placeholder option, does nothing)",
	Long: `This option doesn't create a kubernetes cluster, and needs an existing one. 
It will deploy the Tanuu artifacts to setup the given cluster as a Tanuu managment cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deploy called")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
