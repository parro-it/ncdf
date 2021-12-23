package main

import (
	"fmt"

	"github.com/parro-it/ncdf/read"
)

func main() {
	f, err := read.HeaderFromDisk("fixtures/exampl2.nc")
	if err != nil {
		panic(nil)
	}
	//defer f.Close()
	/*buf, err := json.MarshalIndent(f, "  ", "  ")
	if err != nil {
		panic(nil)
	}
	fmt.Println(string(buf))*/
	fmt.Println(f.CDL())
}
