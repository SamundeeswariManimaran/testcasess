package main

import (
	"os"
	"path/filepath"
	"testing"

	"context"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigValidate() (*rest.Config, error) {
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
func TestNodeFailureDetection(t *testing.T) {
	// Create a Kubernetes client using the default configuration
	config, err := createKubeConfigValidate()
	if err != nil {
		t.Fatalf("Error creating Kubernetes client config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	Convey("Given a Kubernetes cluster with a monitoring system", t, func() {
		nodeName := "minikube"

		Convey("When we simulate a node failure", func() {
			// Simulate a node failure by marking the node unschedulable
			node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
			So(err, ShouldBeNil)

			// Mark the node unschedulable (simulating a failure)
			node.Spec.Unschedulable = true
			_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
			So(err, ShouldBeNil)

			// Wait for the monitoring system to detect the failure (adjust the duration as needed)
			time.Sleep(30 * time.Second)

			Convey("Then the monitoring system should detect the node failure", func() {
				// You should add assertions here to validate that the monitoring system
				// correctly detects the node failure, e.g., by checking alerts or metrics.
				// The specific implementation will depend on your monitoring system.

				// Example: Check for alerts related to node failure
				// alert, err := monitoringSystem.GetNodeFailureAlert(nodeName)
				// So(err, ShouldBeNil)
				// So(alert.Status, ShouldEqual, "Firing")

				// You should adapt this part to match your monitoring system's behavior and APIs.
			})
		})
	})
}
