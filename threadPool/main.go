package main

import (
	"fmt"
	"redisClone/threadPool/server"
)

func main() {
	config := &server.Config{
		Host: "localhost",
		Port: "8000",
	}

	srv := server.New(config)
	fmt.Printf("Starting server on %s:%s\n", config.Host, config.Port)
	srv.Run()
}
