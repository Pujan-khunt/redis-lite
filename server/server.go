package server

import (
	"fmt"
	"io"
	"log"
	"net"

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

	for {
		value, err := respReader.Read()
		if err != nil {
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

		switch command {
		case "SET":
			if len(value.Array) != 3 {
				fmt.Fprintf(conn, "-ERR invalid number of arguments for 'SET' command\r\n")
				continue
			}
			key, val := value.Array[1].Str, value.Array[2].Str
			s.store.Set(key, val)
			fmt.Fprintf(conn, "+OK\r\n")
		case "GET":
			if len(value.Array) != 2 {
				fmt.Fprintf(conn, "-ERR invalid number of arguments for 'GET' command\r\n")
				continue
			}
			key := value.Array[1].Str
			if val, ok := s.store.Get(key); ok {
				fmt.Fprintf(conn, "%s\r\n", val)
			} else {
				fmt.Fprintf(conn, "$-1\r\n")
			}
		case "DEL":
			if len(value.Array) != 2 {
				fmt.Fprintf(conn, "-ERR invalid number of arguments for 'DEL' command\r\n")
				continue
			}
			key := value.Array[1].Str
			if ok := s.store.Del(key); ok {
				fmt.Fprintf(conn, ":1\r\n")
			} else {
				fmt.Fprintf(conn, "-ERR failed to delete key")
			}
		default:
			fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", value.Array[0].Str)
		}
	}
}
