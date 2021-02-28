package handler

import (
	"github.com/digikarya/helper"
	"github.com/digikarya/kepegawaian/app/model"
	"gorm.io/gorm"
	"net/http"
)

func GetToken(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	authService := model.Auth{}
	if err := helper.DecodeJson(r,&authService.Request);err != nil {
		helper.RespondJSONError(w, http.StatusUnauthorized, err)
		return
	}
	data,err := authService.Login(db)
	if err != nil {
		helper.RespondJSONError(w, http.StatusUnauthorized, err)
		return
	}
	helper.RespondJSON(w, "Succes",http.StatusOK, data)
}

