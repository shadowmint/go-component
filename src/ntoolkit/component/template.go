package component

// ObjectTemplate is a simple, flat, serializable object structure that directly converts to and from Objects.
type ObjectTemplate struct {
	Name       string
	Components []ComponentTemplate
	Objects    []ObjectTemplate
}

// ComponentTemplate is a serializable representation of a component
type ComponentTemplate struct {
	Type string
	Data string
}

