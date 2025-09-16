package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

)


var CONFIG_PATH = filepath.Join(FindHomePath(), ".screech.music.json")
var MusicFiles []string
var supportedTypes = []string{".mp3", ".aac", ".wav"}


func FindHomePath() string {
	var homepath, homepathErr = os.UserHomeDir()
	if homepathErr != nil {
		fmt.Println("Unable to locate the home dir! using current dir...")
		homepath = "."
	}

	return homepath
}

var musicMap = map[string][]string{
}

func FullScan(updateBodyChan chan string) bool {
	MusicFiles = []string{}
	var homepath = FindHomePath()

	scanDir(homepath, updateBodyChan)
	storeResponseInPath(CONFIG_PATH)
	

	return true
}

func scanDir(path string, updateBodyChan chan string) bool {
	var dir, err = os.ReadDir(path) 
	if err != nil {
		return false
	}

					updateBodyChan <- fmt.Sprintf("Scanning device for music...\nChecking path: %v\nSongs found: %v", path, len(MusicFiles))

	for _, file := range(dir) {
		var checkLinker = detectLinker(filepath.Join(path, file.Name())) == 2
		if file.IsDir() || checkLinker {
			scanDir(filepath.Join(path, file.Name()), updateBodyChan)
		}

		for _, mtype := range(supportedTypes) {
		if strings.HasSuffix(file.Name(), mtype) {
					MusicFiles = append(MusicFiles, filepath.Join(path, file.Name()))
					updateBodyChan <- fmt.Sprintf("Scanning device for music...\nChecking path: %v\nSongs found: %v", path, len(MusicFiles))
				}
		}	
	}

	return true
}

func storeResponseInPath(filePath string) {
	musicMap["paths"] = MusicFiles
	var mBytes, err = json.Marshal(musicMap)
	if err != nil {
		return
	}

	os.WriteFile(filePath, mBytes, os.ModePerm)

}

func detectLinker(path string) int {
 // not clean, i know, but it should work for  now
 // -1 false file
 // 0 for file 
 // 1 for directory 
 // 2 for linker 	

 var info, err = os.Lstat(path) 
 if err != nil {
	 return -1
 }

 if info.Mode()&os.ModeSymlink != 0 {
	 return 2
 } else if info.IsDir() {
	 return 1
 } 
 return 0
}
