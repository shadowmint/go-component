package component_test

import (
	"ntoolkit/assert"
	"testing"
	c "ntoolkit/component"
	"ntoolkit/iter"
	"reflect"
)

func TestBasicTemplateToObject(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		factory := c.NewObjectFactory()
		factory.Register(&FakeComponent{})
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

func TestComponentDeserialization(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		factory := c.NewObjectFactory()
		factory.Register(&FakeComponent{})
		template := c.ObjectTemplate{
			Components: []c.ComponentTemplate{{Type: "*component_test.FakeComponent"}},
			Objects: []c.ObjectTemplate{
				{Name: "First Child"},
				{Components: []c.ComponentTemplate{{Type: "*component_test.FakeComponent", Data: "Value2,5"}},
					Objects: []c.ObjectTemplate{
						{Components: []c.ComponentTemplate{{Type: "*component_test.FakeComponent", Data: "Value1,1"}}},
						{Name: "Last Child"}}}}}

		instance, err := factory.Deserialize(&template)

		T.Assert(err == nil)
		T.Assert(instance != nil)

		cmps, _ := iter.Collect(instance.GetComponentsInChildren(reflect.TypeOf((*FakeComponent)(nil))))

		T.Assert(cmps[1].(*FakeComponent).Id == "Value2")
		T.Assert(cmps[1].(*FakeComponent).Count == 5)
		T.Assert(cmps[2].(*FakeComponent).Id == "Value1")
		T.Assert(cmps[2].(*FakeComponent).Count == 1)
	})
}

func TestObjectToTemplate(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		factory := c.NewObjectFactory()
		factory.Register(&FakeComponent{})

		template := &c.ObjectTemplate{
			Components: []c.ComponentTemplate{{Type: "*component_test.FakeComponent"}},
			Objects: []c.ObjectTemplate{
				{Name: "First Child"},
				{Components: []c.ComponentTemplate{{Type: "*component_test.FakeComponent"}},
					Objects: []c.ObjectTemplate{
						{Components: []c.ComponentTemplate{{Type: "*component_test.FakeComponent"}}},
						{Name: "Last Child"}}}}}

		instance, _ := factory.Deserialize(template)
		dump1 := instance.Debug()

		template, err := factory.Serialize(instance)
		T.Assert(err == nil)

		instance, err = factory.Deserialize(template)
		T.Assert(err == nil)
		T.Assert(instance != nil)

		dump2 := instance.Debug()
		T.Assert(dump1 == dump2)
	})
}
