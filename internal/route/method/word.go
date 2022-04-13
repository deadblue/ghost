package method

type _WordToken int

const (
	tokenWord _WordToken = iota
	tokenKwAs
	tokenKwBy
	tokenEOL

	// Keyword
	keywordAs = "as"
	keywordBy = "by"
)

type _Word string

func (w _Word) Token() _WordToken {
	switch w {
	case keywordBy:
		return tokenKwBy
	case keywordAs:
		return tokenKwAs
	case "":
		return tokenEOL
	default:
		return tokenWord
	}
}
