package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// MockService represents a mock service that emulates the behavior of the real service
type MockServices struct{}

// Implement a method to emulate the behavior of the service
func (s *MockServices) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Simulate the service's response
	fmt.Fprint(w, "Service Response")
}

func TestPodCommunicationWithService(t *testing.T) {
	Convey("Given a pod and a ClusterIP service", t, func() {
		// Create a mock HTTP server to emulate the service
		mockService := &MockServices{}
		server := httptest.NewServer(mockService)
		defer server.Close()

		// Replace "http://your-service-clusterip:port/path" with the actual ClusterIP
		// and path you want to test
		serviceURL := server.URL + "/path"

		Convey("When the pod communicates with the service", func() {
			response, err := http.Get(serviceURL)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should return a successful response", func() {
				So(response.StatusCode, ShouldEqual, http.StatusOK)
			})

			Convey("The response body should contain the expected message", func() {
				body, _ := ioutil.ReadAll(response.Body)
				So(string(body), ShouldEqual, "Service Response")
			})
		})
	})
}
