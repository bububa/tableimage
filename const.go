package tableimage

const (
	// DefaultLineHeight default row space
	DefaultLineHeight = 1.2
	// DefaultFontSize default font size
	DefaultFontSize = 13
	// DefaultPadding default seperator padding
	DefaultPadding = 10
	// DefaultWrapWords default wrap words count
	DefaultWrapWords = 20
	// DefalutColor default text color
	DefaultColor = "#212121"
	// DefaultBorderWidth default stroke line width
	DefaultBorderWidth = 1
)

// ImageType image type for writer
type ImageType int

const (
	_ ImageType = iota
	// JPEG jpeg image
	JPEG
	// PNG png image
	PNG
)

// Align Alignment
type Align int

const (
	// UnknownAlign unknown alignment
	UnknownAlign Align = iota
	// LEFT align left
	LEFT
	// RIGHT align right
	RIGHT
	// CENTER align center
	CENTER
)

// VAlign vertical alignment
type VAlign int

const (
	// UnknownVAlign unknown alignment
	UnknownVAlign VAlign = iota
	// TOP vertical align top
	TOP
	// BOTTOM vertical align bottom
	BOTTOM
	// MIDDLE vertical align bottom
	MIDDLE
)
