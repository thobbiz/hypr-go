package main

import (
	"log"

	"github.com/thobbiz/hypr-go/layout"
)

func main() {
	// err := layout.SaveWorkspace(4, "reading")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// workspaceID := 5
	err := layout.RestoreWorkspace("coded", nil)
	if err != nil {
		log.Fatal()
	}
}
