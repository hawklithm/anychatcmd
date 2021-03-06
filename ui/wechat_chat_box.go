package ui

import (
	"github.com/hawklithm/anychatcmd/utils"
	"github.com/hawklithm/anychatcmd/wechat"
	"github.com/hawklithm/termui"
	"github.com/hawklithm/termui/widgets"
	"github.com/skratchdot/open-golang/open"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
	groupChan chan SelectEvent          //  聊天成员变更
	messageIn chan wechat.Message       //  聊天成员变更

	conversationBox *widgets.ImageList
	editBox         *widgets.Paragraph
	msgInBox        *widgets.Paragraph
	picked          bool
	userChatLog     map[string]*ChatLogRecord
	memberList      []*UserInfo
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
	l.conversationBox.BorderStyle = termui.NewStyle(termui.ColorRed)
	l.editBox.BorderStyle = termui.NewStyle(termui.ColorRed)
	termui.Render(l.conversationBox, l.editBox)
}

func (l *ChatBox) Unpick() {
	l.picked = false
	l.conversationBox.BorderStyle = termui.NewStyle(termui.ColorMagenta)
	l.editBox.BorderStyle = termui.NewStyle(termui.ColorMagenta)
	termui.Render(l.conversationBox, l.editBox)
}

func (l *ChatBox) showDetail() {
	record := l.userChatLog[l.Id]
	l.logger.Println("current record is= ", record)
	item := record.record[l.conversationBox.SelectedRow]
	l.logger.Println("item detail selected! item=", item)
	if item.ContentImg == nil && item.Url == "" {
		return
	}
	if item.Url != "" {
		if err := Open(item.Url); err != nil {
			panic(err)
		}
	} else if item.ContentImg != nil {
		root := "/tmp"
		key := time.Now().UTC().UnixNano()
		builder := strings.Builder{}
		builder.WriteString(root)
		builder.WriteRune(os.PathSeparator)
		builder.WriteString(strconv.FormatInt(key, 10))
		builder.WriteString(".png")
		out, err := os.Create(builder.String())
		if err != nil {
			l.logger.Fatalln("open file failed! path=", builder.String(), err)
		}
		if err := png.Encode(out, *item.ContentImg); err != nil {
			l.logger.Fatalln("encode image failed! path=", builder.String(), err)
		} else {
			_ = open.Start(builder.String())
		}
	}

}

func (l *ChatBox) appendToConversationBox(msg wechat.MessageRecord) {
	item := widgets.NewImageListItem()
	item.Url = msg.Url
	if msg.ContentImg != nil {
		item.SetImage(*msg.ContentImg)
	}
	item.Text = msg.Text
	if l.memberListMap[msg.Speaker] != nil {
		member := l.memberListMap[msg.Speaker]
		item.Title = utils.If(member.DisplayName != "", member.DisplayName,
			member.Nick).(string)
	} else if msg.Speaker == l.MyId {
		item.Title = "我"
	} else {
		item.Title = "unknow"
	}
	l.conversationBox.Rows = append(l.conversationBox.Rows, item)
	termui.Render(l.conversationBox)
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
	if !l.picked {
		return false
	}
	switch e.ID {
	case "<Enter>":
		if l.editBox.Text != "" {
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
	return t, s[idx+6:]
}

func (l *ChatBox) apendChatLogOut(msg wechat.MessageRecord) *wechat.MessageRecord {
	if l.userChatLog[msg.To] == nil {
		l.userChatLog[msg.To] = &ChatLogRecord{}
	}

	l.userChatLog[msg.To].record = append(l.userChatLog[msg.To].record, msg)
	if msg.To == l.Id {
		l.appendToConversationBox(msg)
	}

	return &msg
}

func (l *ChatBox) apendChatLogIn(msg wechat.MessageRecord) *wechat.MessageRecord {
	if l.userChatLog[msg.From] == nil {
		l.userChatLog[msg.From] = &ChatLogRecord{}
	}

	l.userChatLog[msg.From].record = append(l.userChatLog[msg.From].record, msg)
	if msg.From == l.Id {
		l.appendToConversationBox(msg)
	}

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
					ContentImg: msg.Img,
					MsgId:      msg.MsgId})
			} else {
				var speaker, content string
				if msg.FromUserName[:2] == "@@" {
					speaker, content = getSpeakerIdAndContent(msg.Content)
				} else {
					speaker, content = msg.FromUserName, msg.Content
				}
				newMsg = l.apendChatLogIn(wechat.MessageRecord{To: msg.
					ToUserName,
					Text: content, ContentImg: msg.Img,
					Type: msg.MsgType, From: msg.FromUserName, Speaker: speaker,
					MsgId: msg.MsgId})
			}

			l.logger.Println("message receive = ", newMsg)

		}

	}
	return
}

func (l *ChatBox) resetRows() {
	var rows []*widgets.ImageListItem
	record := l.userChatLog[l.Id]
	if record != nil && record.record != nil {
		for _, i := range record.record {
			item := widgets.NewImageListItem()
			item.Text = i.Text
			var from string
			if i.Speaker == l.MyId {
				from = "我"
			} else if l.memberListMap[i.Speaker] != nil {
				p := l.memberListMap[i.Speaker]
				from = utils.If(p.DisplayName != "", p.DisplayName, p.Nick).(string)
			}
			item.Title = from
			if i.ContentImg != nil {
				item.SetImage(*i.ContentImg)
			} else if i.Url != "" {
				item.Url = i.Url
			}
			rows = append(rows, item)
		}
	}
	l.conversationBox.Rows = rows
	l.conversationBox.SelectedRow = len(l.conversationBox.Rows) - 1
	if l.conversationBox.SelectedRow < 0 {
		l.conversationBox.SelectedRow = 0
	}
}

func NewChatBox(myId, myName string, baseX, baseY, width, height int,
	logger *log.Logger,
	msgIn chan wechat.Message, msgOut chan wechat.MessageRecord,
	groupChan chan SelectEvent) *ChatBox {

	c := &ChatBox{
		baseX:     baseX,
		baseY:     baseY,
		width:     width,
		height:    height,
		logger:    logger,
		messageIn: msgIn,
		msgOut:    msgOut,
		groupChan: groupChan,
		MyId:      myId,
		MyName:    myName,
	}
	go c.displayMsgIn()
	c.Reset()
	go func() {
		for {
			group := <-c.groupChan
			c.Id = group.GetId()
			c.name = group.GetName()
			c.memberList = group.GetUserList()
			c.memberListMap = make(map[string]*UserInfo)
			for _, user := range c.memberList {
				c.memberListMap[user.UserId] = user
			}
			c.resetRows()
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
