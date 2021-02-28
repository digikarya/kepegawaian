package authHelper

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenJWT struct {
	User string
	Username string
	Name string
	Role string
	Agen string
	AgenID string
	Scope string
	Expired string
	Token string
}

func (r *TokenJWT) Init(string2 ...string)  {
	r.User = string2[0]
	r.Username = string2[1]
	r.Name = string2[2]
	r.Role = string2[3]
	r.Agen = string2[4]
	r.AgenID = string2[5]
}

func (r *TokenJWT) CreateJWT() (error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = r.User
	claims["username"] = r.Username
	claims["name"] = r.Name
	claims["role"] = r.Role
	claims["agen"] = r.Agen
	claims["agen_id"] = r.AgenID
	time := time.Now().Add(time.Hour * 3600)
	claims["exp"] = time.Unix()
	r.Expired = time.Format("H:i:s")
	tokenString, err := token.SignedString([]byte("j24g$a@T8#mHN4%"))
	if err != nil {
		return err
	}
	r.Token = tokenString
	return nil
}

func (r *TokenJWT) ValidateJWT() (jwt.MapClaims, error) {
	token, err := jwt.Parse(r.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid token")
		}
		return []byte("j24g$a@T8#mHN4%"), nil
	})
	if token.Valid {
		claims, _ := token.Claims.(jwt.MapClaims)
		return claims,nil
	}else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("Invalid token:Expired")
		}
	}
	return nil, errors.New("Invalid token")
}

