package main

import (
	"encoding/json"
	"fmt"

	"github.com/parro-it/ncdf"
)

func main() {
	f, err := ncdf.Open("fixtures/exampl2.nc")
	if err != nil {
		panic(nil)
	}
	defer f.Close()
	buf, err := json.MarshalIndent(f, "  ", "  ")
	if err != nil {
		panic(nil)
	}
	fmt.Println(string(buf))
}
