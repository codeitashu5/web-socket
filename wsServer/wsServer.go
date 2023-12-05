package wsServer

import (
	"web-socket/models"
)

type WsServer struct {
	Clients    map[*models.Client]bool
	Register   chan *models.Client
	UnRegister chan *models.Client
}

// NewWebSocketServer to crate a new webSocket server
func NewWebSocketServer() *WsServer {
	return &WsServer{
		Clients:    make(map[*models.Client]bool),
		Register:   make(chan *models.Client),
		UnRegister: make(chan *models.Client),
	}
}

// Run creates and infinite loop that waits for chanel input
func (wsServer *WsServer) Run() {
	for {
		// used to wait on chanel operation
		select {
		// get the client that is getting registered
		case client := <-wsServer.Register:
			wsServer.registerClient(client)

		case client := <-wsServer.UnRegister:
			wsServer.unRegisterClient(client)
		}
	}
}

func (wsServer *WsServer) registerClient(client *models.Client) {
	wsServer.Clients[client] = true
}

func (wsServer *WsServer) unRegisterClient(client *models.Client) {
	// you find this client in map and delete it
	delete(wsServer.Clients, client)
}
