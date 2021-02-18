package kubeclient

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/betorvs/secretpublisher/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// lazyInit lazy funcion to init Repository
func lazyInit() *kubernetes.Clientset {
	var clientConfig *rest.Config
	var err error
	if config.LocalKubeconfig {
		// home := os.Getenv("HOME")
		var kubeconfig string
		if home := homeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
		// kubeconfig = fmt.Sprintf("%s/%s", home, ".kube/config")
		// use the current context in kubeconfig
		clientConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// creates the in-cluster config
		clientConfig, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// GetSecrets return all secrets from a namespace and labels
func GetSecrets(namespace, labels string) (*v1.SecretList, error) {
	kube := lazyInit()
	listOptions := metav1.ListOptions{}
	if len(labels) > 0 {
		listOptions.LabelSelector = labels
	}
	secrets, err := kube.CoreV1().Secrets(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return &v1.SecretList{}, fmt.Errorf("Failed to get secrets: %v", err)
	}
	fmt.Printf("Number of kubernetes secrets found: %d \n", len(secrets.Items))
	return secrets, nil
}

// GetConfigMaps return all configMap from a namespace and labels
func GetConfigMaps(namespace, labels string) (*v1.ConfigMapList, error) {
	kube := lazyInit()
	listOptions := metav1.ListOptions{}
	if len(labels) > 0 {
		listOptions.LabelSelector = labels
	}
	cm, err := kube.CoreV1().ConfigMaps(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return &v1.ConfigMapList{}, fmt.Errorf("Failed to get config maps: %v", err)
	}
	fmt.Printf("Number of kubernetes config maps found: %d \n", len(cm.Items))
	return cm, nil
}
