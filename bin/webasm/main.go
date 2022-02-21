package main

import (
	"fmt"
	"log"
)

func main() {
	ch := make(chan int)
	log.Println("Running GoLang Web ASM")
	fmt.Println("This module contains no code at the moment, however, this will eventually be used to optimise JS")

	log.Println("Module by Brandon Plank, Copyright 2022")
	<-ch
}
