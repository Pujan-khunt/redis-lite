package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Pujan-khunt/redis-lite/storage"
)

var (
	ErrClientDisconnected error = errors.New("server: client disconnected")
	ErrInvalidNumArgs     error = errors.New("server: invalid number of arguments")
	ErrKeyNotExists       error = errors.New("server: key doesn't exist")
	ErrInvalidCommand     error = errors.New("server: invalid command")
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

	fmt.Printf("Listening for connections on %s:%d\n", s.addr.IP, s.addr.Port)

	conn, err := listener.AcceptTCP()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	errCh := make(chan error)
	go s.handleConnection(conn, errCh)
	return <-errCh
}

func (s *Server) handleConnection(conn *net.TCPConn, errCh chan<- error) error {
	for {
		r := bufio.NewReader(conn)
		msg, err := r.ReadString('\n')
		if err != nil {
			errCh <- ErrClientDisconnected
		}

		msg = strings.TrimSpace(msg)
		cmd := strings.Split(msg, " ")

		if len(cmd) == 0 || cmd[0] == "" {
			continue
		}

		switch cmd[0] {
		case "SET":
			if len(cmd) < 3 {
				errCh <- ErrInvalidNumArgs
				continue
			}
			key, val := cmd[1], cmd[2]
			s.store.Set(key, val)
			fmt.Fprintf(conn, "OK\n")
		case "GET":
			key := cmd[1]
			if val, ok := s.store.Get(key); ok {
				fmt.Fprintln(conn, val)
			} else {
				errCh <- ErrKeyNotExists
			}
		case "DEL":
			key := cmd[1]
			s.store.Del(key)
		default:
			errCh <- ErrInvalidCommand
		}
	}
}
