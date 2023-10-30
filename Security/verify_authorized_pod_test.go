package main

import (
	"os"
	"path/filepath"
	"testing"

	"context"

	. "github.com/smartystreets/goconvey/convey"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigverify() (*rest.Config, error) {
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

func TestAccessSecretWithDefaultServiceAccount(t *testing.T) {
	// Create a Kubernetes client using the default configuration
	config, err := createKubeConfigverify()
	if err != nil {
		t.Fatalf("Error creating Kubernetes client config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	Convey("Access Secret with Default Service Account", t, func() {
		namespace := "default"
		secretName := "my-secret"
		podName := "pod-with-default-sa"

		Convey("Create a Secret with sensitive data", func() {
			_, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: secretName},
				Data:       map[string][]byte{"password": []byte("sensitive-data")},
			}, metav1.CreateOptions{})
			So(err, ShouldBeNil)

			Convey("Create a Pod with the default Service Account", func() {
				_, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: podName},
					Spec: corev1.PodSpec{
						// Default Service Account is used
						Containers: []corev1.Container{
							{
								Name:  "test-container",
								Image: "nginx",
								// Add container specs here
							},
						},
					},
				}, metav1.CreateOptions{})
				So(err, ShouldBeNil)

				Convey("Verify that the Pod can access the Secret", func() {
					// Add assertions to check that the Pod can access the Secret
					// You may need to check environment variables or mounted volumes for the Secret
				})
			})
		})
	})
}
