package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	window := a.NewWindow("SQLens")

	queryEntry := widget.NewMultiLineEntry()
	queryEntry.SetPlaceHolder("Enter query here")

	executeButton := widget.NewButton("Execute", func() {
		query := queryEntry.Text
		fmt.Println(query)
	})

	window.SetContent(container.NewVBox(
		queryEntry,
		executeButton,
	))

	window.ShowAndRun()
}
