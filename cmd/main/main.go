package main

import (
	"fmt"

	auth "github.com/dionv/spogo/internal"
	"github.com/fatih/color"
)

const lOGINSUCCESSMSG = `
Login Success!
`

func main() {
	code := auth.Authenticate()

	tok, _ := auth.ExchangeToken(code)

	fmt.Print(tok)

	color.Green(lOGINSUCCESSMSG)
}
