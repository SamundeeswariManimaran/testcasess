package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Policy is a simple map that defines which actions are allowed.
var Policy = map[string]bool{
	"read":   true,
	"write":  true,
	"delete": false, // Denied by default
}

// IsActionAllowed checks whether an action is allowed according to the policy.
func IsActionAllowed(action string) bool {
	allowed, ok := Policy[action]
	return ok && allowed
}

func TestDefaultDenyPolicy(t *testing.T) {
	Convey("Given a system with a default-deny policy", t, func() {
		Convey("When a user tries to perform actions", func() {
			Convey("Then unauthorized actions should be denied", func() {
				// Simulate an unauthorized action
				unauthorizedAction := "delete"
				allowed := IsActionAllowed(unauthorizedAction)

				Convey("So, the action should not be allowed", func() {
					So(allowed, ShouldBeFalse)
				})
			})

			Convey("Then authorized actions should be allowed", func() {
				// Simulate an authorized action
				authorizedAction := "read"
				allowed := IsActionAllowed(authorizedAction)

				Convey("So, the action should be allowed", func() {
					So(allowed, ShouldBeTrue)
				})
			})

			Convey("Then actions not in the policy should be denied", func() {
				// Simulate an action that is not defined in the policy
				unknownAction := "unknown"
				allowed := IsActionAllowed(unknownAction)

				Convey("So, the action should be denied by default", func() {
					So(allowed, ShouldBeFalse)
				})
			})
		})
	})
}
