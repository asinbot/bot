package rcon

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Conn represents a connection to a HLL RCON server. A Conn supports multiple thread-safe
// connections.
type Conn struct {
	active  int  // Estimated active connections.
	closing bool // Closing fall to stop polling.

	pool sync.Pool // Collection of sessions.
}

type session struct {
	net.Conn
	key []byte // XOR key.
}

const msglen = 8196

var ErrResultFailed = errors.New("got FAIL response from server")

// New returns a new HLL RCON client to set/get server parameters.
func New(addr string, password string) (*Conn, error) {
	c := &Conn{}

	c.pool = sync.Pool{
		New: func() interface{} {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				return err
			}

			// Retrieve the XOR key used to encrypt communications between client/server.
			key := make([]byte, msglen)

			n, err := conn.Read(key)
			if err != nil {
				return err
			}

			key = key[:n]

			s := &session{
				Conn: conn,
				key:  key,
			}

			err = s.login(password)
			if err != nil {
				return err
			}

			if !c.closing {
				c.active++
			}

			return s
		},
	}

	switch s := c.pool.Get().(type) {
	case error:
		return nil, s
	case *session:
		c.pool.Put(s)
	}

	return c, nil
}

// Close will attempt to close all active connections held by the internal pool.
func (c *Conn) Close() error {
	c.closing = true

	for i := 0; i < c.active; i++ {
		switch s := c.pool.Get().(type) {
		case error:
			return s
		case *session:
			err := s.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Send will send a list of commands to a server and return the response.
// This should only be used when another server function is not explicity defined or implemented.
func (c *Conn) Send(cmds ...string) (string, error) {
	return c.send(cmds...)
}

func (c *Conn) send(cmds ...string) (string, error) {
	switch s := c.pool.Get().(type) {
	case error:
		return "", s
	case *session:
		defer c.pool.Put(s)

		return s.send(cmds...)
	}

	return "", fmt.Errorf("an unknown error has occured")
}

func (s *session) login(password string) error {
	result, err := s.send("login " + password)
	if err != nil {
		return err
	}

	if result != "SUCCESS" {
		return fmt.Errorf("rcon authentication failed for %s with key %x", s.RemoteAddr(), s.key)
	}

	return nil
}

func (s *session) send(cmds ...string) (string, error) {
	_, err := s.Write(s.xor([]byte(strings.Join(cmds, " "))))
	if err != nil {
		return "", err
	}

	b := make([]byte, msglen)

	n, err := s.Read(b)
	if err != nil {
		return "", err
	}

	result := string(s.xor(b)[:n])

	if result == "FAIL" {
		return "", ErrResultFailed
	}

	return result, nil
}

func (s *session) xor(b []byte) []byte {
	d := make([]byte, len(b))

	for i := range b {
		d[i] = b[i] ^ s.key[i%len(s.key)]
	}

	return d
}
