package cdl

import (
	"fmt"
	"log"

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

func (p *Parser) parseDimensions(f *types.File) bool {
	if p.consume() || p.last.Type != TkColon {
		log.Panicf("`:` is required after a `dimensions` directive")
	}
	f.Dimensions = []types.Dimension{}
	for {
		if p.consume() || p.last.Type == TkCurClose {
			return false
		}
		if p.last.Type == TkData || p.last.Type == TkVariables {
			return true
		}
		var d types.Dimension
		if p.last.Type != TkName {
			panic("dimension name expected")
		}
		d.Name = p.last.Text
		if p.consume() || p.last.Type != TkEqual {
			panic("`=` expected")
		}

		if p.consume() || p.last.Type != TkInt {
			panic("dimension name expected")
		}
		d.Len = int32(p.last.NumVal)
		f.Dimensions = append(f.Dimensions, d)
		if p.consume() || p.last.Type != TkSemicolon {
			panic("`;` expected")
		}
	}
}

// TODO: parse attributes
func (p *Parser) parseVariables(f *types.File) bool {
	if p.consume() || p.last.Type != TkColon {
		log.Panicf("`:` is required after a `variables` directive")
	}
	dimensions := map[string]*types.Dimension{}
	for _, d := range f.Dimensions {
		dimensions[d.Name] = &d
	}
	if p.consume() || p.last.Type == TkCurClose {
		return false
	}
	for {

		var v types.Var
		if p.last.Type != TkVarType {
			panic("variable type expected")
		}
		v.Size = 1
		for _, d := range f.Dimensions {
			v.Size *= d.Len
		}
		switch p.last.Text {
		case "float":
			v.Size *= 4
			v.Type = types.Float
		case "byte":
			v.Type = types.Byte
		case "char":
			v.Type = types.Char
		case "short":
			v.Size *= 2
			v.Type = types.Short
		case "int":
			v.Size *= 4
			v.Type = types.Int
		case "double":
			v.Size *= 8
			v.Type = types.Double
		}

		if p.consume() || p.last.Type != TkName {
			panic("variable name expected")
		}
		v.Name = p.last.Text

		if p.consume() || p.last.Type != TkParOpen {
			panic("dimension list expected")
		}

		for {
			if p.consume() || p.last.Type != TkName {
				panic("dimension name expected")
			}
			d, ok := dimensions[p.last.Text]
			if !ok {
				log.Panicf("unknown dimension name `%s`", p.last.Text)
			}
			if p.consume() {
				panic("dimension list not closed by `)`")
			}

			v.Dimensions = append(v.Dimensions, d)

			if p.last.Type == TkParClose {
				if p.consume() || p.last.Type != TkSemicolon {
					panic("`;` expected")
				}

				f.Vars.Set(v.Name, v)
				if p.consume() {
					panic("`}` expected")
				}
				if p.last.Type == TkName {
					return true
				}
				if p.last.Type == TkCurClose || p.last.Type == TkData || p.last.Type == TkVariables {
					return false
				}
				break
			}

		}

	}
}

func (p *Parser) parseData(f *types.File) bool {
	return false
}

func (p *Parser) parseStatement(f *types.File) bool {

	switch p.last.Type {
	case TkDimensions:
		return p.parseDimensions(f)
	case TkVariables:
		return p.parseVariables(f)
	case TkData:
		return p.parseData(f)
	case TkCurClose:
		return false
	default:
		log.Panicf("unexpected token %v", p.last)
	}

	return false
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
	if p.consume() {
		panic("expected }")
	}
	for p.parseStatement(f) {

	}

	if p.last.Type != TkCurClose {
		panic("expected }")
	}

	if !p.consume() {
		log.Panicf("unexpected %v", p.last)
	}
}
