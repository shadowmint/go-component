package component

import (
	"ntoolkit/errors"
	"fmt"
)

type ObjectStorageSetter interface {
	Set(id string, obj *Object, factory *ObjectFactory) error
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
	fmt.Printf("Enter GET\n")
	if s.active.Has(id) {
		fmt.Printf("In active\n")
		return s.active.Get(id)
	} else if s.storage.Has(id) {
		fmt.Printf("Not active but in storage\n")
		if err := s.SetActive(id, true); err != nil {
			return nil, err
		} else {
			return s.Get(id)
		}
	}
	fmt.Printf("No match\n")
	return nil, errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("No match for id %s in storage", id))
}

// Add a new object into the active list.
func (s *ObjectStorage) Add(id string, obj *Object) error {
	return s.active.Set(id, obj)
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
			fmt.Printf("Storage failed to load object\n")
			return err
		}
		fmt.Printf("Activate? --> %s\n", obj)
		if err := s.active.Set(id, obj); err != nil {
			return err
		}
		return nil
	} else {
		fmt.Printf("Set active -> false\n")
		if !s.active.Has(id) {
			fmt.Printf("Not active currently\n")
			return errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("No match for id %s in active storage", id))
		}
		obj, err := s.active.Get(id)
		if err != nil {
			fmt.Printf("Failed to get object\n")
			return err
		}
		if err := s.storage.Set(id, obj); err != nil {
			fmt.Printf("Failed to save object\n")
			return err
		}
		fmt.Printf("Dropping instance!\n")
		if err := s.active.Clear(id); err != nil {
			fmt.Printf("Failed to clear object\n")
			return err
		}
		fmt.Printf("All good?")
		return nil
	}
}

// Drop drops an object entirely.
func (s *ObjectStorage) Drop(id string) error {
	fmt.Printf("Drop!!!\n\n")
	if s.active.Has(id) {
		fmt.Printf("Drop from Active!\n")
		if err := s.active.Clear(id); err != nil {
			fmt.Printf("Nope\n")
			return err
		}
	}
	if s.storage.Has(id) {
		fmt.Printf("Drop from Storage!\n")
		if err := s.storage.Clear(id); err != nil {
			fmt.Printf("Nope\n")
			return err
		}
	}
	fmt.Printf("\n\n")
	return nil
}

// Active checks if the given object is currently in the active list.
// Notice this should be used in combination with `Exists` for full check.
func (s *ObjectStorage) Active(id string) bool {
	return s.active.Has(id)
}

// Exists checks if the given object is currently in the active list or storage.
func (s *ObjectStorage) Exists(id string) bool {
	fmt.Printf("%s -- %s\n", s.active.Has(id), s.storage.Has(id))
	return s.active.Has(id) || s.storage.Has(id)
}

