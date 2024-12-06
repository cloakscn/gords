package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/cloakscn/gords/message"
)



type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) Read() (message.Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return message.Value{}, err
	}

	switch _type {
	case message.ARRAY.Type:
		return r.readArray()
	case message.BULK.Type:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return message.Value{}, nil
	}
}

func (r *Resp) readArray() (message.Value, error) {
	v := message.Value{}
	v.Typ = message.ARRAY.Str

	// read length of array
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the message.value
	v.Array = make([]message.Value, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// add parsed message.value to array
		v.Array[i] = val
	}

	return v, nil
}

func (r *Resp) readBulk() (message.Value, error) {
	v := message.Value{}

	v.Typ = message.BULK.Str

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.Bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v message.Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
