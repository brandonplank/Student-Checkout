package main

import (
	"bufio"
	"fmt"
	"github.com/skip2/go-qrcode"
	"log"
	"os"
)

const outputPath = "codes"

func main() {
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		_ = os.Mkdir(outputPath, os.ModePerm)
	}
	students, err := os.Open("students.txt")
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(students)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		log.Println(fileScanner.Text())
		err := qrcode.WriteFile(fileScanner.Text(), qrcode.Medium, 256, fmt.Sprintf("%s/%s-code.png", outputPath, fileScanner.Text()))
		if err != nil {
			fmt.Printf("Couldn't create qrcode:,%v", err)
		}
	}

	students.Close()
}
