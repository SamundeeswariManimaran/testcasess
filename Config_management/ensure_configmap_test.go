package main

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestConfigMapAppliedToPod(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", "C:/Users/SamundeeswariManimar/.kube/config")
	if err != nil {
		t.Fatalf("Error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating clientset: %v", err)
	}

	convey.Convey("Given a ConfigMap and a Pod", t, func() {
		configMapName := "myy-configmap"
		podName := "myy-pod"
		namespace := "default"

		// Create a ConfigMap
		convey.Convey("Create a ConfigMap", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: configMapName,
				},
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			}
			_, err := clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
			convey.So(err, convey.ShouldBeNil)
		})

		// Create a Pod
		convey.Convey("When the Pod is created", func() {
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: podName,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "my-container",
							Image: "nginx",
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "config-volume",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: configMapName,
									},
								},
							},
						},
					},
				},
			}
			_, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
			convey.So(err, convey.ShouldBeNil)

			// Verify that the Pod has the ConfigMap applied
			retryErr := wait.Poll(2*time.Second, 60*time.Second, func() (bool, error) {
				pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
				if err != nil {
					return false, err
				}
				// Check if the ConfigMap is applied to the Pod
				if len(pod.Spec.Volumes) > 0 && pod.Spec.Volumes[0].Name == "config-volume" {
					return true, nil
				}
				return false, nil
			})
			convey.So(retryErr, convey.ShouldBeNil)
		})

	})
}
