package main

import (
	"encoding/json/v2"
	"flag"
	"io"
	"log"
	"os"
	"path"
)

type Settings struct {
	Generate    bool
	LocalUpload bool
}

var GlobalSettings = &Observable[Settings]{
	Data: Settings{
		Generate:    false,
		LocalUpload: false,
	},
	Subscribers: []*func(Settings){},
}

func applyFlags() {
	generateFlag := flag.Bool("gen", false, "Generate typescript types to ./frontend")
	localUploadFlag := flag.Bool("localUpload", true, "Upload files to local dir instead of sftp server")
	flag.Parse()

	GlobalSettings.Update(func(s Settings) Settings {
		s.Generate = *generateFlag
		s.LocalUpload = *localUploadFlag
		return s
	})
}

var State = &Observable[[]Play]{
	Data:        []Play{},
	Subscribers: []*func([]Play){},
}

func toFileInfo(entries []os.DirEntry) ([]os.FileInfo, error) {
	var result []os.FileInfo

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		result = append(result, info)
	}

	return result, nil
}

func update(client *Client) error {
	plays := []Play{}
	var loc = GlobalSettings.Data.LocalUpload

	var entries, _ []os.FileInfo
	var err error
	if loc {
		entries1, _ := os.ReadDir("plays")
		entries, err = toFileInfo(entries1)
	} else {
		entries, err = client.sftp.ReadDir("plays")
	}
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := path.Join("plays", entry.Name())
		var content []byte
		if loc {
			content, err = os.ReadFile(path.Join(dirPath, "content.json"))
		} else {
			file, err := client.sftp.Open(path.Join(dirPath, "content.json"))
			if err != nil {
				log.Println(err)
				continue
			}

			content, err = io.ReadAll(file)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		if err != nil {
			log.Println(err)
			continue
		}

		newPlay := Play{}

		err = json.Unmarshal(content, &newPlay)
		if err != nil {
			log.Println(err)
			continue
		}

		plays = append(plays, newPlay)

	}

	State.UpdateValue(plays)
	return nil
}
