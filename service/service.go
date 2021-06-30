package service

import (
	"chatroom/model"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// HandlerTalk 处理talk类型的消息
func HandlerTalk(message model.ChatRoomRequest) error {
	logrus.WithFields(logrus.Fields{
		"username": message.UserName,
		"content":  message.Content,
	}).Info("服务端广播消息")
	// Send it out to every client that is currently connected
	for client, _ := range model.Clients {
		data, err := proto.Marshal(&message)
		if err != nil {
			logrus.Error("序列化失败", err)
			return err
		}
		err = client.WriteMessage(1, data)
		if err != nil {
			logrus.Error("发送消息失败", err)
			err := client.Close()
			if err != nil {
				logrus.Error("关闭连接失败", err)
				return err
			}
			delete(model.Clients, client)
		}
	}
	return nil
}

// HandlerExit 处理exit类型的消息
func HandlerExit(conn *websocket.Conn) error {
	delete(model.Clients, conn)
	err := conn.Close()
	if err != nil {
		logrus.Fatalln("关闭连接失败", err)
	}
	return err
}

// HandlerUserList 处理userlist类型的消息
func HandlerUserList(conn *websocket.Conn) error {
	list := make(map[string]string)
	for _, username := range model.Clients {
		list[username] = username
	}
	userList := model.ChatRoomRequest{UserList: list}
	data, err2 := proto.Marshal(&userList)
	if err2 != nil {
		logrus.Error("消息序列化失败", err2)
		return err2
	}
	err3 := conn.WriteMessage(1, data)
	if err3 != nil {
		logrus.Error("发送消息失败", err3)
		return err3
	}
	return nil
}
