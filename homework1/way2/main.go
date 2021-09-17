package main

import (
	"fmt"
	"strings"
)

type Arg = map[string]int
type Opt func(a *Arg)

func char(x int) Opt {
	return func(a *Arg) {
		(*a)["char"] = x
	}
}

func size(x int) Opt {
	return func(a *Arg) {
		(*a)["size"] = x
	}
}

func color(x int) Opt {
	return func(a *Arg) {
		(*a)["color"] = x
	}
}

func getLine(num_line int, char int, size int) string {
	s := make([]string, size)
	for i := 0; i < size; i++ {
		if (i == num_line) || (i == size-1-num_line) {
			s[i] = string(rune(char))
		} else {
			s[i] = " "
		}
	}
	return strings.Join(s, "")
}

func drawSandglass(char int, size int, color int) {
	edge_line := strings.Repeat(string(char), size)
	fmt.Printf("\033[%dm", color)
	fmt.Println(edge_line)
	for i := 1; i < size-1; i++ {
		fmt.Println(getLine(i, char, size))
	}
	fmt.Println(edge_line)
}

func sandglass(opts ...Opt) {
	args := Arg{
		"char":  'X',
		"size":  15,
		"color": 30,
	}
	for _, opt := range opts {
		opt(&args)
	}
	drawSandglass(args["char"], args["size"], args["color"])
}

func main() {
	sandglass(color(34), size(5))
	sandglass(char('X'))
}
