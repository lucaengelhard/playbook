package main

import (
	"embed"
	"encoding/json/v2"
	"errors"
	"log"
	"os"
	"path"

	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
)

//go:embed all:frontend
var frontendFs embed.FS
var devMode = false

type Play struct {
	Title     string
	HeroImage string
	Content   string
	Gallery   []string `json:"Gallery,omitempty"`
}

func ScaffoldOffline(play Play) error {
	if play.Title == "" {
		return errors.New("Cannot upload play with empty title :(")
	}

	dirPath := path.Join("plays", play.Title)

	err := os.MkdirAll(path.Join(dirPath, "media"), 0755)
	if err != nil {
		return err
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
		os.WriteFile(dirPath+"/"+entry.Name(), file, 0755)

	}

	heroImage, err := os.ReadFile(play.HeroImage)
	if err != nil {
		return err
	}

	play.HeroImage = path.Join("media", path.Base(play.HeroImage))
	err = os.WriteFile(path.Join(dirPath, play.HeroImage), heroImage, 0755)
	if err != nil {
		return err
	}

	content, err := json.Marshal(play)
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(dirPath, "content.json"), content, 0755)

}

func Scaffold(client *Client, play Play) error {
	if play.Title == "" {
		return errors.New("Cannot upload play with empty title :(")
	}

	dirPath := path.Join("plays", play.Title)

	if client.sftp == nil {
		client.dialog.Show()
		return errors.New("Connect before uploading")
	}

	err := client.sftp.MkdirAll(path.Join(dirPath, "media"))
	if err != nil {
		return err
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

		err = client.Upload(file, path.Join(dirPath, entry.Name()))
		if err != nil {
			return err
		}
	}

	heroImage, err := os.ReadFile(play.HeroImage)
	if err != nil {
		return err
	}

	play.HeroImage = path.Join("media", path.Base(play.HeroImage))
	err = client.Upload(heroImage, path.Join(dirPath, play.HeroImage))
	if err != nil {
		return err
	}

	content, err := json.Marshal(play)
	if err != nil {
		return err
	}

	return client.Upload(content, path.Join(dirPath, "content.json"))
}

func createEditor(client *Client) *widget.Form {
	play := Play{}

	devModeCheckbox := widget.NewCheck("Stores directory locally", func(b bool) { devMode = b })
	titleEntry := widget.NewEntry()
	heroImageButton := widget.NewButton("Select...", func() {
		file, err := zenity.SelectFile(
			zenity.FileFilters{
				{
					Name:     "Images",
					Patterns: []string{"*.png", "*.jpg"},
					CaseFold: true,
				},
				{
					Name:     "All Files",
					Patterns: []string{"*"},
					CaseFold: true,
				},
			})

		if err != nil {
			// TODO
			return
		}

		play.HeroImage = file
	})

	content := widget.NewMultiLineEntry()

	galleryButton := widget.NewButton("Select Images...", func() {
		files, err := zenity.SelectFileMultiple(
			zenity.FileFilters{
				{
					Name:     "Images",
					Patterns: []string{"*.png", "*.jpg"},
					CaseFold: true,
				},
				{
					Name:     "All Files",
					Patterns: []string{"*"},
					CaseFold: true,
				},
			},
		)

		if err != nil {
			// TODO
			return
		}

		play.Gallery = files
	})

	form := widget.NewForm(
		widget.NewFormItem("Dev Mode", devModeCheckbox),
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Content", content),
		widget.NewFormItem("Hero Image", heroImageButton),
		widget.NewFormItem("Gallery", galleryButton),
	)
	form.OnSubmit = func() {
		play.Title = titleEntry.Text
		play.Content = content.Text

		var err error
		if devMode {
			err = ScaffoldOffline(play)
		} else {
			err = Scaffold(client, play)
		}
		if err != nil {
			log.Println(err)
		}
		// TODO
		// Clear fields
		// reset play struct
		// Show Progress dialog (https://docs.fyne.io/api/v2/dialog/progressdialog/)
		// Show Success Dialog
	}

	return form
}
