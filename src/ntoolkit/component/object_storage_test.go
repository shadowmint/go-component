package component_test

import (
	"ntoolkit/assert"
	"testing"
	"ntoolkit/component"
)

func objectStorageFixture() (*component.ObjectStorage, *component.ObjectStorageMemory, *component.ObjectStorageMemory, *component.ObjectStorageMemory) {
	layer1 := component.NewObjectStorageMemory()
	layer1.Pattern = "^Foo.*"
	layer1.CanGet = false

	layer2 := component.NewObjectStorageMemory()
	layer2.Pattern = "^Bar.*"
	layer2.CanSet = false

	layer3 := component.NewObjectStorageMemory()

	stack := component.NewObjectStorageStack()
	stack.Add(layer1)
	stack.Add(layer2)
	stack.Add(layer3)

	store := component.NewObjectStorage(component.NewObjectFactory(), stack)
	return store, layer1, layer2, layer3
}

func TestCreateObjectStorage(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		store, _, _, _ := objectStorageFixture()
		T.Assert(store != nil)
	})
}

func TestChainedObjectSet(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		store, l1, _, l3 := objectStorageFixture()

		obj := component.NewObject("Foo1")
		err := store.SetObject("Foo1", obj)
		T.Assert(err == nil)

		obj = component.NewObject("Bar1")
		err = store.SetObject("Bar1", obj)
		T.Assert(err == nil)

		_, err = l1.Get("Foo1")
		T.Assert(err == nil)

		_, err = l3.Get("Bar1")
		T.Assert(err == nil)

		_, err = l3.Get("Other")
		T.Assert(err != nil)
	})
}

func TestChainedObjectGet(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		store, _, l2, l3 := objectStorageFixture()

		factory := component.NewObjectFactory()
		obj := component.NewObject("Bar1")
		objtmp, _ := factory.Serialize(obj)
		l2.Set("Bar1", objtmp)

		obj = component.NewObject("Other")
		objtmp, _ = factory.Serialize(obj)
		l3.Set("Other", objtmp)

		_, err := store.GetObject("Bar1")
		T.Assert(err == nil)

		_, err = store.GetObject("Other")
		T.Assert(err == nil)

		_, err = store.GetObject("Other2323")
		T.Assert(err != nil)
	})
}