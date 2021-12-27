package pdf

import (
	"fmt"
	"os"
	"os/exec"
	"ping/extract-data/ocr"
	"ping/extract-data/tools"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	defaultOutput       = "/tmp/pdftotext"
	defaultOutputImages = "/tmp/pdfimages"
)

type pdf struct {
	inputPath       string
	outputPath      string
	outputImagePath string
	myOcr           ocr.OcrReader
}

type PdfData struct {
	Name     string            `json:"name"`
	Metadata map[string]string `json:"metadata"`
	Text     string            `json:"text"`
	Images   []PdfImage        `json:"images"`
}

type PdfImage struct {
	Path   string  `json:"path"`
	Page   int     `json:"page"`
	Num    int     `json:"num"`
	Type   string  `json:"type"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Color  string  `json:"color"`
	Comp   int     `json:"comp"`
	Bpc    int     `json:"bpc"`
	Enc    string  `json:"enc"`
	Interp string  `json:"interp"`
	Object int     `json:"object"`
	Id     int     `json:"id"`
	Xppi   int     `json:"xppi"`
	Yppi   int     `json:"yppi"`
	Size   float64 `json:"size"`
	Ratio  float64 `json:"ratio"`
	Text   string  `json:"text"`
}

func Init(path string, ocr ocr.OcrReader) *pdf {
	myPdf := pdf{
		path,
		defaultOutput + uuid.NewString(),
		defaultOutputImages + uuid.NewString(),
		ocr,
	}
	return &myPdf
}

func (p pdf) Extract() PdfData {
	chanText := p.getText()
	chanInfo := p.getInfo()
	chanImages := p.getImages()
	result := PdfData{
		p.inputPath,
		<-chanInfo,
		<-chanText,
		<-chanImages,
	}
	for i, image := range result.Images {
		text, err := p.myOcr.GetText(image.Path)
		if err != nil {
			fmt.Println(err)
		}
		result.Images[i].Text = text
	}

	return result
}

func (p pdf) getText() chan string {
	r := make(chan string)
	go func() {
		_, err := exec.Command("pdftotext", p.inputPath, p.outputPath).Output()
		if err != nil {
			r <- ""
			return
		}

		bytes, err := os.ReadFile(p.outputPath)
		if err != nil {
			r <- ""
			return
		}
		r <- string(bytes)
	}()
	return r
}

func (p pdf) getInfo() chan map[string]string {
	r := make(chan map[string]string)
	go func() {
		metadata := make(map[string]string)
		stdout, stderr, err := tools.ExecCommandWithOutput("pdfinfo", p.inputPath)
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

func (p pdf) getImages() chan []PdfImage {
	r := make(chan []PdfImage)
	go func() {
		var pdfImages []PdfImage
		cmdExtractImages := exec.Command("pdfimages", "-png", p.inputPath, p.outputImagePath)
		err := cmdExtractImages.Start()
		if err != nil {
			r <- nil
			return
		}

		stdout, _, err := p.getImagesInfo()
		if err != nil {
			r <- nil
			return
		}

		// return lines without header
		if string(stdout) == "" {
			r <- nil
			return
		}
		lines := strings.SplitAfter(string(stdout), "\n")[2:]

		for i, line := range lines {
			if line == "" {
				break
			}
			new_vals := strings.SplitAfter(line, " ")
			var vals []string
			// remove empty items and trim
			for _, val := range new_vals {
				if val != " " {
					vals = append(vals, strings.TrimSpace(val))
				}
			}
			path := p.outputImagePath + "-" + indexToStringFormatted(i) + ".png"
			page, _ := strconv.Atoi(vals[0])
			num, _ := strconv.Atoi(vals[1])
			width, _ := strconv.Atoi(vals[3])
			height, _ := strconv.Atoi(vals[4])
			comp, _ := strconv.Atoi(vals[6])
			bpc, _ := strconv.Atoi(vals[7])
			object, _ := strconv.Atoi(vals[10])
			id, _ := strconv.Atoi(vals[11])
			xppi, _ := strconv.Atoi(vals[12])
			yppi, _ := strconv.Atoi(vals[13])
			size, _ := strconv.ParseFloat(vals[14][:len(vals[14])-1], 64)
			ratio, _ := strconv.ParseFloat(vals[15][:len(vals[15])-1], 64)
			data := PdfImage{
				path,
				page,
				num,
				vals[2],
				width,
				height,
				vals[5],
				comp,
				bpc,
				vals[8],
				vals[9],
				object,
				id,
				xppi,
				yppi,
				size,
				ratio,
				"",
			}

			pdfImages = append(pdfImages, data)
		}
		err = cmdExtractImages.Wait()
		if err != nil {
			r <- nil
			return
		}

		r <- pdfImages
	}()
	return r
}

func (p pdf) getImagesInfo() ([]byte, []byte, error) {
	return tools.ExecCommandWithOutput("pdfimages", "-list", p.inputPath)
}

func indexToStringFormatted(i int) string {
	if i < 10 {
		return "00" + strconv.Itoa(i)
	} else if i < 100 {
		return "0" + strconv.Itoa(i)
	} else {
		return strconv.Itoa(i)
	}
}
