package scanner

import (
	"encoding/json"
	"os"
)

func LoadMusic(path string) []string {
	var fileBytes, err = os.ReadFile(path) 
	if err != nil {
		return []string{}
	}

	var response map[string][]string
	if err := json.Unmarshal(fileBytes, &response); err != nil {
		return []string{}
	}

	if music, ok := response["paths"]; ok {
		return music
	}

	return []string{}
}
