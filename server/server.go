package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

// Server ...
type Server struct {
	host string
	port string
}

// Client ...
type Client struct {
	conn net.Conn
}

// Config ...
type Config struct {
	Host string
	Port string
}

// New ...
func New(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

// Run ...
func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go client.handleRequest()
	}
}

func (client *Client) handleRequest() {
	for {
		log.Println(client.conn.RemoteAddr())
		var buf []byte = make([]byte, 1000)
		_, err := client.conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		// process
		time.Sleep(time.Second * 10)
		client.conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nMessage received.\n"))
	}
}