package server

import (
	"log"
	"golang.org/x/net/websocket"
	"chattoo/user"
)

type Server struct {
	clients      map[int64]*client
	connected    chan *client
	disconnected chan *client
	incoming     chan *message
	done         chan bool
	errors       chan error
}

// Create new ws server.
func NewServer() *Server {
	clients := make(map[int64]*client)
	connected := make(chan *client)
	disconnected := make(chan *client)
	incoming := make(chan *message)
	done := make(chan bool)
	errors := make(chan error)

	return &Server{
		clients,
		connected,
		disconnected,
		incoming,
		done,
		errors,
	}
}

func (s *Server) connect(c *client) {
	s.connected <- c
}

func (s *Server) disconnect(c *client) {
	s.disconnected <- c
}

func (s *Server) send(msg *message) {
	s.incoming <- msg
}

func (s *Server) shutdown() {
	s.done <- true
}

func (s *Server) err(err error) {
	s.errors <- err
}

func (s *Server) sendHistory(c *client) {
	for _, msg := range s.clients[c.id].history {
		c.write(msg)
	}
}

func (s *Server) isConnected(id int64) bool {
	_, ok := s.clients[id]
	return ok
}

func (s *Server) sendToRecipient(msg *message) {
	s.clients[msg.To.Id].write(msg)
}

// listen and serve.
// It serves client connection and sendAll request.
func (s *Server) Listen() {
	log.Println("Listening server...")

	for {
		select {
		case c := <-s.connected:
			s.clients[c.id] = c
			s.sendHistory(c)
			log.Println("Connected new client:", c.id, ". Now", len(s.clients), "clients connected.")

		case c := <-s.disconnected:
			delete(s.clients, c.id)
			log.Println("Disconnected client: ", c.id)

		case msg := <-s.incoming:
			s.sendToRecipient(msg)
			log.Println("New message:", msg)

		case err := <-s.errors:
			log.Println("Error:", err.Error())

		case <-s.done:
			log.Println("shutdown server")
			return
		}
	}
}

func (s *Server) HandleWS(u user.User) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		client := newClient(ws, s, u)
		s.connect(client)
		client.listen()
	})
}
