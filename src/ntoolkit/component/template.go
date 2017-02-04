package component

import (
	"ntoolkit/errors"
	"encoding/json"
)

// ObjectTemplate is a simple, flat, serializable object structure that directly converts to and from Objects.
type ObjectTemplate struct {
	Name       string
	Components []ComponentTemplate
	Objects    []ObjectTemplate
}

// ComponentTemplate is a serializable representation of a component
type ComponentTemplate struct {
	Type string
	Data interface{}
}

// FromJson loads an object template from a json block.
func ObjectTemplateFromJson(raw string) (*ObjectTemplate, error) {
	var data ObjectTemplate
	err := json.Unmarshal([]byte(raw), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// AsObject converts a map[string]interface{} as a typed object.
// This is a helper for serializable components.
func AsObject(target interface{}, raw interface{}) (err error) {
	defer (func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	})()
	if raw == nil {
		return errors.Fail(ErrNullValue{}, nil, "No data (null)")
	}
	if target == nil {
		return errors.Fail(ErrNullValue{}, nil, "No target (null)")
	}

	value := raw.(map[string]interface{})

	bytes, err := json.Marshal(value)
	if err != nil {
		return errors.Fail(ErrBadValue{}, err, "Failed to re-encode data")
	}

	err = json.Unmarshal(bytes, target)
	if err != nil {
		return errors.Fail(ErrBadValue{}, err, "Failed to decode data")
	}

	return nil
}

// AsString converts a map[string]interface{} as a string object.
// This is a helper for serializable components.
func AsString(raw interface{}) (rtn string, err error) {
	defer (func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	})()
	if raw == nil {
		return "", errors.Fail(ErrNullValue{}, nil, "No data (null)")
	}
	value := raw.(string)
	return value, nil
}
