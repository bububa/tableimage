package tableimage

import (
	"image"

	"github.com/llgcode/draw2d"
)

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
	rows, cols, heights := initRows(ti, rows)
	table := &Table{
		caption:    caption,
		footer:     footer,
		rows:       rows,
		rowsHeight: heights,
		colsWidth:  cols,
	}
	table.initCaption(ti.fontCache)
	table.initFooter(ti.fontCache)
	return table, nil
}

func initRows(ti *TableImage, rows []Row) ([]Row, []int, []int) {
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
	return updatedRows, cols, heights
}

func (r *Table) initCaption(cache draw2d.FontCache) {
	if r.caption == nil {
		return
	}
	if r.caption.Style == nil {
		r.caption.Style = DefaultCaptionStyle()
		r.caption.Style.LoadFont(cache)
	} else {
		r.caption.Style.Inherit(DefaultCaptionStyle(), cache)
	}
	if r.caption.Style.MaxWidth == 0 || r.caption.Style.MaxWidth > r.Size().X {
		r.caption.Style.MaxWidth = r.Size().X - r.caption.Style.BorderPadding().Size().X
	}
	r.captionSize = r.caption.Size()
}

func (r *Table) initFooter(cache draw2d.FontCache) {
	if r.footer == nil {
		return
	}
	if r.footer.Style == nil {
		r.footer.Style = DefaultFooterStyle()
		r.footer.Style.LoadFont(cache)
	} else {
		r.footer.Style.Inherit(DefaultFooterStyle(), cache)
	}
	if r.footer.Style.MaxWidth == 0 || r.footer.Style.MaxWidth > r.Size().X {
		r.footer.Style.MaxWidth = r.Size().X - r.footer.Style.BorderPadding().Size().X
	}
	r.footerSize = r.footer.Size()
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
	bounds := image.Rect(pt.X, pt.Y, pt.X+r.Size().X, pt.Y+r.captionSize.Y)
	r.caption.Draw(img, bounds)
}

// DrawFooter draw table footer
func (r Table) DrawFooter(img *image.RGBA, pt image.Point) {
	if r.footer == nil {
		return
	}
	bounds := image.Rect(pt.X, pt.Y, pt.X+r.Size().X, pt.Y+r.footerSize.Y)
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
