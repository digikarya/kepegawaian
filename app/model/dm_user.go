package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type UserPayload struct{
		UserID	     	uint `gorm:"column:user_id; PRIMARY_KEY" json:"-"`
		HashID			string `json:"id"  validate:""`
		Email			string `json:"email"  validate:"required,email"`
		Password		string `json:"password"  validate:"required"`
		Role			string `json:"role"  validate:"required"`
		KaryawanID		string `json:"karyawan_id"  validate:"required,base64url"`
		AgenID			string `json:"agen_id"  validate:"required,base64url"`
		Agen 			string `json:"agen"  validate:""`
		Karyawan 		string `json:"karyawan"  validate:""`
	}
type UserResponse struct{
		UserID	     	uint `gorm:"column:user_id; PRIMARY_KEY" json:"-"`
		HashID			string `json:"id"  validate:""`
		Email			string `json:"email"  validate:"required,email"`
		Password		string `json:"password"  validate:"required"`
		Role			string `json:"role"  validate:"required"`
		KaryawanID		uint `json:"karyawan_id"  validate:"required"`
		AgenID			uint `json:"agen_id"  validate:"required"`
		Agen 			string `json:"agen"  validate:""`
		Karyawan 		string `json:"karyawan"  validate:""`
}

func (UserPayload) TableName() string {
	return "user"
}
func (UserResponse) TableName() string {
	return "user"
}
func (data *UserPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	tmp,err := data.defineValue()
	trx := db.Begin()
	result := trx.Select("email","password","role","karyawan_id","agen_id").Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	if err := data.updateHashId(trx,int(tmp.UserID));err != nil{
		trx.Rollback()
		return nil, err
	}
	data.Password = "-"
	trx.Commit()
	return data,nil
}

func (data *UserPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := UserResponse{}
	if err := db.Where("user_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Where("user_id = ?", id).Save(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return data,nil
}


func (data *UserPayload) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	//tmpData := []UserPayload{}
	sql :=  "SELECT " +
		"	 user.hash_id,user.email,user.role," +
		"    agen.agen_id 'id',agen.hash_id 'agen_id', agen.nama 'agen'," +
		"    karyawan.karyawan_id 'id',karyawan.hash_id 'karyawan_id', karyawan.nama 'karyawan'" +
		"	 FROM user " +
		"    JOIN agen ON user.agen_id=agen.agen_id" +
		"    JOIN karyawan ON user.karyawan_id=karyawan.karyawan_id" +
		"	 WHERE user.user_id = ? LIMIT 1"
	exec := db.Raw(sql, id).Scan(&data)
	if exec.Error != nil {
		return nil,exec.Error
	}
	if exec.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}

func (data *UserPayload) FindByKaryawanID(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	tmpData := []UserPayload{}
	sql :=  "SELECT " +
			"	 user.hash_id,user.email,user.role," +
			"    agen.agen_id 'id',agen.hash_id 'agen_id', agen.nama 'agen'," +
			"    karyawan.karyawan_id 'id',karyawan.hash_id 'karyawan_id', karyawan.nama 'karyawan'" +
		"	 FROM user " +
		"    JOIN agen ON user.agen_id=agen.agen_id" +
		"    JOIN karyawan ON user.karyawan_id=karyawan.karyawan_id" +
		"	 WHERE user.karyawan_id = ?"
	exec := db.Raw(sql, id).Scan(&tmpData)
	if exec.Error != nil {
		return nil,exec.Error
	}
	if exec.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmpData,nil
}


func (data *UserPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("user_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("user_id = ?",id).Delete(&data)
	if response.Error != nil {
		//log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *UserResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	var result []UserPayload
	limit,err := strconv.Atoi(string[1])
	if err != nil {
		return nil, err
	}
	//trans := db.Limit(limit).Find(&result)
	sql :=  "SELECT " +
		"	 user.hash_id,user.email,user.role," +
		"    agen.agen_id 'id',agen.hash_id 'agen_id', agen.nama 'agen'," +
		"    karyawan.karyawan_id 'id',karyawan.hash_id 'karyawan_id', karyawan.nama 'karyawan'" +
		"	 FROM user " +
		"    JOIN agen ON user.agen_id=agen.agen_id" +
		"    JOIN karyawan ON user.karyawan_id=karyawan.karyawan_id"

	hashID := string[0]
	param1 := limit
	param2 := limit
	if hashID != "" {
		id,err := helper.DecodeHash(hashID)
		if err != nil {
			return nil,err
		}
		sql += " WHERE user_id > ?"
		param1 = int(id)
		//trans = trans.Where("agen_id > ?",id).Find(&result)
	}
	exec := db.Raw(sql+" LIMIT ?", param1,param2).Scan(&result)
	if exec.Error != nil {
		return nil,exec.Error
	}
	return result,nil
}




// General Function =================================================================================
// ==================================================================================================
// ==================================================================================================
func (data *UserPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}


func (data *UserPayload) defineValue()  (tmp UserResponse,err error) {
	tmp.Email = data.Email
	tmp.Role = data.Role
	data.setPassword()
	tmp.Password = data.Password
	tmp.KaryawanID,err = helper.DecodeHash(data.KaryawanID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	tmp.AgenID,err = helper.DecodeHash(data.AgenID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	return tmp,nil
}
func (data *UserPayload) setPassword() error{
	passs := strings.TrimSpace(data.Password)
	if passs == ""  || len(passs) < 7 {
		return errors.New("Password harus diisi,minimal harus 7 karakter")
	}
	pass := helper.Password{}
	password,err := pass.HashPassword(passs)
	if err != nil {
		return err
	}
	data.Password = password
	return err
}

func (data *UserResponse) switchValue(tmp *UserResponse) {
	data.Email = tmp.Email
	data.Role = tmp.Role
	data.Password = tmp.Password
	data.AgenID = tmp.AgenID
	data.KaryawanID = tmp.KaryawanID
}

func (data *UserResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *UserPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&UserResponse{}).Where("user_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *UserPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.UserID)
	response := db.Model(&data).Where("user_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal update")
	}
	return nil
}
