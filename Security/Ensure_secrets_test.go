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

func createKubeConfigverifyensuresec() (*rest.Config, error) {
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
func TestSecretSecurity(t *testing.T) {
	// Create a Kubernetes client using the default configuration
	config, err := createKubeConfigverifyensuresec()
	if err != nil {
		t.Fatalf("Error creating Kubernetes client config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	Convey("Secret Security Test", t, func() {
		namespace := "default"
		secretName := "demo-secure"

		Convey("Create a Secret with sensitive data", func() {
			_, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: secretName},
				Data:       map[string][]byte{"password": []byte("sensitive-data")},
			}, metav1.CreateOptions{})
			So(err, ShouldBeNil)

			Convey("Check that Secret is not leaked in Pod logs", func() {
				podName := "pod-demo-secure"
				_, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: podName},
					Spec: corev1.PodSpec{
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

				// Verify that the Pod logs do not contain sensitive information
				// You can fetch Pod logs and check for sensitive data

				// Example:
				// logs, err := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{}).Do()
				// So(err, ShouldBeNil)
				// So(logs, ShouldNotContainSubstring, "sensitive-data")
			})
		})
	})
}
