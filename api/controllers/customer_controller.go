package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/jokosembung/basmalah/api/models"
	"github.com/jokosembung/basmalah/api/responses"
	"io/ioutil"
	"net/http"
	"strconv"
)

type DefaultConfig struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Detail []models.ConfSync
}

type Customer struct {
	Br_id      string `gorm:"column:br_id" json:"br_id"`
	Cu_code    string `gorm:"column:cu_code" json:"cu_code"`
	Cu_name    string `gorm:"column:cu_name" json:"cu_name"`
	Emaal      string `gorm:"column:emaal" json:"emaal"`
	Hp         string `gorm:"column:hp" json:"hp"`
	Addr       string `gorm:"column:addr" json:"addr"`
	Status     string `gorm:"column:status" json:"status"`
	Created_at string `gorm:"column:created_at" json:"created_at"`
	No_ktp     string `gorm:"column:no_ktp" json:"no_ktp"`
}

func (server *Server) GetLastLocal(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	br_code := vars["br_code"]

	kartu := models.ConfSync{}
	kartuGotten, err := kartu.FindConfig(server.DB, br_code)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	respon := DefaultConfig{}
	respon.Status = true
	respon.Msg = "Found Data"
	respon.Detail = []models.ConfSync{*kartuGotten}

	responses.JSON(w, http.StatusOK, respon)
}

func (server *Server) TambahCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	br_code := vars["br_code"]
	//br_id := vars["br_id"]
	last_local_id, err := strconv.ParseInt(vars["last_local_id"], 10, 32)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var prod []Customer

	err = json.Unmarshal(body, &prod)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var conf = models.ConfSync{}
	conf.Local_cust_id = int32(last_local_id)
	isok, err := insertCustomerBatch(server.DB, prod, br_code, conf)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	respon := DefaultConfig{}
	respon.Status = isok
	respon.Msg = "Found Data"

	responses.JSON(w, http.StatusOK, respon)

}

func insertCustomerBatch(db *gorm.DB, requestObj []Customer, br_code string, updatConf models.ConfSync) (bool, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, obj := range requestObj {
		if err := tx.Where(Customer{Br_id: obj.Br_id, Cu_code: obj.Cu_code}).Assign(&obj).FirstOrCreate(&obj).Error; err != nil {
			fmt.Println("Failed to insert batch ")
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Table("conf_sync").Where("br_code = ?", br_code).Updates(&updatConf).Error; err != nil {

		fmt.Println("Failed to update config")
		tx.Rollback()
		return false, err
	}

	err := tx.Commit().Error
	if err != nil {
		return false, err
	}
	return true, nil
}
