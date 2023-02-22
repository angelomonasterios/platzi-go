package webSocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could Not Open WebSocket Connection", http.StatusInternalServerError)
	}
	client := NewClient(hub, socket)
	hub.register <- client
	go client.Write()
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnnect(client)
		case client := <-hub.unregister:
			hub.onDisconect(client)
		}
	}
}

func (hub *Hub) onConnnect(client *Client) {
	log.Println("client connected", client.socket.RemoteAddr())
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	client.id = client.socket.RemoteAddr().String()
	hub.clients = append(hub.clients, client)
}

func (hub *Hub) onDisconect(client *Client) {
	log.Println("client Disconnected", client.socket.RemoteAddr())
	client.socket.Close()
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	client.id = client.socket.RemoteAddr().String()
	hub.clients = append(hub.clients, client)
	i := -1

	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
		}
	}

	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]
}

func (hub *Hub) BroadCast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if client != ignore {
			client.outBound <- data
		}
	}
}
