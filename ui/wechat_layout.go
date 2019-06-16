package ui

import (
	"fmt"
	"github.com/hawklithm/anychatcmd/wechat"
	ui "github.com/hawklithm/termui"
	"github.com/hawklithm/termui/widgets"
	"image"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

const (
	SelectedMark   = "(bg:red)"
	UnSelectedMark = "(bg:blue)"
	PageSize       = 45
)

type Layout struct {
	chatBox         *ChatBox //聊天窗口
	userIDList      []string
	curUserIndex    int
	masterName      string // 主人的名字
	masterID        string //主人的id
	currentMsgCount int
	maxMsgCount     int
	Notify          bool
	userIn          chan []string            // 用户的刷新
	msgIn           chan wechat.Message      // 消息刷新
	imageIn         chan wechat.MessageImage //  图片消息
	closeChan       chan int
	autoReply       chan int
	showUserList    []string
	userCount       int //用户总数，这里有重复,后面会修改
	pageCount       int // page总数。
	userCur         int // 当前page中所选中的用户
	curPage         int // 当前所在页
	pageSize        int // page的size默认是50
	curUserId       string
	userMap         map[string]string
	logger          *log.Logger
	userChatLog     map[string][]*wechat.MessageRecord
	groupMemberMap  map[string]map[string]string
	imageMap        map[string]image.Image
	msgIdList       map[string][]string
	selectedMsgId   string
}

type WidgetPicker interface {
	Pick()
	Unpick()
	Action(e ui.Event) bool
}

func NewLayout(
	recentUserList []UserInfo, recentGroupList []Group, userList []UserInfo,
	groupList []Group, userChangeEvent chan UserChangeEvent,
	selectEvent chan SelectEvent,
	myName, myID string,
	msgIn chan wechat.Message, msgOut chan wechat.MessageRecord,
	logger *log.Logger) {

	//	chinese := false
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	//用户列表框
	userChatLog := make(map[string][]*wechat.MessageRecord)
	groupMemberMap := make(map[string]map[string]string)
	imageMap := make(map[string]image.Image)

	width, height := ui.TerminalDimensions()

	var pickerList []WidgetPicker

	userListWidget := NewUserList(recentUserList, recentGroupList, userList,
		groupList, selectEvent, width*2/10, height, 0, 0, logger)
	userListWidget.Pick()

	pickerList = append(pickerList, userListWidget)

	chatBox := NewChatBox(width*2/10, 0, width*8/10, height, logger, msgIn,
		msgOut, selectEvent)

	pickerList = append(pickerList, chatBox)

	_ = &Layout{
		userCur:         0,
		curPage:         0,
		chatBox:         chatBox,
		currentMsgCount: 0,
		maxMsgCount:     18,
		pageSize:        PageSize,
		curUserIndex:    0,
		masterID:        myID,
		masterName:      myName,
		logger:          logger,
		userChatLog:     userChatLog,
		groupMemberMap:  groupMemberMap,
		imageMap:        imageMap,
		msgIdList:       make(map[string][]string),
		Notify:          true,
	}

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		catched := false
		for _, picker := range pickerList {
			if picker.Action(e) {
				catched = true
				break
			}
		}
		if catched {
			continue
		}
		switch e.ID {
		case "<C-c>", "<C-d>":
			return
		}

	}

}

//func setRows(p *widgets.ImageList, records []*wechat.MessageRecord) {
//	var rows []*widgets.ImageListItem
//	for _, i := range records {
//		item := widgets.NewImageListItem()
//		if i.ContentImg != nil {
//			item.Img = i.ContentImg
//		} else if i.Url != "" {
//			item.Url = i.Url
//			item.Text = i.From + "->" + i.To + ": " + i.Content
//		} else {
//			item.Text = i.From + "->" + i.To + ": " + i.Content
//		}
//		rows = append(rows, item)
//	}
//	p.Rows = rows
//	p.SelectedRow = len(p.Rows) - 1
//	if p.SelectedRow < 0 {
//		p.SelectedRow = 0
//	}
//}

var commands = map[string]string{
	"darwin": "open",
	"linux":  "xdg-open",
}

var notifyCommands = map[string]string{
	"darwin": "osascript",
	"linux":  "notify-send",
}

var notifyCommandsArgs = map[string][]string{
	"darwin": {"-e", `display notification "%s" with title "%s"`},
	"linux":  {`%s`},
}

type formatFunc func(format []string, args ...string) []string

func darwinFormatFunc() formatFunc {
	return func(format []string, args ...string) []string {
		newString := format[1]
		newString = fmt.Sprintf(newString, args[0], args[1])
		return []string{format[0], newString}
	}
}

func linuxFormatFunc() formatFunc {
	return func(format []string, args ...string) []string {
		newString := format[0]
		newString = fmt.Sprintf(newString, args[0])
		return []string{newString}
	}
}

var notifyFormatFunc = map[string]formatFunc{
	"darwin": darwinFormatFunc(),
	"linux":  linuxFormatFunc(),
}

func ShowNotify(message string) error {
	run, ok := notifyCommands[runtime.GOOS]
	args := notifyCommandsArgs[runtime.GOOS]
	notifyFunc := notifyFormatFunc[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to send notify on %s"+
			" platform",
			runtime.GOOS)
	}
	if message == "" {
		message = "收到了一个消息"
	}
	newArgs := notifyFunc(args, message, "wechat")
	cmd := exec.Command(run, newArgs...)
	return cmd.Start()
}

// Open calls the OS default program for uri
func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	if !strings.HasPrefix(uri, "http") {
		uri = "http://" + uri
	}

	cmd := exec.Command(run, uri)
	return cmd.Start()
}

func (l *Layout) getUserIdFromContent(content string,
	userMap map[string]string) string {
	if userMap == nil {
		return content
	}
	s := strings.Split(content, ":")
	if len(s) > 0 && userMap[s[0]] != "" {
		s[0] = userMap[s[0]]
	}
	l.logger.Println("groupMap=", userMap, "s=", s)
	builder := strings.Builder{}
	for i, sub := range s {
		builder.WriteString(sub)
		if i != len(s)-1 {
			builder.WriteString(":")
		}
	}
	return builder.String()
}

func (l *Layout) getUserIdAndConvertImgContent(content string,
	userMap map[string]string) string {
	if userMap == nil {
		return content
	}
	s := strings.Split(content, ":")
	if len(s) > 0 && userMap[s[0]] != "" {
		s[0] = userMap[s[0]]
	}
	l.logger.Println("groupMap=", userMap, "s=", s)
	return s[0] + ":" + AddUnSelectedBg("图片")
}

//func (l *Layout) apendChatLogIn(msg wechat.Message) *wechat.MessageRecord {
//	if l.userChatLog[msg.FromUserName] == nil {
//		l.userChatLog[msg.FromUserName] = []*wechat.MessageRecord{}
//	}
//
//	newMsg := wechat.NewMessageRecordIn(msg)
//
//	if l.groupMemberMap[newMsg.From] != nil {
//		if newMsg.Type == 3 {
//			newMsg.Content = l.getUserIdAndConvertImgContent(newMsg.Content,
//				l.groupMemberMap[newMsg.From])
//		} else {
//			newMsg.Content = l.getUserIdFromContent(newMsg.Content,
//				l.groupMemberMap[newMsg.From])
//		}
//	} else {
//		if newMsg.Type == 3 {
//			newMsg.Content = "图片"
//		}
//	}
//
//	if l.userMap[newMsg.To] != "" {
//		newMsg.To = l.userMap[newMsg.To]
//	}
//
//	if l.userMap[newMsg.From] != "" {
//		newMsg.From = l.userMap[newMsg.From]
//	}
//
//	l.userChatLog[msg.FromUserName] = append(l.userChatLog[msg.
//		FromUserName], newMsg)
//
//	return newMsg
//
//}

//func (l *Layout) apendImageChatLogIn(msg wechat.MessageImage) *wechat.MessageRecord {
//	if l.userChatLog[msg.FromUserName] == nil {
//		l.userChatLog[msg.FromUserName] = []*wechat.MessageRecord{}
//	}
//
//	newMsg := wechat.NewImageMessageRecordIn(msg)
//
//	if l.groupMemberMap[newMsg.From] != nil {
//		if newMsg.Type == 3 {
//			newMsg.Content = l.getUserIdAndConvertImgContent(newMsg.Content,
//				l.groupMemberMap[newMsg.From])
//		} else {
//			newMsg.Content = l.getUserIdFromContent(newMsg.Content,
//				l.groupMemberMap[newMsg.From])
//		}
//	} else {
//		if newMsg.Type == 3 {
//			newMsg.Content = "图片"
//		}
//	}
//
//	if l.userMap[newMsg.To] != "" {
//		newMsg.To = l.userMap[newMsg.To]
//	}
//
//	if l.userMap[newMsg.From] != "" {
//		newMsg.From = l.userMap[newMsg.From]
//	}
//
//	l.userChatLog[msg.FromUserName] = append(l.userChatLog[msg.
//		FromUserName], newMsg)
//
//	return newMsg
//
//}

func AddSelectedBg(msg string) string {
	return AddBgColor(DelBgColor(msg), SelectedMark)
}

func AddUnSelectedBg(msg string) string {
	return AddBgColor(DelBgColor(msg), UnSelectedMark)
}
func AddBgColor(msg string, color string) string {
	if strings.HasPrefix(msg, "[") {
		return msg
	}
	return "[" + msg + "]" + color
}
func DelBgColor(msg string) string {

	if !strings.HasPrefix(msg, "[") {
		return msg
	}
	return msg[1 : len(msg)-9]
}

func appendToPar(p *widgets.Paragraph, k string) {
	if strings.Count(p.Text, "\n") >= 20 {
		subText := strings.Split(p.Text, "\n")
		p.Text = strings.Join(subText[len(subText)-20:], "\n")
	}
	p.Text += k
	ui.Render(p)
}

func appendToList(p *widgets.ImageList, k *wechat.MessageRecord) {
	item := widgets.NewImageListItem()
	item.Text = k.String()
	if k.Url != "" {
		item.Url = k.Url
	}
	p.Rows = append(p.Rows, item)
	ui.Render(p)
}

func appendImageToList(p *widgets.ImageList, k image.Image) {
	item := widgets.NewImageListItem()
	item.WrapText = true
	item.Img = k
	p.Rows = append(p.Rows, item)
	ui.Render(p)
}

func resetPar(p *widgets.Paragraph) {
	p.Text = ""
	ui.Render(p)
}

func setPar(p *widgets.Paragraph) {
	ui.Render(p)
}

//func convertChatLogToText(records []*wechat.MessageRecord) string {
//	var b strings.Builder
//	var start = 0
//	if len(records) > 20 {
//		start = len(records) - 20
//	}
//	for _, i := range records[start:] {
//		_, _ = fmt.Fprint(&b, i.From+"->"+i.To+": "+i.Content+"\n")
//	}
//	return b.String()
//}
