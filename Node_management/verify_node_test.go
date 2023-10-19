package main

import (
	"os/exec"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type Cluster struct {
	Nodes []string
}

func (c *Cluster) AddNode(nodeName string) {
	// Add logic to add a node to the cluster
	c.Nodes = append(c.Nodes, nodeName)
}

func (c *Cluster) RemoveNode(nodeName string) {
	// Add logic to remove a node from the cluster
	for i, node := range c.Nodes {
		if node == nodeName {
			c.Nodes = append(c.Nodes[:i], c.Nodes[i+1:]...)
			break
		}
	}
}

func TestCluster(t *testing.T) {
	Convey("verify node", t, func() { // Run the kubectl command to list nodes and capture the output
		cmd := exec.Command("kubectl", "get", "nodes", "-o=jsonpath='{.items[*].metadata.name}'")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("Error running kubectl: %v\n%s", err, output)
		}

		// Check if the output contains the name of the node (e.g., "minikube")
		expectedNodeName := "minikube"
		if !strings.Contains(string(output), expectedNodeName) {
			t.Fatalf("Node %s not found in the cluster\n%s", expectedNodeName, output)
		}
	})

}
