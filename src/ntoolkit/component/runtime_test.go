package component_test

import (
	"ntoolkit/assert"
	"testing"
	"fmt"
	"ntoolkit/component"
	"reflect"
)

type FakeComponent struct {
}

func (fake FakeComponent) Type() reflect.Type {
	return reflect.TypeOf(fake)
}

func (fake FakeComponent) Update(_ *component.Context) {
	fmt.Printf("Update component\n")
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
		obj.AddComponent(FakeComponent{})
		obj.AddComponent(FakeComponent{})
		obj.AddComponent(FakeComponent{})

		runtime.Root().AddObject(obj)

		runtime.Update(1.0)
	})
}

/*func TestRoom(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		room := room.New()
		room.GameObject.AddComponent(&FakeComponent{})
		runtime := runtime.New(runtime.Config{ThreadPoolSize: 10})
		runtime.AddRoom(room)
		runtime.Step(0.5)
	})
}*/