package component

import (
	"ntoolkit/threadpool"
	"ntoolkit/iter"
)

// Config configures a runtime.
type Config struct {
	ThreadPoolSize int
}

// Runtime is the basic operating unit of the mud.
// A Runtime executes the main game loop on objects.
type Runtime struct {
	locked  bool                   // Is this runtime current locked for updates?
	root    *Object                // The root object for this runtime.
	workers *threadpool.ThreadPool // The thread pool for updating objects
}

// New returns a new Runtime instance
func NewRuntime(config Config) *Runtime {
	runtime := &Runtime{
		locked: false,
		root: NewObject(),
		workers: threadpool.New()}
	runtime.workers.MaxThreads = config.ThreadPoolSize
	return runtime
}

// Return a reference to the root object for the runtime
func (runtime *Runtime) Root() *Object {
	return runtime.root
}

// Return the set of objects as an iterator.
func (runtime *Runtime) Objects() iter.Iter {
	return runtime.root.Objects()
}

// Execute the update step of all components on all objects in worker threads
func (runtime *Runtime) Update(step float32) {
	objects := runtime.Objects()
	for val, err := objects.Next(); err == nil; val, err = objects.Next() {
		obj := val.(*Object)
		runtime.updateObject(step, obj)
	}
	runtime.workers.Wait()
}

// Execute a single object update
func (runtime *Runtime) updateObject(step float32, obj *Object) {
	runtime.workers.Run(func() {
		obj.Update(step, runtime)
	})
}