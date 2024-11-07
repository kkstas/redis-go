package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/kkstas/redis-go/internal/config"
	"github.com/kkstas/redis-go/internal/parser"
)

var _ = net.Listen
var _ = os.Exit

type Store interface {
	Set(key, val string)
	SetWithExpiry(key, val string, expiry int)
	Get(key string) (string, bool)
}

type Server struct {
	config *config.Config
	store  Store
}

func New(config *config.Config, store Store) *Server {
	return &Server{
		config: config,
		store:  store,
	}
}

func Run(s *Server, network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println("error reading from connection:", err.Error())
			}
			break
		}

		s.executeCommand(conn, string(buf[:n]))
	}
}

func (s *Server) executeCommand(conn net.Conn, inputStr string) error {
	input, err := parser.Parse(inputStr)
	if err != nil {
		return fmt.Errorf("parsing command failed: %w", err)
	}
	cmd := parser.GetCommand(input)

	switch cmd {
	case "ping":
		return s.ping(conn)
	case "echo":
		return s.echo(conn, input[1])
	case "set":
		return s.set(conn, input[1], input[2])
	case "set_expiry":
		return s.setWithExpiry(conn, input[1], input[2], input[4])
	case "get":
		return s.get(conn, input[1])
	case "config_get":
		return s.configGet(conn, input[2])
	}

	log.Println("no command found for input", inputStr)
	return nil
}
