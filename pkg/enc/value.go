package enc

import "io"

type Value interface {
	Type() Type
	String() string
	Encode(w io.Writer) error
}

type Type string

const (
	TypeArray        = "array"
	TypeSimpleString = "simple_string"
	TypeBulkString   = "bulk_string"
)
