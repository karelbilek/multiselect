package main

import (
	"encoding/json"
	"fmt"

	"github.com/karelbilek/multiselect"
)

func main() {
	// file, err := dialog.File().Title("Save As").Filter("All Files", "*").Save()
	// fmt.Println(file)
	// fmt.Println("Error:", err)
	// dialog.Message("You chose file: %s", file).Title("Goodbye world!").Error()
	dir, err := multiselect.Fileselect("select CSV", "csv", "comma separated value")
	js, _ := json.MarshalIndent(dir, "", "    ")
	fmt.Printf("%s\n", js)
	fmt.Println("Error:", err)
}
