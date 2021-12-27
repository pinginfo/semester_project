package ocr

import (
	"fmt"
	"os"
	"ping/extract-data/tools"
)

type TesseractCLI struct {
	langs []string
}

func TesseractInit(langs []string) TesseractCLI {
	return TesseractCLI{
		langs,
	}
}

func (c TesseractCLI) GetText(imagePath string) (string, error) {
	stdout, stderr, err := tools.ExecCommandWithOutput("tesseract", imagePath, imagePath, "-l", c.langsParsed())
	if err != nil {
		return "", err
	}
	fmt.Println("stdout: ", string(stdout))
	fmt.Println("stderr: ", string(stderr))

	bytes, err := os.ReadFile(imagePath + ".txt")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c TesseractCLI) langsParsed() string {
	str := ""
	size := len(c.langs)
	for i, lang := range c.langs {
		str += lang
		if i != size-1 {
			str += "+"
		}
	}
	return str
}
