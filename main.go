package main

import "os"
import (
	"fmt"
	"strconv"
	"./modifier"
)

func usageAndExit(exitCode int) {
	println("wrong arguments, usage: ")
	println("subsync FILENAME [+|-] NUMBER_OF_SECONDS")
	os.Exit(exitCode)
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		usageAndExit(1)
	}
	fileName := args[0]
	addOrSub := args[1]
	nbSecondsArg := args[2]
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Printf	("file named  %v does not exist\r\n", fileName)
		usageAndExit(2)
	}
	if addOrSub != "+" && addOrSub != "-" {
		fmt.Printf	("invalid operand %v, should be + or -\r\n", addOrSub)
		usageAndExit(3)
	}
	nbSeconds, err := strconv.Atoi(nbSecondsArg)
	if err != nil {
		fmt.Printf	("invalid number of seconds %v, should be an integer\r\n", nbSecondsArg)
		usageAndExit(4)
	}
	subModifier := modifier.NewModifier(fileName, addOrSub, nbSeconds)
	err = subModifier.Process()
	if err != nil {
		fmt.Printf("runtime error: %v", err)
		os.Exit(5)
	}
}
