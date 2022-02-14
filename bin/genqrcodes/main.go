package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/ledongthuc/pdf"
	"github.com/skip2/go-qrcode"
	"log"
	"os"
	"regexp"
)

const codeOutputPath = "codes"
const cardOutputPath = "cards.pdf"

func readPdf(path string) (string, error) {
	var ret string
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				ret += word.S + "\n"
			}
		}
	}
	return ret, nil
}

func DoesSliceHaveName(name string, slice []string) bool {
	for _, obj := range slice {
		if obj == name {
			return true
		}
	}
	return false
}

func main() {
	if _, err := os.Stat(codeOutputPath); os.IsNotExist(err) {
		_ = os.Mkdir(codeOutputPath, os.ModePerm)
	}

	content, err := readPdf("report.pdf")
	if err != nil {
		log.Fatal(err)
	}

	// this was painful
	var re = regexp.MustCompile(`(?m)^([a-zA-Z\-]+)\s*,\s*([a-zA-Z]+)(\s+([a-zA-Z]+))?$`)

	pdf := gofpdf.New("P", "mm", "A4", "")

	var stored []string

	for _, match := range re.FindAllString(content, -1) {
		if !DoesSliceHaveName(match, stored) {
			qrCodeImg := fmt.Sprintf("%s/%s-code.png", codeOutputPath, match)
			err = qrcode.WriteFile(match, qrcode.Medium, 256, qrCodeImg)
			if err != nil {
				fmt.Printf("Couldn't create qrcode:,%v", err)
			} else {
				pdf.AddPage()
				pdf.SetFont("Arial", "B", 16)

				// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
				pdf.CellFormat(190, 7, match, "0", 0, "CM", false, 0, "")

				// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
				pdf.ImageOptions(
					qrCodeImg, 1, 1,
					0, 0,
					false,
					gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true},
					0,
					"",
				)
			}
			stored = append(stored, match)
		}
	}
	pdf.OutputFileAndClose(cardOutputPath)
}
