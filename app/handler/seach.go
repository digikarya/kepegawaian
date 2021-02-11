package handler

import (
	"github.com/digikarya/helper"
	"github.com/digikarya/kepegawaian/app/model"
	"gorm.io/gorm"
	"net/http"
)

func SearchDaerah(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.SearchRequest{}
	data,err := serv.DaerahSearch(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Found",http.StatusOK, data)
	return
}

func SearchAgen(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.SearchRequest{}
	data,err := serv.AgenSearch(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Found",http.StatusOK, data)
	return
}