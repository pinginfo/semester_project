package ocr

import (
	"github.com/otiai10/gosseract/v2"
)

type Gosseract struct {
	client *gosseract.Client
}

func GosseractInit() Gosseract {
	gosse := Gosseract{
		gosseract.NewClient(),
	}

	return gosse
}

func (g Gosseract) GetText() (string, error) {
	out, err := g.client.Text()
	if err != nil {
		return "", err
	}
	return out, nil
}

func (g Gosseract) SetImage(path string) error {
	g.client.SetImage(path)
	return nil
}
