package tableimage

import (
	"regexp"
	"strconv"
	"unicode"

	"github.com/golang/freetype/truetype"
	"github.com/mattn/go-runewidth"
	"golang.org/x/image/font"
)

var (
	reText = regexp.MustCompile(`<text(\s+?(?P<attrs>.+?))?(\s+)?>(?P<txt>.+?)</text>`)
	reAttr = regexp.MustCompile(`(?P<attr>\w+)=["|'](?P<value>#?\w+)["|']`)
)

// Text string with width
type Text struct {
	Value   string
	Width   int
	Color   string
	BgColor string
	Padding int
	Pos     [2]int
}

// TextFromText create Text from Text
func TextFromText(value string, txt Text, fontFace font.Face) Text {
	txt.Value = value
	txt.Width = int(stringWidth(value, fontFace)) + txt.Padding*2
	return txt
}

// SameStyle check if two Text style is same
func (t Text) SameStyle(t2 Text) bool {
	return t.Color == t2.Color && t.BgColor == t2.BgColor && t.Padding == t2.Padding
}

// Word text array
type Word []Text

// Width calculate word length
func (w Word) Width() int {
	var l int
	for _, t := range w {
		l += t.Width
	}
	return l
}

func wrap(s string, w int, fontFace font.Face, ignoreInlineStyle bool) ([]Word, int) {
	segments := extractTexts(s, ignoreInlineStyle)
	words := separateWords(segments, fontFace)
	if w > 0 {
		words = wrapWords(words, w, fontFace)
	}
	var maxWidth int
	for _, w := range words {
		if w.Width() > maxWidth {
			maxWidth = w.Width()
		}
	}
	return words, maxWidth
}

// wrapWords wrap words to lines in w length
func wrapWords(words []Word, w int, fontFace font.Face) []Word {
	var (
		retWords []Word
		word     Word
	)
	for _, segs := range words {
		for _, txt := range segs {
			ww := int(stringWidth(txt.Value, fontFace))
			if txt.Value == "\n" || word.Width()+ww > w {
				retWords = append(retWords, word)
				word = Word{}
				if txt.Value != "\n" {
					word = append(word, txt)
				}
				continue
			}
			if len(word) > 0 && word[len(word)-1].SameStyle(txt) {
				lastWord := word[len(word)-1]
				word[len(word)-1] = TextFromText(lastWord.Value+txt.Value, word[len(word)-1], fontFace)
			} else {
				word = append(word, txt)
			}
		}
	}
	if len(word) > 0 {
		retWords = append(retWords, word)
	}
	return retWords
}

// separateWords seperate a string into words and not break word
func separateWords(segments []Text, fontFace font.Face) []Word {
	var (
		words     []Word
		wordTexts Word
	)
	for _, seg := range segments {
		var segWord []rune
		for _, r := range []rune(seg.Value) {
			l := len(segWord)
			cw := runewidth.RuneWidth(r)
			segWordStr := string(segWord)
			if cw == 0 {
				if l > 0 {
					wordTexts = append(wordTexts, TextFromText(segWordStr, seg, fontFace))
					words = append(words, wordTexts)
				}
				wordTexts = Word{}
				segWord = []rune{}
				continue
			} else if unicode.IsSpace(r) || unicode.IsPunct(r) || cw == 2 { // \n \t \s or double width add new word
				if l > 0 {
					wordTexts = append(wordTexts, TextFromText(segWordStr, seg, fontFace))
				}
				wordTexts = append(wordTexts, TextFromText(string(r), seg, fontFace))
				words = append(words, wordTexts)
				wordTexts = Word{}
				segWord = []rune{}
				continue
			}
			segWord = append(segWord, r)
		}
		if len(segWord) > 0 {
			segWordStr := string(segWord)
			wordTexts = append(wordTexts, TextFromText(segWordStr, seg, fontFace))
		}
	}
	if len(wordTexts) > 0 {
		words = append(words, wordTexts)
	}
	return words
}

func extractTexts(s string, ignoreInlineStyle bool) []Text {
	txtL := len(s)
	var segments []Text
	matches := reText.FindAllStringSubmatchIndex(s, -1)
	groupNames := reText.SubexpNames()
	totalMatches := len(matches)
	if totalMatches == 0 || ignoreInlineStyle {
		return []Text{{
			Value: s,
		}}
	}
	for idx, locs := range matches {
		if idx == 0 && locs[0] > 0 {
			segments = append(segments, Text{
				Pos:   [2]int{0, locs[0]},
				Value: s[0:locs[0]],
			})
		} else if len(segments) > 0 && locs[0] > segments[idx-1].Pos[1] {
			lastTxt := segments[idx-1]
			segments = append(segments, Text{
				Pos:   [2]int{lastTxt.Pos[1], locs[0]},
				Value: s[lastTxt.Pos[1]:locs[0]],
			})
		}
		text := Text{
			Pos: [2]int{locs[0], locs[1]},
		}
		for j, name := range groupNames {
			if j == 0 || name == "" {
				continue
			}
			content := s[locs[j*2]:locs[j*2+1]]
			if name == "attrs" {
				attrs := extractAttrs(content)
				for k, v := range attrs {
					switch k {
					case "color":
						text.Color = v
					case "bgcolor":
						text.BgColor = v
					case "padding":
						text.Padding, _ = strconv.Atoi(v)
					}
				}
			} else if name == "txt" {
				text.Value = content
			}
		}
		segments = append(segments, text)
		if idx == totalMatches-1 && locs[1] < txtL {
			segments = append(segments, Text{
				Pos:   [2]int{locs[1], txtL},
				Value: s[locs[1]:txtL],
			})
		}
	}
	return segments
}

// extractAttrs extract text tag attributes
func extractAttrs(s string) map[string]string {
	matches := reAttr.FindAllStringSubmatch(s, -1)
	groupNames := reAttr.SubexpNames()
	ret := make(map[string]string, len(matches))
	for _, m := range matches {
		var (
			k string
			v string
		)
		for j, name := range groupNames {
			if j == 0 || name == "" {
				continue
			}
			switch name {
			case "attr":
				k = m[j]
			case "value":
				v = m[j]
			}
		}
		if k != "" && v != "" {
			ret[k] = v
		}
	}
	return ret
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
