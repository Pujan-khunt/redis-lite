package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

var (
	port = flag.Uint("port", 6379, "port for redis-lite to listen on")
	ip   = flag.String("ip", "0.0.0.0", "ip of the interface for the server to listen on")
)

func main() {
	flag.Parse()

	// TCP Address of the server
	addr := net.TCPAddr{
		IP:   net.ParseIP(*ip),
		Port: int(*port),
	}

	// Listen for TCP connections
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening for connections on %s\n", addr.IP.String()+":"+strconv.Itoa(addr.Port))

	// In memory database
	db := make(map[string]string)

	// Accept the connection
	conn, err := listener.AcceptTCP()
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Create reader to read from client
		r := bufio.NewReader(conn)
		msg, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("error big one", err)
			conn.Close()
			break
		}

		// Parse client message into command handling both LF(\n) and CRLF(\r\n)
		msg = strings.TrimSpace(msg)
		cmd := strings.Split(msg, " ")

		// Execute logic (Store, Update, Delete) based on the command
		switch cmd[0] {
		case "SET":
			key, val := cmd[1], cmd[2]
			db[key] = val
			fmt.Fprintf(conn, "SET Operation Performed with Key: %s and Val: %s\n", key, val)
		case "GET":
			key := cmd[1]
			if value, ok := db[key]; ok {
				fmt.Fprint(conn, value)
			} else {
				fmt.Fprintf(conn, "No key exists. Key: %s\n", key)
			}
		case "DEL":
			key := cmd[1]
			delete(db, key)
			fmt.Fprintf(conn, "Deleted key value pair where key: %s\n", key)
		default:
			fmt.Fprintf(conn, "Invaild command used: %s\n", cmd[0])
		}
	}
}
