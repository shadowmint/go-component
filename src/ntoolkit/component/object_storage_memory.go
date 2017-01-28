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
	CanSet  bool
	CanGet  bool
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

func (s *ObjectStorageMemory) Set(id string, obj *ObjectTemplate) error {
	if err := s.rebuild(); err != nil {
		return err
	}
	if !s.regex.Match([]byte(id)) {
		return errors.Fail(ErrBadValue{}, nil, fmt.Sprintf("Key %s does not match storage pattern %d", id, s.Pattern))
	}
	s.data[id] = obj
	return nil
}

func (s *ObjectStorageMemory) Get(id string) (*ObjectTemplate, error) {
	if err := s.rebuild(); err != nil {
		return nil, err
	}
	if !s.regex.Match([]byte(id)) {
		return nil, errors.Fail(ErrBadValue{}, nil, fmt.Sprintf("Key %s does not match storage pattern %d", id, s.Pattern))
	}
	template, ok := s.data[id]
	if !ok {
		return nil, errors.Fail(ErrNoMatch{}, nil, fmt.Sprintf("No id %s in storage", id))
	}
	return template, nil
}

func (s *ObjectStorageMemory) Getter() ObjectStorageGetter {
	return s
}

func (s *ObjectStorageMemory) Setter() ObjectStorageSetter {
	return s
}