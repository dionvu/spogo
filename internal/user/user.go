package user

import "github.com/dionv/spogo/internal/tokens"

type User struct {
	AccessToken  *tokens.AccessToken
	RefreshToken *tokens.RefreshToken
}

func New() *User {
	return &User{
		AccessToken:  &tokens.AccessToken{},
		RefreshToken: &tokens.RefreshToken{},
	}
}
