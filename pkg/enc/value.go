package enc

import "io"

type Value interface {
	Type() ValueType
	String() string
	Encode(w io.Writer) error
}

type ValueType string

const (
	ValueTypeArray        = "array"
	ValueTypeSimpleString = "simple_string"
	ValueTypeBulkString   = "bulk_string"
)
