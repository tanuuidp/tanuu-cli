/*
Copyright (c) 2022 tanuuidp
*/
package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	l "github.com/k3d-io/k3d/v5/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanuuidp/tanuu-cli/cmd/setup"
	"gopkg.in/src-d/go-git.v4"
)

// startCmd represents the demo command
var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demo a local Tanuu setup",
	Long:  `A local Tanuu cluster, based on k3d. ghtoken is required, this is a user token for the tanuudemo gh user, where example apps are templated to.`,
	Run: func(cmd *cobra.Command, args []string) {
		demo = true
		os.Setenv("KUBECONFIG", "kubeconfig.tmp")
		os.RemoveAll(viper.GetString("localrepo"))
		_, err := git.PlainClone(viper.GetString("localrepo"), false, &git.CloneOptions{
			URL: "https://github.com/tanuuidp/tanuu",
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		ghtoken := viper.GetString("ghtoken")
		if ghtoken == "" {
			l.Log().Fatal("missing github token variable")
			return
		}

		kubeconfig, err := NewCmdClusterCreate()
		if err != nil {
			l.Log().Fatal(err)
		}
		l.Log().Info("Waiting for Tanuu demo, this might take 4-5 minutes, grab a cuppa!")
		l.Log().Debug(kubeconfig)

		time.Sleep(45 * time.Second)
		encoded1Text, err := base64.StdEncoding.DecodeString(viper.GetString("ghtoken"))
		if err != nil {
			panic(err)
		}
		encodedText, err := base64.StdEncoding.DecodeString(string(encoded1Text))
		if err != nil {
			panic(err)
		}
		setup.CheckDemo(kubeconfig)
		setup.SetDemoBackstageSecrets(kubeconfig, string(encodedText))
		adminPW := ""
		for adminPW == "" {
			adminPW = setup.GetArgoSecrets(kubeconfig)
		}
		time.Sleep(45 * time.Second)
		token, err := getArgoToken(adminPW)
		if err != nil {
			fmt.Println(err)

		}
		argoResync(token)
		time.Sleep(45 * time.Second)
		setup.CheckBackstage(kubeconfig)

		println("")
		println("READY")
		println("")
		println("ArgoCD URL: http://127.0.0.1:8081")
		println("AdminPW: ", adminPW)
		println("")
		println("ArgoWorkflow URL: http://127.0.0.1:8082")
		println("")
		println("Backstage URL: http://127.0.0.1:7007")
		println("")
		println("Demo app (once deployed): http://127.0.0.1:8084")
		println("")
		os.RemoveAll("/tmp/temp-manifests")
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
	rootCmd.PersistentFlags().StringVarP(&configfile, "ghtoken", "", "", "Github user token used for demo.")
	viper.BindPFlag("ghtoken", rootCmd.PersistentFlags().Lookup("ghtoken"))

}

type Token struct {
	Token string
}

func getArgoToken(password string) (string, error) {
	bodystring := fmt.Sprintf(`{"username": "admin","password": "%s" }`, password)
	jsonstring := []byte(bodystring)
	body := bytes.NewBuffer(jsonstring)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8081/api/v1/session", body)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	respBody, _ := io.ReadAll(resp.Body)
	var argotoken Token
	json.Unmarshal(respBody, &argotoken)
	return argotoken.Token, nil
}

func argoResync(argotoken string) {
	json := []byte(`{"name": "tanuu-setup-app"}`)
	body := bytes.NewBuffer(json)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8081/api/v1/applications/tanuu-setup-app/sync", body)
	if err != nil {
		fmt.Println("Failure : ", err)
	}
	authHeader := fmt.Sprintf("Bearer %s", argotoken)
	req.Header.Add("Authorization", authHeader)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failure : ", err)
	}
	respBody, _ := io.ReadAll(resp.Body)
	l.Log().Debug(respBody)
}
