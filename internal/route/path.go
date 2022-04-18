package route

import (
	"github.com/deadblue/ghost/internal/route/parser"
)

func splitRequestPath(path string) (pieces []string, depth int, ext string) {
	pieces, depth = make([]string, 0), 1
	chars := []rune(path)
	start, length := 1, len(chars)
	dotIndex, slashIndex := -1, -1
	for i := 1; i < length; i++ {
		switch chars[i] {
		case '.':
			dotIndex = i
		case '/':
			slashIndex = i
			pieces = append(pieces, string(chars[start:i]))
			depth += 1
			start = i + 1
		}
	}
	if dotIndex > slashIndex {
		ext = string(chars[dotIndex+1:])
		pieces = append(pieces, string(chars[start:dotIndex]))
	} else {
		pieces = append(pieces, string(chars[start:]))
	}
	return
}

func setPathRoot(rule *parser.Rule, root string) {
	if root == "" || root == "/" {
		return
	}
	chars := []rune(root)
	start, length := 1, len(chars)
	// Trim the tailing slash
	for length >= 1 && chars[length-1] == '/' {
		length -= 1
	}
	// Split root path
	for i := 1; i < length; i++ {
		if chars[i] == '/' {
			if i > start {
				rule.Pieces.Append(&parser.PathPiece{
					IsVar: false,
					Name:  string(chars[start:i]),
				})
			}
			start = i + 1
		}
	}
	if length > start {
		rule.Pieces.Append(&parser.PathPiece{
			IsVar: false,
			Name:  string(chars[start:length]),
		})
	}
	rule.Depth = rule.Pieces.Len()
}
