package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	client := &Client{}
	//defer client.Close()

	hostEntry := widget.NewEntry()
	portEntry := widget.NewEntry()
	userEntry := widget.NewEntry()
	pwdEntry := widget.NewEntry()

	connect_form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Host", Widget: hostEntry},
			{Text: "Port", Widget: portEntry},
			{Text: "User", Widget: userEntry},
			{Text: "Password", Widget: pwdEntry},
		},
		OnSubmit: func() {
			client.Connect(hostEntry.Text, portEntry.Text, userEntry.Text, pwdEntry.Text)
		},
	}

	titleEntry := widget.NewEntry()
	content := widget.NewMultiLineEntry()
	var heroImagePath string

	heroSelector := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		heroImagePath = reader.URI().Path()
	}, w)
	heroImageButton := widget.NewButton("Select...", func() {
		heroSelector.Show()
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: titleEntry},
			{Text: "Content", Widget: content},
			{Widget: heroImageButton},
		},
		OnSubmit: func() {
			Scaffold(client, Play{
				Title:     titleEntry.Text,
				HeroImage: heroImagePath,
				Content:   content.Text,
			})

			// TODO upload bytes & template to server
		},
	}

	w.SetContent(container.NewVBox(connect_form, form))
	w.ShowAndRun()
}
