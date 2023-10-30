package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type NetworkPolicy struct {
	AllowedPods map[string]struct{}
}

func NewNetworkPolicy(allowedPods []string) *NetworkPolicy {
	policy := &NetworkPolicy{
		AllowedPods: make(map[string]struct{}),
	}
	for _, pod := range allowedPods {
		policy.AllowedPods[pod] = struct{}{}
	}
	return policy
}

func (np *NetworkPolicy) AllowsTraffic(sourcePod, destinationPod string) bool {
	_, allowed := np.AllowedPods[destinationPod]
	return allowed
}

func TestNetworkPolicy(t *testing.T) {
	Convey("Given a NetworkPolicy that allows specific ingress traffic", t, func() {
		allowedPods := []string{"pod-1", "pod-2"}
		networkPolicy := NewNetworkPolicy(allowedPods)

		Convey("When a pod tries to connect", func() {
			sourcePod := "pod-3"      // This pod does not match the NetworkPolicy
			destinationPod := "pod-1" // This pod matches the NetworkPolicy

			Convey("Then the policy should allow traffic to allowed pods", func() {
				allowed := networkPolicy.AllowsTraffic(sourcePod, destinationPod)
				So(allowed, ShouldBeTrue)
			})

			Convey("And the policy should deny traffic to non-allowed pods", func() {
				nonAllowedPod := "pod-4" // This pod does not match the NetworkPolicy
				allowed := networkPolicy.AllowsTraffic(sourcePod, nonAllowedPod)
				So(allowed, ShouldBeFalse)
			})
		})
	})
}
