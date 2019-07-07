package main

import (
	"flag"
	"github.com/hawklithm/anychatcmd/ui"
	"github.com/hawklithm/anychatcmd/wechat"
	"github.com/hawklithm/termui"
	"log"
	"os"
	"testing"
)

func Test_chat_box(t *testing.T) {
	flag.Parse()

	fileName := "log.txt"
	var logFile *os.File
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer logFile.Close()

	err = termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	if err != nil {
		print("打开文件失败!\n")
	}

	logger := log.New(logFile, "[*]", log.LstdFlags)

	msgIn := make(chan wechat.Message)
	msgOut := make(chan wechat.MessageRecord)
	groupChan := make(chan ui.SelectEvent)

	chatBox := ui.NewChatBox("", "", 0, 0, 100, 40, logger, msgIn, msgOut,
		groupChan)

	chatBox.Pick()

	groupChan <- ui.Group{UserList: []*ui.UserInfo{{UserId: "1234",
		Nick: "testNick", DisplayName: "TestDisplayName"}, {UserId: "12345",
		Nick: "testNick2", DisplayName: "TestDisplayName2"}},
		Name:    "test1",
		GroupId: "12345678"}

	msgIn <- wechat.Message{FromUserName: "12345678", Content: "123125125",
		MsgType: 1, MsgId: "123124124", ToUserName: "87654321"}

	msgIn <- wechat.Message{FromUserName: "12345678",
		Content: "absdfasdflkasgjdklajl",
		MsgType: 1, MsgId: "123124125", ToUserName: "87654321"}

	uiEvents := termui.PollEvents()

	go func() {
		for {
			<-msgOut
		}
	}()

	for {
		e := <-uiEvents
		if chatBox.Action(e) {
			continue
		}
		switch e.ID {
		case "<C-c>":
			return
		}
	}

}
