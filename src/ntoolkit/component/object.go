package component

import (
	"ntoolkit/iter"
	"reflect"
	"fmt"
	"strings"
)

// Node is a game object type.
type Object struct {
	Name       string
	Runtime    *Runtime
	components []*componentInfo // The set of components attached to this node
	children   []*Object        // The set of child objects attached to this node
}

// New returns a new Node
func NewObject(names ...string) *Object {
	name := ""
	if len(names) > 0 {
		name = names[0]
	}
	return &Object{
		Name: name,
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

// Objects returns an iterator of all the child objects on a game object
func (n *Object) Objects() iter.Iter {
	return fromObject(n)
}

// GetComponents returns an iterator of all components matching the given type.
func (n *Object) GetComponents(T reflect.Type) iter.Iter {
	return fromComponentArray(&n.components, T)
}

// GetComponentsInChildren returns an iterator of all components matching the given type in all children.
func (n *Object) GetComponentsInChildren(T reflect.Type) iter.Iter {
	cIter := fromComponentArray(nil, T)
	objIter := n.Objects()
	var val interface{} = nil
	var err error = nil
	for val, err = objIter.Next(); err == nil; val, err = objIter.Next() {
		componentList := &val.(*Object).components
		if len(*componentList) > 0 {
			cIter.Add(componentList)
		}
	}
	return cIter
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
func (n *Object) addChildren(iterator *ObjectIter) {
	if len(n.children) > 0 {
		iterator.values.PushBack(&n.children)
	}
}

// Debug prints out a summary of the object and its components
func (n *Object) Debug(indents ...int) string {
	indent := 0
	if len(indents) > 0 {
		indent = indents[0]
	}

	name := n.Name
	if len(name) == 0 {
		name = "Untitled"
	}

	rtn := fmt.Sprintf("object: %s (%d / %d)\n", name, len(n.children), len(n.components))
	if len(n.components) > 0 {
		for i := 0; i < len(n.components); i++ {
			rtn += fmt.Sprintf("! %s\n", n.components[i].Type)
		}
	}

	if len(n.children) > 0 {
		for i := 0; i < len(n.children); i++ {
			rtn += n.children[i].Debug(indent + 1) + "\n"
		}
	}

	lines := strings.Split(rtn, "\n")
	prefix := strings.Repeat("  ", indent)
	if indent != 0 {
		prefix += " "
	}
	output := ""
	for i := 0; i < len(lines); i++ {
		if len(strings.Trim(lines[i], " ")) != 0 {
			output += prefix
			output += lines[i]
			if i != (len(lines) - 1) {
				output += "\n"
			}
		}
	}

	return output
}