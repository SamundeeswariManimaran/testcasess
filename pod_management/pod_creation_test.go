package main

import (
	"testing"

	"context"
	"os"
	"path/filepath"

	"github.com/smartystreets/goconvey/convey"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfig() (*rest.Config, error) {
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

func TestCreatePodWithImage(t *testing.T) {
	// Create a Kubernetes client configuration
	config, err := createKubeConfig()
	if err != nil {
		t.Fatalf("Error creating in-cluster config: %v", err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// Define the namespace and pod configuration
	namespace := "default"
	podName := "test-pod"
	containerImage := "nginx:latest"

	// Create the pod configuration
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx-container",
					Image: containerImage,
				},
			},
		},
	}

	// Create the pod in the Kubernetes cluster
	_, err = clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating pod: %v", err)
	}

	// Use GoConvey to assert that the pod exists
	convey.Convey("Create Pod with Specific Image", t, func() {
		convey.So(checkPodExists(clientset, namespace, podName), convey.ShouldBeTrue)
	})
}

func checkPodExists(clientset *kubernetes.Clientset, namespace, podName string) bool {
	_, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	return err == nil
}
