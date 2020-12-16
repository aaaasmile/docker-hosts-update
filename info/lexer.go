package info

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	EOFRune = -1 // serve al parser per riconoscere la fine del source
)

/// RUN Stack used in Rewind
type runeNode struct {
	r    rune
	next *runeNode
}

type runeStack struct {
	start *runeNode
}

func newRuneStack() runeStack {
	return runeStack{}
}

func (s *runeStack) push(r rune) {
	node := &runeNode{r: r}
	if s.start == nil {
		s.start = node
	} else {
		node.next = s.start
		s.start = node
	}
}

func (s *runeStack) pop() rune {
	if s.start == nil {
		return EOFRune
	} else {
		n := s.start
		s.start = n.next
		return n.r
	}
}

func (s *runeStack) clear() {
	s.start = nil
}

// ********** TokenType ******************

type TokenType int

type Token struct {
	Type  TokenType
	Value string
}

func (tk *Token) String() string {
	switch tk.Type {
	case itemEOF:
		return "EOF"
	case itemError:
		return tk.Value
	}
	if len(tk.Value) > 30 {
		return fmt.Sprintf("%.10q...", tk.Value)
	}
	return fmt.Sprintf("%q", tk.Value)
}

// ********** Lexer ******************
type StateFunc func(*L) StateFunc

type PropInfo struct {
	Keyword   string
	TokenType TokenType
}

type L struct {
	source          string
	start, position int
	state           StateFunc
	afterOpenState  StateFunc  // custom igsa
	stackPropInfo   []PropInfo // custom igsa
	tokens          chan Token
	runstack        runeStack
	onlyHeader      bool
}

// NewL creates a returns a lexer ready to parse the given source code.
func NewL(src string, start StateFunc) *L {
	l := L{
		source:   src,
		state:    start,
		start:    0,
		position: 0,
		runstack: newRuneStack(),
	}
	buffSize := len(l.source) / 2
	if buffSize <= 0 {
		buffSize = 1
	}
	l.tokens = make(chan Token, buffSize)

	return &l
}

func (l *L) current() string {
	return l.source[l.start:l.position]
}

func (l *L) emit(t TokenType) {
	tok := Token{
		Type:  t,
		Value: l.current(),
	}
	l.tokens <- tok
	l.start = l.position
	l.runstack.clear()
}

func (l *L) ignore() {
	l.runstack.clear()
	l.start = l.position
}

func (l *L) peek() rune {
	r := l.next()
	l.rewind()
	return r
}

func (l *L) rewind() {
	r := l.runstack.pop()
	if r > EOFRune {
		size := utf8.RuneLen(r)
		l.position -= size
		if l.position < l.start {
			l.position = l.start
		}
	}
}

func (l *L) next() rune {
	var (
		r rune
		s int
	)
	str := l.source[l.position:]
	if len(str) == 0 {
		r, s = EOFRune, 0
	} else {
		r, s = utf8.DecodeRuneInString(str)
	}
	l.position += s
	l.runstack.push(r)

	return r
}

func (l *L) take(chars string) {
	r := l.next()
	for strings.ContainsRune(chars, r) {
		r = l.next()
	}
	l.rewind() // last next wasn't a match
}

func (l *L) nextItem() Token {
	for {
		select {
		case item := <-l.tokens:
			return item
		default:
			if l.state != nil {
				l.state = l.state(l)
			} else {
				return Token{Type: itemEOF, Value: ""}
			}
		}
	}
}

func (l *L) errorf(format string, args ...interface{}) StateFunc {
	l.tokens <- Token{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *L) popPropInfo() (*PropInfo, error) {
	n := len(l.stackPropInfo) - 1
	if n < 0 {
		return nil, fmt.Errorf("popPropInfo stack is empty")
	}
	keyInfoState := l.stackPropInfo[n]
	l.stackPropInfo = l.stackPropInfo[:n] // pop
	return &keyInfoState, nil
}

////////////////////////////////////////////////////////////////////////
const (
	CommentKey = "#"
)

// for String(), use in this dir:  stringer.exe -type TokenType
const (
	itemCommentKey TokenType = iota
	itemComment
	itemText
	itemName
	itemIPPart
	itemDotKey
	itemSpaceSepa
	itemTrail
	itemError
	itemEOF
)

func lexStateAfterName(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case r == '\n':
			l.emit(itemTrail)
			return lexStateInit
		case r == EOFRune:
			l.emit(itemTrail)
			return nil
		default:
			l.ignore()
		}
	}
}

func lexStateName(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '.':
			// nothing to do
		case unicode.IsSpace(r):
			l.rewind()
			if l.position > l.start {
				l.emit(itemName)
				return lexStateAfterName
			}
			return lexStateInit
		case r == EOFRune:
			l.emit(itemName)
			return nil
		default:
			l.emit(itemText)
			return lexStateInit
		}
	}
}

func lexStateBeforeName(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			l.rewind()
			l.emit(itemSpaceSepa)
			return lexStateName
		case r == '\n' || r == '\r':
			l.rewind()
			l.emit(itemText)
			return lexStateInit
		case r == EOFRune:
			l.emit(itemText)
			return nil
		}
	}
}

func lexStateIpD(l *L) StateFunc {
	dcount := 0
	for {
		switch r := l.next(); {
		case unicode.IsDigit(r):
			dcount++
			if dcount > 3 {
				return l.errorf("IP block D is too big")
			}
		case r == '\n' || r == '\r':
			l.rewind()
			l.emit(itemText)
			return lexStateInit
		case unicode.IsSpace(r):
			l.rewind()
			l.emit(itemIPPart)
			l.position++
			return lexStateBeforeName
		case r == EOFRune:
			l.emit(itemText)
			return nil
		default:
			l.emit(itemText)
			return l.errorf("Expect digit or dot")
		}
	}
}

func lexStateIpABC(l *L) StateFunc {
	dcount := 0
	for {
		switch r := l.next(); {
		case unicode.IsDigit(r):
			dcount++
			if dcount > 3 {
				return l.errorf("IP block is too big")
			}
		case r == '.':
			dcount = 0
			l.rewind()
			keyInfoState, err := l.popPropInfo()
			if err != nil {
				return l.errorf("Error: %v", err)
			}
			l.emit(keyInfoState.TokenType)
			l.next()
			l.emit(itemDotKey)
			if len(l.stackPropInfo) == 0 {
				return lexStateIpD
			}
		case r == '\n' || r == '\r':
			l.rewind()
			l.emit(itemText)
			return lexStateInit
		case r == EOFRune:
			l.emit(itemText)
			return nil
		default:
			l.emit(itemText)
			return l.errorf("Expect digit or dot")
		}
	}
}

func lexStateInComment(l *L) StateFunc {
	l.position++
	l.emit(itemCommentKey)
	for {
		switch r := l.next(); {
		case r == '\n' || r == '\r':
			l.rewind()
			l.emit(itemComment)
			return lexStateInit
		case r == EOFRune:
			l.emit(itemText)
			return nil
		}
	}
}

func lexStateInit(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case unicode.IsDigit(r):
			l.rewind()
			l.emit(itemText)
			l.stackPropInfo = append(l.stackPropInfo, PropInfo{TokenType: itemIPPart}, PropInfo{TokenType: itemIPPart}, PropInfo{TokenType: itemIPPart})
			return lexStateIpABC
		case r == '#':
			l.rewind()
			l.emit(itemText)
			return lexStateInComment
		case r == EOFRune:
			l.emit(itemText)
			return nil
		}
	}
}

///////////////////////////////////////////////////////////////////////
type HostsParser struct {
	ChangedSource string
	HasChanges    bool
	DebugParser   bool
	MapIp         map[string]string
	UpdatedHosts  []string
}

func (hp *HostsParser) ParseHosts(source string) error {
	ll := NewL(source, lexStateInit)
	hostIP := ""
	currSpace := ""
	defer close(ll.tokens)
	for {
		item := ll.nextItem()
		if hp.DebugParser {
			fmt.Println("*** type: ", item.Type.String(), item.String())
		}
		switch item.Type {
		case itemIPPart:
			hostIP += item.Value
		case itemDotKey:
			hostIP += item.Value
		case itemSpaceSepa:
			currSpace = item.Value
		case itemName:
			currName := item.Value
			if changedIP, ok := hp.MapIp[currName]; ok {
				if changedIP != hostIP {
					hp.HasChanges = true
					hp.ChangedSource += fmt.Sprintf("%s          %s\n", changedIP, currName)
					hp.UpdatedHosts = append(hp.UpdatedHosts, currName)
				}
			} else {
				hp.ChangedSource += fmt.Sprintf("%s%s%s\n", hostIP, currSpace, currName)
			}
		case itemTrail:
		case itemError:
			return errors.New(item.Value)
		case itemEOF:
			return nil
		default:
			hostIP = ""
			hp.ChangedSource += item.Value

		}
	}
}
