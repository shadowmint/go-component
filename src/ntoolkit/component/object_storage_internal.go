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

// SetObject serializes an object to template and saves it
func (s *objectStorage) SetObject(id string, obj *Object) error {
	template, err := s.factory.Serialize(obj)
	if err != nil {
		return err
	}
	return s.SetObjectTemplate(id, template)
}

// SetObjectTemplate saves a template directly
func (s *objectStorage) SetObjectTemplate(id string, obj *ObjectTemplate) error {
	if s.setter == nil {
		return errors.Fail(ErrNotSupported{}, nil, "Set is not supported by this ObjectStorage")
	}
	return s.setter.Set(id, obj)
}

// GetObject deserializes and returns the object for id
func (s *objectStorage) GetObject(id string) (*Object, error) {
	template, err := s.GetObjectTemplate(id)
	if err != nil {
		return nil, err
	}
	return s.factory.Deserialize(template)
}

// GetObjectTemplate loads a template directly
func (s *objectStorage) GetObjectTemplate(id string) (*ObjectTemplate, error) {
	if s.getter == nil {
		return nil, errors.Fail(ErrNotSupported{}, nil, "Get is not supported by this ObjectStorage")
	}
	return s.getter.Get(id)
}