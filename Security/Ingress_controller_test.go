package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type IngressController struct {
	IngressRules map[string]string
}

func NewIngressController() *IngressController {
	return &IngressController{
		IngressRules: make(map[string]string),
	}
}

func (ic *IngressController) AddIngressRule(serviceName, path string) {
	ic.IngressRules[serviceName] = path
}

func (ic *IngressController) GetIngressPath(serviceName string) string {
	return ic.IngressRules[serviceName]
}

func TestIngressController(t *testing.T) {
	Convey("Given an IngressController managing ingress rules", t, func() {
		ingressController := NewIngressController()

		Convey("When adding ingress rules for services", func() {
			serviceName := "kubernetes"
			ingressPath := "/my-service"

			ingressController.AddIngressRule(serviceName, ingressPath)

			Convey("Then the IngressController should return the correct ingress path for a service", func() {
				path := ingressController.GetIngressPath(serviceName)
				So(path, ShouldEqual, ingressPath)
			})

			Convey("And the IngressController should return an empty string for an unknown service", func() {
				unknownService := "unknown-service"
				path := ingressController.GetIngressPath(unknownService)
				So(path, ShouldBeEmpty)
			})
		})
	})
}
