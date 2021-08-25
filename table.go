package tableimage

import (
	"image"
	"math"
	"net/http"
)

// Cell in table
type Cell struct {
	// Text content of a cell
	Text string `json:"text,omitempty"`
	// Image image for a cell
	Image *Image `json:"image,omitempty"`
	// Style for cell
	Style *Style `json:"style,omitempty"`
}

// Draw render cell to image
func (c Cell) Draw(img *image.RGBA, bounds image.Rectangle) {
	if c.Style == nil {
		return
	}
	if c.Style.BgColor != "" {
		drawRect(img, bounds, "", c.Style.BgColor, 0)
	}
	if c.Style.Border != nil {
		c.Style.Border.Draw(img, bounds)
	}
	if c.Style.Font != nil {
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
		wrapedTexts, _ := c.Wrap(imgXOffset)
		innerBounds := c.InnerBounds(bounds)
		var (
			textStartX int
			textHeight = len(wrapedTexts) * lineHeight
			y          int
			imgX       int
			imgY       int
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
		if c.Image != nil && c.Image.Data != nil {
			imgX = c.Image.PaddingLeft()
			imgY = c.Image.PaddingTop()
			if c.Image.VAlign == TOP || c.Image.VAlign == BOTTOM {
				if c.Image.VAlign == TOP {
					imgY += y
					y += imgSize.Y
				} else {
					imgY += y + textHeight
				}
				switch c.Style.Align {
				case RIGHT:
					imgX += innerBounds.Max.X - imgSize.X
				case CENTER:
					center := (innerBounds.Dx() - imgSize.X) / 2
					imgX += innerBounds.Min.X + center
				default:
					imgX += innerBounds.Min.X
				}
			} else if c.Image.Align == LEFT || c.Image.Align == RIGHT {
				if c.Image.Align == LEFT {
					imgX += innerBounds.Min.X
					textStartX = imgSize.X
				} else {
					imgX += innerBounds.Max.X - imgSize.X
				}
				switch c.Style.VAlign {
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
			}
			pt := image.Pt(imgX, imgY)
			drawImage(img, c.Image, pt)
		}
		for lineIdx, wrapedText := range wrapedTexts {
			y += lineIdx * lineHeight
			var x int
			switch c.Style.Align {
			case RIGHT:
				if textStartX > 0 {
					x = innerBounds.Max.X - wrapedText.Width
				} else {
					x = innerBounds.Max.X - wrapedText.Width - imgXOffset
				}
			case CENTER:
				center := (innerBounds.Dx() - wrapedText.Width - imgXOffset) / 2
				x = innerBounds.Min.X + center
			default:
				x = innerBounds.Min.X + textStartX
			}
			pt := image.Pt(x, y)
			drawText(img, pt, wrapedText.Value, c.Style.Color, c.Style.Font)
		}
	}
}

// Wrap wraps cell content returns paragraphs, and max content width
func (c Cell) Wrap(xOffset int) ([]Text, int) {
	if c.Style == nil || c.Style.Font == nil {
		return []Text{{
			Value: c.Text,
			Width: 0,
		}}, 0
	}
	maxWidth := c.Style.MaxWidth - xOffset
	fontFace := newFontFace(c.Style.Font.Font, c.Style.Font.Size)
	return wrap(c.Text, maxWidth, fontFace)
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
	wrapedTexts, maxWidth := c.Wrap(xOffset)
	if maxWidth < imgW {
		maxWidth = imgW
	}
	x := maxWidth + c.Style.BorderSize().X
	lineHeight := stringHeight(c.Style.Font.Size, c.Style.LineHeight)
	textHeight := len(wrapedTexts) * lineHeight
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

// Table table struct
type Table struct {
	colsWidth   []int
	rowsHeight  []int
	rows        []Row
	caption     *Cell
	footer      *Cell
	captionSize image.Point
	footerSize  image.Point
}

// NewTable create Table instance
func NewTable(ti *TableImage, rows []Row, caption *Cell, footer *Cell) (*Table, error) {
	var maxCols int
	for _, row := range rows {
		cols := len(row.Cells)
		if cols > maxCols {
			maxCols = cols
		}
	}
	cols := make([]int, maxCols)
	heights := make([]int, len(rows))
	updatedRows := make([]Row, 0, len(rows))
	for rowIdx, row := range rows {
		if row.Style == nil {
			row.Style = ti.style
		} else {
			row.Style.Inherit(ti.style, ti.fontCache)
		}
		rowCells := make([]Cell, 0, len(row.Cells))
		for cellIdx, cell := range row.Cells {
			if cell.Style == nil {
				cell.Style = row.Style
			} else {
				cell.Style.Inherit(row.Style, ti.fontCache)
			}
			cell.GetImage(ti.imageCache)
			cellSize := cell.Size()
			if cellSize.X > cols[cellIdx] {
				cols[cellIdx] = cellSize.X
			}
			if cellSize.Y > heights[rowIdx] {
				heights[rowIdx] = cellSize.Y
			}
			rowCells = append(rowCells, cell)
		}
		row.Cells = rowCells
		updatedRows = append(updatedRows, row)
	}
	table := &Table{
		caption:    caption,
		footer:     footer,
		rows:       updatedRows,
		rowsHeight: heights,
		colsWidth:  cols,
	}
	if caption != nil {
		if caption.Style == nil {
			caption.Style = DefaultCaptionStyle()
			caption.Style.LoadFont(ti.fontCache)
		} else {
			caption.Style.Inherit(DefaultCaptionStyle(), ti.fontCache)
		}
		if caption.Style.MaxWidth == 0 || caption.Style.MaxWidth > table.Size().X {
			caption.Style.MaxWidth = table.Size().X - caption.Style.BorderPadding().Size().X
		}
		table.captionSize = caption.Size()
	}
	if footer != nil {
		if footer.Style == nil {
			footer.Style = DefaultCaptionStyle()
			footer.Style.LoadFont(ti.fontCache)
		} else {
			footer.Style.Inherit(DefaultCaptionStyle(), ti.fontCache)
		}
		if footer.Style.MaxWidth == 0 || footer.Style.MaxWidth > table.Size().X {
			footer.Style.MaxWidth = table.Size().X - footer.Style.BorderPadding().Size().X
		}
		table.footerSize = footer.Size()
	}
	return table, nil
}

// Size get bound size
func (r Table) Size() image.Point {
	rowsSize := r.RowsSize()
	if rowsSize.X < r.captionSize.X {
		rowsSize.X = r.captionSize.X
	}
	if rowsSize.X < r.footerSize.X {
		rowsSize.X = r.footerSize.X
	}
	rowsSize.Y += r.captionSize.Y + r.footerSize.Y
	return rowsSize
}

// RowsSize get rows size
func (r Table) RowsSize() image.Point {
	var (
		width  int
		height int
	)
	for _, w := range r.colsWidth {
		width += w
	}
	for _, h := range r.rowsHeight {
		height += h
	}
	return image.Pt(width, height)
}

// RowsStartPoint table rows start point
func (r Table) RowsStartPoint() image.Point {
	return r.captionSize
}

// DrawCaption draw table caption
func (r Table) DrawCaption(img *image.RGBA, pt image.Point) {
	if r.caption == nil {
		return
	}
	bounds := image.Rect(pt.X, pt.Y, pt.X+r.captionSize.X, pt.Y+r.captionSize.Y)
	r.caption.Draw(img, bounds)
}

// DrawFooter draw table footer
func (r Table) DrawFooter(img *image.RGBA, pt image.Point) {
	if r.footer == nil {
		return
	}
	bounds := image.Rect(pt.X, pt.Y, pt.X+r.footerSize.X, pt.Y+r.footerSize.Y)
	r.footer.Draw(img, bounds)
}

// CellBounds get a cell bounds
func (r Table) CellBounds(rowIdx int, cellIdx int) image.Rectangle {
	var (
		x int
		y int
		w = r.colsWidth[cellIdx]
		h = r.rowsHeight[rowIdx]
	)
	for i, v := range r.rowsHeight {
		if i > rowIdx-1 {
			break
		}
		y += v
	}
	for i, v := range r.colsWidth {
		if i > cellIdx-1 {
			break
		}
		x += v
	}
	return image.Rect(x, y, x+w, y+h)
}

// Rows get rows
func (r Table) Rows() []Row {
	return r.rows
}

// Image image setting
type Image struct {
	// URL image link
	URL string `json:"url,omitempty"`
	// Data image data
	Data image.Image
	// Inline display inline
	Inline bool `json:"inline,omitempty"`
	// Size image width/height
	Size image.Point `json:"size,omitempty"`
	// Align image text alignment
	Align Align `json:"align,omitempty"`
	// VAlign image text vertical alignment
	VAlign VAlign `json:"valign,omitempty"`
	// Padding image padding
	Padding *Padding `json:"padding,omitempty"`
}

// PaddingX horizontal padding
func (i Image) PaddingX() int {
	if i.Padding == nil {
		return 0
	}
	return i.Padding.Left + i.Padding.Right
}

// PaddingY vertical padding
func (i Image) PaddingY() int {
	if i.Padding == nil {
		return 0
	}
	return i.Padding.Top + i.Padding.Bottom
}

// PaddingLeft left padding
func (i Image) PaddingLeft() int {
	if i.Padding == nil {
		return 0
	}
	return i.Padding.Left
}

// PaddingTop top padding
func (i Image) PaddingTop() int {
	if i.Padding == nil {
		return 0
	}
	return i.Padding.Top
}

// Download image data
func (i *Image) Download() error {
	if i.Data != nil || i.URL == "" {
		return nil
	}
	resp, err := http.DefaultClient.Get(i.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return err
	}
	i.Data = img
	i.UpdateSize()
	return nil
}

// UpdateSize update Size based on image bounds
func (i *Image) UpdateSize() {
	bounds := i.Data.Bounds()
	if i.Size.X == 0 {
		i.Size.X = bounds.Dx()
	}
	if i.Size.Y == 0 {
		i.Size.Y = bounds.Dy()
	}
	scale := i.Scale()
	i.Size.X = int(math.Round(float64(bounds.Dx()) * scale))
	i.Size.Y = int(math.Round(float64(bounds.Dy()) * scale))
}

// BoundSize get Image width/height
func (i Image) BoundSize() image.Point {
	return image.Pt(i.Size.X+i.PaddingX(), i.Size.Y+i.PaddingY())
}

// Scale get image scale
func (i Image) Scale() float64 {
	bounds := i.Data.Bounds()
	return math.Min(float64(i.Size.X)/float64(bounds.Dx()), float64(i.Size.Y)/float64(bounds.Dy()))
}
