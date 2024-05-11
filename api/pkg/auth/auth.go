package auth

import (
	"slices"
	"strings"
)

type Permission string

func (p Permission) String() string {
	return string(p)
}

type Permissions []Permission

func (ps Permissions) String() string {
	str := ""
	for i, p := range ps {
		switch i {
		case 0:
			str = p.String()
		default:
			str = str + "-" + p.String()
		}
	}

	return str
}

func (ps Permissions) FromString(str string) Permissions {
	parts := strings.Split(str, "-")
	for _, p := range parts {
		ps = append(ps, Permission(p))
	}

	return ps
}

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
