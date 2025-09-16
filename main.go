package main

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"screech/scanner"
	"strings"

	"github.com/dhowden/tag"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Song struct {
	path string 
	title string
	state string
	tempImage string 
	tempImgBytes []byte
}
func main() {
	var currentSong Song
	var app = tview.NewApplication()
	var header = tview.NewTextView().
										SetText("Screech 0.1").
										SetTextAlign(tview.AlignCenter)
	
	var bodyTitle = tview.NewTextView(). 
									SetText("Body!").
									SetTextAlign(tview.AlignCenter).
									SetScrollable(true)

	var bodyImage = tview.NewImage()

	var bodyList = tview.NewList()

	var body = tview.NewFlex().SetDirection(tview.FlexRow)
	body.AddItem(bodyTitle, 3, 0, false) 
	body.AddItem(bodyImage, 0, 1, false)
	body.AddItem(bodyList, 0, 2, true)


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
		case tcell.KeyCtrlC: 
			exec.Command("termux-media-player", "stop").Run()
		if (len(currentSong.tempImage) > 0) {
				cleanTempFile(&currentSong)
			}
			app.Stop()
			case tcell.KeyCtrlF:
				bodyList.Clear()
				DeepSearchDevice(app, bodyTitle, true)
			case tcell.KeyRune:
				switch event.Rune() {
					case 'o':
						ShowMusicMenu(app, bodyTitle, bodyList, &currentSong, bodyImage)
					case 'u':
						PauseSong(bodyTitle, &currentSong)
					case 's':
						StopSong(bodyTitle, &currentSong)
					case 'p':
						PlaySong(bodyTitle, &currentSong)
				}
		}

		return event
	})
	
	
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}

func OpenSong(title *tview.TextView, imageWidget *tview.Image, song *Song, app *tview.Application) {
	if (len(song.tempImage) > 0) {
		cleanTempFile(song)
	}

	getMetadataFromFile(song)
	img, _, err := image.Decode(bytes.NewReader(song.tempImgBytes))
	if err == nil {
	app.QueueUpdateDraw(func() {
	imageWidget.SetImage(img)
})
	}
	
	var formattedString = fmt.Sprintf("%s\n state: playing", song.title)

	title.SetText(formattedString)
	exec.Command("termux-media-player", "play", song.path).Run()
}

func PlaySong(title *tview.TextView, song *Song) {
	if (len(song.title) > 0) {
		var formattedString = fmt.Sprintf("%s\n state: playing", song.title)

			title.SetText(formattedString)
			exec.Command("termux-media-player", "play").Run()

	} else {
		title.SetText("No song is currently open!")
	}
}



func PauseSong(title *tview.TextView, song *Song) {
	var formattedString = fmt.Sprintf("%s\n state: paused", song.title)

	title.SetText(formattedString)
	exec.Command("termux-media-player", "pause").Run()
}


func StopSong(title *tview.TextView, song *Song) {
	var formattedString = fmt.Sprintf("%s\n state: stopped", song.title)

	title.SetText(formattedString)
	exec.Command("termux-media-player", "stop").Run()
}




func ShowMusicMenu(app *tview.Application, title *tview.TextView, list *tview.List, song *Song, imgWidget *tview.Image) {
	var loadMusic = scanner.LoadMusic(scanner.CONFIG_PATH)
	
		title.SetText("Select song: ")
		list.Clear()
		for _, path := range loadMusic {
			var relativeP = strings.Split(path, "/")
			var relative = relativeP[len(relativeP) - 1]

			list.AddItem(relative, "", 0, func() {
				list.Clear()
				song.title = relative
				song.path = path
				OpenSong(title, imgWidget, song, app)
			})
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

func getMetadataFromFile(song *Song) {
	file, err := os.Open(song.path) 
	if err != nil {
		return
	}

	metadata, err2 := tag.ReadFrom(file)
	if err2 != nil {
		return
	}

	var tempImageBytes = metadata.Picture().Data
	var tempImagePath = filepath.Join(scanner.FindHomePath(), "." + strings.TrimSpace(metadata.Title()) + "." + mimeToExt(metadata.Picture().MIMEType))
	song.tempImage = tempImagePath
	song.tempImgBytes = tempImageBytes

	err3 := os.WriteFile(tempImagePath, tempImageBytes, os.ModePerm)
	if err3 != nil {
		return
	}

}

func cleanTempFile(song *Song) {
	os.Remove(song.tempImage)
}

func mimeToExt(mime string) string {
	switch strings.ToLower(mime) {
	case "image/jpeg", "image/jpg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	case "image/bmp":
		return "bmp"
	case "image/tiff":
		return "tiff"
	default:
		return "bin" 
	}
}

