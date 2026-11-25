package main

import (
	"encoding/json/v2"
	"flag"
	"log"
	"os"
	"path"
)

type Tuple[T any] struct {
	first  T
	second T
}
type Settings struct {
	Generate   bool
	WindowSize Tuple[float32]
}

var GlobalSettings = &Observable[Settings]{
	Data: Settings{
		Generate:   false,
		WindowSize: Tuple[float32]{800, 600},
	},
	Subscribers: []*func(Settings){},
}

func applyFlags() {
	generateFlag := flag.Bool("gen", false, "Generate typescript types to ./frontend")
	flag.Parse()

	GlobalSettings.Update(func(s Settings) Settings {
		s.Generate = *generateFlag
		return s
	})
}

var State = &Observable[[]Play]{
	Data:        []Play{},
	Subscribers: []*func([]Play){},
}

func update() error {
	plays := []Play{}

	entries, err := os.ReadDir("plays")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := path.Join("plays", entry.Name())

		content, err := os.ReadFile(path.Join(dirPath, "content.json"))
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
