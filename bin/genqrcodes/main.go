package main

import (
	"fmt"
	"github.com/ledongthuc/pdf"
	"github.com/skip2/go-qrcode"
	"log"
	"os"
	"regexp"
)

const outputPath = "codes"

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

func main() {
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		_ = os.Mkdir(outputPath, os.ModePerm)
	}

	content, err := readPdf("report.pdf")
	if err != nil {
		log.Fatal(err)
	}

	// this was painful
	var re = regexp.MustCompile(`(?m)^([a-zA-Z\-]+)\s*,\s*([a-zA-Z]+)(\s+([a-zA-Z]+))?$`)

	for _, match := range re.FindAllString(content, -1) {
		err = qrcode.WriteFile(match, qrcode.Medium, 256, fmt.Sprintf("%s/%s-code.png", outputPath, match))
		if err != nil {
			fmt.Printf("Couldn't create qrcode:,%v", err)
		}
	}
}
