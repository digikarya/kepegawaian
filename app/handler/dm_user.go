package handler

import (
	"errors"
	"github.com/digikarya/helper"
	"github.com/digikarya/kepegawaian/app/model"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func UserCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.UserPayload{}
	data,err := serv.Create(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func UserAll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.UserResponse{}
	hashID,limit := helper.DecodeURLParam(r)
	data,err := serv.All(db,hashID,limit)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func UserFind(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.UserPayload{}
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
func UserFindByKaryawanID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.UserPayload{}
	vars := mux.Vars(r)
	hashid,exist := vars["hashid"]
	if !exist{
		helper.RespondJSONError(w, http.StatusBadRequest, errors.New("Invalid id"))
		return
	}
	data,err := serv.FindByKaryawanID(db,hashid)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Found",http.StatusOK, data)
	return
}

func UserUpdate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.UserPayload{}
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

func UserDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.UserPayload{}
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


