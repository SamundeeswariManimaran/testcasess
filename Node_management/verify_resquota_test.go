package main

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func setupClient() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func createNamespaces(clientset *kubernetes.Clientset, namespace string) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func createResourceQuota(clientset *kubernetes.Clientset, namespace, cpuLimit, memoryLimit string) {
	resourceQuota := &v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name: "resource-quota",
		},
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceLimitsCPU:    resource.MustParse(cpuLimit),
				v1.ResourceLimitsMemory: resource.MustParse(memoryLimit),
			},
		},
	}
	_, err := clientset.CoreV1().ResourceQuotas(namespace).Create(context.TODO(), resourceQuota, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func TestResourceQuotasEnforcement(t *testing.T) {
	clientset, err := setupClient()
	if err != nil {
		t.Fatalf("Error setting up Kubernetes client: %v", err)
	}

	// Define test cases for resource quota enforcement.
	convey.Convey("Testing resource quota enforcement in a namespace", t, func() {
		namespace := "node-namespace"
		cpuLimit := "200m"
		memoryLimit := "200Mi"

		convey.Convey("Creating a Pod should succeed within quota limits", func() {
			createNamespaces(clientset, namespace)
			createResourceQuota(clientset, namespace, cpuLimit, memoryLimit)

			// Create a Pod that should fit within the resource quota.
			// You need to implement the Pod creation logic based on your cluster setup.

			// Check that the Pod exists in the namespace.

			// Clean up the created Pod, ResourceQuota, and namespace.
		})

		convey.Convey("Creating a Pod should fail when exceeding quota limits", func() {
			createNamespaces(clientset, namespace)
			createResourceQuota(clientset, namespace, cpuLimit, memoryLimit)

			// Create a Pod that exceeds the resource quota.
			// You need to implement the Pod creation logic based on your cluster setup.

			// Check that the Pod does not exist in the namespace.

			// Clean up the created ResourceQuota and namespace.
		})
	})
}
