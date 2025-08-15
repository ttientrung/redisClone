package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Job struct {
	client Client
}

type Worker struct {
	id      int
	jobChan chan Job
}

type Pool struct {
	jobQueue chan Job
	workers  []*Worker
}

func NewWorker(id int, jobChan chan Job) *Worker {
	return &Worker{
		id:      id,
		jobChan: jobChan,
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobChan {
			fmt.Printf("Worker %d is handling job from %s", w.id, job.client.conn.RemoteAddr())
			handleRequest(job.client)
		}
	}()
}

func NewPool(numOfWorker int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers:  make([]*Worker, numOfWorker),
	}
}

func (p *Pool) AddJob(client Client) {
	p.jobQueue <- Job{client: client}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue)
		p.workers[i] = worker
		worker.Start()
	}
}

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

	pool := NewPool(2)
	pool.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := Client{
			conn: conn,
		}

		pool.AddJob(client)
		// go client.handleRequest()
	}
}

func handleRequest(client Client) {
	for {
		log.Println(client.conn.RemoteAddr())
		var buf []byte = make([]byte, 1000)
		_, err := client.conn.Read(buf)
		if err != nil {
			client.conn.Close()
			log.Fatal(err)
		}
		// process
		time.Sleep(time.Second * 10)
		client.conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nMessage received.\n"))
	}
}
