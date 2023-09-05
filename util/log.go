package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Log logs into fmt stdout
func Log(tag, text string) {
	fmt.Printf("%s\r\r", tag)

	fmt.Println(strings.ReplaceAll(text, "\n", "\r"))
}

// LogJson logs json object
func LogJson(tag string, val any) {
	bytes, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		Log(tag, fmt.Sprintf("%#v", val))
		return
	}

	Log(tag, string(bytes))
}
