package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type RespReader struct {
	reader *bufio.Reader
}

func NewRespReader(r io.Reader) *RespReader {
	return &RespReader{
		reader: bufio.NewReader(r),
	}
}

// Read peeks at the first byte to determine the type of the RESP message and then delegates the rest to specific parsing methods.
func (r *RespReader) Read() (RespValue, error) {
	prefix, err := r.reader.ReadByte()
	if err != nil {
		return RespValue{}, err
	}
	switch RespType(prefix) {
	case Array:
		return r.readArray()
	case BulkString:
		return r.readBulkString()
	case Integer:
		return r.readInteger()
	case SimpleString:
		return r.readSimpleString()
	default:
		msg := fmt.Sprintf("unknown RESP type: %c\n", prefix)
		return RespValue{}, errors.New(msg)
	}
}

func (r *RespReader) readSimpleString() (RespValue, error) {
	line, _, err := r.readLine()
	if err != nil {
		return RespValue{}, nil
	}
	return RespValue{
		Type: SimpleString,
		Str:  string(line),
	}, nil
}

func (r *RespReader) readInteger() (RespValue, error) {
	line, _, err := r.readLine()
	if err != nil {
		return RespValue{}, err
	}
	num, err := strconv.Atoi(string(line))
	if err != nil {
		return RespValue{}, err
	}
	return RespValue{
		Type: Integer,
		Num:  num,
	}, nil
}

func (r *RespReader) readBulkString() (RespValue, error) {
	line, _, err := r.readLine()
	if err != nil {
		return RespValue{}, err
	}
	length, err := strconv.Atoi(string(line))
	if err != nil {
		return RespValue{}, err
	}
	// Handle null bulk strings
	if length == -1 {
		return RespValue{
			Type: BulkString,
			Str:  "",
		}, nil
	}
	// Read exactly `length` bytes
	bulkStr := make([]byte, length)
	_, err = r.reader.Read(bulkStr)
	if err != nil {
		return RespValue{}, err
	}
	// Consume trailing `\r\n`
	_, _, err = r.readLine()
	if err != nil {
		return RespValue{}, err
	}
	return RespValue{
		Type: BulkString,
		Str:  string(bulkStr),
	}, nil
}

// readLine is a helper to read until '\r\n'
func (r *RespReader) readLine() (line []byte, n int, err error) {
	for {
		byte, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		line = append(line, byte)
		n++
		if len(line)-2 >= 0 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *RespReader) readArray() (RespValue, error) {
	// Read byte containing length of array
	line, _, err := r.readLine()
	if err != nil {
		return RespValue{}, err
	}
	// Convert length from byte to integer
	length, err := strconv.Atoi(string(line))
	if err != nil {
		return RespValue{}, err
	}
	// Create RESP value with type array of specified length
	val := RespValue{Type: Array}
	val.Array = make([]RespValue, length)
	// Loop exactly `length` times and recursively parse rest of the messages.
	for i := range length {
		element, err := r.Read()
		if err != nil {
			return RespValue{}, err
		}
		val.Array[i] = element
	}
	return val, nil
}
