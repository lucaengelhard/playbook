package main

import (
	"embed"
	"encoding/json/v2"
	"io"
	"log"
	"os"
	"path"
)

//go:embed all:frontend
var template embed.FS

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

	indexHtml, err := template.ReadFile("frontend/index.html")
	if err != nil {
		log.Println("index", err)
		return
	}

	indexHtmlRemote, err := client.sftp.Create(path.Join(dirPath, "index.html"))
	if err != nil {
		log.Println("indRem", err)
		return
	}
	defer indexHtmlRemote.Close()
	indexHtmlRemote.Write(indexHtml)

	script, err := template.ReadFile("frontend/script.js")
	if err != nil {
		log.Println("index", err)
		return
	}

	scriptRemote, err := client.sftp.Create(path.Join(dirPath, "script.js"))
	if err != nil {
		log.Println("indRem", err)
		return
	}
	defer scriptRemote.Close()
	scriptRemote.Write(script)

	heroImage, err := os.Open(play.HeroImage)
	if err != nil {
		log.Println("open hero", err)
		return
	}
	defer heroImage.Close()
	play.HeroImage = path.Join("media", path.Base(play.HeroImage))

	heroImageRemote, err := client.sftp.Create(path.Join(dirPath, play.HeroImage))
	if err != nil {
		log.Println("rem", err)
		return
	}
	defer heroImageRemote.Close()

	_, err = io.Copy(heroImageRemote, heroImage)
	if err != nil {
		log.Println("copy", err)
		return
	}

	content, err := json.Marshal(play)
	if err != nil {
		log.Println(err)
		return
	}
	contentRemote, err := client.sftp.Create(path.Join(dirPath, "content.json"))
	if err != nil {
		log.Println(err)
		return
	}
	defer contentRemote.Close()
	contentRemote.Write(content)

	/* entries, _ := template.ReadDir("template")
	// TODO make recursive
	for _, e := range entries {
		path := "template/" + e.Name()
		_, err := os.Open(path)
		if err != nil {
			fmt.Println("Failed to open local file:", err)
			return
		}

		//client.Upload(file, play.Title+"/"+e.Name())
	} */
}
