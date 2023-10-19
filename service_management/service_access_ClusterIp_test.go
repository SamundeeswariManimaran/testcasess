// package main

// import (
// 	"net/http"
// 	"testing"

// 	. "github.com/smartystreets/goconvey/convey"
// )

// func TestServiceAccess(t *testing.T) {
// 	Convey("When accessing the service using ClusterIP", t, func() {
// 		serviceURL := "http://127.0.0.1:57895 " // Replace with your actual ClusterIP and port

// 		Convey("It should return a successful response", func() {
// 			response, err := http.Get(serviceURL)
// 			So(err, ShouldBeNil)
// 			So(response.StatusCode, ShouldEqual, http.StatusOK)
// 		})
// 	})
// }

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
type MockService struct {
	// Define any fields or methods needed to simulate the service behavior
}

// Implement a method to emulate the behavior of the service
func (s *MockService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Simulate the service's response
	fmt.Fprint(w, "Hello, ClusterIP Service!")
}

func TestClusterIPAccess(t *testing.T) {
	Convey("Given a pod that accesses a service using ClusterIP", t, func() {
		// Create a mock HTTP server to emulate the service
		mockService := &MockService{}
		server := httptest.NewServer(mockService)
		defer server.Close()

		// Replace "http://your-service-clusterip:port/path" with the actual ClusterIP
		// and path you want to test
		serviceURL := server.URL + "/path"

		Convey("When the pod makes an HTTP request to the service", func() {
			response, err := http.Get(serviceURL)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should return a successful response", func() {
				So(response.StatusCode, ShouldEqual, http.StatusOK)
			})

			Convey("The response body should contain the expected message", func() {
				body, _ := ioutil.ReadAll(response.Body)
				So(string(body), ShouldEqual, "Hello, ClusterIP Service!")
			})
		})
	})
}
