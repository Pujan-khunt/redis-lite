package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Pujan-khunt/redis-lite/aof"
	"github.com/Pujan-khunt/redis-lite/resp"
	"github.com/Pujan-khunt/redis-lite/storage"
)

type Server struct {
	addr  *net.TCPAddr
	store storage.Store
	aof   *aof.AOF
}

func NewServer(host string, port int, store storage.Store, period time.Duration) (*Server, error) {
	aof, err := aof.NewAOF(period)
	if err != nil {
		return nil, fmt.Errorf("failed to start server: %w", err)
	}
	return &Server{
		addr: &net.TCPAddr{
			IP:   net.ParseIP(host),
			Port: port,
		},
		store: store,
		aof:   aof,
	}, nil
}

func (s *Server) ListenAndServe() error {
	listener, err := net.ListenTCP("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Printf("Listening for connections on %s:%d\r\n", s.addr.IP, s.addr.Port)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatalf("Failed to accept connection: %v\r\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	respReader := resp.NewRespReader(conn)
	respWriter := resp.NewRespWriter(conn)
	for {
		value, err := respReader.Read()
		s.aof.Append(value)
		if err != nil {
			// Client disconnect
			if err == io.EOF {
				break
			}
			log.Println("Error parsing RESP:", err)
			return
		}
		s.handleCommand(value, respWriter)
	}
}

// handleCommand finds the correct handler for the command and executes it.
func (s *Server) handleCommand(value resp.RespValue, w *resp.RespWriter) {
	// Expect command to only be a RESP array of positive length.
	if value.Type != resp.Array || len(value.Array) == 0 {
		return
	}
	command := value.Array[0].Str
	normalizedCommand := strings.ToUpper(command)
	if handler, exists := commandRegistry[normalizedCommand]; exists {
		handler(value.Array, w, s.store)
		// Append only if handler exists
		if err := s.aof.Append(value); err != nil {
			fmt.Printf("Failed to save to append only file: %s\n", normalizedCommand)
		}
	} else {
		msg := fmt.Sprintf("-ERR unknown command '%s'", value.Array[0].Str)
		w.Write(resp.RespValue{Type: resp.Error, Str: msg})
	}

}
