package component_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"ntoolkit/assert"
	"ntoolkit/component"
	"strings"
	"io/ioutil"
)

// Add remove child adds a new child every second.
// When it has 10 children, it removes itself.
type AddRemoveChild struct {
	parent  *component.Object
	count   int
	elapsed float32
}

func (c *AddRemoveChild) New() component.Component {
	return &AddRemoveChild{}
}

func (c *AddRemoveChild) Type() reflect.Type {
	return reflect.TypeOf(c)
}

func (c *AddRemoveChild) Attach(parent *component.Object) {
	c.parent = parent
}

func (c *AddRemoveChild) Update(context *component.Context) {
	context.Logger.Printf("Update: %s", c.parent.Name())
	c.elapsed += context.DeltaTime
	if c.elapsed > 1.0 {
		c.count += 1
		if c.count >= 3 {
			parent := c.parent.Parent()
			if parent != nil {
				parent.RemoveObject(c.parent)
			}
		} else {
			child := component.NewObject(fmt.Sprintf("Child: %d", c.count))
			c.parent.AddObject(child)
		}
		c.elapsed = 0
	}
}

// DumpState dumps an object tree of the runtime every 1/2 seconds
type DumpState struct {
	elapsed float32
}

func (c *DumpState) New() component.Component {
	return &DumpState{}
}

func (c *DumpState) Type() reflect.Type {
	return reflect.TypeOf(c)
}

func (c *DumpState) Update(context *component.Context) {
	c.elapsed += context.DeltaTime
	context.Logger.Printf("DumpState: %f", c.elapsed)
	if c.elapsed >= 0.5 {
		c.elapsed = 0.0
		root := context.Object.Root()
		structure := root.Debug()
		context.Logger.Printf("Tree: %s", structure)
	}
}

func TestComplexRuntime(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		logger := log.New(os.Stdout, "Runtime: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.SetOutput(ioutil.Discard) // No output thanks

		runtime := component.NewRuntime(component.Config{
			ThreadPoolSize: 10,
			Logger:         logger})

		runtime.Root().AddComponent(&DumpState{})

		o1 := component.NewObject("Container One")
		w1 := component.NewObject("Worker 1")
		w2 := component.NewObject("Worker 2")

		o2 := component.NewObject("Container Two")
		w3 := component.NewObject("Worker 3")
		w4 := component.NewObject("Worker 4")

		o3 := component.NewObject("Container Three")
		w5 := component.NewObject("Worker 5")
		w6 := component.NewObject("Worker 6")

		o1.AddObject(w1)
		o1.AddObject(w2)

		o2.AddObject(w3)
		o2.AddObject(w4)

		w4.AddObject(o3)
		o3.AddObject(w5)
		o3.AddObject(w6)

		w1.AddComponent(&AddRemoveChild{})
		w2.AddComponent(&AddRemoveChild{})
		w3.AddComponent(&AddRemoveChild{})
		w4.AddComponent(&AddRemoveChild{})
		w5.AddComponent(&AddRemoveChild{})
		w6.AddComponent(&AddRemoveChild{})

		runtime.Root().AddObject(o1)
		runtime.Root().AddObject(o2)

		for i := 0; i < 50; i++ {
			runtime.Update(0.25)
		}

		expectedOutput := strings.Trim(`
object: Untitled (2 / 1)
! *ntoolkit/component_test.DumpState
   object: Container One (0 / 0)
   object: Container Two (0 / 0)`, " \n")

		actualOutput := strings.Trim(runtime.Root().Debug(), " \n")
		
		T.Assert(expectedOutput == actualOutput)
	})
}
