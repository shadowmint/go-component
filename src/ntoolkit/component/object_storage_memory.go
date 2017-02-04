package component

import (
	"ntoolkit/errors"
	"regexp"
	"fmt"
)

// ObjectStorageMemory serializes and stores objects in a persistent memory cache.
type ObjectStorageMemory struct {
	data    map[string]*ObjectTemplate
	Pattern string // If set, only accept ids that match this.
	CanSet  bool // Controls the Setter() behaviour, not Set()
	CanGet  bool // Controls the Getter() behaviour, not Get()
	pattern string
	regex   *regexp.Regexp
}

// NewObjectStorageMemory returns a new instance that caches templates in a simple local hash.
func NewObjectStorageMemory() *ObjectStorageMemory {
	return &ObjectStorageMemory{
		data: make(map[string]*ObjectTemplate),
		Pattern: ".*",
		pattern: ".*",
		CanSet: true,
		CanGet: true}
}

// rebuild the regex internally if the pattern changed for some reason
func (s *ObjectStorageMemory) rebuild() error {
	if s.pattern != s.Pattern || s.regex == nil {
		s.pattern = s.Pattern
		expr, err := regexp.Compile(s.pattern)
		if err != nil {
			return err
		}
		s.regex = expr
	}
	return nil
}

func (s *ObjectStorageMemory) Set(id string, obj *Object, factory *ObjectFactory) error {
	if err := s.rebuild(); err != nil {
		return err
	}
	if !s.regex.Match([]byte(id)) {
		return errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("Key %s does not match storage pattern %d", id, s.Pattern))
	}
	tmpl, err := factory.Serialize(obj)
	if err != nil {
		return err
	}
	s.data[id] = tmpl
	return nil
}


func (s *ObjectStorageMemory) Clear(id string) error {
	if err := s.rebuild(); err != nil {
		return err
	}
	if !s.regex.Match([]byte(id)) {
		return errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("Key %s does not match storage pattern %d", id, s.Pattern))
	}
	if _, ok := s.data[id]; ok {
		delete(s.data, id)
	}
	return nil
}

func (s *ObjectStorageMemory) Get(id string, factory *ObjectFactory) (*Object, error) {
	if err := s.rebuild(); err != nil {
		return nil, err
	}
	if !s.regex.Match([]byte(id)) {
		return nil, errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("Key %s does not match storage pattern %d", id, s.Pattern))
	}
	template, ok := s.data[id]
	if !ok {
		return nil, errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("No id %s in storage", id))
	}
	obj, err := factory.Deserialize(template)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *ObjectStorageMemory) Has(id string) bool {
	if err := s.rebuild(); err != nil {
		return false
	}
	if !s.regex.Match([]byte(id)) {
		return false
	}
	_, ok := s.data[id]
	return ok
}

func (s *ObjectStorageMemory) Getter() ObjectStorageGetter {
	if !s.CanGet {
		return nil
	}
	return s
}

func (s *ObjectStorageMemory) Setter() ObjectStorageSetter {
	if !s.CanSet {
		return nil
	}
	return s
}