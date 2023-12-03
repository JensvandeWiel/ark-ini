package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenIdentifier
	tokenString
	tokenNumber
	tokenTrue
	tokenFalse
	tokenComma
	tokenEquals
	tokenOpenParen
	tokenCloseParen
	tokenOpenBrace
	tokenCloseBrace
	tokenEOF
)

type token struct {
	typ  tokenType
	pos  int
	val  string
	text string
}

type lexer struct {
	input string
	pos   int
}

func newLexer(input string) *lexer {
	return &lexer{input: input}
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += size
	return r
}

func (l *lexer) backup() {
	l.pos--
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) emit(typ tokenType) *token {
	return &token{typ: typ, pos: l.pos, text: l.input, val: l.input[:l.pos]}
}

func (l *lexer) skipWhitespace() {
	for {
		r := l.next()
		if !unicode.IsSpace(r) {
			l.backup()
			break
		}
	}
}

func (l *lexer) lexIdentifier() *token {
	for {
		r := l.next()
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != ':' {
			l.backup()
			break
		}
	}
	return l.emit(tokenIdentifier)
}

func (l *lexer) lexString() *token {
	// Assuming strings are enclosed in double quotes.
	for {
		r := l.next()
		if r == '"' {
			break
		}
	}
	return l.emit(tokenString)
}

func (l *lexer) lexNumber() *token {
	for {
		r := l.next()
		if !unicode.IsDigit(r) && r != '.' {
			l.backup()
			break
		}
	}
	return l.emit(tokenNumber)
}

func (l *lexer) lexTrue() *token {
	for i := 0; i < len("true"); i++ {
		if l.next() != rune("true"[i]) {
			l.backup()
			return l.emit(tokenIdentifier)
		}
	}
	return l.emit(tokenTrue)
}

func (l *lexer) lexFalse() *token {
	for i := 0; i < len("false"); i++ {
		if l.next() != rune("false"[i]) {
			l.backup()
			return l.emit(tokenIdentifier)
		}
	}
	return l.emit(tokenFalse)
}

func (l *lexer) Lex() *token {
	l.skipWhitespace()

	switch r := l.next(); {
	case unicode.IsLetter(r):
		l.backup()
		return l.lexIdentifier()
	case r == '"':
		return l.lexString()
	case unicode.IsDigit(r):
		l.backup()
		return l.lexNumber()
	case r == 't':
		return l.lexTrue()
	case r == 'f':
		return l.lexFalse()
	case r == ',':
		return l.emit(tokenComma)
	case r == '=':
		return l.emit(tokenEquals)
	case r == '(':
		return l.emit(tokenOpenParen)
	case r == ')':
		return l.emit(tokenCloseParen)
	case r == '{':
		return l.emit(tokenOpenBrace)
	case r == '}':
		return l.emit(tokenCloseBrace)
	case r == 0:
		return l.emit(tokenEOF)
	default:
		return l.emit(tokenError)
	}
}

func main() {
	input := `ConfigOverrideSupplyCrateItems=(SupplyCrateClassString="SupplyCrate_Cave_QualityTier1_EX",MinItemSets=0,MaxItemSets=1,NumItemSetsPower=0,bSetsRandomWithoutReplacement=True,ItemSets=((SetName="Carcha Saddle",ItemEntries=((ItemEntryName="Planos: Monturas",Items=(%2FScript%2FEngine.BlueprintGeneratedClass'"%2FGame%2FPrimalEarth%2FCoreBlueprints%2FItems%2FArmor%2FSaddles%2FPrimalItemArmor_CarchaSaddle.PrimalItemArmor_CarchaSaddle_C"'),MinQuantity=1.000000,MinQuality=1.000000,bForceBlueprint=False,ChanceToBeBlueprintOverride=0.500000,ChanceToActuallyGiveItem=0.900000)),SetWeight=0.015000,bItemsRandomWithoutReplacement=True)),bAppendItemSets=True,bAppendPreventIncreasingMinMaxItemSets=False)`

	lex := newLexer(input)
	for {
		t := lex.Lex()
		if t.typ == tokenEOF {
			break
		}
		fmt.Printf("Token: %-15s Value: %s\n", t.typ, t.val)
	}
}
