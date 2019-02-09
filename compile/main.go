package main

import (
    "fmt"
    "compile/lexical"
)

func main() {
    list := lexical.Parse("\"foo\\\"bar\"")
    for i := 0; i < len(list); i++ {
        el := list[i]
        fmt.Printf("%d %d %s|\n", el.Cnt, el.Type, el.Value)
    }
}
