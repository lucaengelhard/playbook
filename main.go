package main

import (
	"flag"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tompston/gut"
)

func main() {
	// TODO
	// Flags/Settings global struct?
	// Add Dev mode flag
	generateFlag := flag.Bool("gen", false, "Generate typescript types to ./frontend")
	flag.Parse()

	if *generateFlag {
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
