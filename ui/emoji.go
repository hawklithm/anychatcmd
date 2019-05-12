package ui

import (
	"regexp"
	"strconv"
	"strings"
)

var emojiRegex = `\[\S{1,3}\]`

var emojiTagRegex = `<span\sclass=\"emoji\semoji([0-9]|[a-f]){0,6}\"><\/span>`

var emojiMap = make(map[string]string)

func init() {
	emojiMap["微笑"] = "🙂"
	emojiMap["撇嘴"] = "😟"
	emojiMap["色"] = "😍"
	emojiMap["发呆"] = "😲"
	emojiMap["得意"] = "😎"
	emojiMap["流泪"] = "😢"
	emojiMap["害羞"] = "😊"
	emojiMap["闭嘴"] = "🤐"
	emojiMap["睡"] = "😪"
	emojiMap["大哭"] = "😭"
	emojiMap["尴尬"] = "😅"
	emojiMap["发怒"] = "🤬"
	emojiMap["调皮"] = "😜"
	emojiMap["调皮"] = "😜"
	emojiMap["呲牙"] = "😁"
	emojiMap["惊讶"] = "😮"
	emojiMap["难过"] = "🙁"
	emojiMap["囧"] = "😳"
	emojiMap["抓狂"] = "😩"
	emojiMap["吐"] = "🤮"
	emojiMap["愉快"] = "😊"
	emojiMap["白眼"] = "🙄"
	emojiMap["傲慢"] = "😕"
	emojiMap["困"] = "😴"
	emojiMap["惊恐"] = "😱"
	emojiMap["流汗"] = "😅"
	emojiMap["憨笑"] = "😄"
	emojiMap["悠闲"] = "😏"
	emojiMap["咒骂"] = "😤"
	emojiMap["奋斗"] = "💪"
	emojiMap["疑问"] = "❓"
	emojiMap["晕"] = "😖"
	emojiMap["嘘"] = "🤫"
	emojiMap["衰"] = "🥵"
	emojiMap["骷髅"] = "💀"
	emojiMap["敲打"] = "🔨"
	emojiMap["再见"] = "👋"
	emojiMap["擦汗"] = "😅"
	emojiMap["抠鼻"] = "🌝"
	emojiMap["鼓掌"] = "👏"
	emojiMap["坏笑"] = "👻"
	emojiMap["左哼哼"] = "😾"
	emojiMap["右哼哼"] = "😾"
	emojiMap["哈欠"] = "😪"
	emojiMap["鄙视"] = "👎"
	emojiMap["委屈"] = "😢"
	emojiMap["快哭了"] = "😔"
	emojiMap["阴险"] = "😈"
	emojiMap["亲亲"] = "😚"
	emojiMap["可怜"] = "🥺"
	emojiMap["菜刀"] = "🔪"
	emojiMap["西瓜"] = "🍉"
	emojiMap["啤酒"] = "🍺"
	emojiMap["咖啡"] = "☕️"
	emojiMap["猪头"] = "🐷"
	emojiMap["玫瑰"] = "🌹"
	emojiMap["凋谢"] = "👿"
	emojiMap["嘴唇"] = "👄"
	emojiMap["爱心"] = "❤️"
	emojiMap["心碎"] = "💔️"
	emojiMap["蛋糕"] = "🍰️"
	emojiMap["炸弹"] = "💣️"
	emojiMap["便便"] = "💩"
	emojiMap["月亮"] = "🌜️"
	emojiMap["太阳"] = "☀️️"
	emojiMap["拥抱"] = "🤗"
	emojiMap["强"] = "👍️️"
	emojiMap["弱"] = "👎️️"
	emojiMap["握手"] = "🤝"
	emojiMap["胜利"] = "✌️"
	emojiMap["抱拳"] = "🙏"
	emojiMap["拳头"] = "✊"
	emojiMap["OK"] = "👌"
	emojiMap["跳跳"] = "💃"
	emojiMap["发抖"] = "😖"
	emojiMap["怄火"] = "😡"
	emojiMap["转圈"] = "🤸‍️"
	emojiMap["嘿哈"] = "🤪"
	emojiMap["捂脸"] = "🤦‍"
	emojiMap["奸笑"] = "🥴"
	emojiMap["机智"] = "🥳"
	emojiMap["皱眉"] = "🙍"
	emojiMap["耶"] = "✌️"
	emojiMap["红包"] ="🧧"
	emojiMap["發"] = "🤩"



}

var f = func(s string) string {
	//v, _ := strconv.ParseFloat(s, 32)
	s = strings.Trim(s, "[]")
	//t := s[1 : len(s)-1]
	//fmt.Println("len=", len(s), "sub=", s)
	if emojiMap[s] != "" {
		return emojiMap[s]
	} else {
		return "[" + s + "]"
	}
}

var f2 = func(s string) string {
	t := s[24 : len(s)-9]
	r, _ := strconv.ParseInt(t, 16, 32)
	var ru []rune
	ru = append(ru, rune(r))
	return string(ru)
}

func ConvertToEmoji(sentence string) string {
	re, _ := regexp.Compile(emojiRegex)
	str2 := re.ReplaceAllStringFunc(sentence, f)
	return str2
}

func TranslateEmoji(sentence string) string {
	re, _ := regexp.Compile(emojiTagRegex)

	str2 := re.ReplaceAllStringFunc(sentence, f2)
	return str2
}
