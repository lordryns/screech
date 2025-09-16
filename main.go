package main

import (
	"fmt"
	"os"
	"screech/scanner"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	var app = tview.NewApplication()

	var header = tview.NewTextView().
										SetText("Screech 0.1").
										SetTextAlign(tview.AlignCenter)
	
	var bodyTitle = tview.NewTextView(). 
									SetText("Body!").
									SetTextAlign(tview.AlignCenter).
									SetScrollable(true)

	var bodyList = tview.NewList()

	var body = tview.NewFlex().SetDirection(tview.FlexRow)
	body.AddItem(bodyTitle, 3, 0, false) 
	body.AddItem(bodyList, 0, 1, true)


		DeepSearchDevice(app, bodyTitle, false)
		var footer = tview.NewTextView(). 
								SetText("[p]Play [u]Pause [s]Stop \n[^f]Deep Search [o]Open"). 
								SetTextAlign(tview.AlignCenter)

	var layout = tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(header, 3, 0, false)
	layout.AddItem(body, 0, 3, true)
	layout.AddItem(footer, 2, 0, false)


	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
			case tcell.KeyCtrlF:
				bodyList.Clear()
				DeepSearchDevice(app, bodyTitle, true)
			case tcell.KeyRune:
				switch event.Rune() {
					case 'o':
						ShowMusicMenu(app, bodyTitle, bodyList)
				}
		}

		return event
	})
	
	
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}

func OpenSong(title *tview.TextView, songName string, fullPath string) {

}

func ShowMusicMenu(app *tview.Application, title *tview.TextView, list *tview.List) {
	var loadMusic = scanner.LoadMusic(scanner.CONFIG_PATH)
	
		title.SetText("Select song: ")
		list.Clear()
		for _, path := range loadMusic {
			var relativeP = strings.Split(path, "/")
			var relative = relativeP[len(relativeP) - 1]

			list.AddItem(relative, "", 0, nil)
		}
	}

func DeepSearchDevice(app *tview.Application, body *tview.TextView, force bool) {	
	if  _, err := os.Stat(scanner.CONFIG_PATH); err == nil && !force {
		body.SetText("Use [o] to select a song!")
		return
	}

	updateBodyChan := make(chan string, 10)
	go func () {
		for msg := range updateBodyChan {
 			app.QueueUpdateDraw(
			 func() {
				body.SetText(msg)
			 })

		}		
	}()

	go func() {	
		if scanner.FullScan(updateBodyChan) {
			updateBodyChan <- fmt.Sprintf("Scan finished with %v songs found!", len(scanner.MusicFiles))
			close(updateBodyChan)
		}
	}()
}
