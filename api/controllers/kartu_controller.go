package controllers

import (
	//"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jokosembung/emaal/api/auth"
	"github.com/jokosembung/emaal/api/models"
	"github.com/jokosembung/emaal/api/responses"
	//"io/ioutil"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Default struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Detail []models.Kartu
}

func (server *Server) GetKartu(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["user_id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	nokartu := vars["no_kartu"]

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	kartu := models.Kartu{}
	kartuGotten, err := kartu.FindkartuByNomor(server.DB, nokartu)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	respon := Default{}
	respon.Status = true
	respon.Msg = "Found Data"
	respon.Detail = []models.Kartu{*kartuGotten}

	responses.JSON(w, http.StatusOK, respon)
}

func (server *Server) TambahKartu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["user_id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	paramKartu := models.ParamKartu{}
	err = json.Unmarshal(body, &paramKartu)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = paramKartu.Validate("tambah_kartu")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

}
