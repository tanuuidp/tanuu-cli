/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// nukeCmd represents the nuke command
var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Delete the local temp Tanuu server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		os.Setenv("KUBECONFIG", "kubeconfig.tmp")
		//TODO change this to use k3d package, and not CLI
		oscmd := exec.Command("k3d", "cluster", "delete", "tanuu")
		err := oscmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Tanuu nuked.")
	},
}

func init() {
	rootCmd.AddCommand(nukeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nukeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nukeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
