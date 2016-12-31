package component_test

import (
	"ntoolkit/assert"
	"testing"
	"ntoolkit/component"
)

func TestFindObject(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		o1 := component.NewObject("A")
		o2 := component.NewObject("B")
		o3 := component.NewObject("C")
		o4 := component.NewObject("D1")
		o5 := component.NewObject("D2")

		o1.AddObject(component.NewObject())
		o1.AddObject(o2)
		o1.AddObject(component.NewObject())

		o2.AddObject(o3)

		o3.AddObject(o4)
		o3.AddObject(component.NewObject())
		o3.AddObject(o5)

		r, err := o1.FindObject("B")
		T.Assert(err == nil)
		T.Assert(r == o2)

		r, err = o1.FindObject("B", "C")
		T.Assert(err == nil)
		T.Assert(r == o3)

		r, err = o1.FindObject("B", "C", "D1")
		T.Assert(err == nil)
		T.Assert(r == o4)

		r, err = o1.FindObject("B", "C", "D2")
		T.Assert(err == nil)
		T.Assert(r == o5)
	})
}

func TestFindComponent(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		o1 := component.NewObject("A")
		o2 := component.NewObject("B")
		o3 := component.NewObject("C")
		o4 := component.NewObject("D")
		c1 := &FakeComponent{Id:"Component"}

		o1.AddObject(o2)
		o2.AddObject(o3)
		o3.AddObject(o4)
		o4.AddComponent(c1)

		var c2 *FakeComponent
		err := o1.Find(&c2, "B", "C", "D")

		T.Assert(err == nil)
		T.Assert(c2.Id == "Component")
	})
}