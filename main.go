package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ping/extract-data/image"
	"ping/extract-data/ocr"
	"ping/extract-data/pdf"
	"ping/extract-data/tools"
	"time"

	"github.com/google/uuid"
)

type ApiResponse struct {
	Timestamp time.Time     `json:"timestamp"`
	Data      []interface{} `json:"data"`
}

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HandleError(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		apiResponse := ApiError{
			http.StatusBadRequest,
			err.Error(),
		}
		b, _ := json.Marshal(apiResponse)
		io.WriteString(w, string(b))
		return true
	}
	return false
}

func extractHandler(w http.ResponseWriter, r *http.Request) {
	apiResponse := ApiResponse{
		time.Now(),
		nil,
	}

	r.ParseMultipartForm(32 << 20)

	for _, files := range r.MultipartForm.File {
		for _, val := range files {
			filename := "/tmp/uploaded_file_" + uuid.NewString()
			file, err := val.Open()
			if HandleError(w, err) {
				return
			}
			bytes := make([]byte, val.Size)
			c, err := file.Read(bytes)
			if HandleError(w, err) {
				return
			}
			err = tools.SaveInfile(filename, bytes[:c])
			if HandleError(w, err) {
				return
			}

			if pdf.IsAValidPdf(filename) {
				myPdf := pdf.Init(filename, ocr.TesseractInit([]string{"fra", "deu", "ita", "eng"}))
				values := myPdf.Extract()
				apiResponse.Data = append(apiResponse.Data, values)
			} else if image.IsAValidImage(filename) {
				myImage := image.Init(filename, ocr.TesseractInit([]string{"fra", "deu", "ita", "eng"}))
				values := myImage.Extract()
				apiResponse.Data = append(apiResponse.Data, values)
			} else {
				fmt.Println("unsupported format")
			}

		}
	}

	b, _ := json.Marshal(apiResponse)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(b))
}

func main() {
	http.HandleFunc("/", extractHandler)

	http.ListenAndServe(":8080", nil)
}
