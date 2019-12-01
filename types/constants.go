package types

// Some constants used internally by the lazyledger prototype
const namespaceSize = 8

var codedNamespace [namespaceSize]byte

func init() {
	for i := range codedNamespace {
		codedNamespace[i] = 0xFF
	}
}
