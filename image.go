package tableimage

import (
	"image"
	"math"
	"net/http"
)

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
