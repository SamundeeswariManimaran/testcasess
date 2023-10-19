package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigmap() (*rest.Config, error) {
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

func createConfigMap(clientset *kubernetes.Clientset, namespace, configMapName string, data map[string]string) (*v1.ConfigMap, error) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: configMapName,
		},
		Data: data,
	}
	return clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
}

func getConfigMap(clientset *kubernetes.Clientset, namespace, configMapName string) (*v1.ConfigMap, error) {
	return clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
}

func TestCreateConfigMap(t *testing.T) {
	//kubeconfig := "C:/Users/SamundeeswariManimar/.kube/config" // Update with your kubeconfig path.
	namespace := "default"
	configMapName := "new-configmap"
	data := map[string]string{
		"key1": "1",
		"key2": "2",
	}

	config, err := createKubeConfigmap()
	if err != nil {
		t.Fatalf("Failed to create Kubernetes config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}

	Convey("Given a Kubernetes cluster", t, func() {
		Convey("When a ConfigMap is created in the cluster", func() {
			createdConfigMap, err := createConfigMap(clientset, namespace, configMapName, data)
			So(err, ShouldBeNil)

			Convey("Then the ConfigMap should be created successfully", func() {
				// Sleep for a moment to allow time for the ConfigMap to be created.
				time.Sleep(2 * time.Second)

				// Verify that the ConfigMap exists in the cluster.
				retrievedConfigMap, err := getConfigMap(clientset, namespace, configMapName)
				So(err, ShouldBeNil)
				So(retrievedConfigMap.Name, ShouldEqual, createdConfigMap.Name)
				So(retrievedConfigMap.Data, ShouldResemble, createdConfigMap.Data)
			})
		})

		Convey("When a non-existing ConfigMap is retrieved", func() {
			Convey("Then an error should be returned", func() {
				_, err := getConfigMap(clientset, namespace, "non-existent-configmap")
				So(err, ShouldNotBeNil)
			})
		})
	})
}
