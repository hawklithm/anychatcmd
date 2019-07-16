// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"os"
	"strings"

	. "github.com/hawklithm/termui"
)

type ImageList struct {
	Block
	Rows             []*ImageListItem
	WrapText         bool
	TextStyle        Style
	SelectedRow      int
	topRow           int
	SelectedRowStyle Style
}

type ImageListItem struct {
	Block
	Text             string
	img              image.Image
	imgMatrix        []byte
	Url              string
	TextStyle        Style
	WrapText         bool
	SelectedRowStyle Style
}

func clip(origin image.Image) ([]byte, error) {
	emptyBuff := bytes.NewBuffer(nil)
	if err := png.Encode(emptyBuff, origin); err != nil {
		return nil, err
	}
	//switch origin.(type) {
	//case *image.YCbCr:
	//	if err := jpeg.Encode(emptyBuff, origin, nil); err != nil {
	//		return nil, err
	//	}
	//case *image.NRGBA:
	//	if err := png.Encode(emptyBuff, origin); err != nil {
	//		return nil, err
	//	}
	//case *image.RGBA:
	//	if err := png.Encode(emptyBuff, origin); err != nil {
	//		return nil, err
	//	}
	//case *image.Paletted:
	//	if err := gif.Encode(emptyBuff, origin, &gif.Options{}); err != nil {
	//		return nil, err
	//	}
	//}
	dist := base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
	return []byte(dist), nil
}

func (l *ImageListItem) SetImage(origin image.Image) error {
	if os.Getenv("WECHAT_TERM") == "iterm" {
		var err error
		if l.imgMatrix, err = clip(origin); err != nil {
			return err
		}
	} else {
		l.img = origin
	}
	return nil
}

func NewImageListItem() *ImageListItem {
	block := NewBlock()
	block.Border = false
	return &ImageListItem{
		Block:            *block,
		TextStyle:        Theme.List.Text,
		SelectedRowStyle: Theme.List.Text,
	}
}

func NewImageList() *ImageList {
	return &ImageList{
		Block:            *NewBlock(),
		TextStyle:        Theme.List.Text,
		SelectedRowStyle: Theme.List.Text,
	}
}

func (self *ImageListItem) GetHeight() int {
	if !self.hasImage() {
		rows := strings.Split(self.Text, "\n")
		return len(rows) + 2
	} else {
		_, height, _, _ := self.calcImageSize()
		return height + 2
	}

}

func (self *ImageListItem) calcImageSize() (width, height int, imgScale, whratio float64) {
	if os.Getenv("WECHAT_TERM") == "iterm" {
		return 16, 4, 1, DefaultRatio
	}
	if !self.hasImage() {
		return -1, -1, -1, 0
	}
	img := self.img
	width = 40
	height = 40
	whratio = DefaultRatio

	bounds := img.Bounds()
	imgW, imgH := bounds.Dx(), bounds.Dy()

	imgScale = Scale(imgW, imgH, width, height, whratio)

	// Resize canvas to fit scaled image
	width, height = int(float64(imgW)/imgScale), int(float64(imgH)/(imgScale*whratio))
	return width, height, imgScale, whratio
}

func (self *ImageListItem) drawImage(buf *Buffer, selected bool) {
	if selected {
		self.Border = true
	} else {
		self.Border = false
	}
	defer self.Block.Draw(buf)
	if os.Getenv("WECHAT_TERM") == "iterm" {
		matrix := self.imgMatrix
		if self.hasImage() {
			buf.SetCell(Cell{Bytes: matrix}, image.Point{X: 1, Y: 0}.Add(self.Inner.Min))
		}
		return
	}
	img := self.img
	width, height, imgScale, whratio := self.calcImageSize()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate average color for the corresponding image rectangle
			// fitting in this cell. We use a half-block trick, wherein the
			// lower half of the cell displays the character ▄, effectively
			// doubling the resolution of the canvas.
			startX, startY, endX, endY := ImgArea(x, y, imgScale, whratio)

			r, g, b := AvgRGB(img, startX, startY, endX, (startY+endY)/2)
			colorUp := Color(TermColor(r, g, b))

			r, g, b = AvgRGB(img, startX, (startY+endY)/2, endX, endY)
			colorDown := Color(TermColor(r, g, b))

			buf.SetCell(Cell{Rune: '▄', Style: Style{Fg: colorDown - 1,
				Bg: colorUp - 1}}, image.Point{X: x, Y: y}.Add(self.Inner.Min))
		}
	}
}

func (self *ImageListItem) draw(buf *Buffer, selected bool) {

	if selected {
		self.Border = true
	} else {
		self.Border = false
	}

	self.Block.Draw(buf)

	cells := ParseStyles(self.Text, self.TextStyle)
	if self.WrapText {
		cells = WrapCells(cells, uint(self.Inner.Dx()))
	}

	rows := SplitCells(cells, '\n')

	for y, row := range rows {
		if y+self.Inner.Min.Y >= self.Inner.Max.Y {
			break
		}
		row = TrimCells(row, self.Inner.Dx())
		for _, cx := range BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(self.Inner.Min))
		}
	}
}

func (self *ImageList) SelectedLine() int {
	return self.convertRowToLine(self.SelectedRow)
}

func (self *ImageList) TopLine() int {
	height := 0
	if len(self.Rows) == 0 {
		return height
	}
	if self.topRow > len(self.Rows) {
		self.topRow = len(self.Rows)
	}
	if self.topRow < 0 {
		self.topRow = 0
	}
	for _, row := range self.Rows[:self.topRow] {
		height += row.GetHeight()
	}
	return height
}

func (self *ImageList) convertRowToLine(rowNum int) int {
	height := 0
	//fmt.Println(rowNum)
	if len(self.Rows) == 0 {
		return 0
	}
	for _, row := range self.Rows[:rowNum+1] {
		height += row.GetHeight()
	}
	return height
}

func (self *ImageList) convertLineToRow(line int) int {
	height := 0
	for i, row := range self.Rows {
		height += row.GetHeight()
		if height >= line {
			return i + 1
		}
	}
	return len(self.Rows)
}

func (self *ImageListItem) hasImage() bool {
	return self.img != nil || (self.imgMatrix != nil && len(self.imgMatrix) != 0)
}

func (self *ImageList) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	point := self.Inner.Min

	// adjusts view into widget
	var topLine int
	if self.SelectedLine() >= self.Inner.Dy()+self.TopLine() {
		topLine = self.SelectedLine() - self.Inner.Dy() + 1
		self.topRow = self.convertLineToRow(topLine)
	} else if self.SelectedRow < self.topRow {
		self.topRow = self.SelectedRow
	}

	if len(self.Rows) == 0 {
		return
	}

	// draw rows
	for row := self.topRow; row < len(self.Rows) && point.Y < self.Inner.Max.Y; row++ {
		height := self.Rows[row].GetHeight()
		self.Rows[row].SetRect(self.Inner.Min.X, point.Y, self.Inner.Max.X,
			point.Y+height)
		if self.Rows[row].hasImage() {
			self.Rows[row].drawImage(buf, self.SelectedRow == row)
		} else {
			self.Rows[row].draw(buf, self.SelectedRow == row)
		}
		point.Y += height
	}

	// draw UP_ARROW if needed
	if self.topRow > 0 {
		buf.SetCell(
			NewCell(UP_ARROW, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Max.X-1, self.Inner.Min.Y),
		)
	}

	// draw DOWN_ARROW if needed
	if self.convertRowToLine(len(self.Rows)-1) > int(self.TopLine())+self.Inner.
		Dy() {
		buf.SetCell(
			NewCell(DOWN_ARROW, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Max.X-1, self.Inner.Max.Y-1),
		)
	}
}

// ScrollAmount scrolls by amount given. If amount is < 0, then scroll up.
// There is no need to set self.topRow, as this will be set automatically when drawn,
// since if the selected item is off screen then the topRow variable will change accordingly.
func (self *ImageList) ScrollAmount(amount int) {
	if len(self.Rows)-int(self.SelectedRow) <= amount {
		self.SelectedRow = len(self.Rows) - 1
	} else if int(self.SelectedRow)+amount < 0 {
		self.SelectedRow = 0
	} else {
		self.SelectedRow += amount
	}
}

func (self *ImageList) ScrollUp() {
	self.ScrollAmount(-1)
}

func (self *ImageList) ScrollDown() {
	self.ScrollAmount(1)
}

func (self *ImageList) ScrollPageUp() {
	// If an item is selected below top row, then go to the top row.
	if self.SelectedRow > self.topRow {
		self.SelectedRow = self.topRow
	} else {
		self.ScrollAmount(-self.Inner.Dy())
	}
}

func (self *ImageList) ScrollPageDown() {
	self.ScrollAmount(self.Inner.Dy())
}

func (self *ImageList) ScrollHalfPageUp() {
	self.ScrollAmount(-int(FloorFloat64(float64(self.Inner.Dy()) / 2)))
}

func (self *ImageList) ScrollHalfPageDown() {
	self.ScrollAmount(int(FloorFloat64(float64(self.Inner.Dy()) / 2)))
}

func (self *ImageList) ScrollTop() {
	self.SelectedRow = 0
}

func (self *ImageList) ScrollBottom() {
	self.SelectedRow = len(self.Rows) - 1
}
