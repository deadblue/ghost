package rule

import (
	"strings"
	"unicode"
)

const (
	keywordAs = "as"
	keywordBy = "by"
)

type methodTokenizer struct {
	chars  []rune
	index  int
	length int
}

func (t *methodTokenizer) Init(text string) *methodTokenizer {
	t.chars = []rune(text)
	t.index, t.length = 0, len(t.chars)
	return t
}

func (t *methodTokenizer) Next() (word string, token _ParseToken) {
	for i := t.index + 1; i < t.length+1; i++ {
		if i == t.length || unicode.IsUpper(t.chars[i]) {
			slice := t.chars[t.index:i]
			t.index = i
			// Make word
			slice[0] = unicode.ToLower(slice[0])
			word = string(slice)
			// Determine token type
			switch word {
			case keywordAs:
				token = tokenKwAs
			case keywordBy:
				token = tokenKwBy
			default:
				token = tokenWord
			}
			return
		}
	}
	return "", tokenEOL
}

func (r *Rule) FromMethodName(name string) error {
	// Initialize
	r.IsStatic, r.Segments = true, make([]*Segment, 0)

	tk := (&methodTokenizer{}).Init(name)
	// The first word is used as request method
	if word, _ := tk.Next(); word != "" {
		r.Method = strings.ToUpper(word)
	}
	// Parse path
	if err := parse(tk, &r.Segments); err != nil {
		return err
	}

	buf := &strings.Builder{}
	buf.WriteRune('/')
	for i, s := range r.Segments {
		if s.IsVar {
			r.IsStatic = false
		}
		if i > 0 {
			buf.WriteRune('/')
		}
		buf.WriteString(s.String())
	}
	r.Path = buf.String()
	return nil
}
