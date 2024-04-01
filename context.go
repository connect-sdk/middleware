package middleware

// ContextKey represents a context key.
type ContextKey struct {
	// Name is the name of the context key.
	Name string
}

// String returns the context key as a string.
func (k *ContextKey) String() string {
	return k.Name
}
