package main

import (
	"embed"
	"encoding/json/v2"
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
)

//go:embed all:frontend
var frontendFs embed.FS
var contents []Content

type Play struct {
	Title     string
	HeroImage string
	Contents  []Content
	Gallery   []string `json:"Gallery,omitempty"`
}
type Content struct {
	Type    string
	Content string
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
		err = os.WriteFile(path.Join(dirPath, entry.Name()), file, 0755)

	}
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

	for _, c := range play.Contents {
		if c.Type == "Picture" {
			v, err := os.ReadFile(c.Content)
			if err != nil {
				return err
			}

			c.Content = path.Join("media", path.Base(c.Content))
			err = os.WriteFile(path.Join(dirPath, c.Content), v, 0755)
			if err != nil {
				return err
			}
		}
	}

	content, err := json.Marshal(play)
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(dirPath, "content.json"), content, 0755)

}

func createEditor(play Play) *widget.Form {
	var form *widget.Form
	var entries []*widget.Entry
	//dirPath := path.Join("plays", play.Title)

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
	newTextButton := widget.NewButton("Content", nil)
	newPictureButton := widget.NewButton("Image", nil)

	buttons := []*widget.Button{heroImageButton, galleryButton, newTextButton, newPictureButton}
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

	pictureButtonfunc := func() {

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

			contents = append(contents, Content{"Picture", file})
			fyne.Do(form.Refresh)
		}()

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

	newTextButton.OnTapped = func() {
		var c = widget.NewMultiLineEntry()
		entries = append(entries, c)
		form.Append("Content", c)
		form.Refresh()
	}
	newPictureButton.OnTapped = func() {
		var b = widget.NewButton("Image", pictureButtonfunc)
		buttons = append(buttons, b)
		form.Append("Picture", b)
		form.Refresh()
	}
	b := container.NewHBox(
		newTextButton,
		newPictureButton,
	)
	form = widget.NewForm(
		widget.NewFormItem("New", b),
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Hero Image", heroImageButton),
		widget.NewFormItem("Gallery", galleryButton),
	)

	for _, c := range play.Contents {
		if c.Type == "Text" {
			var v = widget.NewMultiLineEntry()
			v.Text = c.Content
			form.Append("Content", v)
		}
		if c.Type == "Picture" {
			var b = widget.NewButton("Picture", pictureButtonfunc)
			form.Append("Picture", b)
		}
	}

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
		//play.Contents = content.Text

		for _, e := range entries {
			contents = append(contents, Content{"Text", e.Text})
		}
		for _, c := range contents {
			play.Contents = append(play.Contents, c)
		}

		var err error
		err = Scaffold(play)

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
