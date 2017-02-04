package component

import (
	"ntoolkit/errors"
	"fmt"
)

type ObjectStorageSetter interface {
	Set(obj *Object, factory *ObjectFactory) error
	Clear(id string) error
}

type ObjectStorageGetter interface {
	Get(id string, factory *ObjectFactory) (*Object, error)
	Has(id string) bool
}

type ObjectStorageProvider interface {
	Setter() ObjectStorageSetter
	Getter() ObjectStorageGetter
}

type ObjectStorage struct {
	active  *objectStorage
	storage *objectStorage
}

// NewObjectStorage returns a new instance of an ObjectStorage
func NewObjectStorage(factory *ObjectFactory, active ObjectStorageProvider, storage ObjectStorageProvider) *ObjectStorage {
	return &ObjectStorage{
		active: newObjectStorage(factory, active),
		storage: newObjectStorage(factory, storage)}
}

// Get an object; if it doesn't exist in the active list, load it.
func (s *ObjectStorage) Get(id string) (*Object, error) {
	if s.active.Has(id) {
		return s.active.Get(id)
	} else if s.storage.Has(id) {
		if err := s.SetActive(id, true); err != nil {
			return nil, err
		} else {
			return s.Get(id)
		}
	}
	return nil, errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("No match for id %s in storage", id))
}

// Add a new object into the active list.
func (s *ObjectStorage) Add(obj *Object) error {
	return s.active.Set(obj)
}

// SetActive either saves and object and removes it from active or activates an object and adds it to active.
func (s *ObjectStorage) SetActive(id string, active bool) error {
	if (active) {
		if s.active.Has(id) {
			return nil
		}
		if !s.storage.Has(id) {
			return errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("No match for id %s in storage", id))
		}
		obj, err := s.storage.Get(id)
		if err != nil {
			return err
		}
		if err := s.active.Set(obj); err != nil {
			return err
		}
		return nil
	} else {
		if !s.active.Has(id) {
			return errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("No match for id %s in active storage", id))
		}
		obj, err := s.active.Get(id)
		if err != nil {
			return err
		}
		if err := s.storage.Set(obj); err != nil {
			return err
		}
		if err := s.active.Clear(id); err != nil {
			return err
		}
		return nil
	}
}

// Drop drops an object entirely.
func (s *ObjectStorage) Drop(id string) error {
	if s.active.Has(id) {
		if err := s.active.Clear(id); err != nil {
			return err
		}
	}
	if s.storage.Has(id) {
		if err := s.storage.Clear(id); err != nil {
			return err
		}
	}
	return nil
}

// Active checks if the given object is currently in the active list.
// Notice this should be used in combination with `Exists` for full check.
func (s *ObjectStorage) Active(id string) bool {
	return s.active.Has(id)
}

// Exists checks if the given object is currently in the active list or storage.
func (s *ObjectStorage) Exists(id string) bool {
	return s.active.Has(id) || s.storage.Has(id)
}

