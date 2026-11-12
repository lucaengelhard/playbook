package main

import (
	"embed"
	"fmt"
	"os"
)

//go:embed all:frontend
var template embed.FS

type Play struct {
	Title     string
	HeroImage string
	Content   string
}

func Scaffold(client *Client, play Play) {
	entries, _ := template.ReadDir("template")
	// TODO make recursive
	for _, e := range entries {
		path := "template/" + e.Name()
		_, err := os.Open(path)
		if err != nil {
			fmt.Println("Failed to open local file:", err)
			return
		}

		//client.Upload(file, play.Title+"/"+e.Name())
	}
}
