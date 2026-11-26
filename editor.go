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

//go:embed all:frontend/dist
var frontendFs embed.FS
var frontendPath = path.Join("frontend", "dist")

type Play struct {
	Title     string
	HeroImage string
	Content   string
	Gallery   []string `json:"Gallery,omitempty"`
}

func uploadFrontendFiles(entries []fs.DirEntry, dirPath string) error {
	for _, entry := range entries {
		if entry.IsDir() {
			// TODO
			continue
		}
		filePath := path.Join(frontendPath, entry.Name())
		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("Error while reading template file", err)
			continue
		}

		err = os.WriteFile(path.Join(dirPath, entry.Name()), file, 0755)

		if err != nil {
			return err
		}
	}
	return nil
}

func Scaffold(play Play) error {
	if play.Title == "" {
		return errors.New("Cannot upload play with empty title :(")
	}

	dirPath := path.Join("plays", play.Title)

	err := os.MkdirAll(path.Join(dirPath, "media"), 0755)
	if err != nil {
		return err
	}
	entries, err := frontendFs.ReadDir(frontendPath)

	err = uploadFrontendFiles(entries, dirPath)
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

func createEditor(play Play, isNewPlay bool) *widget.Form {
	var form *widget.Form

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

		var heroImagePath string
		if isNewPlay {
			heroImagePath = play.HeroImage
		} else {
			heroImagePath = path.Join("plays", play.Title, play.HeroImage)
		}

		info, err := os.Stat(heroImagePath)
		if err != nil || !info.Mode().IsRegular() {
			return errors.New("Invalid hero image file")
		}
		return nil
	}

	form.OnSubmit = func() {
		play.Title = titleEntry.Text
		play.Content = content.Text

		err := Scaffold(play)
		if err != nil {
			log.Println(err)
		}

		update()
		// TODO
		// Clear fields
		// reset play struct
		// Show Progress dialog (https://docs.fyne.io/api/v2/dialog/progressdialog/)
		// Show Success Dialog
	}

	return form
}
