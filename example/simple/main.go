package main

import (
	"encoding/json"
	"fmt"

	"github.com/karelbilek/multiselect"
)

func main() {
	dir, err := multiselect.Fileselect("select CSV", "csv", "comma separated value")
	js, _ := json.MarshalIndent(dir, "", "    ")
	fmt.Printf("%s\n", js)
	fmt.Println("Error:", err)
}
