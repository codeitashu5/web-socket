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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
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

	// close the connection when the function ends
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
		c.Client.WsServer.BroadCast <- jsonMessage
	}
}

// it will continuously write messages to the client
func (c *Client) writePump() {
	// we create a ticker which will ping the client after the pingTime period is exhausted
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		// after retuning the ticker should stop and the connection is closed
		ticker.Stop()
		c.Client.Conn.Close()
	}()

	// read the message continuously form the
	for {
		select {
		case message, ok := <-c.Client.Send:
			c.Client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				c.Client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// creating a new writer for
			w, err := c.Client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// gets the total no of pending messages
			n := len(c.Client.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				// keep writing to the client until the send channel is empty
				w.Write(<-c.Client.Send)
			}

		// when the time period has elapsed the ticker receives a message
		case <-ticker.C:
			/*
			 so we are saying that we should be able to write our message to
			 the client within wait time provided that's why we are setting this deadline
			 every time we write a message
			*/
			c.Client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
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
