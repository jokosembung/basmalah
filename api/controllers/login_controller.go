package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jokosembung/basmalah/api/auth"
	"github.com/jokosembung/basmalah/api/models"
	"github.com/jokosembung/basmalah/api/responses"
	"golang.org/x/crypto/bcrypt"
)

type ResponseLogin struct {
	Timestamp     string `json:"timestamp"`
	Default_kartu int32  `json:"default_kartu"`
	Customer_id   uint32 `json:"customer_id"`
	Token         string `json:"token"`
	Msg           string `json:"msg"`
}

type Reponse struct {
	Kode        bool   `json:"kode`
	Status_code string `json:"status_code"`
	Result      []ResponseLogin
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User_mobile{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	mpass := user.Password
	mdevice := user.Dev_id
	musername := user.Username

	rowUser, err := user.FindUserByUsername(server.DB, user.Username)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if rowUser.Status != 1 {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Status Username untuk kartu Anda belum di aktivasi, Silahkan aktivasi."))
		return
	}

	if rowUser.Dev_id != mdevice {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Device Tidak DitemukaN, Silahkan registrasi ulang"))
		return
	}

	if server.MD5(mpass) != rowUser.Password {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Password tidak sesuai"))
		return
	}

	kartu := models.Kartu{}
	getKartu, err := kartu.FindkartuByCustomerId(server.DB, rowUser.Customer_id)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = user.UpdateUsermobile(server.DB, uint32(rowUser.Seq_id))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/*
		err = user.UpdateDevice(server.DB, rowUser.Dev_id)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
	*/
	err = models.VerifyPassword(user.Password, rowUser.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := auth.CreateToken(rowUser.Seq_id)

	responseLogin := ResponseLogin{}

	jsonKartu, err := json.Marshal(getKartu)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	unixtime := fmt.Sprintf("%v", time.Now().Unix())
	sdata := fmt.Sprintf("%s:%s:%s", musername, mpass, unixtime)
	passcode := base64.StdEncoding.EncodeToString([]byte(sdata))
	responseLogin.Default_kartu = rowUser.Kartu_id
	responseLogin.Customer_id = rowUser.Customer_id
	responseLogin.Timestamp = unixtime
	responseLogin.Token = token
	responseLogin.Msg = auth.Encrypt(string(jsonKartu), unixtime, passcode)

	resp := Reponse{}
	resp.Kode = true
	resp.Status_code = "000"
	resp.Result = []ResponseLogin{responseLogin}
	responses.JSON(w, http.StatusOK, resp)
}

func (server *Server) MD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}
