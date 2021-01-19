package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	authModel "gitlab.com/quick-count/go/auth/model"
	"gitlab.com/quick-count/go/helper"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func AuthToken(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	authService := authModel.Auth{}
	if err := helper.DecodeJson(r,&authService.Request);err != nil {
		helper.RespondJSONError(w, http.StatusUnauthorized, err)
		return
	}
	data,err := authService.Create(db)
	if err != nil {
		helper.RespondJSONError(w, http.StatusUnauthorized, err)
		return
	}

	helper.RespondJSON(w, "Succes",http.StatusOK, data)
}


func Authorization(db *gorm.DB,r *http.Request,role string) error {
	token := authModel.TokenJWT{}
	authService := authModel.Auth{}
	decode,splitToken,err := DecodeToken(&token,r)
	if err != nil {
		return err
	}
	if role != "*" {
		if decode["role"] != role {
			return errors.New("Unauthorized")
		}
	}

	if err := authService.CheckRevokeToken(db,splitToken);err != nil {
		return err
	}
	return nil
}
func DecodeToken(token *authModel.TokenJWT,r *http.Request) (jwt.MapClaims,[]string,error){
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
	dataToken,err := token.Validate()
	if err != nil{
		return nil,nil,err
	}
	return dataToken,splitToken, nil
}
