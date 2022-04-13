package method

type _ParseState int

const (
	stateInit _ParseState = iota
	stateName
	stateBy
	stateVar
	stateAs
	stateExt
	stateDone
)

type _TransitionKey struct {
	CurrState _ParseState
	Token     _WordToken
}

type _TransitionResult struct {
	NextState _ParseState
	Finish    bool
}

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
	_TransitionTable = map[_TransitionKey]_TransitionResult{
		{stateInit, tokenWord}: {stateName, false},
		{stateInit, tokenKwAs}: {stateAs, false},
		{stateInit, tokenKwBy}: {stateBy, false},
		{stateInit, tokenEOL}:  {stateDone, true},

		{stateName, tokenWord}: {stateName, true},
		{stateName, tokenKwAs}: {stateAs, false},
		{stateName, tokenKwBy}: {stateBy, true},
		{stateName, tokenEOL}:  {stateDone, true},

		{stateBy, tokenWord}: {stateVar, false},

		{stateVar, tokenWord}: {stateName, true},
		{stateVar, tokenKwAs}: {stateAs, false},
		{stateVar, tokenKwBy}: {stateBy, true},
		{stateVar, tokenEOL}:  {stateDone, true},

		{stateAs, tokenWord}: {stateExt, false},

		{stateExt, tokenEOL}: {stateDone, true},
	}
)
