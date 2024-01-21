package server

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	port             string
	ln               net.Listener
	totalConnections int32
	q                chan struct{}
}

func NewServer(port string) *Server {
	return &Server{
		port: ":" + port,
		q:    make(chan struct{}),
	}
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		log.Println("listen error", err.Error())
		return err
	}
	defer listener.Close()

	s.ln = listener

	log.Println("echo server @ localhost" + s.port)

	go s.acceptConnections()
	<-s.q

	return nil
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.ln.Accept() // blocking call
		if err != nil {
			log.Println("conn error", err.Error())
			continue
		}

		updateUsersStatus(conn, "+", &s.totalConnections)

		go s.parseRequest(conn)
	}
}

func (s *Server) parseRequest(c net.Conn) {
	for {
		msg, err := getBuffer(c) // blocking call
		if err != nil {
			if err == io.EOF {
				updateUsersStatus(c, "-", &s.totalConnections)
				c.Close()
				break
			}

			log.Println(err.Error())
			continue
		}

		fmt.Print(c.RemoteAddr(), " -> ", msg)
		if err := sendResponse(c, "-> "+msg); err != nil {
			log.Println(err.Error())
		}
	}
}
