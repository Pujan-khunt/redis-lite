package aof

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Pujan-khunt/redis-lite/resp"
)

type AOF struct {
	file *os.File
	r    *bufio.Reader
	w    *bufio.Writer
	mu   sync.Mutex
}

func NewAOF(period time.Duration) (*AOF, error) {
	f, err := os.OpenFile("appendonly.aof", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	aof := &AOF{
		file: f,
		r:    bufio.NewReader(f),
		w:    bufio.NewWriter(f),
	}
	// Run fsync to flush to disk periodically
	go func() {
		for {
			aof.mu.Lock()
			aof.Flush()
			aof.mu.Unlock()
			time.Sleep(period)
		}
	}()
	return aof, nil
}

// Flush flushes file to disk
func (a *AOF) Flush() {
	// Flush buffer(bufio) onto os.file(io.Writer)
	a.w.Flush()
	// Flush os.File(io.Writer) onto disk.
	a.file.Sync()
}

// Close closes the connection to the file, ensuring one last flush before closing.
func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Flush()
	return a.file.Close()
}

func (a *AOF) Append(value resp.RespValue) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, err := a.w.WriteString(serialize(value))
	return err
}

func serialize(value resp.RespValue) string {
	switch value.Type {
	case resp.SimpleString:
		return fmt.Sprintf("S\r\n%s\r\n", value.Str)
	case resp.Error:
		return fmt.Sprintf("E\r\n%s\r\n", value.Str)
	case resp.Integer:
		return fmt.Sprintf("I\r\n%d\r\n", value.Num)
	case resp.BulkString:
		return fmt.Sprintf("B\r\n%d\r\n%s\r\n", len(value.Str), value.Str)
	case resp.Array:
		var val strings.Builder
		fmt.Fprintf(&val, "A\r\n%d\r\n", len(value.Array))
		for _, respVal := range value.Array {
			val.WriteString(serialize(respVal))
		}
		return val.String()
	}
	return ""
}
