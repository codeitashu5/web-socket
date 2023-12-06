package wsServer

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
	"web-socket/models"
)

const (
	/*
		we are sending our message to the client we wait for the message to be written to the client
		--> we wait for 10 sec if the client is reading the message or not
		--> if the client was not able to read in this time period we will not wait and move ahead
		--> in this case our write operation to client has failed
	*/
	writeWait = 10 * time.Second

	/*
	  we ping the client,so we expect it to receive our ping and respond with a pong
	  --> the time after which next pong message is expected
	*/
	pongWait = 60 * time.Second

	/*
	  interval after which we are sending ping request
	  --> we made the pingPeriod smaller than pong wait
	  --> making sure that pong wait is not violated
	*/
	pingPeriod = (pongWait * 9) / 10

	// the message size should not exceed this limit in byte ig
	maxMessageSize = 10000
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

// a function that will read messages continuously from the client
func (c *Client) readPump() {

	defer func() {
		err := c.Client.Conn.Close()
		if err != nil {
			panic("un-able to close client")
		}
	}()

	c.Client.Conn.SetReadLimit(maxMessageSize)
	c.Client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Client.Conn.SetPongHandler(func(appData string) error {
		// extend the read deadline so that the connection remains alive
		c.Client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// create an endless loop for message other than ping messages
	for {
		_, jsonMessage, err := c.Client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Closed un-expectedly %v", err)
			}
			break
		}

		// check the given message from the client
		fmt.Println(string(jsonMessage))
	}
}

// crating read pump for client

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
