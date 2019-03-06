package main

import (
    "fmt"
    "github.com/watsonserve/maild/compile/lexical"
)

func main() {
    list := lexical.Parse("2 LOGIN++james@watsonserve.com 123456")
    for i := 0; i < len(list); i++ {
        el := list[i]
        fmt.Printf("%d %d %s|\n", el.Cnt, el.Type, el.Value)
    }
}
