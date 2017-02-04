package component

import (
	"ntoolkit/errors"
	"fmt"
)

// ObjectStorageRuntime serializes and stores objects in a persistent memory cache.
type ObjectStorageRuntime struct {
	root   *Object
	CanSet bool // Controls the Setter() behaviour, not Set()
	CanGet bool // Controls the Getter() behaviour, not Get()
}

// NewObjectStorageRuntime returns a new instance that caches templates in a simple local hash.
func NewObjectStorageRuntime(root *Object) *ObjectStorageRuntime {
	return &ObjectStorageRuntime{
		root: root,
		CanSet: true,
		CanGet: true}
}

func (s *ObjectStorageRuntime) Set(id string, obj *Object, _ *ObjectFactory) error {
	if obj == nil {
		return errors.Fail(ErrNullValue{}, nil, "Unable to set null object")
	}
	old, err := s.root.GetObject(id)
	if !errors.Is(err, ErrNoMatch{}) {
		s.root.RemoveObject(old)
	}
	fmt.Printf("Setting name on %s\n", obj)
	obj.Name = id
	return s.root.AddObject(obj)
}

func (s *ObjectStorageRuntime) Clear(id string) error {
	obj, err := s.root.GetObject(id)
	if errors.Is(err, ErrNoMatch{}) {
		return err
	}
	return s.root.RemoveObject(obj)
}

func (s *ObjectStorageRuntime) Get(id string, _ *ObjectFactory) (*Object, error) {
	return s.root.GetObject(id)
}

func (s *ObjectStorageRuntime) Has(id string) bool {
	_, err := s.root.GetObject(id)
	return err == nil
}

func (s *ObjectStorageRuntime) Getter() ObjectStorageGetter {
	if !s.CanGet {
		return nil
	}
	return s
}

func (s *ObjectStorageRuntime) Setter() ObjectStorageSetter {
	if !s.CanSet {
		return nil
	}
	return s
}