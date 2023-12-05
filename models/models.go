package models

import (
	"github.com/gorilla/websocket"
	"web-socket/wsServer"
)

// Client each client has a connection associated with it
type Client struct {
	Conn     *websocket.Conn
	WsServer *wsServer.WsServer
}
