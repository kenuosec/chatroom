package ctrl

import (
	"chatroom/model"
	"log"
	"net/url"

	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// 保存消息记录的切片
var messageRecords = make([]string, 0)

// UserListText 用户列表框
var UserListText *widget.Label

// RecordsText 消息记录框
var RecordsText *widget.Entry

// GlobalClient 全局连接对象
var GlobalClient *websocket.Conn

// CreatConnection 创建连接
func CreatConnection(username string, addr string) *websocket.Conn {
	if username == "" || addr == "" {
		return nil
	}
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	header := map[string][]string{"username": {username}}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Print("建立连接失败", err)
		return nil
	}
	GlobalClient = c
	return GlobalClient
}

// ReadMessageClient 监听消息
func ReadMessageClient() {
	for {
		// 读消息
		_, messgedata, err := GlobalClient.ReadMessage()
		if err != nil {
			logrus.Error("读取消息失败", err)
			break
		}
		message := model.ChatRoomRequest{}
		err = proto.Unmarshal(messgedata, &message)
		if err != nil {
			logrus.Error("反序列化失败", err)
			break
		}
		var list string
		if message.UserList != nil {
			for _, value := range message.UserList {
				list += value + "\n"
			}
			UserListText.Text = list
			UserListText.Refresh()
			continue
		}
		var records string
		if message.Content != "" && message.UserName != "" {
			messageRecords = append(messageRecords, message.UserName+":"+message.Content+"\n")
			for i := 0; i < len(messageRecords); i++ {
				records += messageRecords[i]
			}
			RecordsText.Text = records
			RecordsText.Refresh()
		}
	}
}

// SendTalk 发送talk消息
func SendTalk(content string, username string) {
	if GlobalClient == nil {
		return
	}
	message := model.ChatRoomRequest{Type: "talk", Content: content, UserName: username}
	data, err := proto.Marshal(&message)
	if err != nil {
		logrus.Print("序列化失败", err)
	}
	err = GlobalClient.WriteMessage(1, data)
	if err != nil {
		logrus.Print("发送消息失败", err)
	}
}

// SendExit 发送exit类型消息
func SendExit() {
	if GlobalClient == nil {
		return
	}
	message := model.ChatRoomRequest{Type: "exit"}
	defer GlobalClient.Close()
	data, err := proto.Marshal(&message)
	if err != nil {
		logrus.Print("序列化失败", err)
	}
	err = GlobalClient.WriteMessage(1, data)
	if err != nil {
		log.Print("发送消息失败", err)
	}
}

// SendUserList 发送userlist类型消息
func SendUserList() {
	if GlobalClient == nil {
		return
	}
	message := model.ChatRoomRequest{Type: "userlist"}
	data, err := proto.Marshal(&message)
	if err != nil {
		logrus.Print("序列化失败", err)
	}
	err = GlobalClient.WriteMessage(1, data)
	if err != nil {
		logrus.Print("发送消息失败", err)
	}
}
