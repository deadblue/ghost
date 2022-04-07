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
	r.IsStatic, r.SegHead = true, &Segment{}

	tk := (&methodTokenizer{}).Init(name)
	// The first word is used as request method
	if word, _ := tk.Next(); word != "" {
		r.Method = strings.ToUpper(word)
	}
	// Parse path
	if err := parse(tk, r.SegHead); err != nil {
		return err
	}

	// TODO: Move following login in `parse` function
	buf := &strings.Builder{}
	for seg := r.SegHead; seg != nil; seg = seg.Next {
		r.Depth += 1
		if seg.IsVar {
			r.IsStatic = false
		}
		buf.WriteRune('/')
		buf.WriteString(seg.String())
	}
	r.Path = buf.String()
	return nil
}
