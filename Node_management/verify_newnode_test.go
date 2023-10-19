package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfignode() (*rest.Config, error) {
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

func TestNodeRegistrationAndReadiness(t *testing.T) {
	// Create a Kubernetes client configuration from your kubeconfig or in-cluster configuration.
	config, err := createKubeConfignode()
	if err != nil {
		t.Fatalf("Failed to create Kubernetes client config: %v", err)
	}

	// Create a Kubernetes clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}

	// Replace "NODE_NAME" with the name of the node you want to test.
	nodeName := "minikube"

	convey.Convey("Node Registration and Readiness Test", t, func() {
		convey.Convey("Node should be registered and ready", func() {
			// Get the node by name.
			node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
			convey.So(err, convey.ShouldBeNil)
			convey.So(node, convey.ShouldNotBeNil)

			// Check if the node is marked as Ready.
			for _, condition := range node.Status.Conditions {
				if condition.Type == corev1.NodeReady {
					convey.So(condition.Status, convey.ShouldEqual, corev1.ConditionTrue)
					return
				}
			}

			// Node is not marked as Ready.
			t.Errorf("Node %s is not ready", nodeName)
		})
	})
}
