package component_test

import (
	"ntoolkit/assert"
	"testing"
	"ntoolkit/component"
	"reflect"
	"ntoolkit/iter"
)

type FakeComponent struct {
	Id    string
	Count int
}

func (fake *FakeComponent) Type() reflect.Type {
	return reflect.TypeOf(fake)
}

func (fake *FakeComponent) Update(_ *component.Context) {
	fake.Count += 1
}

func (fake *FakeComponent) New() component.Component {
	return &FakeComponent{}
}

func TestNew(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		component.NewRuntime(component.Config{
			ThreadPoolSize: 3})
	})
}

func TestUpdate(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		runtime := component.NewRuntime(component.Config{
			ThreadPoolSize: 50})

		obj := component.NewObject()
		obj.AddComponent(&FakeComponent{})
		obj.AddComponent(&FakeComponent{})
		obj.AddComponent(&FakeComponent{})

		count, err := iter.Count(obj.GetComponents(reflect.TypeOf((*FakeComponent)(nil))))
		T.Assert(err == nil)
		T.Assert(count == 3)

		components := obj.GetComponents(reflect.TypeOf((*FakeComponent)(nil)))
		for val, err := components.Next(); err == nil; val, err = components.Next() {
			T.Assert(val.(*FakeComponent).Count == 0)
		}

		runtime.Root().AddObject(obj)

		runtime.Update(1.0)
		runtime.Update(1.0)

		components = obj.GetComponents(reflect.TypeOf((*FakeComponent)(nil)))
		for val, err := components.Next(); err == nil; val, err = components.Next() {
			T.Assert(val.(*FakeComponent).Count == 2)
		}
	})
}

func TestComponentsAreUpdated(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		runtime := component.NewRuntime(component.Config{
			ThreadPoolSize: 50})

		obj := component.NewObject()
		obj.AddComponent(&FakeComponent{})
		obj.AddComponent(&FakeComponent{})
		obj.AddComponent(&FakeComponent{})

		root := component.NewObject()
		root.AddObject(obj)

		count, _ := iter.Count(root.Objects())
		T.Assert(count == 2)

		count, err := iter.Count(root.GetComponentsInChildren(reflect.TypeOf((*FakeComponent)(nil))))
		T.Assert(err == nil)
		T.Assert(count == 3)

		runtime.Update(1.0)
		runtime.Update(1.0)

		components := root.GetComponentsInChildren(reflect.TypeOf((*FakeComponent)(nil)).Elem())
		for val, err := components.Next(); err == nil; val, err = components.Next() {
			T.Assert(val.(*FakeComponent).Count == 2)
		}
	})
}