package tableimage

import (
	"image"
	"math"
)

// Cell in table
type Cell struct {
	// Text content of a cell
	Text string `json:"text,omitempty"`
	// Image image for a cell
	Image *Image `json:"image,omitempty"`
	// Style for cell
	Style *Style `json:"style,omitempty"`
	// IgnoreInlineStyle ignore inline text style parsing
	IgnoreInlineStyle bool `json:"ignore_inline_style,omitempty"`
}

// Draw render cell to image
func (c Cell) Draw(img *image.RGBA, bounds image.Rectangle) {
	if c.Style == nil {
		return
	}
	c.drawBorderAndBg(img, bounds)
	if c.Style.Font == nil {
		return
	}
	var (
		imgSize    image.Point
		imgXOffset int
	)
	if c.Image != nil && c.Image.Data != nil {
		imgSize = c.ImageSize()
		if c.Image.Align == LEFT || c.Image.Align == RIGHT {
			imgXOffset = imgSize.X
		}
	}
	lineHeight := stringHeight(c.Style.Font.Size, c.Style.LineHeight)
	lines, _ := c.Wrap(imgXOffset)
	innerBounds := c.InnerBounds(bounds)
	var (
		textStartX int
		textHeight = len(lines) * lineHeight
		y          int
	)
	switch c.Style.VAlign {
	case MIDDLE:
		middle := (innerBounds.Dy() - textHeight - imgSize.Y) / 2
		y = innerBounds.Min.Y + middle
	case BOTTOM:
		y = innerBounds.Max.Y - textHeight - imgSize.Y
	default:
		y = innerBounds.Min.Y
	}
	textStartX, y = c.drawImage(img, y, textHeight, imgSize, innerBounds)
	c.drawText(img, lines, textStartX, imgXOffset, y, lineHeight, innerBounds)
}

func (c Cell) drawBorderAndBg(img *image.RGBA, bounds image.Rectangle) {
	if c.Style.BgColor != "" {
		drawRect(img, bounds, "", c.Style.BgColor, 0)
	}
	if c.Style.Border != nil {
		c.Style.Border.Draw(img, bounds)
	}
}

func (c Cell) drawImage(img *image.RGBA, y int, textHeight int, imgSize image.Point, innerBounds image.Rectangle) (int, int) {
	if c.Image == nil || c.Image.Data == nil {
		return 0, y
	}
	imgX := c.Image.PaddingLeft()
	imgY := c.Image.PaddingTop()
	var textStartX int
	if c.Image.VAlign == TOP || c.Image.VAlign == BOTTOM {
		imgX, imgY, y = calcImageVAlign(c.Image.VAlign, c.Style.Align, imgX, imgY, y, textHeight, imgSize, innerBounds)
	} else if c.Image.Align == LEFT || c.Image.Align == RIGHT {
		if c.Image.Align == LEFT {
			textStartX = imgSize.X
		}
		imgX, imgY, y = calcImageAlign(c.Image.Align, c.Style.VAlign, imgX, imgY, y, textHeight, imgSize, innerBounds)
	}
	pt := image.Pt(imgX, imgY)
	drawImage(img, c.Image, pt)
	return textStartX, y
}

func calcImageVAlign(valign VAlign, align Align, imgX int, imgY int, y int, textHeight int, imgSize image.Point, innerBounds image.Rectangle) (int, int, int) {
	if valign == TOP {
		imgY += y
		y += imgSize.Y
	} else {
		imgY += y + textHeight
	}
	switch align {
	case RIGHT:
		imgX += innerBounds.Max.X - imgSize.X
	case CENTER:
		center := (innerBounds.Dx() - imgSize.X) / 2
		imgX += innerBounds.Min.X + center
	default:
		imgX += innerBounds.Min.X
	}
	return imgX, imgY, y
}

func calcImageAlign(align Align, valign VAlign, imgX int, imgY int, y int, textHeight int, imgSize image.Point, innerBounds image.Rectangle) (int, int, int) {
	if align == LEFT {
		imgX += innerBounds.Min.X
	} else {
		imgX += innerBounds.Max.X - imgSize.X
	}
	switch valign {
	case BOTTOM:
		imgY = innerBounds.Max.Y - imgSize.Y
		y = innerBounds.Max.Y - textHeight
	case MIDDLE:
		middle := (innerBounds.Dy() - imgSize.Y) / 2
		imgY += innerBounds.Min.Y + middle
		middle = (innerBounds.Dy() - textHeight) / 2
		y = innerBounds.Min.Y + middle
	default:
		imgY += innerBounds.Min.Y
	}
	return imgX, imgY, y
}

func (c Cell) drawText(img *image.RGBA, lines []Word, textStartX int, imgXOffset int, y int, lineHeight int, innerBounds image.Rectangle) {
	for _, line := range lines {
		var x int
		switch c.Style.Align {
		case RIGHT:
			if textStartX > 0 {
				x = innerBounds.Max.X - line.Width()
			} else {
				x = innerBounds.Max.X - line.Width() - imgXOffset
			}
		case CENTER:
			center := (innerBounds.Dx() - line.Width() - imgXOffset) / 2
			x = innerBounds.Min.X + center
		default:
			x = innerBounds.Min.X + textStartX
		}
		pt := image.Pt(x, y)
		for _, txt := range line {
			if txt.Color == "" {
				txt.Color = c.Style.Color
			}
			txtBounds := image.Rect(pt.X, pt.Y, pt.X+txt.Width, pt.Y+lineHeight)
			drawText(img, txtBounds, &txt, c.Style.Font)
			pt = pt.Add(image.Pt(txt.Width, 0))
		}
		y += lineHeight
	}
}

// Wrap wraps cell content returns paragraphs, and max content width
func (c Cell) Wrap(xOffset int) ([]Word, int) {
	if c.Style == nil || c.Style.Font == nil {
		return nil, 0
	}
	maxWidth := c.Style.MaxWidth - xOffset
	fontFace := newFontFace(c.Style.Font.Font, c.Style.Font.Size)
	return wrap(c.Text, maxWidth, fontFace, c.IgnoreInlineStyle)
}

// ImageSize get image size, will update Image.Size based on max width setting
func (c Cell) ImageSize() image.Point {
	if c.Style == nil || c.Image == nil || c.Image.Data == nil {
		return image.ZP
	}
	if c.Style.MaxWidth > 0 && c.Style.MaxWidth < c.Image.BoundSize().X {
		maxWidth := float64(c.Style.MaxWidth - c.Image.PaddingX())
		scale := maxWidth / float64(c.Image.Size.X)
		c.Image.Size.X = int(math.Round(float64(c.Image.Size.X) * scale))
		c.Image.Size.Y = int(math.Round(float64(c.Image.Size.Y) * scale))
	}
	return c.Image.BoundSize()
}

// Size returns cell width/height
func (c Cell) Size() image.Point {
	var (
		xOffset int
		yOffset int
		imgW    int
		imgH    int
	)
	if c.Style == nil || c.Style.Font == nil {
		return image.ZP
	}
	if c.Image != nil && c.Image.Data != nil {
		imgSize := c.ImageSize()
		if c.Image.Align == LEFT || c.Image.Align == RIGHT {
			xOffset = imgSize.X
			imgH = imgSize.Y
		} else if c.Image.VAlign == TOP || c.Image.VAlign == BOTTOM {
			yOffset = imgSize.Y
			imgW = imgSize.X
		}
	}
	lines, maxWidth := c.Wrap(xOffset)
	if maxWidth < imgW {
		maxWidth = imgW
	}
	x := maxWidth + c.Style.BorderSize().X
	lineHeight := stringHeight(c.Style.Font.Size, c.Style.LineHeight)
	textHeight := len(lines) * lineHeight
	if textHeight < imgH {
		textHeight = imgH
	}
	y := textHeight + c.Style.BorderSize().Y
	x += xOffset
	y += yOffset
	return image.Pt(x, y)
}

// InnerBounds cell content bounds
func (c Cell) InnerBounds(bounds image.Rectangle) image.Rectangle {
	if c.Style == nil {
		return image.Rectangle{}
	}
	return c.Style.InnerBounds(bounds)
}

// GetImage download cell image
func (c Cell) GetImage(cache ImageCache) error {
	if c.Image == nil || (c.Image.Data != nil && c.Image.URL == "") {
		return nil
	}
	if cache != nil {
		if img, err := cache.Get(c.Image.URL); err == nil {
			c.Image.Data = img
			c.Image.UpdateSize()
			return nil
		}
	}
	if err := c.Image.Download(); err != nil {
		return err
	}
	if cache != nil {
		if err := cache.Set(c.Image.URL, c.Image.Data); err != nil {
			return err
		}
	}
	return nil
}

// Row in table
type Row struct {
	// Cells in row
	Cells []Cell `json:"cells,omitempty"`
	// Style for row
	Style *Style `json:"style,omitempty"`
}
