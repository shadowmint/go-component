package component_test

import (
	"ntoolkit/assert"
	"testing"
	c "ntoolkit/component"
	"fmt"
)

func TestBasicTemplateToObject(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		factory := c.NewObjectFactory()
		factory.Register(&FakeComponent{})
		fmt.Printf("%s\n", factory)
		template := c.ObjectTemplate{
			Components: []c.ComponentTemplate{{Type: "*component_test.FakeComponent"}},
			Objects: []c.ObjectTemplate{
				{Name: "First Child"},
				{Components: []c.ComponentTemplate{{Type: "*component_test.FakeComponent"}},
					Objects: []c.ObjectTemplate{
						{Components: []c.ComponentTemplate{{Type: "*component_test.FakeComponent"}}},
						{Name: "Last Child"}}}}}

		instance, err := factory.Deserialize(&template)

		T.Assert(err == nil)
		T.Assert(instance != nil)

		T.Assert(instance.Debug() == `object: Untitled (2 / 1)
! *component_test.FakeComponent
   object: First Child (0 / 0)
   object: Untitled (2 / 1)
   ! *component_test.FakeComponent
        object: Untitled (0 / 1)
        ! *component_test.FakeComponent
        object: Last Child (0 / 0)
`)
	})
}