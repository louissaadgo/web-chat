package routes

import (
	"fmt"
	"log"
	"sort"

	"github.com/gofiber/websocket/v2"
)

var wsChan = make(chan WsPayload)
var keepAliveChan = make(chan WebSocketConnection)
var clients = make(map[WebSocketConnection]string)

type WebSocketConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

func Ws(c *websocket.Conn) {
	var response = WsJsonResponse{Message: `<em><small>Connected to server</small></em>`}
	conn := WebSocketConnection{Conn: c}
	clients[conn] = ""

	err := c.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
	go ListenForWs(&conn)
	for {
		client := <-keepAliveChan
		if c == client.Conn {
			break
		}
	}
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		recover()
	}()
	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse

	for {
		e := <-wsChan
		switch e.Action {
		case "username":
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			BroadcastToAll(response)

		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			BroadcastToAll(response)

		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			BroadcastToAll(response)
		}
	}
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	sort.Strings(userList)
	return userList
}

func BroadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			client.Close()
			delete(clients, client)
			keepAliveChan <- client
		}
	}
}
