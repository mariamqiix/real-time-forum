package Server

import (
	"RealTimeForum/database"
	"RealTimeForum/structs"
	"strconv"

	"encoding/json"
	"log"
	"net/http"

	// "strconv"
	// "strings"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ID         int
	connection *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections []Connection

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	user := GetUser(r)
	if user == nil {
		log.Println("User not authenticated")
		return
	}

	connection := Connection{
		ID:         user.Id,
		connection: conn,
	}
	connections = append(connections, connection)
	defer func() {
		RemoveConnection(user.Id)
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			println(err)
			return
		}
		MessageRequest := BodyToMessage(p)
		if MessageRequest == nil {
			panic("Message is nil.")
		}

		message := structs.Message{
			SenderId:   MessageRequest.SenderId,
			ReceiverId: MessageRequest.RecipientId,
			Message:    MessageRequest.Message,
			Time:       MessageRequest.Time,
		}

		err = database.AddMessage(message)
		if err != nil {
			errorServer(w, r, http.StatusInternalServerError)
			continue
		}

		reciverConnections, ok := GetConnectionByID(message.ReceiverId)
		if ok {
			SendMessage(*reciverConnections, &message)
		} else {
			log.Println("No connection found for the user with id: ", message.ReceiverId)
		}

	}
}

func SendMessage(conn Connection, message *structs.Message) {
	b, err := json.Marshal(message)
	if err != nil {
		log.Println("Error wrapping the message to bytes. " + err.Error())
		conn.connection.Close()
		RemoveConnection(conn.ID)
	}
	err = conn.connection.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("Error writting the message into the Web Socket. ", err.Error())
		conn.connection.Close()
		RemoveConnection(conn.ID)
	}
}

func BodyToMessage(body []byte) *structs.MessageRequest {
	if len(body) == 0 {
		return nil
	}

	var message structs.MessageRequest
	err := json.Unmarshal(body, &message)
	if err != nil {
		println("Error: ", err.Error())
		return nil
	}

	return &message
}

// RemoveConnection removes a connection by its ID
func RemoveConnection(userID int) {

	for i, conn := range connections {
		if conn.ID == userID {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

// GetConnectionByID returns the connection where the ID matches the recipientID
func GetConnectionByID(recipientID int) (*Connection, bool) {

	for _, conn := range connections {
		if conn.ID == recipientID {
			return &conn, true
		}
	}
	return nil, false
}

func IsUserOnline(userID int) bool {
	_, ok := GetConnectionByID(userID)
	return ok
}

func checkUserOnlineHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	if IsUserOnline(userID) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User is online"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User is offline"))
	}
}
