package model

import (
	"gitlab.com/quick-count/go/helper"
	"time"
)

type Auth struct {
	Request struct {
		Username string `json:"username" validate:"required" `
		Password string `json:"password" validate:"required" `
	}
	Response struct{
		Credential struct{
			UserID string `json:"user_id"`
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
			//ExpiresIn   string `json:"expire_in"`
			Scope        string `json:"scope"`
		} `json:"credential"`
	}
	revokeToken
}
type revokeToken struct{
	Token string
	UserID string
	LogoutAt time.Time
	helper.TimeModel
}
func (revokeToken) TableName() string {
	return "revoke_tokens"
}