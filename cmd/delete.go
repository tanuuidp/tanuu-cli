/*
Copyright (c) 2022 tanuuidp
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the local temp Tanuu server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		os.Setenv("KUBECONFIG", "kubeconfig.tmp")
		DeleteCluster()
		os.RemoveAll("kubeconfig.tmp")
		fmt.Println("Tanuu deleted.")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
