package component_test

import (
	"ntoolkit/assert"
	"testing"
	"ntoolkit/component"
//	"fmt"
	//"strings"
	"strings"
	"fmt"
)

func objectStorageRuntimeFixture(T *assert.T) (*component.Runtime, *component.ObjectFactory, *component.ObjectStorage, *component.ObjectStorageRuntime, *component.ObjectStorageMemory, *component.ObjectStorageMemory) {
	config := component.Config{ThreadPoolSize: 10}
	runtime := component.NewRuntime(config)

	rooms := component.NewObject("Rooms")
	runtime.Root().AddObject(rooms)

	active := component.NewObjectStorageRuntime(rooms)

	templates := component.NewObjectStorageMemory()
	templates.CanSet = false

	db := component.NewObjectStorageMemory()

	storage := component.NewObjectStorageStack()
	storage.Add(db)
	storage.Add(templates)

	factory := component.NewObjectFactory()

	// Push in a few templates to the hard coded template block
	T.Assert(templates.Set("Room1", component.NewObject("Room1"), factory) == nil)
	T.Assert(templates.Set("Room2", component.NewObject("Room2"), factory) == nil)
	T.Assert(templates.Set("Room3", component.NewObject("Room3"), factory) == nil)

	rtn := component.NewObjectStorage(factory, active, storage)
	return runtime, factory, rtn, active, db, templates
}

func TestRuntimeStorageFixture(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		objectStorageRuntimeFixture(T)
	})
}

func TestRuntimeStorageInitialState(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		_, factory, fixture, active, db, templates := objectStorageRuntimeFixture(T)

		T.Assert(fixture.Exists("Room1"))
		T.Assert(fixture.Exists("Room2"))
		T.Assert(fixture.Exists("Room3"))
		T.Assert(!fixture.Exists("Room4"))

		_, err := active.Get("Room1", factory); T.Assert(err != nil)
		_, err = db.Get("Room1", factory); T.Assert(err != nil)
		_, err = templates.Get("Room1", factory); T.Assert(err == nil)
	})
}

func TestRuntimeStorageStateAfterGet(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		_, factory, fixture, active, db, templates := objectStorageRuntimeFixture(T)

		fixture.Get("Room1")
		T.Assert(fixture.Active("Room1"))

		_, err := active.Get("Room1", factory); T.Assert(err == nil)
		_, err = db.Get("Room1", factory); T.Assert(err != nil)
		_, err = templates.Get("Room1", factory); T.Assert(err == nil)
	})
}

func TestRuntimeStorageStateAfterSave(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		_, factory, fixture, active, db, templates := objectStorageRuntimeFixture(T)

		fixture.Get("Room1")
		fixture.SetActive("Room1", false)

		_, err := active.Get("Room1", factory); T.Assert(err != nil)
		_, err = db.Get("Room1", factory); T.Assert(err == nil)
		_, err = templates.Get("Room1", factory); T.Assert(err == nil)
	})
}

func TestRuntimeStorageStateAfterDrop(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		_, factory, fixture, active, db, templates := objectStorageRuntimeFixture(T)

		fixture.Get("Room1")
		fixture.SetActive("Room1", false)
		fixture.Drop("Room1")

		_, err := active.Get("Room1", factory); T.Assert(err != nil)
		_, err = db.Get("Room1", factory); T.Assert(err != nil)
		_, err = templates.Get("Room1", factory); T.Assert(err == nil)
	})
}


func TestRuntimeStorageCanMoveBetweenStates(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		runtime, factory, fixture, _, db, templates := objectStorageRuntimeFixture(T)

		room1, err := fixture.Get("Room1")
		T.Assert(err == nil)
		T.Assert(room1 != nil)

		room2, err := fixture.Get("Room2")
		T.Assert(err == nil)
		T.Assert(room2 != nil)

		room3 := component.NewObject()
		fixture.Add("Room3", room3)

		// Adding shouldn't have saved objects
		_, err = db.Get("Room1", factory)
		T.Assert(err != nil)

		_, err = templates.Get("Room1", factory)
		T.Assert(err == nil)

		T.Assert(strings.TrimSpace(runtime.Root().Debug()) == strings.TrimSpace(`
object: Untitled (1 / 0)
   object: Rooms (3 / 0)
        object: Room1 (0 / 0)
        object: Room2 (0 / 0)
        object: Room3 (0 / 0)
`))

		T.Assert(fixture.Active("Room1"))

		fixture.SetActive("Room1", false)

		T.Assert(!fixture.Active("Room1"))
		T.Assert(fixture.Exists("Room1"))

		T.Assert(strings.TrimSpace(runtime.Root().Debug()) == strings.TrimSpace(`
object: Untitled (1 / 0)
   object: Rooms (2 / 0)
        object: Room2 (0 / 0)
        object: Room3 (0 / 0)
`))

		fixture.SetActive("Room1", true)

		T.Assert(strings.TrimSpace(runtime.Root().Debug()) == strings.TrimSpace(`
object: Untitled (1 / 0)
   object: Rooms (3 / 0)
        object: Room2 (0 / 0)
        object: Room3 (0 / 0)
        object: Room1 (0 / 0)
`))

		// Object should be saved now
		_, err = db.Get("Room1", factory)
		T.Assert(err == nil)

		_, err = templates.Get("Room1", factory)
		T.Assert(err == nil)
	})
}

func TestRuntimeStorageDrop(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		runtime, factory, fixture, _, db, _ := objectStorageRuntimeFixture(T)

		room1, err := fixture.Get("Room1")
		T.Assert(err == nil)
		T.Assert(room1 != nil)

		room2, err := fixture.Get("Room2")
		T.Assert(err == nil)
		T.Assert(room2 != nil)

		room3 := component.NewObject()
		fixture.Add("Room4", room3)

		T.Assert(strings.TrimSpace(runtime.Root().Debug()) == strings.TrimSpace(`
object: Untitled (1 / 0)
   object: Rooms (3 / 0)
        object: Room1 (0 / 0)
        object: Room2 (0 / 0)
        object: Room4 (0 / 0)
`))

		fixture.SetActive("Room1", false)

		_, err = db.Get("Room1", factory)
		T.Assert(err == nil)

		fixture.Drop("Room1")
		T.Assert(!fixture.Active("Room1"))
		T.Assert(fixture.Exists("Room1"))

		_, err = db.Get("Room1", factory)
		T.Assert(err != nil)

		fixture.Drop("Room4")
		T.Assert(!fixture.Active("Room4"))
		T.Assert(!fixture.Exists("Room4"))

		fmt.Printf("%s\n", runtime.Root().Debug())
		T.Assert(strings.TrimSpace(runtime.Root().Debug()) == strings.TrimSpace(`
object: Untitled (1 / 0)
   object: Rooms (1 / 0)
        object: Room2 (0 / 0)
`))

		fixture.SetActive("Room1", true)

		fmt.Printf("%s\n", runtime.Root().Debug())
		T.Assert(strings.TrimSpace(runtime.Root().Debug()) == strings.TrimSpace(`
object: Untitled (1 / 0)
   object: Rooms (2 / 0)
        object: Room2 (0 / 0)
        object: Room1 (0 / 0)
`))
	})
}