package main

import (
	"RuntimeError/server"
	"fmt"
)

func main() {
	fmt.Println("Hello world")

	s := server.NewServer()
	s.Init()
	s.Run()
}