package main

import (
	"fmt"
	"io"
	"os"
)

func checkDoesWantToPlay(reader io.Reader) func(int) {
	var play string
	for {
		fmt.Println("Do you want to play a game?")
		fmt.Fscanf(reader, "%s", &play)
		switch play {
		case "y", "Y", "yes", "Yes":
			fmt.Println("Ok, great!")
			return nil
		case "n", "N", "no", "No":
			fmt.Println("Ok, maybe later.")
			return func(i int) {
				os.Exit(i)
			}
		default:
			fmt.Println("You must enter yes or no")
		}
	}
}
