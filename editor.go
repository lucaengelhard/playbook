package main

import (
	"embed"
	"encoding/json/v2"
	"errors"
	"io/fs"
	"log"
	"os"
	"path"

	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
)

//go:embed all:frontend
var frontendFs embed.FS

type Play struct {
	Title     string
	HeroImage string
	Content   string
	Gallery   []string `json:"Gallery,omitempty"`
}

func uploadFrontendFiles(entries []fs.DirEntry, dirPath string, client *Client) error {
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
		if GlobalSettings.Data.LocalUpload {
			err = os.WriteFile(path.Join(dirPath, entry.Name()), file, 0755)
		} else {
			err = client.Upload(file, path.Join(dirPath, entry.Name()))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func Scaffold(client *Client, play Play) error {
	//loc is true if we want to create the playbook locally
	var loc = GlobalSettings.Data.LocalUpload
	//Check if title is set
	if play.Title == "" {
		return errors.New("Cannot upload play with empty title :(")
	}

	dirPath := path.Join("plays", play.Title)

	if !loc {
		//Check if client is connected
		if client.sftp == nil {
			client.dialog.Show()
			return errors.New("Connect before uploading")
		}
	}

	var err error
	//Create media folder
	if loc {
		err = os.MkdirAll(path.Join(dirPath, "media"), 0755)
	} else {
		err = client.sftp.MkdirAll(path.Join(dirPath, "media"))
	}
	if err != nil {
		return err
	}

	//Read frontend files
	entries, err := frontendFs.ReadDir("frontend")
	err = uploadFrontendFiles(entries, dirPath, client)
	if err != nil {
		return err
	}
	//Read HeroImage
	heroImage, err := os.ReadFile(play.HeroImage)
	if err != nil {
		return err
	}

	//Upload HeroImage
	play.HeroImage = path.Join("media", path.Base(play.HeroImage))
	if loc {
		err = os.WriteFile(path.Join(dirPath, play.HeroImage), heroImage, 0755)
	} else {
		err = client.Upload(heroImage, path.Join(dirPath, play.HeroImage))
	}
	if err != nil {
		return err
	}

	//save play as json file
	content, err := json.Marshal(play)
	if err != nil {
		return err
	}
	if loc {
		return os.WriteFile(path.Join(dirPath, "content.json"), content, 0755)
	} else {
		return client.Upload(content, path.Join(dirPath, "content.json"))
	}
}

func createEditor(client *Client, play Play) *widget.Form {
	localUploadCheckbox := widget.NewCheck("Stores directory locally", func(b bool) {
		GlobalSettings.Update(func(s Settings) Settings {
			s.LocalUpload = b
			return s
		})
	})
	localUploadCheckbox.Checked = GlobalSettings.Data.LocalUpload
	titleEntry := widget.NewEntry()
	titleEntry.Text = play.Title
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
	content.Text = play.Content

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
		widget.NewFormItem("Local Upload", localUploadCheckbox),
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Content", content),
		widget.NewFormItem("Hero Image", heroImageButton),
		widget.NewFormItem("Gallery", galleryButton),
	)
	form.OnSubmit = func() {
		play.Title = titleEntry.Text
		play.Content = content.Text

		var err error
		err = Scaffold(client, play)
		if err != nil {
			log.Println(err)
		}

		update(client)
		// TODO
		// Clear fields
		// reset play struct
		// Show Progress dialog (https://docs.fyne.io/api/v2/dialog/progressdialog/)
		// Show Success Dialog
	}

	return form
}
