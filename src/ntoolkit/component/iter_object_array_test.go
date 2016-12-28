package component_test

import (
	"ntoolkit/assert"
	"testing"
	"ntoolkit/component"
	"ntoolkit/iter"
)

func TestRecursiveIterator(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		obj := component.NewObject()
		obj2 := component.NewObject()
		obj3 := component.NewObject()
		obj.AddObject(obj2)
		obj2.AddObject(obj3)

		results, err := iter.Collect(obj.Objects())
		T.Assert(err == nil)
		T.Assert(len(results) == 2)
	})
}
