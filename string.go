package tableimage

import (
	"strings"
	"unicode"

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
	var (
		width int
		lines []Text
		line  strings.Builder
		words []string
		word  []rune
	)

	// seperate words
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
	var maxWidth float64
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
