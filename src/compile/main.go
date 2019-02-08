package main

import (
    "fmt"
    "compile/lexical"
)

func main() {
    list := lexical.Parse("\"foo\\\"bar\"")
    for item := list.Front(); item != nil; item = item.Next() {
        el := item.Value.(lexical.Lexical_t)
        fmt.Printf("%d %d %s|\n", el.Cnt, el.Type, el.Value)
    }
}
