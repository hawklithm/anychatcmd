package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/daviddengcn/go-colortext"
	"github.com/hawklithm/anychatcmd/jike"
	"github.com/hawklithm/anychatcmd/ui"
	chat "github.com/hawklithm/anychatcmd/wechat"
)

const (
	maxChanSize = 50
)

var (
	logger *log.Logger
)

type Config struct {
	SaveToFile   bool     `json:"save_to_file"`
	AutoReply    bool     `json:"auto_reply"`
	AutoReplySrc bool     `json:"auto_reply_src"`
	ReplyMsg     []string `json:"reply_msg"`
}

func startWechat(wxLogger *log.Logger) {
	wechat := chat.NewWechat(wxLogger)

	if err := wechat.WaitForLogin(); err != nil {
		logger.Fatalf("等待失败：%s\n", err.Error())
		return
	}
	srcPath, err := os.Getwd()
	if err != nil {
		logger.Printf("获得路径失败:%#v\n", err)
	}
	configFile := path.Join(path.Clean(srcPath), "config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logger.Fatalln("请提供配置文件：config.json")
		return
	}

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Fatalln("读取文件失败：%#v", err)
		return
	}
	var config *Config
	err = json.Unmarshal(b, &config)

	logger.Printf("登陆...\n")

	wechat.AutoReplyMode = config.AutoReply
	wechat.ReplyMsgs = config.ReplyMsg
	wechat.AutoReplySrc = config.AutoReplySrc

	if err := wechat.Login(); err != nil {
		logger.Printf("登陆失败：%v\n", err)
		return
	}
	logger.Printf("配置文件:%+v\n", config)

	logger.Println("成功!")

	logger.Println("微信初始化成功...")

	logger.Println("开启状态栏通知...")
	if err := wechat.StatusNotify(); err != nil {
		return
	}
	if err := wechat.GetContacts(); err != nil {
		logger.Fatalf("拉取联系人失败:%v\n", err)
		return
	}

	if err := wechat.TestCheck(); err != nil {
		logger.Fatalf("检查状态失败:%v\n", err)
		return
	}

	var recentUserList []ui.UserInfo
	var recentGroupList []*ui.Group
	var userInfos []ui.UserInfo
	var groupInfos []*ui.Group

	for _, member := range wechat.InitContactList {
		if strings.HasPrefix(member.UserName, "@@") {
			recentGroupList = append(recentGroupList, &ui.Group{GroupId: member.
				UserName, Name: member.NickName,
				LastChatTime: time.Now()})
		} else {
			recentUserList = append(recentUserList, ui.UserInfo{UserId: member.
				UserName, Nick: member.NickName, DisplayName: member.RemarkName,
				LastChatTime: time.Now()})
		}
	}

	for _, member := range wechat.ContactList {
		userInfos = append(userInfos, ui.UserInfo{UserId: member.
			UserName, Nick: member.NickName, DisplayName: member.RemarkName,
			LastChatTime: time.Now()})
	}

	for _, member := range wechat.PublicUserList {
		userInfos = append(userInfos, ui.UserInfo{UserId: member.
			UserName, Nick: member.NickName, DisplayName: member.RemarkName,
			LastChatTime: time.Now()})
	}

	for _, member := range wechat.GroupMemberList {
		groupInfos = append(groupInfos, &ui.Group{GroupId: member.
			UserName, Name: member.NickName,
			LastChatTime: time.Now()})
	}
	//groupIdList := []string{}
	//for _, user := range userIDList {
	//	if strings.HasPrefix(user, "@@") {
	//		groupIdList = append(groupIdList, user)
	//	}
	//}

	////群成员列表
	//groupMemberList, err := wechat.GetContactsInBatch(groupIdList)
	//if err != nil {
	//	logger.Fatal("get batch contact error=", err)
	//	return
	//}

	ui.InitTalkInfo(wechat, logger, groupInfos)

	msgIn := make(chan chat.Message, maxChanSize)
	msgOut := make(chan chat.MessageRecord, maxChanSize)
	selectEvent := make(chan ui.SelectEvent, maxChanSize)
	autoChan := make(chan int, 1)

	go wechat.SyncDaemon(msgIn)

	go wechat.MsgDaemon(msgOut, autoChan)

	//logger.Println("recentUserList size=", len(recentUserList))
	//logger.Println("recentGroupList size=", len(recentGroupList))
	//logger.Println("userInfos size=", len(userInfos))
	//logger.Println("groupInfos size=", len(groupInfos))

	ui.NewLayout(recentUserList, recentGroupList, userInfos, groupInfos,
		nil, selectEvent,
		wechat.User.NickName,
		wechat.User.UserName, msgIn, msgOut,
		wxLogger, wechat)
}

func startJike(wxLogger *log.Logger) {
	wechat := jike.NewJike(wxLogger)

	if err := wechat.WaitForLogin(); err != nil {
		logger.Fatalf("等待失败：%s\n", err.Error())
		return
	}
	//srcPath, err := os.Getwd()
	//if err != nil {
	//	logger.Printf("获得路径失败:%#v\n", err)
	//}
	//configFile := path.Join(path.Clean(srcPath), "config.json")
	//if _, err := os.Stat(configFile); os.IsNotExist(err) {
	//	logger.Fatalln("请提供配置文件：config.json")
	//	return
	//}
	//
	//b, err := ioutil.ReadFile(configFile)
	//if err != nil {
	//	logger.Fatalln("读取文件失败：%#v", err)
	//	return
	//}
	//var config *Config
	//err = json.Unmarshal(b, &config)
	//
	//logger.Printf("登陆...\n")
	//
	//wechat.AutoReplyMode = config.AutoReply
	//wechat.ReplyMsgs = config.ReplyMsg
	//wechat.AutoReplySrc = config.AutoReplySrc

	//if err := wechat.Login(); err != nil {
	//	logger.Printf("登陆失败：%v\n", err)
	//	return
	//}
	//logger.Printf("配置文件:%+v\n", config)
	//
	//logger.Println("成功!")
	//
	//logger.Println("微信初始化成功...")

	//logger.Println("开启状态栏通知...")
	//if err := wechat.StatusNotify(); err != nil {
	//	return
	//}
	//if err := wechat.GetContacts(); err != nil {
	//	logger.Fatalf("拉取联系人失败:%v\n", err)
	//	return
	//}
	//
	//if err := wechat.TestCheck(); err != nil {
	//	logger.Fatalf("检查状态失败:%v\n", err)
	//	return
	//}
	//
	//var recentUserList []ui.UserInfo
	//var recentGroupList []*ui.Group
	//var userInfos []ui.UserInfo
	//var groupInfos []*ui.Group
	//
	//for _, member := range wechat.InitContactList {
	//	if strings.HasPrefix(member.UserName, "@@") {
	//		recentGroupList = append(recentGroupList, &ui.Group{GroupId: member.
	//			UserName, Name: member.NickName,
	//			LastChatTime: time.Now()})
	//	} else {
	//		recentUserList = append(recentUserList, ui.UserInfo{UserId: member.
	//			UserName, Nick: member.NickName, DisplayName: member.RemarkName,
	//			LastChatTime: time.Now()})
	//	}
	//}
	//
	//for _, member := range wechat.ContactList {
	//	userInfos = append(userInfos, ui.UserInfo{UserId: member.
	//		UserName, Nick: member.NickName, DisplayName: member.RemarkName,
	//		LastChatTime: time.Now()})
	//}
	//
	//for _, member := range wechat.PublicUserList {
	//	userInfos = append(userInfos, ui.UserInfo{UserId: member.
	//		UserName, Nick: member.NickName, DisplayName: member.RemarkName,
	//		LastChatTime: time.Now()})
	//}
	//
	//for _, member := range wechat.GroupMemberList {
	//	groupInfos = append(groupInfos, &ui.Group{GroupId: member.
	//		UserName, Name: member.NickName,
	//		LastChatTime: time.Now()})
	//}
	////groupIdList := []string{}
	////for _, user := range userIDList {
	////	if strings.HasPrefix(user, "@@") {
	////		groupIdList = append(groupIdList, user)
	////	}
	////}
	//
	//////群成员列表
	////groupMemberList, err := wechat.GetContactsInBatch(groupIdList)
	////if err != nil {
	////	logger.Fatal("get batch contact error=", err)
	////	return
	////}
	//
	//ui.InitTalkInfo(wechat, logger, groupInfos)
	//
	//msgIn := make(chan chat.Message, maxChanSize)
	//msgOut := make(chan chat.MessageRecord, maxChanSize)
	//selectEvent := make(chan ui.SelectEvent, maxChanSize)
	//autoChan := make(chan int, 1)
	//
	//go wechat.SyncDaemon(msgIn)
	//
	//go wechat.MsgDaemon(msgOut, autoChan)

}

func main() {

	mode := flag.String("mode", "wechat", "anychatcmd start mode")
	flag.Parse()

	ct.Foreground(ct.Green, true)
	flag.Parse()
	logger = log.New(os.Stdout, "[*🤔 *]->:", log.LstdFlags)

	logger.Println("启动...")
	fileName := "log.txt"
	var logFile *os.File
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	defer logFile.Close()
	if err != nil {
		logger.Printf("打开文件失败!\n")
	}

	wxlogger := log.New(logFile, "[*]", log.LstdFlags)

	if *mode == "wechat" {
		startWechat(wxlogger)
	} else {
		startJike(wxlogger)
	}

}
