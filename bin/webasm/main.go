package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("Running GoLang Web ASM")
	log.Println("This module contains no code at the moment, however, this will eventually be used to optimise JS")
	log.Println("Module by Brandon Plank, Copyright 2022")
	for i := 0; i < 10; i++ {
		fmt.Print("*")
	}
	<-make(chan bool)
}
