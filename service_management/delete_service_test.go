package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigdelete() (*rest.Config, error) {
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
func deleteService(clientset *kubernetes.Clientset, serviceName, namespace string) error {
	// Delete the service
	deletePolicy := metav1.DeletePropagationForeground
	return clientset.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func isServiceAccessible(serviceURL string) bool {
	response, err := http.Get(serviceURL)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	return response.StatusCode == http.StatusOK
}

func TestDeleteService(t *testing.T) {
	Convey("Given a Kubernetes cluster configuration", t, func() {
		// Load the Kubernetes configuration from in-cluster
		config, err := createKubeConfigdelete()
		So(err, ShouldBeNil)

		// Create a Kubernetes clientset
		clientset, err := kubernetes.NewForConfig(config)
		So(err, ShouldBeNil)

		// Define the service name and namespace
		serviceName := "new-service" // Replace with the actual service name
		namespace := "default"       // Replace with the actual namespace

		Convey("When deleting the service", func() {
			err := deleteService(clientset, serviceName, namespace)
			So(err, ShouldBeNil)

			Convey("It should no longer be accessible", func() {
				// Construct the service URL based on the ClusterIP and port
				serviceURL := " http://127.0.0.1:58840" // Replace with your service details

				// Check if the service is still accessible
				accessible := isServiceAccessible(serviceURL)
				So(accessible, ShouldBeFalse)
			})
		})
	})
}
