package handler

import (
	"errors"
	"github.com/digikarya/helper"
	"github.com/digikarya/kepegawaian/app/model"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func KaryawanCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KaryawanPayload{}
	data,err := serv.Create(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func KaryawanAll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KaryawanResponse{}
	hashID,limit := helper.DecodeURLParam(r)
	data,err := serv.All(db,hashID,limit)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func KaryawanFind(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KaryawanResponse{}
	vars := mux.Vars(r)
	hashid,exist := vars["hashid"]
	if !exist{
		helper.RespondJSONError(w, http.StatusBadRequest, errors.New("Invalid id"))
		return
	}
	data,err := serv.Find(db,hashid)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Found",http.StatusOK, data)
	return
}

func KaryawanUpdate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KaryawanPayload{}
	vars := mux.Vars(r)
	hashid,exist := vars["hashid"]
	if !exist{
		helper.RespondJSONError(w, http.StatusBadRequest, errors.New("Invalid id"))
		return
	}
	data,err := serv.Update(db,r,hashid)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}

	helper.RespondJSON(w, "Updated",http.StatusOK, data)
	return
}

func KaryawanDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KaryawanPayload{}
	vars := mux.Vars(r)
	hashid,exist := vars["hashid"]
	if !exist{
		helper.RespondJSONError(w, http.StatusBadRequest, errors.New("Invalid id"))
		return
	}
	data,err := serv.Delete(db,hashid)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}

	helper.RespondJSON(w, "Deleted",http.StatusOK, data)
	return
}


