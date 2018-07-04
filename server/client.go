package server

import (
	"fmt"
	"io"
	"log"
	"golang.org/x/net/websocket"
	"chattoo/user"
	"errors"
)

const channelBufSize = 100

type client struct {
	id       int64
	ip       string
	user     user.User
	ws       *websocket.Conn
	server   *Server
	messages chan *message
	done     chan bool
}

func newClient(ws *websocket.Conn, server *Server, u user.User) *client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	return &client{
		id:       u.Id,
		ip:       ws.Request().RemoteAddr,
		user:     u,
		ws:       ws,
		server:   server,
		messages: make(chan *message, channelBufSize),
		done:     make(chan bool),
	}
}

func (c *client) write(msg *message) {
	select {
	case c.messages <- msg:

	default:
		c.server.disconnect(c)
		err := fmt.Errorf("client %d is disconnected", c.user.Username)
		c.server.err(err)
	}
}

func (c *client) listen() {
	go c.listenReceiver()
	c.listenSender()
}

// Listening messages to client
func (c *client) listenReceiver() {
	for {
		select {
		case msg := <-c.messages:
			websocket.JSON.Send(c.ws, msg)
			log.Println("New incoming msg:", msg)

		case <-c.done:
			c.server.disconnect(c)
			c.done <- true
			return
		}
	}
}

// Listening messages from client
func (c *client) listenSender() {
	for {
		select {
		case <-c.done:
			c.server.disconnect(c)
			c.done <- true
			return

		default:
			var msg message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.done <- true
				return
			}

			if err != nil {
				c.server.err(err)
				return
			}

			if err = c.validate(msg); err != nil {
				c.server.err(err)
				break
			}

			c.server.send(&msg)
		}
	}
}

func (c *client) validate(msg message) error {
	if err := msg.validate(); err != nil {
		return err
	}

	if msg.Type == privateMessageType {
		_, found := c.server.clients.Load(msg.To.Id)
		if !found {
			return errors.New("client not found on server")
		}
	}

	return nil
}
