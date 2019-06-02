package framework

import (
	"strings"
)

// A Binding is a textual representation of a specific method of a type within a package.
type Binding struct {
	Type    string
	Package string
	Method  string
}

// Service returns the fully qualified name of the type (package + type name)
func (b Binding) Service() string {
	return b.Package + "." + b.Type
}

func (b Binding) String() string {
	return b.Service() + "::" + b.Method
}

// NewBinding returns a new instance of a Binding given a method definition in of the form <package>.<type>::<method>
func NewBinding(methodDef string) Binding {
	parts := strings.Split(methodDef, "::")
	nameparts := strings.Split(parts[0], ".")

	return Binding{
		Method:  parts[1],
		Package: nameparts[0],
		Type:    nameparts[1],
	}
}
