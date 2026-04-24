package enc

import "io"

type Value interface {
	Type() Type
	String() string
	Encode(w io.Writer) error
}

type Type string

const (
	TypeArray        Type = "array"
	TypeSimpleString Type = "simple_string"
	TypeSimpleError  Type = "simple_error"
	TypeBulkString   Type = "bulk_string"
	TypeInteger      Type = "integer"
)
