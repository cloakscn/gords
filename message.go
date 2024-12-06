package main

import "strconv"

type message struct {
	Type byte
	Str  string
}

var (
	STRING  = message{Type: '+', Str: "string"}
	ERROR   = message{Type: '-', Str: "error"}
	INTEGER = message{Type: ':', Str: "integer"}
	BULK    = message{Type: '$', Str: "bulk"}
	ARRAY   = message{Type: '*', Str: "array"}
	NULL    = message{Type: '.', Str: "null"}
)

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

func (v Value) Marshal() []byte {
	switch v.Typ {
	case ARRAY.Str:
		return v.marshalArray()
	case BULK.Str:
		return v.marshalBulk()
	case STRING.Str:
		return v.marshalString()
	case NULL.Str:
		return v.marshallNull()
	case ERROR.Str:
		return v.marshallError()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING.Type)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK.Type)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalArray() []byte {
	len := len(v.Array)
	var bytes []byte
	bytes = append(bytes, ARRAY.Type)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.Array[i].Marshal()...)
	}

	return bytes
}

func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR.Type)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}