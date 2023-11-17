package main

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDeleteNamespace(t *testing.T) {
	// Load the Kubernetes configuration
	config, err := clientcmd.BuildConfigFromFlags("", "C:/Users/SamundeeswariManimar/.kube/config")
	if err != nil {
		t.Fatalf("Failed to load Kubernetes config: %v", err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}

	// Define the namespace to delete
	namespace := "module-namespace"

	// Delete the namespace
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, v1.DeleteOptions{})
	if err != nil {
		t.Fatalf("Failed to delete namespace: %v", err)
	}

	// Wait for the namespace to be deleted
	time.Sleep(15 * time.Second)

	// Verify that the namespace and associated resources are removed
	convey.Convey("Verify namespace deletion", t, func() {
		// Check if the namespace still exists
		_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, v1.GetOptions{})
		convey.So(errors.IsNotFound(err), convey.ShouldBeTrue)

		// Check if associated resources are removed
		// Add more assertions here based on the resources you want to verify
	})
}
