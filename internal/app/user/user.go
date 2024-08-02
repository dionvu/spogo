package user

import "github.com/dionv/spogo/internal/app/user/tokens"

type User struct {
	AccessToken  *tokens.AccessToken
	RefreshToken *tokens.RefreshToken
}
