package resp

import (
	"io"
	"strconv"
)

type RespWriter struct {
	writer io.Writer
}

func NewRespWriter(w io.Writer) *RespWriter {
	return &RespWriter{
		writer: w,
	}
}

func (w *RespWriter) Write(value RespValue) error {
	var bytes []byte

	switch value.Type {
	case SimpleString:
		bytes = w.writeSimpleString(value.Str)
	case Error:
		bytes = w.writeError(value.Str)
	case Integer:
		bytes = w.writeInteger(value.Num)
	case BulkString:
		bytes = w.writeBulkString(value.Str)
	case Array:
		bytes = w.writeArray(value.Array)
	default:
		bytes = w.writeError("unknown value type")
	}
	_, err := w.writer.Write(bytes)
	return err
}

func (w *RespWriter) writeArray(array []RespValue) []byte {
	if array == nil {
		return []byte("*-1\r\n")
	}
	bytes := []byte("*" + strconv.Itoa(len(array)) + "\r\n")
	for _, val := range array {
		var childBytes []byte
		switch val.Type {
		case SimpleString:
			childBytes = w.writeSimpleString(val.Str)
		case Error:
			childBytes = w.writeError(val.Str)
		case Integer:
			childBytes = w.writeInteger(val.Num)
		case BulkString:
			childBytes = w.writeBulkString(val.Str)
		case Array:
			childBytes = w.writeArray(val.Array)
		}
		bytes = append(bytes, childBytes...)
	}
	return bytes
}

func (w *RespWriter) writeBulkString(str string) []byte {
	// Null bulk string
	if str == "" {
		return []byte("$-1\r\n")
	}
	return []byte("$" + strconv.Itoa(len(str)) + "\r\n" + str + "\r\n")
}

func (w *RespWriter) writeInteger(num int) []byte {
	return []byte(":" + strconv.Itoa(num) + "\r\n")
}

func (w *RespWriter) writeError(str string) []byte {
	return []byte("-" + str + "\r\n")
}

func (w *RespWriter) writeSimpleString(str string) []byte {
	return []byte("+" + str + "\r\n")
}
