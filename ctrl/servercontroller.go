package ctrl

import (
	"chatroom/model"
	"chatroom/service"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{}

// Broadcast 消息管道
var Broadcast = make(chan model.ChatRoomRequest)

const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client 连接对象
type Client struct {
	Conn *websocket.Conn
}

// HandleConnections 处理连接
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("升级WebSocket失败")
	}
	username := r.Header.Get("username")
	model.Clients[ws] = username
	client := &Client{Conn: ws}
	go client.ReadMessage()
	client.WriteMessage()
}

// ReadMessage 读取消息
func (c *Client) ReadMessage() {
	defer func() {
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// 读消息
		_, messgedata, err := c.Conn.ReadMessage()
		if err != nil {
			delete(model.Clients, c.Conn)
			logrus.Error("读取消息失败", err)
			break
		}
		message := model.ChatRoomRequest{}
		err = proto.Unmarshal(messgedata, &message)
		if err != nil {
			logrus.Error("反序列化失败", err)
			break
		}
		message.UserName = model.Clients[c.Conn]
		// 将监听到的消息放入管道
		Broadcast <- message
	}
}

// WriteMessage 发送消息
func (c *Client) WriteMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message := <-Broadcast:
			if message.Type == "talk" {
				err := service.HandlerTalk(message)
				if err != nil {
					return
				}
			}
			if message.Type == "exit" {
				err := service.HandlerExit(c.Conn)
				if err != nil {
					return
				}
			}
			if message.Type == "userlist" {
				err := service.HandlerUserList(c.Conn)
				if err != nil {
					return
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
