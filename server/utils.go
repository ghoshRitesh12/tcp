package server

import (
	"log"
	"net"
	"sync/atomic"
)

// Gets the buffer from the socket.
// Default buffer size is 3048
func getBuffer(c net.Conn) (string, error) {
	buffer := make([]byte, 3048)

	n, err := c.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil
}

// Sends the response for a subsequent request
func sendResponse(c net.Conn, msg string) error {
	if _, err := c.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}

// Atomically updates the number of concurrent users
func updateUsersStatus(c net.Conn, status string, allUsers *int32) {
	if status == "-" {
		status = "LEFT"
		atomic.AddInt32(allUsers, -1)
	} else if status == "+" {
		atomic.AddInt32(allUsers, 1)
		status = "JOINED"
	}

	log.Println("USER", c.RemoteAddr(), status, "\nTotal concurrent users: ", atomic.LoadInt32(allUsers))
}
