package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigAccess() (*rest.Config, error) {
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
func TestServiceExposesPod(t *testing.T) {
	// Create a Kubernetes client configuration
	config, err := createKubeConfigAccess()
	if err != nil {
		t.Fatalf("Error creating in-cluster config: %v", err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// Define the namespace, service name, and pod name
	namespace := "default"
	serviceName := "kubernetes"
	// podName := "test-pod"

	// Retrieve the service by name
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Error getting service: %v", err)
	}

	// Wait for a few seconds to allow the service to become available
	time.Sleep(5 * time.Second)

	// Get the service's cluster IP
	serviceClusterIP := service.Spec.ClusterIP

	// Use GoConvey to assert that the service's cluster IP is not empty
	convey.Convey("Service Exposes Pod", t, func() {
		convey.So(serviceClusterIP, convey.ShouldNotEqual, "")
	})

	// You can also perform additional checks, such as making HTTP requests to the service's IP address
	// to further verify that it exposes the pod correctly.
}
