package component

import "ntoolkit/errors"

// ObjectStorageStack abstracts over a stack of storage options.
// Both get and set operations are chained through until a match is hit.
type ObjectStorageStack struct {
	get []*ObjectStorage
	set []*ObjectStorage
}

// NewObjectStorageMemory returns a new instance that caches templates in a simple local hash.
func NewObjectStorageStack() *ObjectStorageStack {
	return &ObjectStorageStack{
		get: make([]*ObjectStorage, 0),
		set: make([]*ObjectStorage, 0)}
}

// Add a new storage tier to the supported channels
func (s *ObjectStorageStack) Add(storage *ObjectStorage) {
	if storage.CanGet() {
		s.get = append(s.get, storage)
	}
	if storage.CanSet() {
		s.set = append(s.set, storage)
	}
}

func (s *ObjectStorageStack) Set(id string, obj *ObjectTemplate) error {
	var err error
	for i := 0; i < len(s.set); i++ {
		if err = s.set[i].SetObjectTemplate(id, obj); err == nil {
			return nil
		}
	}
	return errors.Fail(ErrBadObject{}, err, "Unable to save object in any object storage instances")
}

func (s *ObjectStorageStack) Get(id string) (*ObjectTemplate, error) {
	var err error
	var rtn *ObjectTemplate
	for i := 0; i < len(s.set); i++ {
		if rtn, err = s.get[i].GetObjectTemplate(id); err == nil {
			return rtn, nil
		}
	}
	return nil, errors.Fail(ErrNoMatch{}, err, "Unable to get object from any object storage instances")
}
