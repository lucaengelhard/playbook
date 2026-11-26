package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tompston/gut"
)

func main() {
	applyFlags()

	if GlobalSettings.Data.Generate {
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

	var width float32 = GlobalSettings.Data.WindowSize.first
	var height float32 = GlobalSettings.Data.WindowSize.second
	w.Resize(fyne.NewSize(width, height))

	var nav *container.Navigation

	playList := container.NewVBox()

	State.OnUpdate(func(p []Play) {
		playList.Objects = []fyne.CanvasObject{}
		for _, play := range p {
			playList.Objects = append(playList.Objects, widget.NewButton(play.Title, func() {
				nav.PushWithTitle(createEditor(&play), play.Title)
			}))
		}
	})

	sidebar := container.NewVScroll(
		container.NewVBox(
			playList,
			widget.NewButton("New", func() {
				nav.PushWithTitle(createEditor(&Play{}), "New Play")
			}),
		))

	content := container.NewVBox(widget.NewLabel("Select a Play :)"))
	nav = container.NewNavigationWithTitle(content, "Home")

	split := container.NewHSplit(sidebar, nav)
	split.SetOffset(0.2)

	err := update()
	if err != nil {
		log.Println(err)
	}

	w.SetContent(split)
	w.ShowAndRun()
}
