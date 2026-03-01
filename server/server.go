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

		switch command {
		case "SET":
			if len(value.Array) != 3 {
				respWriter.Write(resp.RespValue{Type: resp.Error, Str: "-ERR invalid number of arguments for 'SET' command"})
				continue
			}
			key, val := value.Array[1].Str, value.Array[2].Str
			s.store.Set(key, val)
			respWriter.Write(resp.RespValue{Type: resp.SimpleString, Str: "OK"})
		case "GET":
			if len(value.Array) != 2 {
				respWriter.Write(resp.RespValue{Type: resp.Error, Str: "-ERR invalid number of arguments for 'GET' command\r\n"})
				continue
			}
			key := value.Array[1].Str
			if val, ok := s.store.Get(key); ok {
				respWriter.Write(resp.RespValue{Type: resp.BulkString, Str: val})
			} else {
				respWriter.Write(resp.RespValue{Type: resp.BulkString, Str: ""})
			}
		case "DEL":
			if len(value.Array) != 2 {
				respWriter.Write(resp.RespValue{Type: resp.Error, Str: "-ERR invalid number of arguments for 'DEL' command\r\n"})
				continue
			}
			key := value.Array[1].Str
			if ok := s.store.Del(key); ok {
				respWriter.Write(resp.RespValue{Type: resp.Integer, Num: 1})
			} else {
				respWriter.Write(resp.RespValue{Type: resp.Error, Str: "-ERR failed to delete key"})
			}
		default:
			msg := fmt.Sprintf("-ERR unknown command '%s'", value.Array[0].Str)
			respWriter.Write(resp.RespValue{Type: resp.Error, Str: msg})
		}
	}
}
