package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func createKubeConfigmapps() (*rest.Config, error) {
	// Check if running inside a Kubernetes cluster
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		// Use in-cluster config
		return rest.InClusterConfig()
	}

	// Create out-of-cluster config (e.g., for local development)
	userHomeDir, _ := os.UserHomeDir()
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func TestDeleteConfigMap(t *testing.T) {
	// Initialize Kubernetes clientset
	clientset := getKubernetesClientset()

	Convey("Given a ConfigMap exists", t, func() {
		// Create a sample ConfigMap
		configMap := createSampleConfigMap()
		_, err := clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), configMap, metav1.CreateOptions{})
		So(err, ShouldBeNil)

		Convey("When deleting the ConfigMap", func() {
			err := deleteConfigMap(clientset, "sample-configmap")

			Convey("It should delete the ConfigMap without errors", func() {
				So(err, ShouldBeNil)
			})

			Convey("The ConfigMap should no longer exist", func() {
				exists, _ := configMapExists(clientset, "sample-configmap")
				So(exists, ShouldBeFalse)
			})
		})
	})
}

func getKubernetesClientset() *kubernetes.Clientset {
	config, err := createKubeConfigmapps()
	if err != nil {
		kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func createSampleConfigMap() *v1.ConfigMap {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sample-configmap",
		},
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	return configMap
}

func configMapExists(clientset *kubernetes.Clientset, configMapName string) (bool, error) {
	_, err := clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func deleteConfigMap(clientset *kubernetes.Clientset, configMapName string) error {
	return clientset.CoreV1().ConfigMaps("default").Delete(context.TODO(), configMapName, metav1.DeleteOptions{})
}
