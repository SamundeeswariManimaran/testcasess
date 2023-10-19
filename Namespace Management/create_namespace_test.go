package main

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createNamespaceIfNotExists(clientset *kubernetes.Clientset, namespaceName string) error {
	_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// Namespace does not exist, create it
			newNamespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: namespaceName},
			}
			_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), newNamespace, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			// Other error occurred, return it
			return err
		}
	}
	return nil
}

func TestCreateNamespaceIfNotExists(t *testing.T) {
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

	Convey("Test CreateNamespaceIfNotExists", t, func() {
		namespaceName := "test-namespace"

		err := createNamespaceIfNotExists(clientset, namespaceName)
		So(err, ShouldBeNil)

		// Verify that the namespace exists
		_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
		So(err, ShouldBeNil)
	})
}
