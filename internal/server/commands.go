package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/kkstas/redis-go/internal/parser"
)

const (
	nullBulkString = "$-1\r\n"
)

func (s *Server) ping(conn net.Conn) error {
	err := writeRaw(conn, parser.ToSimpleString("PONG"))
	if err != nil {
		return fmt.Errorf("PING failed: %w", err)
	}
	return nil
}

func (s *Server) echo(conn net.Conn, str string) error {
	err := writeRaw(conn, parser.ToSimpleString(str))
	if err != nil {
		return fmt.Errorf("ECHO failed: %w", err)
	}
	return nil
}

func (s *Server) set(conn net.Conn, key, val string) error {
	s.store.Set(key, val)
	err := writeRaw(conn, parser.ToSimpleString("OK"))
	if err != nil {
		return fmt.Errorf("writing after SET failed: %w", err)
	}
	return nil
}

func (s *Server) setWithExpiry(conn net.Conn, key, val string, expiry string) error {
	expiryInMS, err := strconv.Atoi(expiry)
	if err != nil {
		return fmt.Errorf("error while parsing expiry amount: %w", err)
	}

	s.store.SetWithExpiry(key, val, expiryInMS)

	err = writeRaw(conn, parser.ToSimpleString("OK"))
	if err != nil {
		return fmt.Errorf("writing after SET failed: %w", err)
	}
	return nil
}

func (s *Server) get(conn net.Conn, key string) error {
	val, found := s.store.Get(key)
	if !found {
		err := writeRaw(conn, nullBulkString)
		if err != nil {
			return fmt.Errorf("writing after GET ( not found ) failed: %w", err)
		}
		return nil
	}
	err := writeRaw(conn, parser.ToBulkString(val))
	if err != nil {
		return fmt.Errorf("writing after GET ( found ) failed: %w", err)
	}
	return nil
}

func (s *Server) configGet(conn net.Conn, key string) error {
	if key == "dir" {
		err := writeRaw(conn, parser.ToRESPArray([]string{"dir", s.config.Dir}))
		if err != nil {
			return fmt.Errorf("writing after GET ( dir ) failed: %w", err)
		}
		return nil
	}
	if key == "dbfilename" {
		err := writeRaw(conn, parser.ToRESPArray([]string{"dbfilename", s.config.Dir}))
		if err != nil {
			return fmt.Errorf("writing after GET ( dir ) failed: %w", err)
		}
		return nil
	}
	return nil
}
