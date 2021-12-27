package ocr

type OcrReader interface {
	GetText(string) (string, error)
}
