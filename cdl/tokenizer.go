package cdl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// TokenType describe the type of tokens
// of a CDL file.
type TokenType int

const (
	// TkEmpty represents a non valid token
	TkEmpty TokenType = iota

	// TkName represents the name of variables, dimensions etc.
	TkName

	// TkVarType represents the type of variables (byte,short,int,char,float,double)
	TkVarType

	// TkCurOpen - { char
	TkCurOpen

	// TkCurClose - } char
	TkCurClose

	// TkParOpen - ( char
	TkParOpen

	// TkParClose - ) char
	TkParClose

	// TkEqual - = char
	TkEqual

	// TkColon - : char
	TkColon

	// TkSemicolon - ; char
	TkSemicolon

	// TkComma - , char
	TkComma

	// TkComment ...
	TkComment

	// TkNetCdf - netcdf string
	TkNetCdf

	// TkStr ...
	TkStr
	// TkInt ...
	TkInt
	// TkDec ...
	TkDec
)

// CodePoint represent a single point in a source file,
// either with column/index numer of character index.
type CodePoint struct {
	Col uint
	Row uint
	Idx uint
}

// CodePosition represent a the position of a chunk of text
// in a source file, using a start and end CodePoint
// (inclusives)
type CodePosition struct {
	Start CodePoint
	End   CodePoint
}

// Token ...
type Token struct {
	Pos    CodePosition
	Type   TokenType
	Text   string
	NumVal float64
}

type tokenizer struct {
	r     *bufio.Reader
	res   chan Token
	curr  rune
	atEnd bool

	curpos CodePoint
}

// Tokenize ...
func Tokenize(r io.Reader) (ch chan Token, errs chan error) {
	errs = make(chan error)

	tkn := tokenizer{
		r:   bufio.NewReader(r),
		res: make(chan Token),
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				close(tkn.res)
				var err error
				switch e := r.(type) {
				case error:
					err = e
				default:
					s := fmt.Sprint(r)
					err = errors.New(s)
				}
				errs <- fmt.Errorf("Tokenization failed: %w", err)
			}
			close(errs)
		}()

		tkn.read()
	}()

	return tkn.res, errs

}

// TODO: name must support special charaters aftere first letter
func (tkn *tokenizer) readIdent() {
	var buf strings.Builder
	for !tkn.atEnd && (unicode.IsLetter(tkn.curr) || unicode.IsDigit(tkn.curr)) {
		buf.WriteRune(tkn.curr)
		tkn.readRune()
	}
	val := buf.String()
	var tkType TokenType
	switch val {
	case "byte", "short", "int", "float", "double", "char":
		tkType = TkVarType
	case "netcdf":
		tkType = TkNetCdf
	default:
		tkType = TkName
	}
	tkn.res <- Token{
		Type: tkType,
		Text: val,
	}
}

// TODO: implements comments parsing
func (tkn *tokenizer) skipComment() {
	tkn.readRune()
	if tkn.curr != '/' {
		panic("unexpected char `/`")
	}

	tkn.readRune()
	for !tkn.atEnd && tkn.curr != '\n' {
		tkn.readRune()
	}
}

func (tkn *tokenizer) read() {
	tkn.readRune()
	for !tkn.atEnd {
		//fmt.Printf("%c %t\n", tkn.curr, unicode.IsSymbol(tkn.curr) || unicode.IsPunct(tkn.curr))
		switch true {
		case tkn.curr == '/':
			tkn.skipComment()
		case tkn.curr == '"':
			tkn.readString()
		case unicode.IsDigit(tkn.curr):
			tkn.readNumber()

		case unicode.IsLetter(tkn.curr):
			tkn.readIdent()
		default:
			tk := Token{
				Pos:  CodePosition{tkn.curpos, tkn.curpos},
				Text: string([]rune{tkn.curr}),
			}
			// single char tokens
			switch tkn.curr {
			case '{':
				tk.Type = TkCurOpen
			case '}':
				tk.Type = TkCurClose
			case '(':
				tk.Type = TkParOpen
			case ')':
				tk.Type = TkParClose
			case '=':
				tk.Type = TkEqual
			case ':':
				tk.Type = TkColon
			case ';':
				tk.Type = TkSemicolon
			case ',':
				tk.Type = TkComma
			}

			if tk.Type != TkEmpty {
				tkn.res <- tk
			}

			tkn.readRune()
		}
	}
	close(tkn.res)
}

// TODO: parse negative numbers
func (tkn *tokenizer) readNumber() {
	tokType := TkInt
	text := ""
	for !tkn.atEnd && isNumChar(tkn.curr) {
		if tkn.curr == '.' {
			if tokType == TkInt {
				tokType = TkDec
			} else {
				panic("unexpected dot")
			}
		}
		text += string(tkn.curr)
		tkn.readRune()
	}
	// TODO: parse integers with dedicated func
	num, _ := strconv.ParseFloat(text, 64)

	tkn.res <- Token{
		Type: tokType,
		// TODO: use different field for integers
		NumVal: num,
		Text:   text,
	}
}

func (tkn *tokenizer) readRune() {
	r, _, err := tkn.r.ReadRune()
	if err == io.EOF {
		tkn.curr = rune(0)
		tkn.atEnd = true
		return
	}
	if err != nil {
		panic(err)
	}

	tkn.curr = r
}

func isNumChar(ch rune) bool {
	return unicode.IsDigit(ch) || ch == '.'
}

func (tkn *tokenizer) readString() {
	text := ""
	tkn.readRune()

	for !tkn.atEnd && tkn.curr != '"' {
		text += string(tkn.curr)
		tkn.readRune()
	}

	if tkn.curr != '"' {
		panic("unclosed string")
	}
	tkn.readRune()

	tkn.res <- Token{
		Type: TkStr,
		Text: text,
	}
}
