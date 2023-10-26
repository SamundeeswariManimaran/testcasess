package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigenvironment() (*rest.Config, error) {
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

func TestPodWithEnvironmentalVariable(t *testing.T) {
	config, err := createKubeConfigenvironment() // Use in-cluster config if applicable
	if err != nil {
		t.Fatalf("Error creating Kubernetes in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Define your pod specifications, including the environmental variable.
	podName := "env-pod"
	namespace := "default" // Replace with your target namespace
	envVarName := "podenv"
	envVarValue := "1"

	podSpec := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test-container",
					Image: "nginx",
					Env: []corev1.EnvVar{
						{
							Name:  envVarName,
							Value: envVarValue,
						},
					},
				},
			},
		},
	}

	// Create the pod.
	createpod, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), podSpec, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating pod: %v", err)
	}

	// Ensure the pod is running.
	err = wait.PollImmediate(1*time.Second, 30*time.Second, func() (bool, error) {
		pod, getErr := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if getErr != nil {
			return false, getErr
		}
		return pod.Status.Phase == corev1.PodRunning, nil
	})

	if err != nil {
		t.Fatalf("Error waiting for pod to become running: %v", err)
	}

	// Run the test cases with GoConvey.
	convey.Convey("Testing pod with environmental variable", t, func() {
		convey.Convey("The environmental variable should have the expected value", func() {
			// Retrieve the pod to ensure it's running.
			pod, getErr := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
			if getErr != nil {
				t.Fatalf("Error getting pod: %v", getErr)
			}

			// Find the environmental variable in the pod's containers.
			var actualValue string
			for _, container := range pod.Spec.Containers {
				for _, env := range container.Env {
					if env.Name == envVarName {
						actualValue = env.Value
						break
					}
				}
			}

			// Set your expected value.
			expectedValue := envVarValue

			// Define your assertion.
			convey.So(actualValue, convey.ShouldEqual, expectedValue)
			fmt.Println(createpod)
		})
	})
}
