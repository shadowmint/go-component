package component

import "ntoolkit/errors"

type ObjectStorageSetter interface {
	Set(id string, obj *ObjectTemplate) error
}

type ObjectStorageGetter interface {
	Get(id string) (*ObjectTemplate, error)
}

// ObjectStorage is a common high level interface for getting object templates and objects.
type ObjectStorage struct {
	Factory *ObjectFactory
	getter  ObjectStorageGetter
	setter  ObjectStorageSetter
}

// NewObjectStorage returns a new instance of an ObjectStorage
func NewObjectStorage(factory *ObjectFactory, getter ObjectStorageGetter, setter ObjectStorageSetter) *ObjectStorage {
	return &ObjectStorage{
		Factory: factory,
		setter: setter,
		getter: getter}
}

// CanSet returns true if an object setter is set
func (s *ObjectStorage) CanSet() bool {
	return s.setter != nil
}

// CanGet returns true if an object getter is set
func (s *ObjectStorage) CanGet() bool {
	return s.getter != nil
}

// SetObject serializes an object to template and saves it
func (s *ObjectStorage) SetObject(id string, obj *Object) error {
	template, err := s.Factory.Serialize(obj)
	if err != nil {
		return err
	}
	return s.SetObjectTemplate(id, template)
}

// SetObjectTemplate saves a template directly
func (s *ObjectStorage) SetObjectTemplate(id string, obj *ObjectTemplate) error {
	if !s.CanSet() {
		return errors.Fail(ErrNotSupported{}, nil, "Set is not supported by this ObjectStorage")
	}
	return s.setter.Set(id, obj)
}

// GetObject deserializes and returns the object for id
func (s *ObjectStorage) GetObject(id string) (*Object, error) {
	template, err := s.GetObjectTemplate(id)
	if err != nil {
		return nil, err
	}
	return s.Factory.Deserialize(template)
}

// GetObjectTemplate loads a template directly
func (s *ObjectStorage) GetObjectTemplate(id string) (*ObjectTemplate, error) {
	if !s.CanGet() {
		return nil, errors.Fail(ErrNotSupported{}, nil, "Get is not supported by this ObjectStorage")
	}
	return s.getter.Get(id)
}