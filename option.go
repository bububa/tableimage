package tableimage

import (
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
)

// Option tableimage option interface
type Option interface {
	apply(*TableImage)
}

type optionFunc func(*TableImage)

func (fn optionFunc) apply(ti *TableImage) {
	fn(ti)
}

// WithFontSize set font size
func WithFontSize(size float64) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		if ti.style.Font == nil {
			ti.style.Font = &Font{}
		}
		ti.style.Font.Size = size
	})
}

// WithLineHeight set line height
func WithLineHeight(height float64) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		ti.style.LineHeight = height
	})
}

// WithPadding set padding
func WithPadding(padding *Padding) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		ti.style.Padding = padding
	})
}

// WithMargin set margin
func WithMargin(margin *Padding) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		ti.style.Margin = margin
	})
}

// WithColor set text color
func WithColor(color string) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		ti.style.Color = color
	})
}

// WithBorderColor set border color
func WithBorderColor(color string) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		if ti.style.Border == nil {
			ti.style.Border = DefaultBorder()
		}
		ti.style.Border.ChangeColor(color)
	})
}

// WithBorderWidth set border width
func WithBorderWidth(width int) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		if ti.style.Border == nil {
			ti.style.Border = DefaultBorder()
		}
		ti.style.Border.ChangeWidth(width)
	})
}

// WithBorder set border setting
func WithBorder(border *Border) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		ti.style.Border = border
	})
}

// WithBgColor set background color
func WithBgColor(bgColor string) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		ti.style.BgColor = bgColor
	})
}

// WithAlign set alignment
func WithAlign(align Align) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		ti.style.Align = align
	})
}

// WithVAlign set vertical alignment
func WithVAlign(align VAlign) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		ti.style.VAlign = align
	})
}

// WithFontFolder set font folder
func WithFontFolder(fontFolder string) Option {
	return optionFunc(func(ti *TableImage) {
		ti.fontFolder = fontFolder
	})
}

// WithFontData set font data
func WithFontData(font *draw2d.FontData) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		if ti.style.Font == nil {
			ti.style.Font = &Font{}
		}
		ti.style.Font.Data = font
	})
}

// WithFont set font setting
func WithFont(font *truetype.Font) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		if ti.style.Font == nil {
			ti.style.Font = &Font{}
		}
		ti.style.Font.Font = font
	})
}

// WidthDPI set font dpi
func WithDPI(dpi int) Option {
	return optionFunc(func(ti *TableImage) {
		if ti.style == nil {
			ti.style = &Style{}
		}
		if ti.style.Font == nil {
			ti.style.Font = &Font{}
		}
		ti.style.Font.DPI = dpi
	})
}

// WithFontCache set font cache
func WithFontCache(cache draw2d.FontCache) Option {
	return optionFunc(func(ti *TableImage) {
		ti.fontCache = cache
	})
}

// WithStyle set style setting
func WithStyle(style *Style) Option {
	return optionFunc(func(ti *TableImage) {
		ti.style = style
	})
}
