package main

import (
	"github.com/dionv/spogo/internal/auth"
	"github.com/joho/godotenv"
	// "github.com/fatih/color"
)

func main() {
	godotenv.Load()

	token, refreshToken, err := auth.GetTokens()
	if err != nil {
		code := auth.Authenticate()
		token, refreshToken = auth.ExchangeToken(code)
		auth.SaveToken(token, refreshToken)
	}

	token, valid := auth.EnsureValidTokens(token, refreshToken)
	if !valid {
		auth.SaveToken(token, refreshToken)
	}
}
