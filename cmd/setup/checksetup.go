package setup

import (
	"context"
	"fmt"
	"time"

	l "github.com/k3d-io/k3d/v5/pkg/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CheckSetup(kubeconfig string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		crossplane, err := clientset.AppsV1().Deployments("crossplane-system").Get(context.TODO(), "crossplane", metav1.GetOptions{})

		if errors.IsNotFound(err) {
			fmt.Println("still working on it... please be patient...")
		} else {
			fmt.Println(crossplane.Name + " is now ready,")
			return
		}
		time.Sleep(10 * time.Second)
	}
}
func CheckBackstage(kubeconfig string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		crossplane, err := clientset.AppsV1().Deployments("backstage").Get(context.TODO(), "backstage", metav1.GetOptions{})

		if errors.IsNotFound(err) {
			fmt.Println("still working on it... please be patient...")
		} else {
			fmt.Println(crossplane.Name + " is now ready.")
			return
		}
		time.Sleep(10 * time.Second)
	}
}
func CheckDemo(kubeconfig string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		crossplane, err := clientset.AppsV1().Deployments("argocd").Get(context.TODO(), "argocd-server", metav1.GetOptions{})

		if errors.IsNotFound(err) {
			fmt.Println("still working on it... please be patient...")
		} else {
			fmt.Println(crossplane.Name + " is now ready, setting up the tools now...")
			return
		}
		time.Sleep(10 * time.Second)
	}
}

func SetDemoBackstageSecrets(kubeconfig string, secretdata string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	secret := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "backstage",
			Name:      "backstage-secrets",
		},
		Data: map[string][]byte{
			"GITHUB_TOKEN": []byte(secretdata),
		},
		Type: coreV1.SecretTypeOpaque,
	}
	createdSecret, err := clientset.CoreV1().Secrets("backstage").Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	l.Log().Trace("Secrets data: %+v", createdSecret.Data)

}
func GetArgoSecrets(kubeconfig string) string {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	otherSecret, _ := clientset.CoreV1().Secrets("argocd").Get(context.TODO(), "argocd-initial-admin-secret", metav1.GetOptions{})
	for _, value := range otherSecret.Data {
		return string(value)
	}
	return ""

}
func SetCredentialSecrets(kubeconfig string, secretdata string, namespace string, name string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	secret := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		StringData: map[string]string{
			"creds": secretdata,
		},
		Type: coreV1.SecretTypeOpaque,
	}
	createdSecret, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	l.Log().Trace("Secrets data: %+v", createdSecret.Data)
}

//TODO Add a 'create cluster for AWS/Azure from claim' section here

//TODO apply all the yamls also to the newly created cluster

//TODO dump all the yamls to a temp directory, so they can be added to repo

//TODO setup some bootstrap user, in case of issues in ziti

//TODO export everything, make some kind of test and then delete
