package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigSecret() (*rest.Config, error) {
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

// Function to apply a secret to a pod
func applySecretToPod(pod *v1.Pod, secret *v1.Secret) {
	pod.Spec.Volumes = append(pod.Spec.Volumes, v1.Volume{
		Name: "my-secret-volume",
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: secret.Name,
			},
		},
	})
}

func TestApplySecretToPod(t *testing.T) {
	convey.Convey("Given a Kubernetes pod and a Secret", t, func() {
		// Load in-cluster config
		config, err := createKubeConfigSecret()
		if err != nil {
			t.Fatal(err)
		}

		// Create a Kubernetes clientset using the in-cluster config
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			t.Fatal(err)
		}

		// Create a pod
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "se-pod"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "my-container",
						Image: "nginx",
					},
				},
			},
		}

		// Create a secret
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "my-secret"},
			Data: map[string][]byte{
				"username": []byte("my-username"),
				"password": []byte("my-password"),
			},
		}

		// Create the pod in the cluster
		_, err = clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {
			t.Fatal(err)
		}

		// Apply the secret to the pod
		applySecretToPod(pod, secret)

		// convey.Convey("When applying the Secret to the pod", func() {
		// 	// Update the pod in the cluster with the new configuration
		// 	_, err = clientset.CoreV1().Pods("default").Update(context.TODO(), pod, metav1.UpdateOptions{})
		// 	if err != nil {
		// 		t.Fatal(err)
		// 	}

		// convey.Convey("The pod should have the Secret correctly applied", func() {
		// 	// Verify that the pod's volume has the correct SecretName
		// 	updatedPod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), "secrett-pod", metav1.GetOptions{})
		// 	if err != nil {
		// 		t.Fatal(err)
		// 	}

		// 	volume := updatedPod.Spec.Volumes[0].VolumeSource.Secret
		// 	convey.So(volume.SecretName, convey.ShouldEqual, secret.Name)

		// 	// You can add more assertions to test other aspects of the pod's configuration
		// })
	})

	// // Clean up the resources by deleting the pod and the secret (not shown in this example)
	// err = clientset.CoreV1().Pods("default").Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// 	})
}
