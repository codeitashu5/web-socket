package wsServer

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"web-socket/models"
)

// Client each client has a connection associated with it
type Client struct {
	Client *models.Client
}

// NewClient each client is created with a connection
func NewClient(wsServer *WsServer, conn *websocket.Conn) *Client {
	return &Client{
		Client: &models.Client{
			Conn:     conn,
			WsServer: wsServer,
		},
	}
}

func (c *Client) Reader() {
	for {
		// read message from the client
		messageType, p, err := c.Client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		// print the message that is read
		fmt.Println(string(p))

		// send conformation message back to the client
		err = c.Client.Conn.WriteMessage(messageType, []byte("Message is been read"))
		if err != nil {
			log.Println(err)
		}
	}
}
