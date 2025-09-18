package cli

import (
	"fmt"
	"os/exec"
	"strings"
)


func PlayAudio(path string) {
	exec.Command("termux-media-player", "play", path).Run()
}

func ResumeAudio() {
	exec.Command("termux-media-player", "play").Run()
}

func PauseAudio() {
	exec.Command("termux-media-player", "pause").Run()
}

func StopAudio(){
	exec.Command("termux-media-player", "stop").Run()
}

func NowPlayingNotification(songName string, coverPath string) {
    nowPlaying := fmt.Sprintf("Now playing: %v", songName)
    cmd := exec.Command(
        "termux-notification",
        "--id", "now_playing",
        "--title", nowPlaying,
        "--type", "media",
        "--media-play", "termux-media-player play",
        "--media-pause", "termux-media-player pause",
        "--media-next", "termux-media-player next",
        "--media-previous", "termux-media-player prev",
        "--image-path", coverPath, 
        "--ongoing",
    )

    cmd.CombinedOutput()
}
	func RemoveNowPlayingNotification () {
		exec.Command("termux-notification-remove", "now_playing").Run()
	}


// the struct and functio n are ai generated 
// ie MediaInfo and GetMediaInfo 
// sorry i got lazy 

func main() {

}

// MediaInfo represents the track info
type MediaInfo struct {
	Status   string `json:"status"`
	Track    string `json:"track"`
	Position string `json:"position"`
	Duration string `json:"duration"`
}

// GetMediaInfo runs `termux-media-player info` and parses it
func GetMediaInfo() (*MediaInfo, error) {
	cmd := exec.Command("termux-media-player", "info")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	info := &MediaInfo{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "Status:"):
			info.Status = strings.TrimSpace(strings.TrimPrefix(line, "Status:"))
		case strings.HasPrefix(line, "Track:"):
			info.Track = strings.TrimSpace(strings.TrimPrefix(line, "Track:"))
		case strings.HasPrefix(line, "Current Position:"):
			parts := strings.Split(line, " ")
			if len(parts) >= 5 {
				info.Position = parts[2] // current time
				info.Duration = parts[4] // total duration
			}
		}
	}

	return info, nil
}
