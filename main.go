package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tompston/gut"
)

func main() {
	applyFlags()
	if GlobalSettings.Generate {
		gut.Generate(
			"./frontend/types.gen.ts",
			gut.Convert(Play{}),
			gut.Settings{
				FirstLine: "// DONT EDIT THIS FILE! It is auto-generated",
				DateType:  "Date",
			},
		)
		return
	}

	a := app.New()
	w := a.NewWindow("Hello World")

	client := NewClient(w)
	defer client.Close()

	editor := createEditor(client)

	w.SetContent(container.NewVBox(widget.NewButton("Connect", func() {
		client.dialog.Show()
	}), editor))
	w.ShowAndRun()
}
