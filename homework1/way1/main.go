package main

import (
	"fmt"
	"strings"
)

type Args = map[string]int

func getLine(numLine int, char int, size int) string {
	s := make([]string, size)
	for i := 0; i < size; i++ {
		if (i == numLine) || (i == size-1-numLine) {
			s[i] = string(rune(char))
		} else {
			s[i] = " "
		}
	}
	return strings.Join(s, "")
}

func drawSandglass(char int, size int, color int) {
	edgeLine := strings.Repeat(string(char), size)
	fmt.Printf("\033[%dm", color)
	fmt.Println(edgeLine)
	for i := 1; i < size-1; i++ {
		fmt.Println(getLine(i, char, size))
	}
	fmt.Println(edgeLine)
}

func sandglass(args Args) {
	defArgs := Args{
		"char":  'X',
		"size":  15,
		"color": 30,
	}
	for k, v := range args {
		defArgs[k] = v
	}
	drawSandglass(defArgs["char"], defArgs["size"], defArgs["color"])
}

func main() {
	sandglass(Args{"color": 35, "size": 11})
}
