package main

import (
  "fmt"
  "strings"
)

type Args = map[string]int


func get_line(num_line int, char int, size int) string {
  s := make([]string, size)
  for i := 0; i < size; i ++ {
      if (i == num_line) || (i == size-1-num_line) {
        s[i] = string(char)
      } else {
        s[i] = " "
      }
  }
  return strings.Join(s, "")
}


func draw_sandglass(char int, size int, color int) {
  edge_line := strings.Repeat(string(char), size)
  fmt.Printf("\033[%dm", color)
  fmt.Println(edge_line)
  for i := 1; i < size-1; i ++ {
    fmt.Println(get_line(i, char, size))
  }
  fmt.Println(edge_line)
}


func sandglass(args Args) {
  def_args := Args{
    "char" : 'X',
    "size" : 15,
    "color" : 30,
  }
  for k, v := range args {
    def_args[k] = v
  }
  draw_sandglass(def_args["char"], def_args["size"], def_args["color"])
}



func main() {
  sandglass(Args{"color" : 35, "size": 11})
}
