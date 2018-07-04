package server

import (
	"log"
	"golang.org/x/net/websocket"
	"chattoo/user"
	"sync"
)

const historySize = 100

type Server struct {
	history      []message
	clients      sync.Map
	connected    chan *client
	disconnected chan *client
	broadcast    chan *message
	incoming     chan *message
	done         chan bool
	errors       chan error
}

// Create new ws server.
func NewServer() *Server {
	return &Server{
		history:      make([]message, 0),
		clients:      sync.Map{},
		connected:    make(chan *client),
		disconnected: make(chan *client),
		broadcast:    make(chan *message),
		incoming:     make(chan *message),
		done:         make(chan bool),
		errors:       make(chan error),
	}
}

func (s *Server) connect(c *client) {
	s.connected <- c
}

func (s *Server) disconnect(c *client) {
	s.disconnected <- c
}

func (s *Server) send(msg *message) {
	switch msg.Type {
	case privateMessageType:
		s.incoming <- msg
	case publicMessageType:
		s.broadcast <- msg
	}
}

func (s *Server) shutdown() {
	s.done <- true
}

func (s *Server) err(err error) {
	s.errors <- err
}

// Listen and serve.
func (s *Server) Listen() {
	log.Println("Listening server...")

	for {
		select {
		case c := <-s.connected:
			s.clients.Store(c.id, c)
			for _, msg := range s.history {
				websocket.JSON.Send(c.ws, msg)
			}
			log.Println("Connected new client:", c.id)

		case c := <-s.disconnected:
			s.clients.Delete(c.id)
			log.Println("Disconnected client:", c.id)

		case msg := <-s.broadcast:
			s.clients.Range(func(key, value interface{}) bool {
				c := value.(*client)
				c.write(msg)
				return true
			})

			// Append and trim history to `historySize`
			s.history = append(s.history, *msg)
			offset := len(s.history)-historySize
			if offset > 0 {
				s.history = s.history[offset:]
			}

			log.Println("New Broadcast message:", msg)

		case msg := <-s.incoming:
			r, receiverFound := s.clients.Load(msg.To.Id)
			s, senderFound := s.clients.Load(msg.From.Id)
			if receiverFound && senderFound {
				receiver := r.(*client)
				receiver.write(msg)

				sender := s.(*client)
				sender.write(msg)

				log.Println("New private message:", msg)
			}

		case err := <-s.errors:
			log.Println("Error:", err.Error())

		case <-s.done:
			log.Println("Shutdown server")
			return
		}
	}
}

// WebSocket handler.
func (s *Server) HandleWS(u user.User) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		client := newClient(ws, s, u)
		s.connect(client)
		client.listen()
	})
}
