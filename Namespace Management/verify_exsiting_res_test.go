package main

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfigverifyres() (*rest.Config, error) {
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
func TestUpdateWithoutDisruption(t *testing.T) {
	Convey("Given an existing map of data", t, func() {
		existingData := map[string]string{"key1": "value1", "key2": "value2"}
		newData := map[string]string{"key2": "updatedValue", "key3": "newValue"}

		Convey("When we update the map", func() {
			// Simulate updating the map by merging new data into existing data
			for key, value := range newData {
				existingData[key] = value
			}

			Convey("Then the map should have the updated values", func() {
				So(existingData["key2"], ShouldEqual, "updatedValue")
				So(existingData["key3"], ShouldEqual, "newValue")
			})

			Convey("And it should not disrupt the existing data", func() {
				So(existingData["key1"], ShouldEqual, "value1")
			})
		})
	})
}
