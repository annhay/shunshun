package utils

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 消息推送

// Upgrader 持有 WebSocket 升级的配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有来源的跨域请求，生产环境中应根据需要进行配置
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端管理器
type ClientManager struct {
	clients   map[*websocket.Conn]bool
	mutex     sync.Mutex
	broadcast chan []byte
}

var manager = ClientManager{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan []byte),
}

// 启动广播器，它会持续监听 broadcast 通道并将消息发送给所有客户端
func (manager *ClientManager) startBroadcaster() {
	for {
		// 从广播通道中接收消息
		message := <-manager.broadcast

		manager.mutex.Lock()
		// 将消息发送给所有连接的客户端
		for client := range manager.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(manager.clients, client)
			}
		}
		manager.mutex.Unlock()
	}
}

// 处理客户端连接
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	log.Println("New client connected")

	// 将新客户端添加到管理器
	manager.mutex.Lock()
	manager.clients[ws] = true
	manager.mutex.Unlock()

	// 循环读取客户端发送的消息
	for {
		// ReadMessage 会阻塞，直到收到一条消息或连接关闭
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			manager.mutex.Lock()
			delete(manager.clients, ws)
			manager.mutex.Unlock()
			break
		}

		log.Printf("Received message: %s", message)
		// 将收到的消息发送到广播通道
		manager.broadcast <- message
	}
}
