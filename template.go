package main

import (
	"embed"
	"fmt"
	"os"
)

//go:embed all:template
var template embed.FS

type Template struct{}

type Play struct {
	Title     string
	HeroImage string
	Content   string
}

func (t *Template) Scaffold(play Play) {
	entries, _ := template.ReadDir("template")
	// TODO make recursive
	for _, e := range entries {
		path := "template/" + e.Name()
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Failed to open local file:", err)
			return
		}

		client.Upload(file, play.Title+"/"+e.Name())
	}
}
