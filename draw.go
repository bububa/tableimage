package tableimage

import (
	"image"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"golang.org/x/image/font"
)

func drawText(img *image.RGBA, point image.Point, label string, color string, ft *truetype.Font, fontSize float64) {
	fontCtx := newFreetypeContext()
	fontCtx.SetFont(ft)
	fontCtx.SetFontSize(fontSize)
	fontCtx.SetClip(img.Bounds())
	fontCtx.SetDst(img)
	fontCtx.SetSrc(image.NewUniform(ColorFromHex(color)))
	pt := freetype.Pt(point.X, point.Y+int(fontCtx.PointToFixed(fontSize)>>6))
	fontCtx.DrawString(label, pt)
}

func drawRect(img *image.RGBA, bounds image.Rectangle, borderColor string, bgColor string, strokeWidth float64) {
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetStrokeColor(ColorFromHex(borderColor))
	if bgColor != "" {
		gc.SetFillColor(ColorFromHex(bgColor))
	}
	gc.SetLineWidth(strokeWidth)
	draw2dkit.Rectangle(gc, float64(bounds.Min.X), float64(bounds.Min.Y), float64(bounds.Max.X), float64(bounds.Max.Y))
	gc.FillStroke()
	if bgColor != "" {
		gc.Fill()
	}
}

func drawImage(dst *image.RGBA, img *Image, pt image.Point) {
	if img == nil || img.Data == nil {
		return
	}
	scale := img.Scale()
	scaledImage := scaleImage(img.Data, scale)
	gc := draw2dimg.NewGraphicContext(dst)
	gc.Translate(float64(pt.X), float64(pt.Y))
	gc.DrawImage(scaledImage)
}

func scaleImage(img image.Image, scale float64) *image.RGBA {
	i := image.NewRGBA(img.Bounds())
	gc := draw2dimg.NewGraphicContext(i)
	gc.Scale(scale, scale)
	gc.DrawImage(img)
	return i
}

// MeasureString returns the rendered width and height of the specified text
// given the current font face.
func stringWidth(s string, fontFace font.Face) float64 {
	d := &font.Drawer{
		Face: fontFace,
	}
	a := d.MeasureString(s)
	return float64(a >> 6)
}

func stringHeight(fontSize float64, lineHeight float64) int {
	return int(fontSize * lineHeight)
}

func newFontFace(ft *truetype.Font, fontSize float64) font.Face {
	return truetype.NewFace(ft, &truetype.Options{
		Size: fontSize,
	})
}

func newFreetypeContext() *freetype.Context {
	fontCtx := freetype.NewContext()
	fontCtx.SetDPI(72)
	return fontCtx
}
