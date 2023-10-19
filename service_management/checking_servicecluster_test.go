package main

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestServiceExistence(t *testing.T) {
	Convey("Given a Kubernetes cluster configuration", t, func() {
		kubeconfig := "C:/Users/SamundeeswariManimar/.kube/config" // Replace with the path to your kubeconfig file

		// Load the Kubernetes configuration from the kubeconfig file
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		So(err, ShouldBeNil)

		// Create a Kubernetes clientset
		clientset, err := kubernetes.NewForConfig(config)
		So(err, ShouldBeNil)

		Convey("When checking the existence of a service in the cluster", func() {
			serviceName := "kubernetes" // Replace with the actual service name
			namespace := "default"      // Replace with the actual namespace

			// Use the clientset to retrieve the service
			//_, err := clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
			_, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
			if err != nil {
				t.Fatalf("Error getting service: %v", err)
			}

			Convey("It should not return an error, indicating the service exists", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
