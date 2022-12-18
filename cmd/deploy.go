/*
Copyright (c) 2022 tanuuidp
*/
package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Tanuu to existing k8s cluster (non-functional placeholder option, does nothing)",
	Long: `This option doesn't create a kubernetes cluster, and needs an existing one. 
It will deploy the Tanuu artifacts to setup the given cluster as a Tanuu managment cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deploy called")
		input, err := ioutil.ReadFile("aws.tmpl")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		output := bytes.Replace(input, []byte("placeholder"), []byte(viper.GetString("clustername")), -1)

		if err = ioutil.WriteFile("aws.yaml", output, 0666); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// rootCmd.AddCommand(deployCmd)
	rootCmd.PersistentFlags().StringVarP(&clustername, "clustername", "n", "tanuu", "Name of the management cluster, default to 'tanuu'")
	viper.BindPFlag("clustername", rootCmd.PersistentFlags().Lookup("clustername"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
