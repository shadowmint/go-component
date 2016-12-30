package component

// ErrNullValue is raised when a null value is passed in as a pointer
type ErrNullValue struct{}

// ErrUnknownComponent is raised when trying to deserialize an unknown component type
type ErrUnknownComponent struct{}

// ErrBadValue is raised when an invalid value is used, eg. for serialization
type ErrBadValue struct{}