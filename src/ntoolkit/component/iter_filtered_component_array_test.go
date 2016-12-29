package component_test

import (
	"ntoolkit/assert"
	"testing"
	"ntoolkit/component"
	"ntoolkit/iter"
	"reflect"
)

func TestGetComponents(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		obj := component.NewObject("Object 1")
		obj.AddComponent(&FakeComponent{Id: "1"})

		ci, err := iter.Collect(obj.GetComponents(reflect.TypeOf((*FakeComponent)(nil))))
		T.Assert(err == nil)
		T.Assert(len(ci) == 1)
		T.Assert(ci[0].(*FakeComponent).Id == "1")
	})
}

func TestGetComponentsInChildren(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		obj := component.NewObject()
		obj2 := component.NewObject()
		obj3 := component.NewObject()
		obj.AddObject(obj2)
		obj2.AddObject(obj3)

		obj3.AddComponent(&FakeComponent{Id: "1"})
		obj3.AddComponent(&FakeComponent{Id: "2"})

		ci, err := iter.Collect(obj.GetComponents(reflect.TypeOf((*FakeComponent)(nil))))
		T.Assert(err == nil)
		T.Assert(len(ci) == 0)

		ci, err = iter.Collect(obj.GetComponentsInChildren(reflect.TypeOf((*FakeComponent)(nil))))
		T.Assert(err == nil)
		T.Assert(len(ci) == 2)
		T.Assert(ci[0].(*FakeComponent).Id == "1")
		T.Assert(ci[1].(*FakeComponent).Id == "2")
	})
}