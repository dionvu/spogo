package main

import (
	"fmt"

	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/session"
	"github.com/dionv/spogo/internal/user"
	"github.com/fatih/color"
)

func main() {
	c, err := config.New()
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error"), err)
	}

	err = c.Load()
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error"), err)
	}

	s, err := session.New(c)
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error"), err)
	}

	u, _ := user.New(s)
	u.Print()
}
