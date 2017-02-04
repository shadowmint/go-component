package component

import "ntoolkit/errors"

type objectStorage struct {
	factory *ObjectFactory
	getter  ObjectStorageGetter
	setter  ObjectStorageSetter
}

// NewObjectStorage returns a new instance of an ObjectStorage
func newObjectStorage(factory *ObjectFactory, provider ObjectStorageProvider) *objectStorage {
	return &objectStorage{
		factory: factory,
		setter: provider.Setter(),
		getter: provider.Getter()}
}

// Check if the given id is in the storage
func (s *objectStorage) Has(id string) bool {
	return s.getter.Has(id)
}

// Remove the given id from the storage
func (s *objectStorage) Clear(id string) error {
	return s.setter.Clear(id)
}

// Set saves an object if possible.
func (s *objectStorage) Set(id string, obj *Object) error {
	if s.setter == nil {
		return errors.Fail(ErrNotSupported{}, nil, "Set is not supported on this instance")
	}
	return s.setter.Set(id, obj, s.factory)
}

// Get saves an object if possible.
func (s *objectStorage) Get(id string) (*Object, error) {
	if s.getter == nil {
		return nil, errors.Fail(ErrNotSupported{}, nil, "Get is not supported on this instance")
	}
	return s.getter.Get(id, s.factory)
}