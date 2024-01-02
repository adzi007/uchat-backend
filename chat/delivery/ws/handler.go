package delivery

// import (
// 	"fmt"
// 	"log"
// 	"sync"

// 	"github.com/gofiber/contrib/websocket"
// )

// type Room struct {
// 	ID       string
// 	Clients  map[*websocket.Conn]bool
// 	Messages chan string
// 	mu       sync.Mutex
// }

// func NewRoom(id string) *Room {
// 	return &Room{
// 		ID:       id,
// 		Clients:  make(map[*websocket.Conn]bool),
// 		Messages: make(chan string),
// 	}
// }

// func (r *Room) Join(client *websocket.Conn) {
// 	r.mu.Lock()
// 	r.Clients[client] = true
// 	r.mu.Unlock()

// 	go func() {
// 		for {
// 			message := <-r.Messages
// 			if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
// 				log.Printf("Error sending message to client: %v", err)
// 				return
// 			}
// 		}
// 	}()
// }

// func (r *Room) Leave(client *websocket.Conn) {
// 	r.mu.Lock()
// 	delete(r.Clients, client)
// 	r.mu.Unlock()
// }

// func (r *Room) Broadcast(message string) {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	for client := range r.Clients {

// 		client.WriteJSON(message)

// 	}
// }

// func HandleWebSocket(c *websocket.Conn) {
// 	vars := c.Params("roomID")
// 	room := GetRoomx(vars)
// 	if room == nil {
// 		log.Printf("Room not found")
// 		c.Close()
// 		return
// 	}

// 	room.Join(c)
// 	defer room.Leave(c)

// 	for {
// 		message := ""
// 		if err := c.ReadJSON(&message); err != nil {
// 			log.Printf("Error reading message from client: %v", err)
// 			break
// 		}
// 		room.Broadcast(fmt.Sprintf("[%s] %s", vars, message))
// 	}
// }

// func GetRoomx(roomID string) *Room {
// 	// Implement your logic to retrieve or create a room based on roomID
// 	// For simplicity, let's create a new room for each roomID.
// 	return NewRoom(roomID)
// }
