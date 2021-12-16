package cdl

import (
	"fmt"

	"github.com/parro-it/ncdf/types"
)

// Parser ...
type Parser struct {
	Tokens chan Token
	buf    *Token
	last   Token
}

func (p *Parser) consume() bool {
	if p.buf != nil {
		p.last = *p.buf
		p.buf = nil
		return p.last.Type == TkEmpty
	}
	p.last = <-p.Tokens
	return p.last.Type == TkEmpty
}
func (p *Parser) parseStatement(f *types.File) bool {
	eof := p.consume()
	if eof {
		return false
	}
	return p.last.Type != TkCurClose
}

func (p *Parser) peek() Token {
	tk := <-p.Tokens
	p.buf = &tk
	return *p.buf
}

// Parse ...

func (p *Parser) Parse() (f *types.File, err error) {
	f = new(types.File)
	defer func() {
		if e := recover(); e != nil {
			f = nil
			err = fmt.Errorf("Parse failed: %v", e)
		}
	}()
	p.parseProgram(f)
	return
}

func (p *Parser) parseProgram(f *types.File) {
	p.consume()

	if p.last.Type != TkNetCdf {
		panic("expected netcdf word")
	}

	p.consume()
	if p.last.Type != TkName {
		panic("expected file name")
	}

	p.consume()
	if p.last.Type != TkCurOpen {
		panic("expected {")
	}

	for p.parseStatement(f) {

	}

	if p.last.Type != TkCurClose {
		panic("expected }")
	}

}
