package component

import (
	"reflect"
	"fmt"
	"ntoolkit/errors"
)
// ComponentProvider maps between component instances and component templates
type ComponentProvider interface {
	Type() reflect.Type
	New() Component
}

// ObjectFactory is the overseer that can be used to convert between objects and object templates
type ObjectFactory struct {
	handlers map[string]ComponentProvider
}

// NewObjectFactory returns a new object factory
func NewObjectFactory() *ObjectFactory {
	return &ObjectFactory{handlers: make(map[string]ComponentProvider)}
}

// Register a ComponentProvider that can be used to serialize and deserialize objects
func (factory *ObjectFactory) Register(provider ComponentProvider) {
	factory.handlers[factory.typeName(provider.Type())] = provider
}

// Serialize converts an object into an ObjectTemplate
func (factory *ObjectFactory) Serialize(object *Object) (*ObjectTemplate, error) {
	return nil, nil
}

// Deserialize converts an ObjectTemplate into an object
func (factory *ObjectFactory) Deserialize(template *ObjectTemplate) (*Object, error) {
	obj := NewObject(template.Name)

	// Add components
	for i := 0; i < len(template.Components); i++ {
		c, err := factory.deserializeComponent(&template.Components[i])
		if err != nil {
			return nil, err
		}
		obj.AddComponent(c)
	}

	// Add children
	for i := 0; i < len(template.Objects); i++ {
		child, err := factory.Deserialize(&template.Objects[i])
		if err != nil {
			return nil, err
		}
		obj.AddObject(child)
	}

	return obj, nil
}

// deserializeComponent turns a component template into a component
func (factory *ObjectFactory) deserializeComponent(template *ComponentTemplate) (Component, error) {
	for k, v := range factory.handlers {
		if k == template.Type {
			component := v.New()
			// TODO: Deserialize the data here
			return component, nil
		}
	}
	return nil, errors.Fail(ErrUnknownComponent{}, nil, fmt.Sprintf("Component type %s is not registered with the factory", template.Type))
}

// typeName returns the name for a specific type
func (factory *ObjectFactory) typeName(T reflect.Type) string {
	return fmt.Sprintf("%s", T)
}
