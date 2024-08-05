package main

import (
	"fmt"

	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/session"
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

	_, err = session.New(c)
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error"), err)
	}
}
