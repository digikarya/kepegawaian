package model

import (
	"errors"
	"github.com/digikarya/helper"
	"github.com/digikarya/kepegawaian/authHelper"
	"gorm.io/gorm"
	"log"
)

type Auth struct {
	Request struct {
		Email string `json:"username" validate:"required,email" `
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
		Data struct{
			User UserAuth
			Agen interface{}
			Karyawan interface{}
		} `json:"data"`
	}
}

type UserAuth struct {
	HashID 			string
	Email			string `json:"email"  validate:"required,email"`
	Password		string `json:"password"  validate:"required"`
	Role			string `json:"role"  validate:"required"`
	KaryawanID		uint `json:"karyawan_id"  validate:"required"`
	AgenID			uint `json:"agen_id"  validate:"required"`
}

func (UserAuth) TableName() string {
	return "user"
}


func (data *Auth) initResponse(string2 ...string)  {
	data.Response.Credential.AccessToken = string2[0]
	data.Response.Credential.TokenType = "Bearer"
	data.Response.Credential.UserID = string2[1]
	data.Response.Credential.Scope = string2[2]
}

func (payload *Auth) Login(db *gorm.DB,param ...string) (interface{},error){
	users := UserAuth{}
	if err := helper.ValidateData(payload.Request);err != nil {
		return payload,err
	}
	result := db.Where("email = ?", payload.Request.Email).Find(&users)
	if err := result.Error; err != nil {
		log.Print(err.Error())
		return nil,errors.New("username atau password salah")
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("username atau password salah")
	}
	pass := helper.Password{}
	if pass.Check(payload.Request.Password,users.Password) != true{
		return nil,errors.New("password salah")
	}
	karyawan := KaryawanResponse{}
	agen := AgenPayload{}
	karyawanID,_ := helper.EncodeHash(int(users.KaryawanID))
	agenID,_ := helper.EncodeHash(int(users.AgenID))
	dataKaryawan,err := karyawan.Find(db,karyawanID)
	if err != nil {
		return nil,errors.New("username atau password salah")
	}
	dataAgen,err := agen.Find(db,agenID)
	if err != nil {
		return nil,errors.New("username atau password salah")
	}
	payload.Response.Data.User = users
	payload.Response.Data.Karyawan = dataKaryawan
	payload.Response.Data.Agen = dataAgen
	tokenJWT := authHelper.TokenJWT{}
	tokenJWT.Init(users.HashID,users.Email,dataKaryawan.(*KaryawanResponse).Nama,users.Role,dataAgen.(*AgenPayload).Nama,dataAgen.(*AgenPayload).HashID)
	if err := tokenJWT.CreateJWT();err != nil {
		return nil,err
	}
	payload.initResponse(tokenJWT.Token,users.HashID,tokenJWT.Expired,tokenJWT.Scope)
	return payload.Response,nil
}