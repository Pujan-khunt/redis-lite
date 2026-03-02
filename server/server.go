package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/Pujan-khunt/redis-lite/resp"
	"github.com/Pujan-khunt/redis-lite/storage"
)

type Server struct {
	addr  *net.TCPAddr
	store storage.Store
}

func NewServer(host string, port int, store storage.Store) *Server {
	return &Server{
		addr: &net.TCPAddr{
			IP:   net.ParseIP(host),
			Port: port,
		},
		store: store,
	}
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
		if err != nil {
			// Client disconnect
			if err == io.EOF {
				break
			}
			log.Println("Error parsing RESP:", err)
			return
		}
		// Expect command to only be a RESP positive length array
		if value.Type != resp.Array || len(value.Array) == 0 {
			continue
		}
		command := value.Array[0].Str
		normalizedCommand := strings.ToUpper(command)
		if handler, exists := commandRegistry[normalizedCommand]; exists {
			handler(value.Array, respWriter, s.store)
		} else {
			msg := fmt.Sprintf("-ERR unknown command '%s'", value.Array[0].Str)
			respWriter.Write(resp.RespValue{Type: resp.Error, Str: msg})
		}
	}
}
