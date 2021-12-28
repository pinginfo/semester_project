package image

import (
	"fmt"
	"ping/extract-data/ocr"
	"ping/extract-data/tools"
	"strings"
)

type image struct {
	inputFile string
	myOcr     ocr.OcrReader
}

type ImageData struct {
	Metadata map[string]string `json:"metadata"`
	Text     string            `json:"text"`
}

func Init(path string, ocr ocr.OcrReader) *image {
	return &image{
		path,
		ocr,
	}
}

func IsAValidImage(path string) bool {
	_, stderr, err := tools.ExecCommandWithOutput("exiftool", path)
	if err != nil || len(stderr) != 0 {
		return false
	}
	return true
}

func (i image) Extract() ImageData {
	chanInfo := i.getInfo()
	text, err := i.myOcr.GetText(i.inputFile)
	if err != nil {
		fmt.Println(err)
	}
	return ImageData{
		<-chanInfo,
		text,
	}
}

func (i image) getInfo() chan map[string]string {
	r := make(chan map[string]string)
	go func() {
		metadata := make(map[string]string)
		stdout, stderr, err := tools.ExecCommandWithOutput("exiftool", i.inputFile)
		if err != nil {
			r <- metadata
			return
		}
		if string(stderr) != "" {
			r <- metadata
			return
		}

		lines := strings.Split(string(stdout), "\n")
		for _, line := range lines {
			if line == "" {
				break
			}
			hashmap := strings.Split(line, ":")
			keyLowerCamelCase := strings.ToLower(hashmap[0][:1]) + hashmap[0][1:]
			// metadata[strings.TrimSpace(keyLowerCamelCase)] = strings.TrimSpace(hashmap[1])
			res := "" // TOOD: WHy strings.TrimSpace or strings.Trim doesn't works ?
			for _, t := range strings.Split(keyLowerCamelCase, " ") {
				res += t
			}
			metadata[res] = strings.TrimSpace(hashmap[1])
		}
		r <- metadata
	}()
	return r
}
