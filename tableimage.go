package tableimage

import (
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/llgcode/draw2d"
)

// TableImage core struct
type TableImage struct {
	fontFolder string
	fontCache  draw2d.FontCache
	imageCache ImageCache
	style      *Style
}

// New init a TableImage object
func New(options ...Option) (*TableImage, error) {
	ti := &TableImage{
		style:      DefaultStyle(),
		imageCache: make(DefaultImageCache),
	}
	for _, opt := range options {
		opt.apply(ti)
	}
	if ti.fontFolder != "" {
		ti.fontCache = draw2d.NewSyncFolderFontCache(ti.fontFolder)
		if ti.style != nil {
			if err := ti.style.LoadFont(ti.fontCache); err != nil {
				return nil, err
			}
		}
	}
	return ti, nil
}

// Draw draw table image
func (ti *TableImage) Draw(rows []Row, caption *Cell, footer *Cell) (*image.RGBA, error) {
	table, err := NewTable(ti, rows, caption, footer)
	if err != nil {
		return nil, err
	}
	bounds := ti.Size(table)
	img := image.NewRGBA(image.Rect(0, 0, bounds.X, bounds.Y))
	if ti.style != nil && ti.style.BgColor != "" {
		draw.Draw(img, img.Bounds(), &image.Uniform{ColorFromHex(ti.style.BgColor)}, image.ZP, draw.Src)
	}
	ti.draw(img, table)
	return img, nil
}

// Write witer image to io Writer
func Write(w io.Writer, img *image.RGBA, imageType ImageType) error {
	switch imageType {
	case JPEG:
		return jpeg.Encode(w, img, nil)
	case PNG:
		return png.Encode(w, img)
	}
	return errors.New("unknown image type")
}

// Save an image to file
func Save(filepath string, img *image.RGBA, imageType ImageType) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	return Write(f, img, imageType)
}

// Size get tableimage width/height
func (ti *TableImage) Size(table *Table) image.Point {
	rowsBounds := table.Size()
	if ti.style == nil {
		return rowsBounds
	}
	return rowsBounds.Add(ti.style.BorderSize())
}

// BorderSize get border width of tableimage
func (ti *TableImage) BorderSize() image.Point {
	border := image.ZP
	if ti.style != nil {
		border = ti.style.BorderSize()
	}
	return border
}

func (ti *TableImage) innerStartPoint() image.Point {
	if ti.style == nil {
		return image.ZP
	}
	return ti.style.InnerStart()
}

func (ti *TableImage) draw(img *image.RGBA, table *Table) {
	startPoint := ti.innerStartPoint()
	table.DrawCaption(img, startPoint)
	rowsStartPoint := table.RowsStartPoint()
	rowsPt := image.Pt(startPoint.X, startPoint.Y+rowsStartPoint.Y)
	for rowIdx, row := range table.Rows() {
		for cellIdx, cell := range row.Cells {
			bounds := table.CellBounds(rowIdx, cellIdx)
			bounds = bounds.Add(rowsPt)
			cell.Draw(img, bounds)
		}
	}
	rowsSize := table.RowsSize()
	footerPt := image.Pt(startPoint.X, rowsPt.Y+rowsSize.Y)
	table.DrawFooter(img, footerPt)
}

// CacheImage cache image
func (ti *TableImage) CacheImage(k string, img image.Image) error {
	if ti.imageCache == nil {
		return errors.New("no cache")
	}
	return ti.imageCache.Set(k, img)
}

// GetImage get image from cache
func (ti *TableImage) GetImage(k string) (image.Image, error) {
	if ti.imageCache == nil {
		return nil, errors.New("no cache")
	}
	return ti.imageCache.Get(k)
}

// ImageCache image cache interface
type ImageCache interface {
	// Get get image from cache
	Get(k string) (image.Image, error)
	// Set set image to cache
	Set(k string, img image.Image) error
}

// DefaultImageCache default ImageCache implement
type DefaultImageCache map[string]image.Image

// Get implement ImageCache
func (c DefaultImageCache) Get(k string) (image.Image, error) {
	if img, found := c[k]; found {
		return img, nil
	}
	return nil, errors.New("missing cache")
}

// Set implement ImageCache
func (c DefaultImageCache) Set(k string, img image.Image) error {
	c[k] = img
	return nil
}
