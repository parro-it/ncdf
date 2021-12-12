package main

import (
	"encoding/json"
	"fmt"

	"github.com/parro-it/ncdf"
)

func main() {
	f, err := ncdf.Open("fixtures/example.nc")
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
