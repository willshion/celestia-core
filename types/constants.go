package types

// Some constants used internally by the lazyledger prototype

// TODO(LL): this should rather be configurable:
const NamespaceSize = 8

const DefaultMessageSize = 512

var codedNamespace [NamespaceSize]byte

func init() {
	for i := range codedNamespace {
		codedNamespace[i] = 0xFF
	}
}
