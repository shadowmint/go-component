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
	layer1.CanSet = false

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

func TestChainedObjectStorage(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		store, _, _, _ := objectStorageFixture()

		obj := component.NewObject("Foo1")
		err := store.SetObject("Foo1", obj)
		T.Assert(err == nil)

		// Create fixture object
		// Insert a record that should go into layer 1
		// Get a record from layer 2
		// Insert a record that should go into layer 3
		// Get a record from layer 3
	})
}