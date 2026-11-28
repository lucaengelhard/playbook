package main

import (
	"embed"
	"encoding/json/v2"
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"uuid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
)

//go:embed all:frontend/dist
var frontendFs embed.FS
var contents []*Content
var frontendPath = path.Join("frontend", "dist")

type Play struct {
	Title     string
	HeroImage string
	Contents  []*Content
	Gallery   []string `json:"Gallery,omitempty"`
}
type Content struct {
	Type    string
	Content string
}

func Scaffold(play *Play) error {
	if play.Title == "" {
		return errors.New("Cannot upload play with empty title :(")
	}

	dirPath := path.Join("plays", play.Title)

	err := os.MkdirAll(path.Join(dirPath, "media"), 0755)
	if err != nil {
		return err
	}
	entries, err := frontendFs.ReadDir(frontendPath)

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

	}
	if err != nil {
		return err
	}

	heroImage, err := os.ReadFile(play.HeroImage)
	if err != nil {
		return err
	}

	heroImagePath := NewMediaPath(play.HeroImage)
	play.HeroImage = heroImagePath
	err = os.WriteFile(path.Join(dirPath, heroImagePath), heroImage, 0755)
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
	if len(play.Gallery) > 0 {
		for i, fsPath := range play.Gallery {
			image, err := os.ReadFile(fsPath)
			if err != nil {
				log.Println(err)
				continue
			}

			mediaPath := NewMediaPath(fsPath)
			err = os.WriteFile(path.Join(dirPath, mediaPath), image, 0755)

			if err != nil {
				log.Println(err)
				continue
			}
			play.Gallery[i] = mediaPath
		}
	}

	content, err := json.Marshal(play)
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(dirPath, "content.json"), content, 0755)

}

func createEditor(play *Play, isNewPlay bool) *widget.Form {
	var form *widget.Form
	var entries []*widget.Entry

	//Create all standard widgets
	titleEntry := createTitleEntry(*play)

	heroImageButton := widget.NewButton("Select...", nil)
	galleryButton := widget.NewButton("Select Images...", nil)
	newTextButton := widget.NewButton("Content", nil)
	newPictureButton := widget.NewButton("Image", nil)

	//TODO
	//create a list with all the buttons in the form
	buttons := []*widget.Button{heroImageButton, galleryButton, newTextButton, newPictureButton}

	//Implementation of the button functions
	heroImageButton.OnTapped = func() { imageSelector(0, buttons, play, form) }
	pictureButtonfunc := func() { imageSelector(1, buttons, play, form) }
	galleryButton.OnTapped = func() { imageSelector(2, buttons, play, form) }

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

	//Put all items on the screen
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
	//Update the array with the new buttons
	buttons = GetAllButtons(form)

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
		for _, e := range entries {
			contents = append(contents, &Content{"Text", e.Text})
		}
		play.Contents = contents

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

func createTitleEntry(play Play) *widget.Entry {
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
	return titleEntry
}

func disableButtons(buttons []*widget.Button) {
	for _, b := range buttons {
		b.Disable()
	}
}
func enableButtons(buttons []*widget.Button) {
	for _, b := range buttons {
		b.Enable()
	}
}

func imageSelector(i int, buttons []*widget.Button, play *Play, form *widget.Form) {
	// i = 0: HeroImage, i = 1: image in text, i = 3 gallery

	imageTypeFilter := zenity.FileFilters{
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
	}

	disableButtons(buttons)

	go func() {
		defer fyne.Do(func() { enableButtons(buttons) })
		var err error
		var file string
		var files []string
		if i == 2 {
			files, err = zenity.SelectFileMultiple(imageTypeFilter)
		} else {
			file, err = zenity.SelectFile(imageTypeFilter)
		}
		if err != nil {
			// TODO
			return
		}

		if i == 0 {
			play.HeroImage = file
		}
		if i == 1 {
			contents = append(contents, &Content{"Picture", file})
		}
		if i == 2 {
			play.Gallery = files
			fyne.Do(func() {
				buttons[2].Text = strings.Join(play.Gallery, "\n")
			})
		}
		fyne.Do(form.Refresh)
	}()
}

func GetAllButtons(form *widget.Form) []*widget.Button {
	var buttons []*widget.Button
	var traverse func(obj fyne.CanvasObject)

	traverse = func(obj fyne.CanvasObject) {
		if button, ok := obj.(*widget.Button); ok {
			buttons = append(buttons, button)
		}
		if container, ok := obj.(*fyne.Container); ok {
			for _, child := range container.Objects {
				traverse(child)
			}
		}
	}

	// Traverse form fields
	for _, item := range form.Items {
		traverse(item.Widget)
	}

	return buttons
}

func NewMediaPath(p string) string {
	name := path.Base(p)
	return path.Join("media", uuid.New().String()+"_"+name)
}
