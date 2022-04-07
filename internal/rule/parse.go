package rule

import (
	"fmt"
)

type _ParseState int
type _ParseToken int

type _TransitionKey struct {
	CurrState _ParseState
	Token     _ParseToken
}

type _TransitionResult struct {
	NextState  _ParseState
	ShallFlush bool
}

type tokenizer interface {
	Next() (string, _ParseToken)
}

const (
	tokenWord _ParseToken = iota
	tokenKwAs
	tokenKwBy
	tokenEOL

	stateInit _ParseState = iota
	stateName
	stateAs
	stateExt
	stateBy
	stateVar
	stateDone
)

var (
	// Transition table:
	//   | S\T  | Word | KwAs | KwBy | EOL  |
	//   |------|------|------|------|------|
	//   | Init | Name | KwAs | KwBy | Done |
	//   | Name | Name | KwAs | KwBy | Done |
	//   | As   | Ext  | -    | -    | -    |
	//   | Ext  | -    | -    | -    | Done |
	//   | By   | Var  | -    | -    | -    |
	//   | Var  | Name | KwAs | KwBy | Done |
	//   | Done | -    | -    | -    | -    |
	//   |------|------|------|------|------|
	transitionTable = map[_TransitionKey]_TransitionResult{
		{stateInit, tokenWord}: {stateName, false},
		{stateInit, tokenKwAs}: {stateAs, false},
		{stateInit, tokenKwBy}: {stateBy, false},
		{stateInit, tokenEOL}:  {stateDone, true},
		{stateName, tokenWord}: {stateName, true},
		{stateName, tokenKwAs}: {stateAs, false},
		{stateName, tokenKwBy}: {stateBy, true},
		{stateName, tokenEOL}:  {stateDone, true},
		{stateAs, tokenWord}:   {stateExt, false},
		{stateExt, tokenEOL}:   {stateDone, true},
		{stateBy, tokenWord}:   {stateVar, false},
		{stateVar, tokenWord}:  {stateName, true},
		{stateVar, tokenKwAs}:  {stateAs, false},
		{stateVar, tokenKwBy}:  {stateBy, true},
		{stateVar, tokenEOL}:   {stateDone, true},
	}
)

func parse(tk tokenizer, segments *[]*Segment) error {
	// Split request path
	ss, s := make([]*Segment, 0), &Segment{}
	key := _TransitionKey{CurrState: stateInit}
	for key.CurrState != stateDone {
		// Read next word
		var word string
		word, key.Token = tk.Next()
		// Get transition result
		if res, ok := transitionTable[key]; !ok {
			return fmt.Errorf("unexpected word: %s", word)
		} else {
			if res.ShallFlush && s.IsValid() {
				ss = append(ss, s)
				s = &Segment{}
			}
			// Update state
			key.CurrState = res.NextState
			// Update segment
			switch key.CurrState {
			case stateName, stateVar:
				s.IsVar = key.CurrState == stateVar
				s.Name = word
			case stateExt:
				s.Ext = word
			}
		}
	}
	*segments = ss
	return nil
}
