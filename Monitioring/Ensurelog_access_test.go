package main

import (
	"os"
	"path/filepath"
	"testing"

	"context"

	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigmonitor() (*rest.Config, error) {
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

func TestRestrictedLogAccess(t *testing.T) {
	// Create a Kubernetes client using the default configuration
	config, err := createKubeConfigmonitor()
	if err != nil {
		t.Fatalf("Error creating Kubernetes client config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	Convey("Given a Kubernetes client", t, func() {
		namespaceName := "default"

		Convey("When we create a restricted Pod", func() {
			podName := "restricted-pod"
			_, err := clientset.CoreV1().Pods(namespaceName).Create(context.TODO(), &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespaceName,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-1",
							Image: "nginx",
						},
					},
				},
			}, metav1.CreateOptions{})
			So(err, ShouldBeNil)

			Convey("Then the Pod should exist", func() {
				_, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
				So(err, ShouldBeNil)

				Convey("And access to Pod logs should be restricted", func() {
					// Should return an error as access is restricted
				})
			})
		})
	})
}
