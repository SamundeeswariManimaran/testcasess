package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfignodee() (*rest.Config, error) {
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

func TestNodeInformationUpdate(t *testing.T) {
	// Create a Kubernetes client configuration from your kubeconfig or in-cluster configuration.
	config, err := createKubeConfignodee()
	if err != nil {
		t.Fatalf("Failed to create Kubernetes client config: %v", err)
	}

	// Create a Kubernetes clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}

	convey.Convey("Node Information Update Test", t, func() {
		convey.Convey("Check Node Information Before Update", func() {
			nodeListBefore, err := clientset.CoreV1().Nodes().List(context.TODO(),metav1.ListOptions{})
			convey.So(err, convey.ShouldBeNil)
			convey.So(len(nodeListBefore.Items), convey.ShouldBeGreaterThan, 0)

			// Choose a Node to update (in this example, the first Node)
			nodeToUpdate := &nodeListBefore.Items[0]

			// Assuming you want to add a label to the Node
			nodeToUpdate.Labels["my-label"] = "updated-value"

			// Update the Node information in the cluster
			_, err = clientset.CoreV1().Nodes().Update(context.TODO(),nodeToUpdate, metav1.UpdateOptions{})
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("Check Node Information After Update", func() {
				nodeListAfter, err := clientset.CoreV1().Nodes().List(context.TODO(),metav1.ListOptions{})
				convey.So(err, convey.ShouldBeNil)
				convey.So(len(nodeListAfter.Items), convey.ShouldBeGreaterThan, 0)

				// Ensure that the number of nodes remains the same after the update
				convey.So(len(nodeListAfter.Items), convey.ShouldEqual, len(nodeListBefore.Items))

				// Add checks to verify the updated information, for example, checking the updated label
				updatedNode := &nodeListAfter.Items[0]
				convey.So(updatedNode.Labels["my-label"], convey.ShouldEqual, "updated-value")
			})
		})
	})
}
