package component_test

import (
	"ntoolkit/assert"
	"testing"
	"ntoolkit/component"
)

func objectStorageFixture() (*component.ObjectStorage, *component.ObjectFactory, *component.ObjectStorageMemory, *component.ObjectStorageMemory, *component.ObjectStorageMemory) {
	layer1 := component.NewObjectStorageMemory()
	layer1.Pattern = "^Foo.*"
	layer1.CanGet = false

	layer2 := component.NewObjectStorageMemory()
	layer2.Pattern = "^Bar.*"
	layer2.CanSet = false

	layer3 := component.NewObjectStorageMemory()

	layer4 := component.NewObjectStorageMemory()

	stack := component.NewObjectStorageStack()
	stack.Add(layer1)
	stack.Add(layer2)
	stack.Add(layer3)

	factory := component.NewObjectFactory()
	store := component.NewObjectStorage(factory, stack, layer4)
	return store, factory, layer1, layer2, layer3
}

func TestCreateObjectStorage(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		store, _, _, _, _ := objectStorageFixture()
		T.Assert(store != nil)
	})
}

func TestChainedObjectSet(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		store, f, l1, _, l3 := objectStorageFixture()

		obj := component.NewObject("Foo1")
		err := store.Add("Foo1", obj)
		T.Assert(err == nil)

		obj = component.NewObject("Bar1")
		err = store.Add("Bar1", obj)
		T.Assert(err == nil)

		_, err = l1.Get("Foo1", f)
		T.Assert(err == nil)

		_, err = l3.Get("Bar1", f)
		T.Assert(err == nil)

		_, err = l3.Get("Other", f)
		T.Assert(err != nil)
	})
}

func TestChainedObjectGet(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		store, f, _, l2, l3 := objectStorageFixture()

		obj := component.NewObject("Bar1")
		l2.Set("Bar1", obj, f)

		obj = component.NewObject("Other")
		l3.Set("Other", obj, f)

		_, err := store.Get("Bar1")
		T.Assert(err == nil)

		_, err = store.Get("Other")
		T.Assert(err == nil)

		_, err = store.Get("Other2323")
		T.Assert(err != nil)
	})
}


