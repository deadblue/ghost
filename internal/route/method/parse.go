package method

import (
	"errors"
	"github.com/deadblue/ghost/internal/route/parser"
	"strings"
)

var (
	errParseFailed = errors.New("illegal syntax")
)

func Parse(text string, rule *parser.Rule) (err error) {
	// Initialize tokenizer
	tk := (&tokenizer{}).Init(text)
	// Method name
	rule.Method = strings.ToUpper(tk.Next())
	// Path buffer
	buf := strings.Builder{}
	if rule.Depth > 0 {
		for ok := rule.Pieces.GoFirst(); ok; ok = rule.Pieces.Forward() {
			_, piece := rule.Pieces.Get()
			buf.WriteRune('/')
			buf.WriteString(piece.String())
		}
	}
	// Split path
	key, piece := _TransitionKey{CurrState: stateInit}, &parser.PathPiece{}
	for key.CurrState != stateDone {
		word := strings.ToLower(tk.Next())
		key.Token = _Word(word).Token()
		result, ok := _TransitionTable[key]
		if !ok {
			return errParseFailed
		}
		key.CurrState = result.NextState
		if result.Finish {
			// Increase depth
			rule.Depth += 1
			// Append to path buffer
			buf.WriteRune('/')
			buf.WriteString(piece.String())
			// Append to rule
			rule.Pieces.Append(piece)
			if key.CurrState != stateDone {
				piece = &parser.PathPiece{}
			}
		}
		switch key.CurrState {
		case stateName:
			piece.Name = word
		case stateVar:
			piece.Name, piece.IsVar = word, true
			rule.IsStrict = false
		case stateExt:
			rule.Ext = word
		}
	}
	// Append extension to path buffer
	if rule.Ext != "" {
		buf.WriteRune('.')
		buf.WriteString(rule.Ext)
	}
	rule.Path = buf.String()
	return
}
