package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"web-socket/wsServer"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndPoint(ws *wsServer.WsServer, w http.ResponseWriter, r *http.Request) {
	// it is used to validate the domain but as of now we will allow any host that wants to connect
	upGrader.CheckOrigin = func(r *http.Request) bool { return true }

	// we will give out http w,r to the upGrader, and it will provide us with a web-socket connection
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// we have created a connection for the client using which the client can communicate
	log.Println("Client Connected")

	// now the client is being associated with a server and a connection
	client := wsServer.NewClient(ws, conn)

	// register client on server
	ws.Register <- client.Client

	err = client.Client.Conn.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	// so using this we are reading messages form this client --> the messages that are being send on this socket form the front end
	client.Reader()
}

func setUpRoutes() {
	wsLocalServer := wsServer.NewWebSocketServer()
	// this function will run parallel without stoping the further execution
	go wsLocalServer.Run()

	// associate the server with the request
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		wsEndPoint(wsLocalServer, writer, request)
	})
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
