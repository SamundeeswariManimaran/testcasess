package main

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubernetesClient() (*kubernetes.Clientset, error) {
	kubeconfig := "C:/Users/SamundeeswariManimar/.kube/config" // Update with your kubeconfig path
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

func createConfigMapup(clientset *kubernetes.Clientset, name string, data map[string]string) error {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	_, err := clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), cm, metav1.CreateOptions{})
	return err
}

func getConfigMapup(clientset *kubernetes.Clientset, name string) (*v1.ConfigMap, error) {
	return clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), name, metav1.GetOptions{})
}

func updateConfigMap(clientset *kubernetes.Clientset, name string, data map[string]string) error {
	cm, err := getConfigMapup(clientset, name)
	if err != nil {
		return err
	}
	cm.Data = data
	_, err = clientset.CoreV1().ConfigMaps("default").Update(context.TODO(), cm, metav1.UpdateOptions{})
	return err
}

func TestUpdateConfigMap(t *testing.T) {
	clientset, err := getKubernetesClient()
	if err != nil {
		t.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	Convey("Update a ConfigMap in Kubernetes", t, func() {
		cmName := "new-configmap"
		initialData := map[string]string{"key1": "1", "key2": "2"}
		updatedData := map[string]string{"key1": "3", "key2": "4"}

		Convey("Create the initial ConfigMap", func() {
			So(createConfigMapup(clientset, cmName, initialData), ShouldBeNil)

			Convey("Update the ConfigMap", func() {
				So(updateConfigMap(clientset, cmName, updatedData), ShouldBeNil)

				Convey("Check if the ConfigMap was updated", func() {
					cm, err := getConfigMapup(clientset, cmName)
					So(err, ShouldBeNil)
					So(cm.Name, ShouldEqual, cmName)
					So(cm.Data, ShouldResemble, updatedData)
				})
			})
		})
	})
}
