package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"net/http"
)

type SearchRequest struct {
	Condition []struct{
		Column string `json:"column"  validate:"required,alpha"`
		Value string `json:"value"  validate:"required"`
	} `json:"condition"  validate:"required"`
}

func (payload *SearchRequest) DaerahSearch(db *gorm.DB,r *http.Request)  (interface{},error) {
	err := payload.setPayload(r)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	var tmpData = []DaerahResponse{}
	result := db.Where("kabupaten LIKE ?", "%"+payload.Condition[0].Value+"%").Or("kecamatan LIKE ?", "%"+payload.Condition[0].Value+"%").Find(&tmpData)
	result = result.Order("kabupaten asc, kecamatan asc").Find(&tmpData)
	if err := result.Error; err != nil {
		return nil,errors.New("data tidak ditemukan")
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmpData,nil
}

func (payload *SearchRequest) AgenSearch(db *gorm.DB,r *http.Request)  (interface{},error) {
	err := payload.setPayload(r)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	var tmpData = []AgenResponse{}
	result := db.Where("nama LIKE ?", "%"+payload.Condition[0].Value+"%").Find(&tmpData)
	result = result.Order("nama asc").Find(&tmpData)
	if err := result.Error; err != nil {
		return nil,errors.New("data tidak ditemukan")
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmpData,nil
}

func (payload *SearchRequest) PegawaiSearch(db *gorm.DB,r *http.Request)  (interface{},error) {
	err := payload.setPayload(r)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	var tmpData = []KaryawanPayload{}
	result := db.Where("nama LIKE ?","%"+payload.Condition[0].Value+"%").Or("nomor_telepon LIKE ?","%"+payload.Condition[0].Value+"%").Find(&tmpData)
	result = result.Order("nama asc").Find(&tmpData)
	if err := result.Error; err != nil {
		return nil,errors.New("data tidak ditemukan")
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmpData,nil
}
func (payload *SearchRequest) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&payload);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(payload);err != nil {
		return err
	}
	if len(payload.Condition) < 1 {
		return errors.New("invalid payload")
	}
	return nil
}