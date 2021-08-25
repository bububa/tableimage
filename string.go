package tableimage

import (
	"strings"
	"unicode"

	"github.com/golang/freetype/truetype"
	"github.com/mattn/go-runewidth"
	"golang.org/x/image/font"
)

// Text string with width
type Text struct {
	Value string
	Width int
}

func wrap(s string, w int, fontFace font.Face) ([]Text, int) {
	if w <= 0 {
		width := int(stringWidth(s, fontFace))
		return []Text{
			{
				Value: s,
				Width: width,
			},
		}, width
	}
	words := separateWords(s)
	return wrapWords(words, w, fontFace)
}

// wrapWords wrap words to lines in w length
func wrapWords(words []string, w int, fontFace font.Face) ([]Text, int) {
	var (
		width    int
		lines    []Text
		line     strings.Builder
		maxWidth float64
	)
	for _, word := range words {
		ww := int(stringWidth(word, fontFace))
		if word == "\n" {
			lineWidth := stringWidth(line.String(), fontFace)
			if lineWidth > maxWidth {
				maxWidth = lineWidth
			}
			lines = append(lines, Text{
				Value: line.String(),
				Width: int(lineWidth),
			})
			line.Reset()
			width = 0
			continue
		} else if width+ww > w {
			lineWidth := stringWidth(line.String(), fontFace)
			if lineWidth > maxWidth {
				maxWidth = lineWidth
			}
			lines = append(lines, Text{
				Value: line.String(),
				Width: int(lineWidth),
			})
			line.Reset()
			line.WriteString(word)
			width = 0
			continue
		}
		line.WriteString(word)
		width += ww
	}
	if line.Len() > 0 {
		lineWidth := stringWidth(line.String(), fontFace)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
		lines = append(lines, Text{
			Value: line.String(),
			Width: int(lineWidth),
		})
	}
	return lines, int(maxWidth)
}

// separateWords seperate a string into words and not break word
func separateWords(s string) []string {
	var (
		words []string
		word  []rune
	)
	for _, r := range []rune(s) {
		cw := runewidth.RuneWidth(r)
		l := len(word)
		if cw == 0 {
			if l > 0 {
				words = append(words, string(word))
				word = []rune{}
			}
			continue
		} else if unicode.IsSpace(r) { // \n \t \s add new word
			if l > 0 {
				words = append(words, string(word))
				word = []rune{}
			}
			words = append(words, string(r))
			continue
		} else if l > 0 && unicode.IsPunct(r) && !unicode.IsPunct(word[l-1]) { // if is punct and previous last char is not punct will add a new word
			words = append(words, string(word))
			word = []rune{}
		} else if l > 0 && !unicode.IsPunct(r) && unicode.IsPunct(word[l-1]) { // if is punct and previous last char is not punct will add a new word
			words = append(words, string(word))
			word = []rune{}
		} else if cw == 2 { // double width is a new word
			if l > 0 {
				words = append(words, string(word))
				word = []rune{}
			}
			words = append(words, string(r))
			continue
		}
		word = append(word, r)
	}
	if len(word) > 0 {
		words = append(words, string(word))
	}
	return words
}

// MeasureString returns the rendered width and height of the specified text
// given the current font face.
func stringWidth(s string, fontFace font.Face) float64 {
	if fontFace == nil {
		return 0
	}
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
	if ft == nil {
		return nil
	}
	return truetype.NewFace(ft, &truetype.Options{
		Size: fontSize,
	})
}
