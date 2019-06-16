package main

import (
	"flag"
	"github.com/hawklithm/anychatcmd/ui"
	"github.com/hawklithm/anychatcmd/wechat"
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func Test_UI(t *testing.T) {
	flag.Parse()
	maxChanSize := 10000

	//log.SetLevel(log.DebugLevel)
	msgIn := make(chan wechat.Message, maxChanSize)
	textOut := make(chan wechat.MessageRecord, maxChanSize)
	imageIn := make(chan wechat.MessageImage, maxChanSize)

	closeChan := make(chan int, 1)
	autoChan := make(chan int, 1)

	fileName := "log.txt"
	var logFile *os.File
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	defer logFile.Close()
	if err != nil {
		print("打开文件失败!\n")
	}

	wxLogger := log.New(logFile, "[*]", log.LstdFlags)

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

	groupArray := []ui.Group{
		{GroupId: "12312126", Name: "TestGroup1",
			LastChatTime: time.Now()},
		{GroupId: "12312127", Name: "TestGroup2",
			LastChatTime: time.Now()},
	}

	rgroupArray := []ui.Group{
		{GroupId: "12312130", Name: "TestGroup3",
			LastChatTime: time.Now()},
		{GroupId: "12312131", Name: "TestGroup4",
			LastChatTime: time.Now()},
	}

	userChangeEvent := make(chan ui.UserChangeEvent, maxChanSize)
	selectEvent := make(chan ui.SelectEvent, maxChanSize)

	go func() {
		for {
			t := <-selectEvent
			wxLogger.Println(t)
		}
	}()

	ui.NewLayout(ruserInfoArray, rgroupArray, userInfoArray, groupArray,
		userChangeEvent, selectEvent,
		"test1", "test2",
		msgIn, textOut, autoChan, wxLogger)
}

func Test_NOTIFY(t *testing.T) {

	for i := 0; i < 10; i++ {
		if e := exec.Command("osascript", "-e",
			`display notification "test" with title "test_title"`).Run(); e != nil {
			println("error happen", e.Error())
		}
	}

	for i := 0; i < 10; i++ {
		time.Sleep(100)
		if e := ui.ShowNotify("test msg" + strconv.Itoa(i)); e != nil {
			print(e.Error())
		}
	}
}
