package main

import (
	"flag"
	"github.com/hawklithm/anychatcmd/ui"
	"github.com/hawklithm/termui"
	"log"
	"os"
	"testing"
	"time"
)

func Test_UserList(t *testing.T) {
	flag.Parse()

	userInfoArray := []ui.UserInfo{
		{UserId: "12312124", Nick: "TestNick1",
			LastChatTime: time.Now()},
		{UserId: "12312125", Nick: "TestNick2",
			DisplayName:  "TestDisplay2",
			LastChatTime: time.Now()},
	}

	ruserInfoArray := []ui.UserInfo{
		{UserId: "12312128", Nick: "TestNick3",
			LastChatTime: time.Now()},
		{UserId: "12312129", Nick: "TestNick4",
			DisplayName:  "TestDisplay4",
			LastChatTime: time.Now()},
	}

	groupArray := []*ui.Group{
		{GroupId: "12312126", Name: "TestGroup1",
			LastChatTime: time.Now()},
		{GroupId: "12312127", Name: "TestGroup2",
			LastChatTime: time.Now()},
	}

	rgroupArray := []*ui.Group{
		{GroupId: "12312130", Name: "TestGroup3",
			LastChatTime: time.Now()},
		{GroupId: "12312131", Name: "TestGroup4",
			LastChatTime: time.Now()},
	}

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

	userList := ui.NewUserList(ruserInfoArray, rgroupArray, userInfoArray,
		groupArray,
		nil, 40, 40, 0, 0, logger, nil)

	userList.Pick()

	uiEvents := termui.PollEvents()

	for {
		e := <-uiEvents
		if userList.Action(e) {
			continue
		}
		switch e.ID {
		case "<C-c>":
			return
		}
	}

}
