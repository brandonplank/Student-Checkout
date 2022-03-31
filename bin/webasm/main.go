package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	apiUrl = "https://flappybird.brandonplank.org/v1/"
)

func main() {
	log.Println("Running GoLang Web ASM")
	log.Println("This module contains no code at the moment, however, this will eventually be used to optimise JS")
	log.Println("Module by Brandon Plank, Copyright 2022")
	for i := 0; i < 10; i++ {
		fmt.Print("*")
	}
	fmt.Println()
	fmt.Println("Testing call to my FlappyBird server")
	go func() {
		request, err := http.NewRequest("GET", apiUrl+"users", nil)
		if err != nil {
			log.Fatal(err)
		}
		timeout := 10 * time.Second
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Do(request)

		if resp.StatusCode == 401 {
			log.Fatal("You are not authorized to perform this action")
			os.Exit(-1)
		}

		if resp.StatusCode == 500 {
			log.Fatal("This was not supposed to happen, internal server error")
			os.Exit(-1)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(body))
	}()
	//document := js.Global().Get("document")
	//log.Println(document)
	<-make(chan bool)
}
