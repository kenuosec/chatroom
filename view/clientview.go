package view

import (
	"chatroom/ctrl"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

// Start 启动应用
func Start() {
	myApp := app.New()
	myWindow := myApp.NewWindow("ChatRoom")
	UserNameEntry := widget.NewEntry()
	ServerAddressEntry := widget.NewEntry()
	StateLabel := widget.NewLabel("false")
	UserList := widget.NewLabel("")
	ctrl.UserListText = UserList
	Records := widget.NewMultiLineEntry()
	ctrl.RecordsText = Records
	Message := widget.NewMultiLineEntry()
	connectButton := widget.NewButton("connect", func() {
		client := ctrl.CreatConnection(UserNameEntry.Text, ServerAddressEntry.Text)
		if client != nil {
			StateLabel.Text = "OK"
			// 建立连接后监听消息
			go ctrl.ReadMessageClient()
		}
	})
	sendButton := widget.NewButton("Send", func() {
		// 发送消息获取最新列表以及广播消息
		ctrl.SendUserList()
		ctrl.SendTalk(Message.Text, UserNameEntry.Text)
	})
	exitButton := widget.NewButton("exit", func() {
		ctrl.SendExit()
		StateLabel.Text = "Fasle"
	})
	form := widget.NewForm(
		&widget.FormItem{Text: "UserName", Widget: UserNameEntry},
		&widget.FormItem{Text: "ServerAddress", Widget: ServerAddressEntry},
		&widget.FormItem{Text: "", Widget: connectButton},
		&widget.FormItem{Text: "", Widget: exitButton},
		&widget.FormItem{Text: "status", Widget: StateLabel},
		&widget.FormItem{Text: "UserList", Widget: UserList},
		&widget.FormItem{Text: "Records", Widget: Records},
		&widget.FormItem{Text: "Message", Widget: Message},
		&widget.FormItem{Text: "", Widget: sendButton},
	)
	myWindow.SetContent(form)
	myWindow.Resize(fyne.NewSize(300, 300))
	myWindow.ShowAndRun()
}
