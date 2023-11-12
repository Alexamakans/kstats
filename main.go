package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const windowTitle = "K-Stats"
const width = 1000
const height = 640

func main() {
	println("creating app")
	a := app.New()
	font, err := fyne.LoadResourceFromPath("MonospaceTypewriter.ttf")
	if err != nil {
		panic(err)
	}
	a.Settings().SetTheme(&kTheme{
		font: font,
	})
	println("creating window")
	w := a.NewWindow(windowTitle)
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(width, height))
	println("creating game")
	g, err := newGame("words.txt")
	if err != nil {
		panic(err)
	}
	w.Canvas().SetOnTypedRune(g.onTypedRune)

	cont := container.NewVBox()

	cont.Add(g.wordDisplay.first)
	cont.Add(g.wordDisplay.second)
	cont.Add(g.wordDisplay.third)
	cont.Add(g.statsDisplay)
	w.SetContent(cont)

	// go g.updateStatsLoop(statsDisplay)

	println("showing and running")
	w.ShowAndRun()
}
