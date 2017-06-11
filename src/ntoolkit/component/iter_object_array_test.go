package component_test

import (
	"ntoolkit/assert"
	"testing"
	"ntoolkit/component"
	"ntoolkit/iter"
)

func TestSingleChildIterator(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		obj := component.NewObject("Object 1")
		obj2 := component.NewObject("Object 2")
		obj.AddObject(obj2)

		results, err := iter.Collect(obj.ObjectsInChildren())
		T.Assert(err == nil)
		T.Assert(len(results) == 1)
		T.Assert(results[0] == obj2)
	})
}

func TestDepth3ChildIterator(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		obj := component.NewObject()
		obj2 := component.NewObject()
		obj3 := component.NewObject()
		obj.AddObject(obj2)
		obj2.AddObject(obj3)

		results, err := iter.Collect(obj.ObjectsInChildren())
		T.Assert(err == nil)
		T.Assert(len(results) == 2)
		T.Assert(results[0] == obj2)
		T.Assert(results[1] == obj3)
	})
}