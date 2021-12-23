package main

import (
	"fmt"
	"os"

	"github.com/parro-it/ncdf/cdl"
	"github.com/parro-it/ncdf/write"
)

func main() {
	in, err := os.Open("/mnt/repos/parro-it/ncdf/fixtures/simple.cdl")
	if err != nil {
		panic(err)
	}
	defer in.Close()
	tks, _ := cdl.Tokenize(in)

	p := cdl.Parser{Tokens: tks}
	f, err := p.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", f.ByteSize())

	out, err := os.OpenFile("/mnt/repos/parro-it/ncdf/simple2.nc", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	f.ComputeSizes()

	if err := write.Header(f, out); err != nil {
		panic(err)
	}
	for _, v := range f.Vars {
		write.VarData(v, make([]float32, 100), out)
	}

}
