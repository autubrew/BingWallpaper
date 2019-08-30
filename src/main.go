package main

import (
	"ui"
)

func main() {
	err := Init()
	if err != nil {
		panic(err)
	}

	go UpdateTask()

	ui.Run()
}
