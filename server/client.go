package server

import (
	"fmt"
	"io"
	"log"
	"golang.org/x/net/websocket"
	"chattoo/user"
)

const channelBufSize = 100

type client struct {
	id       int64
	ip       string
	user     user.User
	ws       *websocket.Conn
	server   *Server
	history  []*message
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

	ch := make(chan *message, channelBufSize)
	doneCh := make(chan bool)
	ip := ws.Request().RemoteAddr
	var history []*message
	return &client{u.Id, ip, u, ws, server, history, ch, doneCh}
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
	go c.listenOut()
	c.listenIn()
}

// Listening write to client
func (c *client) listenOut() {
	for {
		select {
		case msg := <-c.messages:
			websocket.JSON.Send(c.ws, msg)
			log.Println("Sent msg:", msg)

		case <-c.done:
			c.server.disconnect(c)
			c.done <- true
			return
		}
	}
}

// Listening read from client
func (c *client) listenIn() {
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

			if c.server.isConnected(msg.To.Id) {
				c.write(&msg)
				c.history = append(c.history, &msg)
				c.server.send(&msg)
			}
		}
	}
}
