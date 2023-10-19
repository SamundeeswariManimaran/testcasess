package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfignetworkpolicy() (*rest.Config, error) {
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

func createPodd(clientset *kubernetes.Clientset, pod *corev1.Pod) error {
	_, err := clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	return err
}

func deletePodd(clientset *kubernetes.Clientset, namespace, podName string) error {
	return clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
}

func TestNetworkPolicyExternalAccess(t *testing.T) {
	// Create a Kubernetes client using the default configuration
	config, err := createKubeConfignetworkpolicy()
	if err != nil {
		t.Fatalf("Error creating Kubernetes client config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	Convey("Given a Kubernetes client", t, func() {
		namespace := "test-namespace"
		podName := "test-pod"

		// Define a Pod for testing
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      podName,
				Namespace: namespace,
				Labels:    map[string]string{"app": "test-app"},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "nginx",
						Image:           "nginx:latest",
						ImagePullPolicy: corev1.PullIfNotPresent,
					},
				},
			},
		}

		Convey("When we create a Pod in a namespace with network policies", func() {
			err := createPodd(clientset, pod)
			So(err, ShouldBeNil)

			// Sleep for a moment to allow network policies to take effect
			time.Sleep(5 * time.Second)

			Convey("Then the Pod should have restricted access to external services or endpoints", func() {
				// Here, you can perform tests to verify network policy enforcement.
				// For example, you can use the Pod to make HTTP requests to external services
				// or endpoints and check that the requests are blocked or restricted as per your network policies.

				// In a real-world scenario, you might use libraries like "curl" or make HTTP requests
				// from within the Pod to external services or endpoints and expect the requests to fail or be restricted.

				// For simplicity, we'll just print a message here.
				fmt.Println("Perform network policy enforcement tests for external access here.")

				Convey("When we delete the Pod", func() {
					err := deletePodd(clientset, namespace, podName)
					So(err, ShouldBeNil)
				})
			})
		})
	})
}
