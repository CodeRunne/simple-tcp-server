package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	listener   net.Listener
	quitch     chan struct{}
	msgChannel chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgChannel: make(chan Message, 10),,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	ErrorHandler(err)
	defer listener.Close()

	s.listener = listener

	go s.AcceptConnection()
	<-s.quitch
	close(s.msgChannel)

	return nil
}

func (s *Server) AcceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Accept error", err)
			continue
		}

		fmt.Println("New connection to the server", conn.RemoteAddr())

		go s.ReadConnection(conn)
	}
}

func (s *Server) ReadConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read error", err)
			continue
		}

		s.msgChannel <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}
	}
}

func ErrorHandler(err error) error {
	if err != nil {
		return err
	}

	return nil
}

func main() {
	server := NewServer(":3000")
	go func() {
		for msg := range server.msgChannel {
			fmt.Printf("Received message from conection (%s): %s\n", msg.from, string(msg.payload))
		}
	}()
	log.Fatal(server.Start())
}
