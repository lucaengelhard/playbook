package main

import (
	"embed"
	"encoding/json/v2"
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"fyne.io/fyne/v2"
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

	err = uploadFrontendFiles(entries, dirPath, nil)
	if err != nil {
		return err
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

	err = uploadFrontendFiles(entries, dirPath, client)
	if err != nil {
		return err
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

func createEditor(client *Client, play Play) *widget.Form {
	var form *widget.Form

	localUploadCheckbox := widget.NewCheck("Stores directory locally", func(b bool) {
		GlobalSettings.Update(func(s Settings) Settings {
			s.LocalUpload = b
			return s
		})
	})
	localUploadCheckbox.Checked = GlobalSettings.Data.LocalUpload
	titleEntry := widget.NewEntry()
	titleEntry.Text = play.Title
	titleEntry.Validator = func(value string) error {
		title := strings.TrimSpace(value)

		if title == "" {
			return errors.New("Title is required")
		}
		if title == "." || title == ".." || strings.ContainsAny(title, "/\\") {
			return errors.New("Invalid title")
		}
		return nil
	}

	heroImageButton := widget.NewButton("Select...", nil)
	galleryButton := widget.NewButton("Select Images...", nil)

	buttons := []*widget.Button{heroImageButton, galleryButton}

	disableButtons := func() {
		for _, b := range buttons {
			b.Disable()
		}
	}
	enableButtons := func() {
		for _, b := range buttons {
			b.Enable()
		}
	}

	heroImageButton.OnTapped = func() {
		{
			disableButtons()
			go func() {
				defer fyne.Do(enableButtons)

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
				fyne.Do(form.Refresh)
			}()
		}
	}

	galleryButton.OnTapped = func() {
		disableButtons()

		go func() {
			defer fyne.Do(enableButtons)
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
		}()
	}

	content := widget.NewMultiLineEntry()
	content.Text = play.Content

	form = widget.NewForm(
		widget.NewFormItem("Local Upload", localUploadCheckbox),
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Content", content),
		widget.NewFormItem("Hero Image", heroImageButton),
		widget.NewFormItem("Gallery", galleryButton),
	)
	// check if hero image exists
	form.Validator = func() error {
		if play.HeroImage == "" {
			return errors.New("Hero image is required")
		}
		info, err := os.Stat(play.HeroImage)
		if err != nil || !info.Mode().IsRegular() {
			return errors.New("Invalid hero image file")
		}
		return nil
	}

	form.OnSubmit = func() {
		play.Title = titleEntry.Text
		play.Content = content.Text

		var err error
		if GlobalSettings.Data.LocalUpload {
			err = ScaffoldOffline(play)
		} else {
			err = Scaffold(client, play)
		}
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
