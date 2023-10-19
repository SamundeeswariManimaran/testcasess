package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigg() (*rest.Config, error) {
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

func TestPodAccessibility(t *testing.T) {
	// Create a Kubernetes client configuration
	config, err := createKubeConfigg()
	if err != nil {
		t.Fatalf("Error creating in-cluster config: %v", err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// Define the namespace and pod name
	namespace := "default"
	podName := "test-pod" // Replace with the actual pod name

	// Define a function to check pod accessibility
	checkPodAccessibility := func() bool {
		// Use clientset to retrieve the pod
		pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			return false
		}

		// Check the pod's status or other criteria to determine accessibility
		// For example, you can check if the pod is running and ready
		return pod.Status.Phase == v1.PodRunning && pod.Status.ContainerStatuses[0].Ready
	}

	// Use GoConvey to assert that the pod is accessible
	convey.Convey("Check if Pod is Accessible", t, func() {
		convey.So(checkPodAccessibility(), convey.ShouldBeTrue)
	})
}
