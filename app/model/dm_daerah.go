package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type DaerahPayload struct{
	DaerahID        uint `gorm:"column:daerah_id; PRIMARY_KEY" json:"-"`
	HashID 			string `json:"id"  `
	Kabupaten		string `json:"kabupaten"  validate:"required"`
	Kecamatan		string `json:"Kecamatan"  validate:"required"`
	Provinsi		string `json:"provinsi"  validate:"required"`
}
type DaerahResponse struct{
	DaerahID        uint `gorm:"column:daerah_id; PRIMARY_KEY" json:"-"`
	HashID 			string `json:"id"`
	Kabupaten		string `json:"kabupaten"  validate:"required"`
	Kecamatan		string `json:"Kecamatan"  validate:"required"`
	Provinsi		string `json:"provinsi"  validate:"required"`
}

func (DaerahPayload) TableName() string {
	return "daerah"
}
func (DaerahResponse) TableName() string {
	return "daerah"
}



func (data *DaerahPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	tmp,err := data.defineValue()
	trx := db.Begin()
	result := trx.Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1{
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	if err := data.updateHashId(trx,int(tmp.DaerahID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *DaerahPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}

	if err := data.setPayload(r);err != nil {
		return nil, err
	}
	if _,err := data.countData(db,id);err != nil {
		return nil, err
	}
	tmp,err := data.defineValue()
	tmpUpdate := DaerahResponse{}
	if err := db.Where("daerah_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Where("daerah_id = ?", id).Save(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return data,nil
}


func (data *DaerahResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("daerah_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}



func (data *DaerahPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("daerah_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("daerah_id = ?",id).Delete(&data)
	if response.Error != nil {
		//log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *DaerahResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	var result []DaerahResponse
	limit,err := strconv.Atoi(string[1])
	if err != nil {
		return nil, err
	}
	trans := db.Limit(limit).Find(&result)
	hashID := string[0]
	if hashID != "" {
		id,err := helper.DecodeHash(hashID)
		if err != nil {
			return nil,err
		}
		trans = trans.Where("daerah_id > ?",id).Find(&result)
	}
	exec := trans.Find(&result)
	if exec.Error != nil {
		return result,exec.Error
	}
	return result,nil
}


// General Function =================================================================================
// ==================================================================================================
// ==================================================================================================
func (data *DaerahPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}


func (data *DaerahPayload) defineValue()  (tmp DaerahResponse,err error) {
	tmp.Kabupaten = data.Kabupaten
	tmp.Kecamatan = data.Kecamatan
	tmp.Provinsi = data.Provinsi

	return tmp,nil
}

func (data *DaerahResponse) switchValue(tmp *DaerahResponse) {
	data.Kabupaten = tmp.Kabupaten
	data.Kecamatan = tmp.Kecamatan
	data.Provinsi = tmp.Provinsi
}

func (data *DaerahResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *DaerahPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&DaerahResponse{}).Where("daerah_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *DaerahPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	response := db.Model(&data).Where("daerah_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal update")
	}
	return nil
}
