package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Node represents the node in the system.
type Node struct {
	IsDeleted bool
}

// System represents the system that manages nodes.
type System struct {
	Nodes map[string]*Node
}

// CancelNodeDeletion cancels the deletion of a node in the system.
func (s *System) CancelNodeDeletion(nodeID string) error {
	node, exists := s.Nodes[nodeID]
	if !exists {
		return nil // Node not found, return success.
	}

	if node.IsDeleted {
		node.IsDeleted = false
	}

	return nil
}

func TestCancelNodeDeletion(t *testing.T) {
	// Initialize the test system and node.
	sys := &System{
		Nodes: map[string]*Node{
			"node1": {IsDeleted: true},
			"node2": {IsDeleted: false},
		},
	}

	Convey("Given a system with nodes", t, func() {
		Convey("When a node deletion is canceled", func() {
			nodeID := "node1"
			err := sys.CancelNodeDeletion(nodeID)

			Convey("The system should cancel the node deletion", func() {
				So(err, ShouldBeNil)
				node, exists := sys.Nodes[nodeID]
				So(exists, ShouldBeTrue)
				So(node.IsDeleted, ShouldBeFalse)
			})

			Convey("When canceling a non-existent node deletion", func() {
				nodeID := "nonexistent"
				err := sys.CancelNodeDeletion(nodeID)

				Convey("The system should not return an error and do nothing", func() {
					So(err, ShouldBeNil)
					_, exists := sys.Nodes[nodeID]
					So(exists, ShouldBeFalse)
				})
			})
		})
	})
}
