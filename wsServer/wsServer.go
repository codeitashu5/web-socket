package wsServer

import (
	"web-socket/models"
)

type WsServer struct {
	clients    map[*models.Client]bool
	register   chan *models.Client
	unRegister chan *models.Client
}

// NewWebSocketServer to crate a new webSocket server
func NewWebSocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*models.Client]bool),
		register:   make(chan *models.Client),
		unRegister: make(chan *models.Client),
	}
}

// Run creates and infinite loop that waits for chanel input
func (wsServer *WsServer) Run() {
	for {
		// used to wait on chanel operation
		select {
		// get the client that is getting registered
		case client := <-wsServer.register:
			wsServer.registerClient(client)

		case client := <-wsServer.unRegister:
			wsServer.unRegisterClient(client)
		}
	}
}

func (wsServer *WsServer) registerClient(client *models.Client) {
	wsServer.clients[client] = true
}

func (wsServer *WsServer) unRegisterClient(client *models.Client) {
	// you find this client in map and delete it
	delete(wsServer.clients, client)
}
