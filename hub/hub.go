package hub

import (
	"github.com/gorilla/websocket"
)

type Connection struct {
	// The websocket Connection.
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan interface{}
}

func (c *Connection) Reader() {
	for {
		_, _, err := c.Ws.ReadMessage()
		if err != nil {
			break
		}
	}
	c.Ws.Close()
}

func (c *Connection) Writer() {
	for change := range c.Send {
		err := c.Ws.WriteJSON(change)
		if err != nil {
			break
		}
	}
	c.Ws.Close()
}

type hub struct {
	// Registered Connections.
	Connections map[*Connection]bool

	// Inbound messages from the Connections.
	Broadcast chan interface{}

	// Register requests from the Connections.
	Register chan *Connection

	// Unregister requests from Connections.
	Unregister chan *Connection
}

func NewHub() hub {
	return hub{
		Broadcast:   make(chan interface{}),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
		Connections: make(map[*Connection]bool),
	}
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Connections[c] = true
		case c := <-h.Unregister:
			if _, ok := h.Connections[c]; ok {
				delete(h.Connections, c)
				close(c.Send)
			}
		case m := <-h.Broadcast:
			for c := range h.Connections {
				select {
				case c.Send <- m:
				default:
					delete(h.Connections, c)
					close(c.Send)
				}
			}
		}
	}
}
