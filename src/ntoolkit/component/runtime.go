package component

import (
	"log"
	"os"
	"sync"

	"ntoolkit/iter"
	"ntoolkit/threadpool"
)

// Config configures a runtime.
type Config struct {
	ThreadPoolSize int
	Logger         *log.Logger
}

// Runtime is the basic operating unit of the mud.
// A Runtime executes the main game loop on objects.
type Runtime struct {
	root       *Object                // The root object for this runtime.
	workers    *threadpool.ThreadPool // The thread pool for updating objects
	logger     *log.Logger            // The logger for this runtime, if any.
	updateLock *sync.Mutex            // The thread safe lock for updates.
}

// New returns a new Runtime instance
func NewRuntime(config Config) *Runtime {
	validateConfig(&config)
	runtime := &Runtime{
		root:       NewObject(),
		logger:     config.Logger,
		updateLock: &sync.Mutex{},
		workers:    threadpool.New()}
	runtime.root.runtime = runtime
	runtime.workers.MaxThreads = config.ThreadPoolSize
	return runtime
}

// Configure sensible defaults if none are provided
func validateConfig(config *Config) {
	if config.ThreadPoolSize <= 0 {
		config.ThreadPoolSize = 10
	}
	if config.Logger == nil {
		config.Logger = log.New(os.Stdout, "runtime: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// Return a reference to the root object for the runtime
func (runtime *Runtime) Root() *Object {
	return runtime.root
}

// Return the set of objects as an iterator, including root.
func (runtime *Runtime) Objects() iter.Iter {
	return runtime.root.Objects()
}

// Schedules a task to execute between the next update loops.
// Return immediately, but the task will only be executed after
// the current Update() loop finishes.
// This effectively blocks until the current loop ends, then runs; then
// finally returns.
func (runtime *Runtime) ScheduleTask(task func()) {
	go (func() {
		defer (func() {
			r := recover()
			if r != nil {
				err, ok := r.(error)
				if ok {
					runtime.logger.Printf("Failed to execute scheduled task: %s", err.Error())
				} else {
					runtime.logger.Printf("Failed to execute scheduled task: %s", r)
				}
			}
		})()
		runtime.updateLock.Lock()
		task()
		defer runtime.updateLock.Unlock()
	})()
}

// Execute the update step of all components on all objects in worker threads
func (runtime *Runtime) Update(step float32) {
	runtime.updateLock.Lock()
	defer runtime.updateLock.Unlock()
	runtime.updateObject(step, runtime.root)
	objects := runtime.Objects()
	if objects != nil {
		var val interface{}
		var err error
		for val, err = objects.Next(); err == nil; val, err = objects.Next() {
			obj := val.(*Object)
			runtime.updateObject(step, obj)
		}
		runtime.workers.Wait()
	}
}

// Execute a single object update
func (runtime *Runtime) updateObject(step float32, obj *Object) {
	runtime.workers.Run(func() {
		if obj.runtime == nil {
			obj.runtime = runtime
		}
		obj.Update(step, runtime)
	})
}
