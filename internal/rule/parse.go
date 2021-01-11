package rule

import (
	"fmt"
	"regexp"
	"strings"
)

type _WordToken int
type _ParseState int

const (
	keywordAs = "as"
	keywordBy = "by"

	tokenGeneral _WordToken = iota
	tokenKwAs
	tokenKwBy
	tokenEos

	stateInit _ParseState = iota
	stateName
	stateKwAs
	stateKwBy
	stateVar
)

type _TransitionTuple struct {
	state _ParseState
	token _WordToken
}
type _TransitionFunc func() (_ParseState, bool)

var (
	_TransitionTable = map[_TransitionTuple]_TransitionFunc{
		{stateInit, tokenGeneral}: func() (_ParseState, bool) {
			return stateName, false
		},
		{stateInit, tokenKwAs}: func() (_ParseState, bool) {
			return stateKwAs, false
		},
		{stateInit, tokenKwBy}: func() (_ParseState, bool) {
			return stateKwBy, false
		},
		{stateName, tokenGeneral}: func() (_ParseState, bool) {
			return stateName, true
		},
		{stateName, tokenKwAs}: func() (_ParseState, bool) {
			return stateKwAs, false
		},
		{stateName, tokenKwBy}: func() (_ParseState, bool) {
			return stateKwBy, true
		},
		{stateName, tokenEos}: func() (_ParseState, bool) {
			return stateInit, true
		},
		{stateKwAs, tokenGeneral}: func() (_ParseState, bool) {
			return stateName, false
		},
		{stateKwAs, tokenEos}: func() (_ParseState, bool) {
			return stateInit, true
		},
		{stateKwBy, tokenGeneral}: func() (_ParseState, bool) {
			return stateVar, false
		},
		{stateVar, tokenGeneral}: func() (_ParseState, bool) {
			return stateName, true
		},
		{stateVar, tokenKwAs}: func() (_ParseState, bool) {
			return stateKwAs, false
		},
		{stateVar, tokenKwBy}: func() (_ParseState, bool) {
			return stateKwBy, true
		},
		{stateVar, tokenEos}: func() (_ParseState, bool) {
			return stateInit, true
		},
	}
)

type _Parser struct {
	words []string
}

func (p *_Parser) Parse() (rule *Rule, err error) {
	count := len(p.words)
	if count == 1 {
		return nil, nil
	}
	// Initial variables
	segments := make([]*_SegmentRule, 0)
	last, state := 0, stateInit
	// Scan words
	for i := 0; i <= count; i++ {
		// Select word and token
		word, token := "", tokenEos
		if i < count {
			word = p.words[i]
			switch word {
			case keywordAs:
				token = tokenKwAs
			case keywordBy:
				token = tokenKwBy
			default:
				token = tokenGeneral
			}
		}
		// Analysis word
		if tf, ok := _TransitionTable[_TransitionTuple{state, token}]; !ok {
			err = fmt.Errorf("unexpected word: %s", word)
			break
		} else {
			flushBuf := false
			if state, flushBuf = tf(); flushBuf {
				segments = append(segments, p.makeSegment(last, i))
				last = i
			}
		}
	}
	if err == nil {
		rule = &Rule{
			depth: len(segments),
			srs:   segments,
		}
	}
	return
}

func (p *_Parser) makeSegment(from, to int) (seg *_SegmentRule) {
	seg, buf := &_SegmentRule{}, &strings.Builder{}
	hasVar := false
	for i := from; i < to; i++ {
		switch word := p.words[i]; word {
		case keywordAs:
			if hasVar {
				buf.WriteString("\\.")
			} else {
				buf.WriteString(".")
			}
		case keywordBy:
			hasVar = true
			i += 1
			seg.varName = p.words[i]
			buf.WriteString("(\\w+)")
		default:
			buf.WriteString(word)
		}
	}
	// Always set value for debugging
	seg.value = buf.String()
	if hasVar {
		seg.pattern = regexp.MustCompile(seg.value)
	}
	return
}
