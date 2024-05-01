package auth

import "slices"

type Permission string

const (
	NamespaceRead   Permission = "namespace:read"
	NamespaceWrite  Permission = "namespace:write"
	NamespaceDelete Permission = "namespace:delete"
)

// All returns an array with all [Permission] values.
func All() []Permission {
	return []Permission{
		NamespaceRead,
		NamespaceWrite,
		NamespaceDelete,
	}
}

// Report reports whether an array of permissions i has a permission t.
func Report(i []Permission, t Permission) bool {
	return slices.Contains(i, t)
}
