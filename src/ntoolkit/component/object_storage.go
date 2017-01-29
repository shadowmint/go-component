package component

type ObjectStorageSetter interface {
	Set(id string, obj *ObjectTemplate) error
	Clear(id string) error
}

type ObjectStorageGetter interface {
	Get(id string) (*ObjectTemplate, error)
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
	return nil, nil
}

// Set an object into the active list.
func (s *ObjectStorage) Set(id string, obj *Object) error {
	return nil
}

// Save pushes the given object from active into storage and drop it from the active list.
func (s *ObjectStorage) Save(id string) error {
	return nil
}

// Active checks if the given object is currently in the active list.
func (s *ObjectStorage) Active(id string) bool {
	return false
}

// Exists checks if the given object is currently in the active list or storage.
func (s *ObjectStorage) Exists(id string) bool {
	return false
}

// Clear completely removes the given object.
func (s *ObjectStorage) Clear(id string) error {
	return nil
}

