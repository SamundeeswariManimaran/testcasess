package main

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const kubeconfigPath = "C:/Users/SamundeeswariManimar/.kube/config"

// deploySampleResource deploys a sample Kubernetes resource
func deploySampleResource(clientset *kubernetes.Clientset, namespace string) error {
	// For example, deploy a ConfigMap
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pipe-configmap",
			Namespace: namespace,
		},
		Data: map[string]string{
			"key": "value",
		},
	}

	_, err := clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	return err
}

func TestDeployAndManageK8sResources(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		t.Fatalf("Error building Kubernetes config: %v", err)
	}

	// For testing only, skip TLS verification
	config.Insecure = true

	convey.Convey("Given a Kubernetes cluster", t, func() {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		convey.So(err, convey.ShouldBeNil)

		clientset, err := kubernetes.NewForConfig(config)
		convey.So(err, convey.ShouldBeNil)

		namespace := "default"

		convey.Convey("When deploying a sample resource", func() {
			err := deploySampleResource(clientset, namespace)

			convey.Convey("Then the deployment should succeed", func() {
				convey.So(err, convey.ShouldBeNil)

				convey.Convey("When checking the existence of the deployed resource", func() {
					// You can check the existence of the deployed resource, e.g., ConfigMap
					_, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), "pipeline-configmap", metav1.GetOptions{})

					convey.Convey("Then the resource should exist", func() {
						convey.So(err, convey.ShouldBeNil)

					})
				})
			})
		})
	})
}
