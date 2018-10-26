package auth

import (
	"fmt"
)

type Author struct {
}

func New() *Author {
    return &Author{}
}

func (this *Author) Auth(username string, password string) string {
    fmt.Printf("%s %s\n", username, password)
    return "null"
}
