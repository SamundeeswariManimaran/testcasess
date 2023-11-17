package main

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestNestedNamespaces(t *testing.T) {
	// Load Kubernetes configuration from a file (you can customize this)
	config, err := clientcmd.BuildConfigFromFlags("", "C:/Users/SamundeeswariManimar/.kube/config")
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	Convey("Test Nested Namespaces", t, func() {
		// Define the parent and child namespaces
		parentNamespace := "parent-namespace"
		childNamespace := "parent-namespace/child-namespace"

		// Function to check if a namespace exists
		namespaceExists := func(namespaceName string) bool {
			_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
			return err == nil
		}

		// Create the parent namespace if it doesn't exist
		if !namespaceExists(parentNamespace) {
			createNamespace(parentNamespace)
		}

		Convey("Create Child Namespace", func() {
			// Check if the child namespace already exists
			if !namespaceExists(childNamespace) {
				// Attempt to create the child namespace
				childNamespaceObj := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{Name: childNamespace},
				}
				_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), childNamespaceObj, metav1.CreateOptions{})

				// Ensure that creating a child namespace within a non-existent parent namespace results in an error
				So(err, ShouldNotBeNil)
			}
		})

		// Convey("Test Nested Namespaces", t, func() {
		// 	// Define the parent and child namespaces
		// 	parentNamespace := "pai-namespace"
		// 	childNamespace := "pai-namespace/child-namespace"

		// 	// Create the parent namespace
		// 	createNamespace := func(namespaceName string) {
		// 		newNamespace := &corev1.Namespace{
		// 			ObjectMeta: metav1.ObjectMeta{Name: namespaceName},
		// 		}
		// 		_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), newNamespace, metav1.CreateOptions{})
		// 		So(err, ShouldBeNil)
		// 	}

		// 	createNamespace(parentNamespace)

		// 	Convey("Create Child Namespace", func() {
		// 		// Attempt to create the child namespace
		// 		childNamespaceObj := &corev1.Namespace{
		// 			ObjectMeta: metav1.ObjectMeta{Name: childNamespace},
		// 		}
		// 		_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), childNamespaceObj, metav1.CreateOptions{})

		// 		// Ensure that creating a child namespace within a non-existent parent namespace results in an error
		// 		So(err, ShouldNotBeNil)
		// 	})

		Convey("Create Pod in Child Namespace", func() {
			// Attempt to create a Pod in the child namespace
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "new-pod", Namespace: childNamespace},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "test-container", Image: "nginx"}},
				},
			}
			_, err := clientset.CoreV1().Pods(childNamespace).Create(context.TODO(), pod, metav1.CreateOptions{})

			// Ensure that creating a Pod in a child namespace within a non-existent parent namespace results in an error
			So(err, ShouldNotBeNil)
		})
	})
}

func createNamespace(parentNamespace string) {
	panic("unimplemented")
}
