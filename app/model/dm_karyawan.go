package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type KaryawanPayload struct{
	KaryawanID     	uint `gorm:"column:karyawan_id; PRIMARY_KEY" json:"-"`
	HashID 			string `json:"id"  `
	Nama			string `json:"nama"  validate:"required"`
	NomorIdentitas  string `json:"no_identitas"  validate:"required,number"`
	NomorTelepon	string `json:"telepon"  validate:"required,number"`
	Alamat			string `json:"alamat"  validate:"required"`
	JenisKelamin	string `json:"jenis_kelamin"  validate:"required,alpha"`
	Jabatan			string `json:"jabatan"  validate:"required,alpha"`
}

type KaryawanResponse struct{
	KaryawanID     	uint `gorm:"column:karyawan_id; PRIMARY_KEY" json:"-"`
	HashID 			string `json:"id"  `
	Nama			string `json:"nama"  validate:"required,alpha"`
	NomorIdentitas  string `json:"no_identitas"  validate:"required,number"`
	NomorTelepon	string `json:"telepon"  validate:"required,number"`
	Alamat			string `json:"alamat"  validate:"required"`
	JenisKelamin	string `json:"jenis_kelamin"  validate:"required,alpha"`
	Jabatan			string `json:"jabatan"  validate:"required,alpha"`
	//Users           interface{}
}


func (KaryawanPayload) TableName() string {
	return "karyawan"
}

func (KaryawanResponse) TableName() string {
	return "karyawan"
}


func (data *KaryawanPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
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
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	if err := data.updateHashId(trx,int(tmp.KaryawanID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *KaryawanPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := KaryawanResponse{}
	if err := db.Where("karyawan_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Where("karyawan_id = ?", id).Save(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return data,nil
}


func (data *KaryawanResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("karyawan_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}



func (data *KaryawanPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("karyawan_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("karyawan_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *KaryawanResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	var result []KaryawanResponse
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
		trans = trans.Where("karyawan_id > ?",id).Find(&result)
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
func (data *KaryawanPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}


func (data *KaryawanPayload) defineValue()  (tmp KaryawanResponse,err error) {
	tmp.Nama = data.Nama
	tmp.Nama = data.Nama
	tmp.NomorIdentitas = data.NomorIdentitas
	tmp.NomorTelepon = data.NomorTelepon
	tmp.Alamat = data.Alamat
	tmp.NomorTelepon = data.NomorTelepon
	tmp.JenisKelamin = data.JenisKelamin
	tmp.Jabatan = data.Jabatan

	return tmp,nil
}

func (data *KaryawanResponse) switchValue(tmp *KaryawanResponse) {
	data.Nama = tmp.Nama
	data.Nama = tmp.Nama
	data.NomorIdentitas = tmp.NomorIdentitas
	data.NomorTelepon = tmp.NomorTelepon
	data.Alamat = tmp.Alamat
	data.NomorTelepon = tmp.NomorTelepon
	data.JenisKelamin = tmp.JenisKelamin
	data.Jabatan = tmp.Jabatan

}

func (data *KaryawanResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *KaryawanPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&KaryawanResponse{}).Where("karyawan_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *KaryawanPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.KaryawanID)
	response := db.Model(&data).Where("karyawan_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal update")
	}
	return nil
}
