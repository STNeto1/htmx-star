package main

import (
	"fmt"
)

func main() {
	container := NewContainer(25, 25)

	for true {
		// container.Print()
		// fmt.Println("===")

		if err := container.Tick(); err != nil {
			fmt.Println(err.Error())
			break
		}

		// time.Sleep(time.Millisecond * 500)
	}

	container.Print()
}
