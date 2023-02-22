package webSocket

import "C"
import "github.com/gorilla/websocket"

type Client struct {
	hub      *Hub
	id       string
	socket   *websocket.Conn
	outBound chan []byte
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outBound: make(chan []byte),
	}
}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.outBound:
			if !ok {
				c.socket.Close()
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
