package tableimage

import (
	"errors"
	"image"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
)

// DefaultStyle default style setting
var DefaultStyle = func() *Style {
	return &Style{
		Color:      DefaultColor,
		Border:     DefaultBorder(),
		LineHeight: DefaultLineHeight,
		Padding:    NewPadding(DefaultPadding),
		Align:      LEFT,
		VAlign:     MIDDLE,
		Font: &Font{
			Size: DefaultFontSize,
		},
	}
}

// DefaultCaptionStyle default caption style setting
var DefaultCaptionStyle = func() *Style {
	return &Style{
		Color:      DefaultColor,
		Border:     NoBorder(),
		LineHeight: DefaultLineHeight,
		Padding:    NewPaddingY(DefaultPadding),
		Align:      LEFT,
		VAlign:     TOP,
		Font: &Font{
			Size: DefaultFontSize,
		},
	}
}

// DefaultFooterStyle default table footer style setting
var DefaultFooterStyle = func() *Style {
	return &Style{
		Color:      DefaultColor,
		Border:     NoBorder(),
		LineHeight: DefaultLineHeight,
		Padding:    NewPaddingY(DefaultPadding),
		Align:      RIGHT,
		VAlign:     TOP,
		Font: &Font{
			Size: DefaultFontSize,
		},
	}
}

// DefaultLine default line setting
var DefaultLine = func() Line {
	return Line{
		Color: DefaultColor,
		Width: DefaultBorderWidth,
	}
}

// DefaultBorder default border setting
var DefaultBorder = func() *Border {
	return &Border{
		Top:    DefaultLine(),
		Right:  DefaultLine(),
		Bottom: DefaultLine(),
		Left:   DefaultLine(),
	}
}

// NoBorder no border setting
var NoBorder = func() *Border {
	return &Border{}
}

// Style for drawing
type Style struct {
	// Color text color
	Color string `json:"color,omitempty"`
	// Border cell border setting
	Border *Border `json:"border,omitempty"`
	// BgColor cell background color
	BgColor string `json:"bg_color,omitempty"`
	// Lineheight lineheight for paragraph
	LineHeight float64 `json:"line_height,omitempty"`
	// Margin cell margin
	Margin *Padding `json:"margin,omitempty"`
	// Padding cell padding
	Padding *Padding `json:"padding,omitempty"`
	// MaxWidth max width
	MaxWidth int `json:"max_width,omitempty"`
	// Align alignment
	Align Align `json:"align,omitempty"`
	// VAlign vertical alignment
	VAlign VAlign `json:"valign,omitempty"`
	// Font font setting
	Font *Font `json:"font,omitempty"`
}

// Inherit from other style
func (s *Style) Inherit(s1 *Style, cache draw2d.FontCache) error {
	if s1 == nil {
		return nil
	}
	if s.Color == "" {
		s.Color = s1.Color
	}
	if s.Border == nil {
		s.Border = s1.Border
	}
	if s.BgColor == "" {
		s.BgColor = s1.BgColor
	}
	if s.Margin != nil {
		s.Margin = s1.Margin
	}
	if s.Padding == nil {
		s.Padding = s1.Padding
	}
	if s.LineHeight < 1e-15 {
		s.LineHeight = s1.LineHeight
	}
	if s.MaxWidth == 0 {
		s.MaxWidth = s1.MaxWidth
	}
	if s.Align == UnknownAlign {
		s.Align = s1.Align
	}
	if s.VAlign == UnknownVAlign {
		s.VAlign = s1.VAlign
	}
	if s.Font == nil {
		s.Font = s1.Font
		if err := s.LoadFont(cache); err != nil {
			return err
		}
	} else {
		if err := s.LoadFont(cache); err != nil {
			return err
		}
		if s1.Font != nil {
			if s.Font.Size < 1e-15 {
				s.Font.Size = s1.Font.Size
			}
			if s.Font.Data == nil {
				s.Font.Data = s1.Font.Data
			}
			if s.Font.Font == nil {
				s.Font.Font = s1.Font.Font
			}
			if s.Font.DPI <= 0 {
				s.Font.DPI = s1.Font.DPI
			}
		}
	}
	return nil
}

// LoadFont load font with fontCache
func (s *Style) LoadFont(cache draw2d.FontCache) error {
	if s.Font == nil {
		return nil
	}
	return s.Font.Load(cache)
}

// BorderSize outer bound size
func (s Style) BorderSize() image.Point {
	var (
		x int
		y int
	)
	if s.Margin != nil {
		pt := s.Margin.Size()
		x += pt.X
		y += pt.Y
	}
	if s.Border != nil {
		pt := s.Border.Size()
		x += pt.X
		y += pt.Y
	}
	if s.Padding != nil {
		pt := s.Padding.Size()
		x += pt.X
		y += pt.Y
	}
	return image.Pt(x, y)
}

// InnerStart content start point
func (s Style) InnerStart() image.Point {
	var (
		x int
		y int
	)
	if s.Margin != nil {
		x += s.Margin.Left
	}
	if s.Border != nil {
		x += s.Border.Left.Width
	}
	if s.Padding != nil {
		x += s.Padding.Left
	}
	if s.Margin != nil {
		y += s.Margin.Top
	}
	if s.Border != nil {
		y += s.Border.Top.Width
	}
	if s.Padding != nil {
		y += s.Padding.Top
	}
	return image.Pt(x, y)
}

// InnerEnd content end point
func (s Style) InnerEnd() image.Point {
	var (
		x int
		y int
	)
	if s.Margin != nil {
		x += s.Margin.Right
	}
	if s.Border != nil {
		x += s.Border.Right.Width
	}
	if s.Padding != nil {
		x += s.Padding.Right
	}
	if s.Margin != nil {
		y += s.Margin.Bottom
	}
	if s.Border != nil {
		y += s.Border.Bottom.Width
	}
	if s.Padding != nil {
		y += s.Padding.Bottom
	}
	return image.Pt(x, y)
}

// InnerBounds content bounds
func (s Style) InnerBounds(bounds image.Rectangle) image.Rectangle {
	return image.Rectangle{
		Min: bounds.Min.Add(s.InnerStart()),
		Max: bounds.Max.Sub(s.InnerEnd()),
	}
}

// BorderPadding border padding
func (s Style) BorderPadding() Padding {
	padding := ZeroPadding
	if s.Margin != nil {
		padding.Add(*s.Margin)
	}
	if s.Border != nil {
		padding.Add(s.Border.Padding())
	}
	if s.Padding != nil {
		padding.Add(*s.Padding)
	}
	return padding
}

// Border border setting
type Border struct {
	Top    Line `json:"top,omitempty"`
	Right  Line `json:"right,omitempty"`
	Bottom Line `json:"bottom,omitempty"`
	Left   Line `json:"left,omitempty"`
}

// ChangeColor change border color
func (b *Border) ChangeColor(color string) {
	b.Top = b.Top.ChangeColor(color)
	b.Right = b.Right.ChangeColor(color)
	b.Bottom = b.Bottom.ChangeColor(color)
	b.Left = b.Left.ChangeColor(color)
}

// ChangeWidth change border width
func (b *Border) ChangeWidth(width int) {
	b.Top = b.Top.ChangeWidth(width)
	b.Right = b.Right.ChangeWidth(width)
	b.Bottom = b.Bottom.ChangeWidth(width)
	b.Left = b.Left.ChangeWidth(width)
}

// Padding border padding
func (b Border) Padding() Padding {
	return Padding{
		Top:    b.Top.Width,
		Right:  b.Right.Width,
		Bottom: b.Bottom.Width,
		Left:   b.Left.Width,
	}
}

// Size border border size
func (b Border) Size() image.Point {
	return image.Pt(b.Left.Width+b.Right.Width, b.Top.Width+b.Bottom.Width)
}

// Draw a border
func (b Border) Draw(img *image.RGBA, bounds image.Rectangle) {
	b.Top.Draw(img, image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y))
	b.Right.Draw(img, image.Rect(bounds.Max.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	b.Bottom.Draw(img, image.Rect(bounds.Min.X, bounds.Max.Y, bounds.Max.X, bounds.Max.Y))
	b.Left.Draw(img, image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X, bounds.Max.Y))
}

// Line border line
type Line struct {
	Color string `json:"color,omitempty"`
	Width int    `json:"width,omitempty"`
}

// ChangeColor return a line with new color
func (l Line) ChangeColor(color string) Line {
	l.Color = color
	return l
}

// ChangeWidth return a line with new width
func (l Line) ChangeWidth(width int) Line {
	l.Width = width
	return l
}

// Draw a new line
func (l Line) Draw(img *image.RGBA, bounds image.Rectangle) {
	if l.Width == 0 {
		return
	}
	drawRect(img, bounds, l.Color, l.Color, float64(l.Width))
}

// ZeroPadding zero padding object
var ZeroPadding = Padding{}

// Padding cell padding
type Padding struct {
	// Top padding
	Top int `json:"top,omitempty"`
	// Right padding
	Right int `json:"right,omitempty"`
	// Bottom padding
	Bottom int `json:"bottom,omitempty"`
	// Left padding
	Left int `json:"left,omitempty"`
}

// NewPadding init padding with same value
func NewPadding(padding int) *Padding {
	return &Padding{
		Top:    padding,
		Right:  padding,
		Bottom: padding,
		Left:   padding,
	}
}

// NewPaddingY init padding with top/bottom
func NewPaddingY(y int) *Padding {
	return &Padding{
		Top:    y,
		Bottom: y,
	}
}

// NewPaddingX init padding with left/right
func NewPaddingX(x int) *Padding {
	return &Padding{
		Right: x,
		Left:  x,
	}
}

// NewPaddingXY init padding for both x / y
func NewPaddingXY(x int, y int) *Padding {
	return &Padding{
		Top:    y,
		Right:  x,
		Bottom: y,
		Left:   x,
	}
}

// Add merget tow paddings
func (p Padding) Add(p2 Padding) Padding {
	p.Left += p2.Left
	p.Top += p2.Top
	p.Right += p2.Right
	p.Bottom += p2.Bottom
	return p
}

// Size padding size
func (p Padding) Size() image.Point {
	return image.Pt(p.Left+p.Right, p.Top+p.Bottom)
}

// Font font info
type Font struct {
	// Size font size
	Size float64 `json:"size,omitempty"`
	// Data font setting
	Data *draw2d.FontData `json:"data,omitempty"`
	// Font
	Font *truetype.Font `json:"-"`
	// DPI
	DPI int `json:"dpi,omitempty"`
}

// Load font from font cache
func (f *Font) Load(cache draw2d.FontCache) error {
	if f.Font != nil {
		return nil
	}
	if f.Data == nil {
		return nil
	}
	if cache == nil {
		return errors.New("missing font cache")
	}
	ft, err := cache.Load(*f.Data)
	if err != nil {
		return err
	}
	f.Font = ft
	return nil
}
