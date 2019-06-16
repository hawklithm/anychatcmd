package ui

import (
	"github.com/hawklithm/anychatcmd/wechat"
	"github.com/hawklithm/termui"
	"github.com/hawklithm/termui/widgets"
	"log"
	"strings"
)

type ChatBox struct {
	MyId   string
	MyName string

	logger    *log.Logger
	baseX     int
	baseY     int
	width     int
	height    int
	msgOut    chan wechat.MessageRecord //  消息输出
	groupChan chan Group                //  聊天成员变更
	messageIn chan wechat.Message       //  聊天成员变更

	conversationBox *widgets.ImageList
	editBox         *widgets.Paragraph
	msgInBox        *widgets.Paragraph
	picked          bool
	userChatLog     map[string]*ChatLogRecord
	memberList      []UserInfo
	Id              string
	name            string

	memberListMap map[string]*UserInfo
	Notify        bool
}

type ChatLogRecord struct {
	record []wechat.MessageRecord
	length int
	id     string
}

func (l *ChatBox) Pick() {
	l.picked = true
}

func (l *ChatBox) Unpick() {
	l.picked = false
}

func (l *ChatBox) showDetail() {
	//item := l.chatBox.Rows[l.chatBox.SelectedRow]
	//l.logger.Println("item detail selected! item=", item)
	//if item.Img == nil && item.Url == "" {
	//	return
	//}
	//if item.Url != "" {
	//	if err := Open(item.Url); err != nil {
	//		panic(err)
	//	}
	//} else if item.Img != nil {
	//	root := "/tmp"
	//	key := time.Now().UTC().UnixNano()
	//	builder := strings.Builder{}
	//	builder.WriteString(root)
	//	builder.WriteRune(os.PathSeparator)
	//	builder.WriteString(strconv.FormatInt(key, 10))
	//	builder.WriteString(".png")
	//	out, err := os.Create(builder.String())
	//	if err != nil {
	//		l.logger.Fatalln("open file failed! path=", builder.String(), err)
	//	}
	//	if err := png.Encode(out, item.Img); err != nil {
	//		l.logger.Fatalln("encode image failed! path=", builder.String(), err)
	//	} else {
	//		_ = open.Start(builder.String())
	//	}
	//}

}

func (l *ChatBox) appendToConversationBox(msg wechat.MessageRecord) {
	item := widgets.NewImageListItem()
	item.Url = msg.Url
	item.Img = msg.ContentImg
	item.Text = msg.Text
	l.conversationBox.Rows = append(l.conversationBox.Rows, item)
}

func (l *ChatBox) NextSelect() {
	l.conversationBox.ScrollDown()
	termui.Render(l.conversationBox)
}

func (l *ChatBox) PrevSelect() {
	l.conversationBox.ScrollUp()
	termui.Render(l.conversationBox)
}

func (l *ChatBox) Action(e termui.Event) bool {
	switch e.ID {
	case "<Enter>":
		if l.editBox.Text != "" {
			appendTextToList(l.conversationBox, l.MyName+"->"+l.name+
				":"+l.editBox.Text+"\n")
			l.SendText(l.editBox.Text)
		}
		resetPar(l.editBox)
		return true
	case "<C-w>":
		l.showDetail()
		return true
	case "<C-j>":
		l.NextSelect()
		return true
	case "<C-k>":
		l.PrevSelect()
		return true
	case "<C-a>":
		l.Notify = !l.Notify
		l.logger.Println("notify state", l.Notify)
		return true
	case "<Space>":
		appendToPar(l.editBox, " ")
		return true
	case "<Backspace>":
		if l.editBox.Text != "" {
			runslice := []rune(l.editBox.Text)
			if len(runslice) != 0 {
				l.editBox.Text = string(runslice[0 : len(runslice)-1])
				setPar(l.editBox)
			}
		}
		return true
	default:
		if e.Type == termui.KeyboardEvent {
			k := e.ID
			appendToPar(l.editBox, k)
		} else if e.Type == termui.ResizeEvent {
			l.logger.Println("resize event received, payload=", e.Payload,
				"id=", e.ID)
		} else if e.Type == termui.MouseEvent {
			l.logger.Println("mouse event received, payload=", e.Payload,
				"id=", e.ID)
		}
		return false
	}
}

func appendTextToList(p *widgets.ImageList, k string) {
	item := widgets.NewImageListItem()
	item.Text = k
	p.Rows = append(p.Rows, item)
	termui.Render(p)
}

func (l *ChatBox) SendText(text string) {
	msg := wechat.MessageRecord{}
	msg.Text = text
	msg.To = l.Id
	msg.Speaker = l.MyId
	msg.From = l.MyId
	l.apendChatLogOut(msg)
	l.msgOut <- msg
}

func getSpeakerIdAndContent(content string) (string, string) {
	s := strings.Trim(content, " ")
	idx := strings.Index(s, ":")
	if idx < 0 {
		return "", s
	}
	t := strings.Trim(s[:idx], " ")
	if len(t) <= 0 {
		return "", s
	}
	return t, s[idx+1:]
}

func (l *ChatBox) apendChatLogOut(msg wechat.MessageRecord) *wechat.MessageRecord {
	if l.userChatLog[msg.To] == nil {
		l.userChatLog[msg.To] = &ChatLogRecord{}
	}

	l.userChatLog[msg.To].record = append(l.userChatLog[msg.To].record, msg)
	l.appendToConversationBox(msg)

	return &msg
}

func (l *ChatBox) apendChatLogIn(msg wechat.MessageRecord) *wechat.MessageRecord {
	if l.userChatLog[msg.From] == nil {
		l.userChatLog[msg.From] = &ChatLogRecord{}
	}

	l.userChatLog[msg.From].record = append(l.userChatLog[msg.From].record, msg)
	l.appendToConversationBox(msg)

	return &msg

}

func (l *ChatBox) displayMsgIn() {
	var (
		msg wechat.Message
	)

	for {
		select {
		case msg = <-l.messageIn:

			var newMsg *wechat.MessageRecord
			msg.Content = TranslateEmoji(ConvertToEmoji(msg.Content))

			if l.MyId == msg.FromUserName {
				newMsg = l.apendChatLogOut(wechat.MessageRecord{To: msg.
					ToUserName, Text: msg.Content, Type: msg.MsgType,
					From: msg.FromUserName, Speaker: l.MyId,
					MsgId: msg.MsgId})
			} else {
				speaker, content := getSpeakerIdAndContent(msg.Content)
				newMsg = l.apendChatLogIn(wechat.MessageRecord{To: msg.
					ToUserName,
					Text: content,
					Type: msg.MsgType, From: msg.FromUserName, Speaker: speaker,
					MsgId: msg.MsgId})
			}

			l.logger.Println("message receive = ", newMsg.String())

		}

	}
	return
}

func NewChatBox(baseX, baseY, width, height int, logger *log.Logger,
	msgIn chan wechat.Message, msgOut chan wechat.MessageRecord, groupChan chan Group) *ChatBox {

	c := &ChatBox{
		baseX:     baseX,
		baseY:     baseY,
		width:     width,
		height:    height,
		logger:    logger,
		messageIn: msgIn,
		msgOut:    msgOut,
		groupChan: groupChan,
	}
	go c.displayMsgIn()
	c.Reset()
	go func() {
		for {
			group := <-c.groupChan
			c.Id = group.GroupId
			c.name = group.Name
			c.memberList = group.UserList
			c.memberListMap = make(map[string]*UserInfo)
			for _, user := range c.memberList {
				c.memberListMap[user.UserId] = &user
			}
			termui.Render(c.conversationBox)
		}
	}()
	return c
}
func (l *ChatBox) Reset() {

	conversationBox := widgets.NewImageList()
	conversationBox.SetRect(l.baseX, l.baseY, l.baseX+l.width*6/8,
		l.baseY+l.height*8/10)

	conversationBox.TextStyle = termui.NewStyle(termui.ColorRed)
	//chatBox.Title = "to:" + userNickList[0]
	conversationBox.BorderStyle = termui.NewStyle(termui.ColorMagenta)

	l.conversationBox = conversationBox

	editBox := widgets.NewParagraph()
	editBox.SetRect(l.baseX, l.baseY+l.height*8/10, l.baseX+l.width*6/8,
		l.height)

	editBox.TextStyle = termui.NewStyle(termui.ColorWhite)
	editBox.Title = "输入框"
	editBox.BorderStyle = termui.NewStyle(termui.ColorCyan)

	l.editBox = editBox

	msgInBox := widgets.NewParagraph()

	msgInBox.SetRect(l.baseX+l.width*6/8, l.baseY, l.baseX+l.width,
		l.baseY+l.height)

	msgInBox.TextStyle = termui.NewStyle(termui.ColorWhite)
	msgInBox.Title = "消息窗"
	msgInBox.BorderStyle = termui.NewStyle(termui.ColorCyan)

	l.msgInBox = msgInBox

	l.userChatLog = make(map[string]*ChatLogRecord)

	termui.Render(conversationBox, editBox, msgInBox)

}
