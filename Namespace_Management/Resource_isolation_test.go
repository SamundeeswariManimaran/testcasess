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

func TestResourceIsolationBetweenNamespaces(t *testing.T) {
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

	Convey("Test Resource Isolation Between Namespaces", t, func() {
		// Create two namespaces
		namespaceA := "namespace-a"
		namespaceB := "namespace-b"

		createNamespace := func(namespaceName string) {
			newNamespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: namespaceName},
			}
			_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), newNamespace, metav1.CreateOptions{})
			So(err, ShouldBeNil)
		}

		createNamespace(namespaceA)
		createNamespace(namespaceB)

		// Create a ConfigMap in namespaceA
		configMapA := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "configmap-a", Namespace: namespaceA},
			Data:       map[string]string{"key": "value"},
		}
		_, err := clientset.CoreV1().ConfigMaps(namespaceA).Create(context.TODO(), configMapA, metav1.CreateOptions{})
		So(err, ShouldBeNil)

		// Attempt to access the ConfigMap from namespaceB
		_, err = clientset.CoreV1().ConfigMaps(namespaceB).Get(context.TODO(), "configmap-a", metav1.GetOptions{})
		So(err, ShouldNotBeNil) // Expect an error because of resource isolation
	})
}
