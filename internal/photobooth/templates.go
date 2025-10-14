package photobooth

import (
	"image"
	"sort"
)

const (
	// CanvasWidth and CanvasHeight are sized for a 4x6 inch photo at 300 DPI.
	CanvasWidth  = 1200
	CanvasHeight = 1800

	margin = 40
	gap    = 20
)

// Template describes a collage layout.
type Template struct {
	Name        string
	Description string
	Slots       []image.Rectangle
}

var templates = []Template{
	singleFull(),
	singlePolaroid(),
	doubleVertical(),
	doubleHorizontal(),
	tripleStack(),
	tripleColumn(),
	tripleLarge(),
	quadGrid(),
	quadStrip(),
	fiveCenterpiece(),
	fivePanorama(),
	sixGrid(),
	sixStrip(),
	sevenMosaic(),
	sevenStory(),
	eightGrid(),
	eightPinwheel(),
	nineGrid(),
	nineHighlight(),
}

// Templates returns the list of templates sorted by name.
func Templates() []Template {
	list := make([]Template, len(templates))
	copy(list, templates)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}

func singleFull() Template {
	return Template{
		Name:        "single-full",
		Description: "ภาพเดี่ยวเต็มเฟรม พร้อมขอบขาวบางๆ",
		Slots: []image.Rectangle{
			image.Rect(margin, margin, CanvasWidth-margin, CanvasHeight-margin),
		},
	}
}

func singlePolaroid() Template {
	bottomMargin := margin + 160
	return Template{
		Name:        "single-polaroid",
		Description: "ภาพเดี่ยวสไตล์โพลารอยด์ มีพื้นที่เขียนข้อความด้านล่าง",
		Slots: []image.Rectangle{
			image.Rect(margin+60, margin, CanvasWidth-margin-60, CanvasHeight-bottomMargin),
		},
	}
}

func doubleVertical() Template {
	slotHeight := (CanvasHeight - 2*margin - gap) / 2
	return Template{
		Name:        "double-vertical",
		Description: "สองภาพบน-ล่าง ขนาดเท่ากัน",
		Slots: []image.Rectangle{
			image.Rect(margin, margin, CanvasWidth-margin, margin+slotHeight),
			image.Rect(margin, margin+slotHeight+gap, CanvasWidth-margin, CanvasHeight-margin),
		},
	}
}

func doubleHorizontal() Template {
	slotWidth := (CanvasWidth - 2*margin - gap) / 2
	return Template{
		Name:        "double-horizontal",
		Description: "สองภาพซ้าย-ขวา ขนาดเท่ากัน",
		Slots: []image.Rectangle{
			image.Rect(margin, margin, margin+slotWidth, CanvasHeight-margin),
			image.Rect(margin+slotWidth+gap, margin, CanvasWidth-margin, CanvasHeight-margin),
		},
	}
}

func tripleStack() Template {
	slotHeight := (CanvasHeight - 2*margin - 2*gap) / 3
	rects := make([]image.Rectangle, 3)
	for i := 0; i < 3; i++ {
		top := margin + i*(slotHeight+gap)
		rects[i] = image.Rect(margin, top, CanvasWidth-margin, top+slotHeight)
	}
	return Template{
		Name:        "triple-stack",
		Description: "สามภาพวางซ้อนแนวตั้ง",
		Slots:       rects,
	}
}

func tripleColumn() Template {
	slotWidth := (CanvasWidth - 2*margin - 2*gap) / 3
	rects := make([]image.Rectangle, 3)
	for i := 0; i < 3; i++ {
		left := margin + i*(slotWidth+gap)
		rects[i] = image.Rect(left, margin, left+slotWidth, CanvasHeight-margin)
	}
	return Template{
		Name:        "triple-column",
		Description: "สามภาพเรียงแนวนอน",
		Slots:       rects,
	}
}

func tripleLarge() Template {
	topHeight := (CanvasHeight - 2*margin - gap) * 2 / 3
	bottomHeight := CanvasHeight - 2*margin - topHeight - gap
	smallWidth := (CanvasWidth - 2*margin - gap) / 2
	return Template{
		Name:        "triple-large",
		Description: "ภาพหลักใหญ่ด้านบน พร้อมสองภาพเล็กด้านล่าง",
		Slots: []image.Rectangle{
			image.Rect(margin, margin, CanvasWidth-margin, margin+topHeight),
			image.Rect(margin, margin+topHeight+gap, margin+smallWidth, margin+topHeight+gap+bottomHeight),
			image.Rect(CanvasWidth-margin-smallWidth, margin+topHeight+gap, CanvasWidth-margin, CanvasHeight-margin),
		},
	}
}

func quadGrid() Template {
	slotWidth := (CanvasWidth - 2*margin - gap) / 2
	slotHeight := (CanvasHeight - 2*margin - gap) / 2
	return Template{
		Name:        "quad-grid",
		Description: "สี่ภาพจัดเป็นตาราง 2x2",
		Slots: []image.Rectangle{
			image.Rect(margin, margin, margin+slotWidth, margin+slotHeight),
			image.Rect(margin+slotWidth+gap, margin, CanvasWidth-margin, margin+slotHeight),
			image.Rect(margin, margin+slotHeight+gap, margin+slotWidth, CanvasHeight-margin),
			image.Rect(margin+slotWidth+gap, margin+slotHeight+gap, CanvasWidth-margin, CanvasHeight-margin),
		},
	}
}

func quadStrip() Template {
	slotHeight := (CanvasHeight - 2*margin - 3*gap) / 4
	rects := make([]image.Rectangle, 4)
	for i := 0; i < 4; i++ {
		top := margin + i*(slotHeight+gap)
		rects[i] = image.Rect(margin, top, CanvasWidth-margin, top+slotHeight)
	}
	return Template{
		Name:        "quad-strip",
		Description: "สี่ภาพเรียงยาวแนวตั้ง",
		Slots:       rects,
	}
}

func fiveCenterpiece() Template {
	cornerWidth := (CanvasWidth - 2*margin - gap) / 2
	cornerHeight := (CanvasHeight - 2*margin - gap) / 2
	left := margin
	top := margin
	right := CanvasWidth - margin
	bottom := CanvasHeight - margin

	return Template{
		Name:        "five-centerpiece",
		Description: "ภาพใหญ่ตรงกลาง รายล้อมด้วยภาพเล็กที่มุม",
		Slots: []image.Rectangle{
			image.Rect(left, top, left+cornerWidth, top+cornerHeight),
			image.Rect(right-cornerWidth, top, right, top+cornerHeight),
			image.Rect(left, bottom-cornerHeight, left+cornerWidth, bottom),
			image.Rect(right-cornerWidth, bottom-cornerHeight, right, bottom),
			image.Rect(left+cornerWidth+gap/2, top+cornerHeight+gap/2, right-cornerWidth-gap/2, bottom-cornerHeight-gap/2),
		},
	}
}

func fivePanorama() Template {
	topHeight := (CanvasHeight - 2*margin - gap) / 2
	bottomHeight := CanvasHeight - 2*margin - topHeight - gap
	smallWidth := (CanvasWidth - 2*margin - 3*gap) / 4
	rects := []image.Rectangle{
		image.Rect(margin, margin, CanvasWidth-margin, margin+topHeight),
	}
	for i := 0; i < 4; i++ {
		left := margin + i*(smallWidth+gap)
		rects = append(rects, image.Rect(left, margin+topHeight+gap, left+smallWidth, margin+topHeight+gap+bottomHeight))
	}
	return Template{
		Name:        "five-panorama",
		Description: "ภาพพาโนรามาด้านบน พร้อมสี่ภาพเรียงด้านล่าง",
		Slots:       rects,
	}
}

func sixGrid() Template {
	slotWidth := (CanvasWidth - 2*margin - 2*gap) / 3
	slotHeight := (CanvasHeight - 2*margin - gap) / 2
	return Template{
		Name:        "six-grid",
		Description: "หกภาพจัดเป็นตาราง 3x2",
		Slots:       createGrid(3, 2, slotWidth, slotHeight),
	}
}

func sixStrip() Template {
	slotHeight := (CanvasHeight - 2*margin - 5*gap) / 6
	rects := make([]image.Rectangle, 6)
	for i := 0; i < 6; i++ {
		top := margin + i*(slotHeight+gap)
		rects[i] = image.Rect(margin, top, CanvasWidth-margin, top+slotHeight)
	}
	return Template{
		Name:        "six-strip",
		Description: "หกภาพแนวตั้งสไตล์ Photo Booth",
		Slots:       rects,
	}
}

func sevenMosaic() Template {
	topHeight := (CanvasHeight - 2*margin - 2*gap) / 3
	middleHeight := topHeight
	bottomHeight := CanvasHeight - 2*margin - topHeight - middleHeight - 2*gap
	smallWidth := (CanvasWidth - 2*margin - 2*gap) / 3
	largeWidth := (CanvasWidth - 2*margin - gap) / 2

	return Template{
		Name:        "seven-mosaic",
		Description: "เจ็ดภาพสลับขนาดแบบโมเสก",
		Slots: []image.Rectangle{
			image.Rect(margin, margin, margin+largeWidth, margin+topHeight),
			image.Rect(margin+largeWidth+gap, margin, CanvasWidth-margin, margin+topHeight),
			image.Rect(margin, margin+topHeight+gap, margin+smallWidth, margin+topHeight+gap+middleHeight),
			image.Rect(margin+smallWidth+gap, margin+topHeight+gap, margin+2*smallWidth+gap, margin+topHeight+gap+middleHeight),
			image.Rect(CanvasWidth-margin-smallWidth, margin+topHeight+gap, CanvasWidth-margin, margin+topHeight+gap+middleHeight),
			image.Rect(margin, CanvasHeight-margin-bottomHeight, margin+smallWidth, CanvasHeight-margin),
			image.Rect(CanvasWidth-margin-smallWidth, CanvasHeight-margin-bottomHeight, CanvasWidth-margin, CanvasHeight-margin),
		},
	}
}

func sevenStory() Template {
	topHeight := (CanvasHeight - 2*margin - 3*gap) / 4
	middleHeight := topHeight
	bottomHeight := CanvasHeight - 2*margin - topHeight - middleHeight - 2*gap
	middleWidth := (CanvasWidth - 2*margin - gap) / 2
	bottomWidth := (CanvasWidth - 2*margin - 3*gap) / 4

	rects := []image.Rectangle{
		image.Rect(margin, margin, CanvasWidth-margin, margin+topHeight),
		image.Rect(margin, margin+topHeight+gap, margin+middleWidth, margin+topHeight+gap+middleHeight),
		image.Rect(CanvasWidth-margin-middleWidth, margin+topHeight+gap, CanvasWidth-margin, margin+topHeight+gap+middleHeight),
	}

	for i := 0; i < 4; i++ {
		left := margin + i*(bottomWidth+gap)
		rects = append(rects, image.Rect(left, CanvasHeight-margin-bottomHeight, left+bottomWidth, CanvasHeight-margin))
	}

	return Template{
		Name:        "seven-story",
		Description: "หนึ่งภาพใหญ่ด้านบน ตามด้วยสองแถวกลางและแถวล่างสี่ภาพ",
		Slots:       rects,
	}
}

func eightGrid() Template {
	slotWidth := (CanvasWidth - 2*margin - 3*gap) / 4
	slotHeight := (CanvasHeight - 2*margin - gap) / 2
	return Template{
		Name:        "eight-grid",
		Description: "แปดภาพจัดเป็นกริด 4x2",
		Slots:       createGrid(4, 2, slotWidth, slotHeight),
	}
}

func eightPinwheel() Template {
	slotWidth := (CanvasWidth - 2*margin - 3*gap) / 4
	slotHeight := (CanvasHeight - 2*margin - 3*gap) / 4
	left := margin
	top := margin
	right := CanvasWidth - margin
	bottom := CanvasHeight - margin

	return Template{
		Name:        "eight-pinwheel",
		Description: "กรอบแปดภาพล้อมรอบช่องว่างตรงกลาง",
		Slots: []image.Rectangle{
			image.Rect(left, top, left+slotWidth, top+slotHeight),
			image.Rect(left+slotWidth+gap, top, left+2*slotWidth+gap, top+slotHeight),
			image.Rect(right-slotWidth, top, right, top+slotHeight),
			image.Rect(right-slotWidth, top+slotHeight+gap, right, top+2*slotHeight+gap),
			image.Rect(right-slotWidth, bottom-slotHeight, right, bottom),
			image.Rect(left+slotWidth+gap, bottom-slotHeight, left+2*slotWidth+gap, bottom),
			image.Rect(left, bottom-slotHeight, left+slotWidth, bottom),
			image.Rect(left, top+slotHeight+gap, left+slotWidth, top+2*slotHeight+gap),
		},
	}
}

func nineGrid() Template {
	slotWidth := (CanvasWidth - 2*margin - 2*gap) / 3
	slotHeight := (CanvasHeight - 2*margin - 2*gap) / 3
	return Template{
		Name:        "nine-grid",
		Description: "เก้าภาพจัดเป็นตาราง 3x3",
		Slots:       createGrid(3, 3, slotWidth, slotHeight),
	}
}

func nineHighlight() Template {
	leftWidth := (CanvasWidth - 2*margin - gap) * 2 / 3
	rightWidth := CanvasWidth - 2*margin - leftWidth - gap
	slotHeight := (CanvasHeight - 2*margin - 2*gap) / 3
	smallWidth := (rightWidth - gap) / 2

	slots := []image.Rectangle{
		image.Rect(margin, margin, margin+leftWidth, margin+slotHeight),
		image.Rect(margin, margin+slotHeight+gap, margin+leftWidth, margin+2*slotHeight+gap),
		image.Rect(margin, CanvasHeight-margin-slotHeight, margin+leftWidth, CanvasHeight-margin),
	}

	rightStart := margin + leftWidth + gap
	for row := 0; row < 3; row++ {
		top := margin + row*(slotHeight+gap)
		for col := 0; col < 2; col++ {
			left := rightStart + col*(smallWidth+gap)
			slots = append(slots, image.Rect(left, top, left+smallWidth, top+slotHeight))
		}
	}

	return Template{
		Name:        "nine-highlight",
		Description: "สามภาพใหญ่ทางซ้าย พร้อมหกภาพเล็กทางขวา",
		Slots:       slots,
	}
}

func createGrid(columns, rows, slotWidth, slotHeight int) []image.Rectangle {
	rects := make([]image.Rectangle, 0, columns*rows)
	for y := 0; y < rows; y++ {
		top := margin + y*(slotHeight+gap)
		for x := 0; x < columns; x++ {
			left := margin + x*(slotWidth+gap)
			rects = append(rects, image.Rect(left, top, left+slotWidth, top+slotHeight))
		}
	}
	return rects
}
