package component

import (
	"ntoolkit/errors"
	"ntoolkit/iter"
	"container/list"
)

// ObjectArrayIter implements Iterator for []Object
type ObjectArrayIter struct {
	values  *list.List
	current *[]*Object
	offset  int
	err     error
}

// fromObjectArray returns a new list iterator for a list
func fromObjectArray(values *[]*Object) *ObjectArrayIter {
	rtn := &ObjectArrayIter{values: list.New(), offset: -1}
	if rtn.values == nil {
		rtn.err = errors.Fail(ErrNullValue{}, nil, "Invalid object array")
	} else {
		rtn.values.PushBack(values)
	}
	return rtn
}

// Next increments the iterator cursor
func (iterator *ObjectArrayIter) Next() (interface{}, error) {
	if iterator.err != nil {
		return nil, iterator.err
	}

	if iterator.current == nil {
		if iterator.nextGroup() {
			return nil, iterator.err
		}
	} else {
		iterator.offset += 1
	}

	if iterator.offset >= len(*iterator.current) {
		if iterator.nextGroup() {
			return nil, iterator.err
		}
	}

	nextValue := (*iterator.current)[iterator.offset]
	nextValue.addChildren(iterator)

	return (*iterator.current)[iterator.offset], nil
}

// Attempt to get the next set if its null.
// Each set is the set of objects from a different parent.
// Returns true if an error is set.
func (iterator *ObjectArrayIter) nextGroup() bool {
	el := iterator.values.Front()
	if el != nil {
		iterator.values.Remove(el)
		iterator.current = el.Value.(*[]*Object)
		iterator.offset = 0
	} else {
		iterator.err = errors.Fail(iter.ErrEndIteration{}, nil, "No more values")
		return true
	}
	return false
}