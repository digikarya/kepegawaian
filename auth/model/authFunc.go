package model

import (
	"errors"
	"gitlab.com/quick-count/go/helper"
	mainModel "gitlab.com/quick-count/go/app/model"
	"gorm.io/gorm"
)

func (data *Auth) initResponse(string2 ...string)  {
	data.Response.Credential.AccessToken = string2[0]
	data.Response.Credential.TokenType = "Bearer"
	data.Response.Credential.Scope = string2[2]
}

func (data *Auth) Create(db *gorm.DB) (interface{},error) {
	saksi := mainModel.Saksi{}
	err := saksi.Users.Auth(db,data.Request.Username,data.Request.Password)
	if err != nil{
		return nil, err
	}
	id, _ := helper.EncodeHash(int(saksi.Users.SaksiID))
	saksi.Find(db,id)
	tokenJWT := TokenJWT{}
	tokenJWT.Init(saksi.Users.HashID,saksi.Users.Username,saksi.Name,saksi.Users.Role)
	if err := tokenJWT.Create();err != nil {
		return nil,err
	}
	data.initResponse(tokenJWT.Token,tokenJWT.Expired,tokenJWT.Scope)
	return data.Response,nil
}

func (data *Auth) CheckRevokeToken(db *gorm.DB,token []string) error {
	isExist := db.Where("token = ?",token[2]).Find(&data.revokeToken)
	if err := isExist.Error; err != nil {
		return err
	}
	if isExist.RowsAffected > 0 {
		return errors.New("Token is expired")
	}
	return nil
}