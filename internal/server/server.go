package server

import (
	"io"
	"log"
	"net"
	"redisClone/internal/config"
	"redisClone/internal/core/io_multiplexing"
	"syscall"
)

func readCommand(fd int) (string, error) {
	var buf = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", io.EOF
	}
	return string(buf[:n]), nil
}

func respond(data string, fd int) error {
	if _, err := syscall.Write(fd, []byte(data)); err != nil {
		return err
	}
	return nil
}

func RunIoMultiplexingServer() {
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Print("Server is running on ", config.Port)
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("Failed to cast to TCPListener")
	}

	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal("Failed to get listener file")
	}
	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())

	ioMultiplexer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal("Failed to create IO Multiplexer:", err)
	}
	defer ioMultiplexer.Close()

	err = ioMultiplexer.Monitor(io_multiplexing.Event{
		Fd: serverFd,
		Op: io_multiplexing.OpRead,
	})
	if err != nil {
		log.Fatal("Failed to monitor server fd:", err)
	}

	var events = make([]io_multiplexing.Event, config.MaxConnection)
	for {
		events, err = ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		for i := 0; i < len(events); i++ {
			if events[i].Fd == serverFd {
				if events[i].Fd == serverFd {
					log.Printf("new client is trying to connect")
					connFd, _, err := syscall.Accept(serverFd)
					if err != nil {
						log.Println("Failed to accept new connection:", err)
						continue
					}
					log.Printf("set up new connection")
					err = ioMultiplexer.Monitor(io_multiplexing.Event{
						Fd: connFd,
						Op: io_multiplexing.OpRead,
					})
					if err != nil {
						log.Fatal("Failed to monitor client fd:", err)
					}
				}
			} else {
				cmd, err := readCommand(events[i].Fd)
				if err != nil {
					if err == io.EOF {
						log.Printf("Client disconnected: fd %d", events[i].Fd)
						_ = syscall.Close(events[i].Fd)
						continue
					}
					log.Printf("read error")
					continue
				}
				err = respond(cmd, events[i].Fd)
				if err != nil {
					log.Printf("respond error")
				}
			}
		}
	}
}
