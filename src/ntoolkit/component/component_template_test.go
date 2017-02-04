package component_test

import (
	"reflect"
	"ntoolkit/component"
)

type ComponentTemplate struct {
	parent *component.Object
}

func (c *ComponentTemplate) New() component.Component {
	return &ComponentTemplate{}
}

func (c *ComponentTemplate) Type() reflect.Type {
	return reflect.TypeOf(c)
}

func (c *ComponentTemplate) Attach(parent *component.Object) {
	c.parent = parent
}

func (c *ComponentTemplate) Update(context *component.Context) {
}

func (c *ComponentTemplate) Serialize() (interface{}, error) {
	return "", nil
}

func (c *ComponentTemplate) Deserialize(raw interface{}) error {
	return nil
}
