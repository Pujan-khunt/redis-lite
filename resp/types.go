package resp

// RespType defines the RESP data type using the prefix byte
type RespType byte

const (
	SimpleString RespType = '+'
	Error        RespType = '-'
	Integer      RespType = ':'
	BulkString   RespType = '$'
	Array        RespType = '*'
)

// RespValue represents a generic RESP value.
// It holds data depending upon what the value of 'Type'.
type RespValue struct {
	Type  RespType
	Str   string      // Used for `SimpleString`, `Error` and `BulkString`
	Num   int         // Used for `Integer`
	Array []RespValue // Used for `Arrays` (recursive)
}
