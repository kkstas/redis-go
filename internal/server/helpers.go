package server

import (
	"fmt"
	"log"
	"net"
)

func writeRaw(conn net.Conn, msg string) error {
	log.Printf("writing %q", msg)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("writing failed: %w", err)
	}
	return nil
}
