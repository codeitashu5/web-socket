package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// creating home page for web socket
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "working fine")
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client each client has a connection associated with it
type Client struct {
	conn *websocket.Conn
}

// each client is created with a connection
func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	// it is used to validate the domain but as of now we will allow any host that wants to connect
	upGrader.CheckOrigin = func(r *http.Request) bool { return true }

	// we will give out http w,r to the upGrader, and it will provide us with a web-socket connection
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// we have created a connection for the client using which the client can communicate
	log.Println("Client Connected")

	// associate your connection with the client
	client := newClient(conn)

	err = client.conn.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	// so using this we are reading messages form this client --> the messages that are being send on this socket form the front end
	reader(client.conn)
}

// this takes a web-socket connection and continuously reads from it until the connection is broken
func reader(conn *websocket.Conn) {
	for {
		// read message from the client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		// print the message that is read
		fmt.Println(string(p))

		// send conformation message back to the client
		err = conn.WriteMessage(messageType, []byte("Message is been read"))
		if err != nil {
			log.Println(err)
		}
	}
}

func setUpRoutes() {
	http.HandleFunc("/health", homePage)
	http.HandleFunc("/ws", wsEndPoint)
}

func main() {
	fmt.Println("hello world")
	setUpRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
 we get a connection --> let's associate that connection with a user
 when the user tries to connect let's add them to a map
 -- > userId -> {
              conn, --> will get updated as we move around
              registered, --> if the user is connected or not
           }

-->  we have fn-> con
             sr-> con

using this conn we communicate

*/

/*
 --/ what we have made today is called simple ping-pong using web-socket
 --* what we want our server to be such that it can connect multiple users

  each user will have some cards with them
  -- each user throws there card
  -- and user with the largest card wins
  -- so the game will have status (started/in playing area)
  -- if the game is in playing area new user can be added (upto 4)
  -- if the game is started no further players can be added
  -- the player with the largest card wins
  -- there should be an api to through card for each player
  -- in the end the player with most no of points wins
  -- todo: add communication for multiple users

*/
