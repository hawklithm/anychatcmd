package ui

import (
	"github.com/hawklithm/anychatcmd/utils"
	ui "github.com/hawklithm/termui"
	"github.com/hawklithm/termui/widgets"
	"log"
	"time"
)

type UserList struct {
	currentId       string
	curUserIndex    int
	logger          *log.Logger
	baseX           int
	baseY           int
	width           int
	height          int
	recentUserList  []UserInfo
	recentGroupList []Group
	userList        []UserInfo
	groupList       []Group
	userChangeEvent chan UserChangeEvent

	tabPane          *widgets.TabPane
	userNickListBox  *widgets.List
	groupNickListBox *widgets.List
	recentListBox    *widgets.List
	tabWidgets       []*widgets.List
	picked           bool
	currentTab       *widgets.List
}

func (this *UserList) Pick() {
	this.picked = true
}

func (this *UserList) Unpick() {
	this.picked = false
}

type UserChangeEvent struct {
}

type UserInfo struct {
	UserId       string
	Nick         string
	DisplayName  string
	LastChatTime time.Time
}

//用户群组
type Group struct {
	UserList     []UserInfo
	Name         string
	GroupId      string
	LastChatTime time.Time
}

type SortItem struct {
	Nick         string
	Id           string
	LastChatTime time.Time
}

type SortItems []SortItem

func (p SortItems) Len() int { return len(p) }

// 根据元素的年龄降序排序 （此处按照自己的业务逻辑写）
func (p SortItems) Less(i, j int) bool {
	return p[i].LastChatTime.After(p[j].LastChatTime)
}

func merge(userInfo []UserInfo, groupInfo []Group) *SortItems {
	z := make(SortItems, len(userInfo)+len(groupInfo))
	for i, user := range userInfo {
		z[i] = SortItem{
			Nick:         utils.If(user.DisplayName != "", user.DisplayName, user.Nick).(string),
			Id:           user.UserId,
			LastChatTime: user.LastChatTime,
		}
	}
	return &z
}

func (this *UserList) renderTab() {
	this.currentTab = this.tabWidgets[this.tabPane.ActiveTabIndex]
	ui.Render(this.tabPane, this.currentTab)
}

func (this *UserList) Reset() {
	this.userNickListBox = widgets.NewList()
	//userNickListBox.Title = "用户列表"
	//userNickListBox.BorderStyle = ui.NewStyle(ui.ColorMagenta)
	//userNickListBox.Border = true
	this.userNickListBox.TextStyle = ui.NewStyle(ui.ColorYellow)
	this.userNickListBox.WrapText = false
	this.userNickListBox.SelectedRowStyle = ui.NewStyle(ui.ColorWhite,
		ui.ColorRed)

	this.userNickListBox.SetRect(this.baseX, this.baseY+2,
		this.baseX+this.width,
		this.baseY+this.height)

	groupNickListBox := widgets.NewList()
	//userNickListBox.Title = "用户列表"
	//userNickListBox.BorderStyle = ui.NewStyle(ui.ColorMagenta)
	//userNickListBox.Border = true
	groupNickListBox.TextStyle = ui.NewStyle(ui.ColorYellow)
	groupNickListBox.WrapText = false
	groupNickListBox.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorRed)

	groupNickListBox.SetRect(this.baseX, this.baseY+2, this.baseX+this.width,
		this.baseY+this.height)

	this.groupNickListBox = groupNickListBox

	this.recentListBox = widgets.NewList()
	this.recentListBox.TextStyle = ui.NewStyle(ui.ColorYellow)
	this.recentListBox.WrapText = false
	this.recentListBox.SelectedRowStyle = ui.NewStyle(ui.ColorWhite,
		ui.ColorRed)

	this.recentListBox.SetRect(this.baseX, this.baseY+2,
		this.baseX+this.width,
		this.baseY+this.height)

	nickList := make([]string, len(this.userList))
	groupList := make([]string, len(this.groupList))

	for i, user := range this.userList {
		if user.DisplayName != "" {
			nickList[i] = user.DisplayName
		} else {
			nickList[i] = user.Nick
		}
	}
	this.userNickListBox.Rows = nickList

	for i, group := range this.groupList {
		groupList[i] = group.Name
	}
	this.groupNickListBox.Rows = groupList

	recentList := merge(this.recentUserList, this.recentGroupList)
	recentNickList := make([]string, recentList.Len())
	for i, r := range *recentList {
		recentNickList[i] = r.Nick
	}
	this.recentListBox.Rows = recentNickList

	this.tabPane = widgets.NewTabPane("聊天列表", "好友列表", "群列表")
	this.tabPane.SetRect(this.baseX, this.baseY, this.baseX+this.width,
		this.baseY+3)
	this.tabPane.Border = true
	tabWidgets := make([]*widgets.List, 3)
	tabWidgets[0] = this.recentListBox
	tabWidgets[1] = this.userNickListBox
	tabWidgets[2] = this.groupNickListBox
	this.tabWidgets = tabWidgets

	this.currentTab = this.userNickListBox

	ui.Render(this.tabPane, this.currentTab)
}

func NewUserList(recentUserList []UserInfo, recentGroupList []Group, userList []UserInfo,
	groupList []Group,
	userChangeEvent chan UserChangeEvent, width, height, baseX, baseY int,
	logger *log.Logger) *UserList {

	//	chinese := false

	l := &UserList{
		curUserIndex:    0,
		logger:          logger,
		baseX:           baseX,
		baseY:           baseY,
		width:           width,
		height:          height,
		recentUserList:  recentUserList,
		recentGroupList: recentGroupList,
		userList:        userList,
		groupList:       groupList,
		userChangeEvent: userChangeEvent,
	}

	l.Reset()

	return l

}

func (l *UserList) nextUser() {
	l.currentTab.ScrollDown()
	l.curUserIndex = l.currentTab.SelectedRow
	ui.Render(l.currentTab)
}

func (l *UserList) prevUser() {
	l.currentTab.ScrollUp()
	l.curUserIndex = l.currentTab.SelectedRow
	ui.Render(l.currentTab)
}

func (l *UserList) focuseLeft() {
	l.tabPane.FocusLeft()
	l.renderTab()
}

func (l *UserList) focuseRight() {
	l.tabPane.FocusRight()
	l.renderTab()
}

func (this *UserList) Action(e ui.Event) bool {
	if !this.picked {
		return false
	}
	switch e.ID {
	case "j":
		this.nextUser()
		return true
	case "k":
		this.prevUser()
		return true
	case "h":
		this.focuseLeft()
		return true
	case "l":
		this.focuseRight()
		return true
	}
	return false

}
