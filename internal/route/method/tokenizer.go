package method

import "unicode"

type tokenizer struct {
	chars  []rune
	index  int
	length int
}

func (t *tokenizer) Init(text string) *tokenizer {
	t.chars = []rune(text)
	t.index, t.length = 0, len(t.chars)
	return t
}

func (t *tokenizer) Next() string {
	for i := t.index + 1; i < t.length+1; i++ {
		if i == t.length || unicode.IsUpper(t.chars[i]) {
			slice := t.chars[t.index:i]
			t.index = i
			return string(slice)
		}
	}
	return ""
}
