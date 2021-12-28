package cdl

import (
	"fmt"
	"log"

	"github.com/parro-it/ncdf/types"
)

// Parser ...
type Parser struct {
	Tokens chan Token
	last   Token
}

func (p *Parser) consume() bool {
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

	if p.consume() || p.last.Type == TkCurClose {
		return false
	}

	dimensions := mapDimensions(f)

	for {
		if p.last.Type == TkVarType {
			v := p.parseVariable(dimensions)
			f.Vars.Set(v.Name, v)
		} else if p.last.Type == TkName || p.last.Type == TkColon {
			a, v := p.parseAttribute(f, dimensions)
			if v == nil {
				f.Attrs.Set(a.Name, a)
			} else {
				v.Attrs.Set(a.Name, a)
				f.Vars.Set(v.Name, *v)
			}

		} else {
			log.Panicf("unexpected token %v", p.last)
		}

		if p.last.Type == TkCurClose {
			return false
		}

		if p.last.Type == TkDimensions || p.last.Type == TkData || p.last.Type == TkVariables {
			return true
		}

	}
}

func (p *Parser) parseAttribute(f *types.File, dimensions map[string]*types.Dimension) (types.Attr, *types.Var) {
	var a types.Attr
	var v *types.Var
	if p.last.Type == TkName {
		tmp := f.Vars.Get(p.last.Text)
		v = &tmp

		if p.consume() || p.last.Type != TkColon {
			panic(": expected")
		}
	}

	if p.consume() || p.last.Type != TkName {
		panic("attribute name expected")
	}
	a.Name = p.last.Text

	if p.consume() || p.last.Type != TkEqual {
		panic("attribute name expected")
	}

	if p.consume() {
		panic("attribute value expected")
	}

	if p.last.Type == TkDec {
		a.Val = []float32{float32(p.last.NumVal)}
		a.Type = types.Float
	} else if p.last.Type == TkInt {
		a.Val = []int16{int16(p.last.NumVal)}
		a.Type = types.Short
	} else if p.last.Type == TkStr {
		a.Val = []byte(p.last.Text)
		a.Type = types.Char
	} else {
		panic("unsupported type")
	}
	if p.consume() || p.last.Type != TkSemicolon {
		panic("; expected")
	}
	if p.consume() {
		panic("} expected")
	}
	return a, v
}

func (p *Parser) parseVariable(dimensions map[string]*types.Dimension) types.Var {
	var v types.Var
	if p.last.Type != TkVarType {
		log.Panicf("variable type expected")
	}
	size := 1
	for _, d := range dimensions {
		size *= int(d.Len)
	}
	v.Type = types.FromCDLName(p.last.Text)
	v.Size = int32(v.Type.ArraySize(size))

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
			if p.consume() {
				panic("`}` expected")
			}

			return v
		}

	}
}

func mapDimensions(f *types.File) map[string]*types.Dimension {
	dimensions := map[string]*types.Dimension{}
	for _, d := range f.Dimensions {
		dimensions[d.Name] = &d
	}
	return dimensions
}

// TODO: parse CDL data
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
