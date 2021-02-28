package authHelper

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func AuthorizeRole(db *gorm.DB,r *http.Request,role string) error {
	token := TokenJWT{}
	decode,_,err := DecodeJWT(&token,r)
	if err != nil {
		return err
	}
	if role != "*" {
		if decode["role"] != role {
			return errors.New("Unauthorized s")
		}
	}
	return nil
}

func DecodeJWT(token *TokenJWT,r *http.Request) (jwt.MapClaims,[]string,error){
	reqToken := r.Header.Get("Authorization")
	splitBearer := strings.Split(reqToken, "Bearer")
	if len(splitBearer) != 2 {
		return nil,nil,errors.New("Invalid token")
	}
	token.Token = strings.TrimSpace(splitBearer[1])
	splitToken := strings.Split(token.Token, ".")
	if len(splitToken) != 3 {
		return nil,nil,errors.New("Unauthorized")
	}
	dataToken,err := token.ValidateJWT()
	if err != nil{
		return nil,nil,err
	}
	return dataToken,splitToken, nil
}
