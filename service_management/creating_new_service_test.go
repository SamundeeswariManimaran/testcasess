package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigservice() (*rest.Config, error) {
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
func createService(clientset *kubernetes.Clientset, serviceName, namespace string) error {
	// Define the service
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app": "my-app", // Adjust the selector to match your application labels
			},
			Ports: []v1.ServicePort{
				{
					Protocol:   v1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(8080), // Adjust the target port as needed
				},
			},
		},
	}

	// Create the service
	_, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	return err
}

func TestServiceExistenceservice(t *testing.T) {
	Convey("Given a Kubernetes cluster configuration", t, func() {
		// Load the Kubernetes configuration from in-cluster
		config, err := createKubeConfigservice()
		So(err, ShouldBeNil)

		// Create a Kubernetes clientset
		clientset, err := kubernetes.NewForConfig(config)
		So(err, ShouldBeNil)

		Convey("When checking the existence of the service in the cluster", func() {
			serviceName := "new-servicekube" // Adjust the service name if needed
			namespace := "default"           // Adjust the namespace if needed

			// Create the service
			err := createService(clientset, serviceName, namespace)
			So(err, ShouldBeNil)

			// Wait for the service to exist (may require a few seconds)
			err = waitForService(clientset, serviceName, namespace)
			So(err, ShouldBeNil)

			Convey("It should not return an error, indicating the service exists", func() {
				// Check the existence of the service
				_, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
				So(err, ShouldBeNil)
			})
		})
	})
}

func waitForService(clientset *kubernetes.Clientset, serviceName, namespace string) error {
	// Implement a function to wait for the service to exist
	// You can use a loop with retries and a timeout to check the service's existence
	// For simplicity, this example uses a single retry with a short timeout
	for i := 0; i < 5; i++ {
		_, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err == nil {
			return nil // Service exists
		}
		time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
	}
	return fmt.Errorf("service %s in namespace %s does not exist", serviceName, namespace)
}
