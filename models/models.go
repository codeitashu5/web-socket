package models

import "github.com/gorilla/websocket"

// Client each client has a connection associated with it
type Client struct {
	Conn *websocket.Conn
}
