package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigdeploy() (*rest.Config, error) {
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

func deployResourcesInNamespace(clientset *kubernetes.Clientset, namespaceName string) error {
	// Deploy your Kubernetes resources here.
	// Example: Deploy a simple Nginx pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: namespaceName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "nginx-container",
					Image:           "nginx:latest",
					ImagePullPolicy: corev1.PullIfNotPresent,
				},
			},
		},
	}

	_, err := clientset.CoreV1().Pods(namespaceName).Create(context.TODO(), pod, metav1.CreateOptions{})
	return err
}

func TestDeployResourcesInNamespace(t *testing.T) {
	// Load Kubernetes configuration from a file (you can customize this)
	config, err := clientcmd.BuildConfigFromFlags("", "C:/Users/SamundeeswariManimar/.kube/config")
	if err != nil {
		config, err = createKubeConfigdeploy()
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	Convey("Test DeployResourcesInNamespace", t, func() {
		namespaceName := "new-namespace"

		// Create a new namespace
		newNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: namespaceName},
		}
		_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), newNamespace, metav1.CreateOptions{})
		So(err, ShouldBeNil)

		// Deploy your resources in the new namespace
		err = deployResourcesInNamespace(clientset, namespaceName)
		So(err, ShouldBeNil)

		// Verify that the resources have been deployed successfully
		pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), "test-pod", metav1.GetOptions{})
		So(err, ShouldBeNil)
		So(pod.Name, ShouldEqual, "test-pod")
	})
}
