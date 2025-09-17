package main

import (
	"fmt"
	"os"
	"path/filepath"
	"screech/cli"
	"screech/helper"
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
	lyrics string
	artist string
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


	var bodyList = tview.NewList()

	var body = tview.NewFlex().SetDirection(tview.FlexRow)
	body.AddItem(bodyTitle, 7, 0, false) 
	body.AddItem(bodyList, 0, 2, true)


		DeepSearchDevice(app, bodyTitle, false)
		var footer = tview.NewTextView(). 
								SetText("[p]Play [u]Pause [s]Stop \n[^f]Deep Search [o]Open"). 
								SetTextAlign(tview.AlignCenter)

	var layout = tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(header, 4, 0, false)
	layout.AddItem(body, 0, 3, true)
	layout.AddItem(footer, 2, 0, false)



	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			cli.StopAudio()
			cli.RemoveNowPlayingNotification()
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
						ShowMusicMenu(app, bodyTitle, bodyList, &currentSong)
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



func OpenSong(title *tview.TextView, song *Song, app *tview.Application) {
	if (len(song.tempImage) > 0) {
		cleanTempFile(song)
	}

	getMetadataFromFile(song)

		
	var formattedString = fmt.Sprintf("Song: %s\nArtist: %s\nstate: playing", song.title, song.artist)

	title.SetText(formattedString)
	go cli.PlayAudio(song.path)
	go cli.NowPlayingNotification(song.title, song.tempImage)
}

func PlaySong(title *tview.TextView, song *Song) {
	if (len(song.title) > 0) {


	var formattedString = fmt.Sprintf("Song: %s\nArtist: %s\nstate: playing", song.title, song.artist)

			title.SetText(formattedString)
			cli.ResumeAudio()
	} else {
		title.SetText("No song is currently open!")
	}
}



func PauseSong(title *tview.TextView, song *Song) {

	var formattedString = fmt.Sprintf("Song: %s\nArtist: %s\nstate: paused", song.title, song.artist)
	title.SetText(formattedString)
	cli.PauseAudio()
}


func StopSong(title *tview.TextView, song *Song) {

	var formattedString = "Use [o] to select a song!"
	title.SetText(formattedString)
	cli.StopAudio()
	song.artist = ""
	song.title = ""
	song.state = "off" 

	cli.RemoveNowPlayingNotification()
	cleanTempFile(song)
}




func ShowMusicMenu(app *tview.Application, title *tview.TextView, list *tview.List, song *Song) {
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
				OpenSong(title, song, app)
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

	var pic = metadata.Picture()
	if pic != nil {
		var tempImageBytes = pic.Data
		var trimTitle = strings.ReplaceAll(song.title, " ", "")
		trimTitle = strings.ReplaceAll(trimTitle, "\uFEFF", "")
			var tempImagePath = filepath.Join(scanner.FindHomePath(), "." + trimTitle + "." + mimeToExt(pic.MIMEType))

			if mimeToExt(pic.MIMEType) == "jpg" {
				var newPath =  filepath.Join(scanner.FindHomePath(), "." + strings.TrimSpace(metadata.Title()) + "." + "png")
				if err := helper.ConvertJpgToPng(tempImageBytes, newPath); err == nil {
					tempImagePath = newPath
				}
			}

			song.lyrics = metadata.Lyrics()
			song.tempImage = tempImagePath
			song.tempImgBytes = tempImageBytes
			song.artist = metadata.Artist()

			// err3 := os.WriteFile(tempImagePath, tempImageBytes, os.ModePerm)
			// if err3 != nil {
				// return
			// }

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

