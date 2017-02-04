package component

import (
	"ntoolkit/iter"
	"reflect"
	"fmt"
	"strings"
	"ntoolkit/errors"
	"sync"
)

// Node is a game object type.
type Object struct {
	Name        string
	Runtime     *Runtime
	components  []*componentInfo // The set of components attached to this node
	children    []*Object        // The set of child objects attached to this node
	parent      *Object
	writeLock   *sync.Mutex
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
		children: make([]*Object, 0),
		writeLock: &sync.Mutex{}}
}

// Add a behaviour to a node
func (n *Object) AddComponent(component Component) {
	n.lock()
	info := newComponentInfo(component)
	n.components = append(n.components, info)
	if info.Attach != nil {
		info.Attach.Attach(n)
	}
	n.unlock()
}

// Add a child object
func (n *Object) AddObject(object *Object) error {
	var err error
	n.lock()
	if n == object || n.HasParent(object) {
		err = errors.Fail(ErrBadObject{}, nil, "Circular object references are not permitted")
	} else {
		n.children = append(n.children, object)
		object.lock()
		object.parent = n
		object.unlock()
	}
	n.unlock()
	return err
}

// Remove a child object
func (n *Object) RemoveObject(object *Object) error {
	n.lock()
	offset := -1
	for i := 0; i < len(n.children); i++ {
		if n.children[i] == object {
			offset = i
			break
		}
	}
	if offset >= 0 {
		n.children = append(n.children[:offset], n.children[offset+1:]...)
	}
	object.parent = nil
	n.unlock()
	return nil
}


// Check if an object has a parent
func (n *Object) HasParent(object *Object) bool {
	root := n
	for root != nil {
		root = root.Parent()
		if root == object {
			return true
		}
	}
	return false
}

// Return the parent of this object
func (n *Object) Parent() *Object {
	return n.parent
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

// Find returns the first matching component on the object tree given by the name sequence or nil
// component should be a pointer to store the output component into.
// eg. If *FakeComponent implements Component, pass **FakeComponent to Find.
func (n *Object) Find(component interface{}, query ...string) error {
	componentType := reflect.TypeOf(component).Elem()

	obj := n
	var err error
	if len(query) != 0 {
		obj, err = n.FindObject(query...)
		if err != nil {
			return err
		}
	}

	cmp, err := obj.GetComponents(componentType).Next()
	if err != nil {
		return err
	}

	reflect.ValueOf(component).Elem().Set(reflect.ValueOf(cmp))
	return nil
}

// FindObject returns the first matching child object on the object tree given by the name sequence or nil
func (n *Object) FindObject(query ...string) (*Object, error) {
	if len(query) == 0 {
		return nil, errors.Fail(ErrBadValue{}, nil, "Invalid query length of zero")
	}

	cursor := n
	queryCursor := 0

	var rtn *Object = nil
	for rtn == nil {
		next, err := cursor.GetObject(query[queryCursor])
		if err != nil {
			return nil, err
		} else {
			cursor = next
		}

		queryCursor += 1
		if queryCursor == len(query) {
			rtn = cursor
		}
	}

	return rtn, nil
}

func (n *Object) GetObject(name string) (*Object, error) {
	for i := 0; i < len(n.children); i++ {
		if n.children[i].Name == name {
			return n.children[i], nil
		}
	}
	return nil, errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("No match for object '%s' on parent '%s'", name, n.Name))
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
			rtn += fmt.Sprintf("! %s\n", typeName(n.components[i].Type))
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

// Unlock for write
func (o *Object) unlock() {
	o.writeLock.Unlock()
}

// In order to lock we must first wait for any other write subtrees to finish.
// We lock up to the root, then lock ourself, then unlock the parent chain.
func (o *Object) lock() {
	i := o.parent
	for {
		if i == nil {
			break
		}
		i.writeLock.Lock()
		i = i.parent
	}
	o.writeLock.Lock()
	i = o.parent
	for {
		if i == nil {
			break
		}
		i.writeLock.Unlock()
		i = i.parent
	}
}