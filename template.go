package main

import (
	"embed"
	"encoding/json/v2"
	"log"
	"os"
	"path"
)

//go:embed all:frontend
var frontendFs embed.FS

type Play struct {
	Title     string
	HeroImage string
	Content   string
}

func Scaffold(client *Client, play Play) {
	dirPath := path.Join("plays", play.Title)

	err := client.sftp.MkdirAll(path.Join(dirPath, "media"))
	if err != nil {
		log.Println("dir", err)
		return
	}

	entries, err := frontendFs.ReadDir("frontend")

	for _, entry := range entries {
		if entry.IsDir() {
			// TODO
			continue
		}
		filePath := path.Join("frontend", entry.Name())
		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("Error while reading template file", err)
			continue
		}

		client.Upload(file, path.Join(dirPath, entry.Name()))
	}

	heroImage, err := os.ReadFile(play.HeroImage)
	if err != nil {
		log.Println("Error while reading hero image", err)
		return
	}

	play.HeroImage = path.Join("media", path.Base(play.HeroImage))
	client.Upload(heroImage, path.Join(dirPath, play.HeroImage))

	content, err := json.Marshal(play)
	if err != nil {
		log.Println(err)
		return
	}

	client.Upload(content, path.Join(dirPath, "content.json"))
}
