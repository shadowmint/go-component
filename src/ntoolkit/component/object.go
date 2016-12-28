package component

import (
	"ntoolkit/iter"
)

// Node is a game object type.
type Object struct {
	Name       string
	Runtime    *Runtime
	components []*componentInfo // The set of components attached to this node
	children   []*Object        // The set of child objects attached to this node
}

// New returns a new Node
func NewObject() *Object {
	return &Object{
		Runtime : nil,
		components: make([]*componentInfo, 0),
		children: make([]*Object, 0)}
}

// Add a behaviour to a node
func (n *Object) AddComponent(component Component) {
	info := newComponentInfo(component)
	n.components = append(n.components, info)
	if info.Attach != nil {
		info.Attach.Attach(n)
	}
}

// Add a child object
func (n *Object) AddObject(object *Object) {
	n.children = append(n.children, object)
}

// Components returns an iterator of all the child objects on a game object
func (n *Object) Objects() iter.Iter {
	if len(n.children) == 0 {
		return nil
	}
	return fromObjectArray(&n.children)
}

// Update all components in this object
func (n *Object) Update(step float32, runtime *Runtime) {
	clone := n.components
	context := Context{Object: n, DeltaTime: step}
	for i := 0; i < len(clone); i++ {
		clone[i].updateComponent(step, runtime, &context)
	}
}

// Extend an existing iterator with more objects
func (n *Object) addChildren(iterator *ObjectArrayIter) {
	if len(n.children) > 0 {
		iterator.values.PushBack(&n.children)
	}
}