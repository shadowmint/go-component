package component

import (
	"ntoolkit/errors"
	"fmt"
)

// ObjectStorageStack abstracts over a stack of storage options.
// Both get and set operations are chained through until a match is hit.
type ObjectStorageStack struct {
	get []ObjectStorageGetter
	set []ObjectStorageSetter
}

// NewObjectStorageMemory returns a new instance that caches templates in a simple local hash.
func NewObjectStorageStack() *ObjectStorageStack {
	return &ObjectStorageStack{
		get: make([]ObjectStorageGetter, 0),
		set: make([]ObjectStorageSetter, 0)}
}

// Add a new storage tier to the supported channels
func (s *ObjectStorageStack) Add(storage ObjectStorageProvider) {
	if storage.Getter() != nil {
		s.get = append(s.get, storage.Getter())
	}
	if storage.Setter() != nil {
		s.set = append(s.set, storage.Setter())
	}
}

func (s *ObjectStorageStack) Set(id string, obj *Object, factory *ObjectFactory) error {
	var err error
	for i := 0; i < len(s.set); i++ {
		if err = s.set[i].Set(id, obj, factory); err == nil {
			return nil
		} else {
			if !errors.Is(err, ErrNotSupported{}) && !errors.Is(err, ErrNoMatch{}) {
				return err
			}
		}
	}
	return errors.Fail(ErrBadObject{}, err, "Unable to save object in any object storage instances")
}

func (s *ObjectStorageStack) Clear(id string) error {
	var err error
	for i := 0; i < len(s.set); i++ {
		fmt.Printf("Dropped instance for %s\n", s.set[i])
		if err = s.set[i].Clear(id); err == nil {
			fmt.Printf("Successful DROP -> %s\n", s)
			return nil
		}
	}
	return errors.Fail(ErrBadObject{}, err, "Unable to clear object in any object storage instances")
}

func (s *ObjectStorageStack) Get(id string, factory *ObjectFactory) (*Object, error) {
	var err error
	var rtn *Object
	for i := 0; i < len(s.get); i++ {
		if rtn, err = s.get[i].Get(id, factory); err == nil {
			return rtn, nil
		} else {
			if !errors.Is(err, ErrNotSupported{}) && !errors.Is(err, ErrNoMatch{}) {
				return nil, err
			}
		}
	}
	return nil, errors.Fail(ErrNoMatch{}, err, "Unable to get object from any object storage instances")
}

func (s *ObjectStorageStack) Has(id string) bool {
	fmt.Printf("Looking for match in %d instance\n", len(s.get))
	for i := 0; i < len(s.get); i++ {
		if has := s.get[i].Has(id); has {
			fmt.Printf("Found match in %s!\n", s.get[i])
			return true
		} else {
			fmt.Printf("No match in %s\n", s.get[i])
		}
	}
	return false
}

func (s *ObjectStorageStack) Getter() ObjectStorageGetter {
	return s
}

func (s *ObjectStorageStack) Setter() ObjectStorageSetter {
	return s
}