package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetNodeDetails(clientset *kubernetes.Clientset, nodeName string) (*NodeDetails, error) {
	// Replace this with actual code to retrieve node details
	// Example:
	// node, err := clientset.CoreV1().Nodes().Get(nodeName, v1.GetOptions{})
	// if err != nil {
	//     return nil, err
	// }
	// details := &NodeDetails{
	//     Name:     node.Name,
	//     Capacity: node.Status.Capacity,
	//     Labels:   node.Labels,
	//     // Add other relevant information
	// }
	// return details, nil

	// For this example, return mock data
	return &NodeDetails{
		Name:     "mock-node-name",
		Capacity: map[string]string{"cpu": "4", "memory": "8Gi"},
		Labels:   map[string]string{"app": "example-app"},
		// Add other mock relevant information
	}, nil
}

type NodeDetails struct {
	Name     string
	Capacity map[string]string
	Labels   map[string]string
	// Add other relevant fields
}

func TestDisplayNodeDetails(t *testing.T) {
	// Load Kubernetes configuration
	kubeconfig := "C:/Users/SamundeeswariManimar/.kube/config" // Update with the actual path to your kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Create a Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	Convey("Given a Node Management System is displaying node details", t, func() {
		nodeName := "minikube" // Replace with the name of the node you want to test

		Convey("When fetching node details", func() {
			nodeDetails, err := GetNodeDetails(clientset, nodeName)
			So(err, ShouldBeNil)

			Convey("Then the displayed node details should match the expected values", func() {
				expectedDetails := &NodeDetails{
					Name:     "mock-node-name",
					Capacity: map[string]string{"cpu": "4", "memory": "8Gi"},
					Labels:   map[string]string{"app": "example-app"},
					// Add expected values for other relevant fields
				}

				So(nodeDetails, ShouldResemble, expectedDetails)
			})
		})
	})
}
