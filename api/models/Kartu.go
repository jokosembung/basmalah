package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
)

type Kartu struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	No_kartu string `gorm:"size:255;not null;unique" json:"no_kartu"`
	Nama     string `gorm:"size:100;not null;unique" json:"nama"`
	Saldo    string `gorm:"size:100;not null;" json:"saldo"`
}

type ParamKartu struct {
	No_kartu    string `json:"no_kartu"`
	Pin         string `json:"pin"`
	Hp          string `json:"hp"`
	Device_id   string `json:"device_id"`
	Customer_id string `json:"customer_id"`
}

func (u *ParamKartu) Validate(action string) error {
	switch strings.ToLower(action) {
	case "trambah_kartu":
		if u.No_kartu == "" {
			return errors.New("Required No Kartu")
		}
		if u.Hp == "" {
			return errors.New("Required HP")
		}
		if u.Pin == "" {
			return errors.New("Required PIN")
		}
		return nil

	default:
		if u.No_kartu == "" {
			return errors.New("Required No Kartu")
		}
		if u.Hp == "" {
			return errors.New("Required HP")
		}
		if u.Pin == "" {
			return errors.New("Required PIN")
		}
		return nil
	}
}

func (u *Kartu) FindkartuByCustomerId(db *gorm.DB, uid uint32) (*Kartu, error) {
	var err error
	err = db.Debug().Model(Kartu{}).Where("customer_id = ?", uid).Take(&u).Error
	if err != nil {
		return &Kartu{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Kartu{}, errors.New("Kartu tidak ditemukan silahkan aktivasi ulang")
	}
	return u, err
}

func (u *Kartu) FindkartuByNomor(db *gorm.DB, uid string) (*Kartu, error) {
	var err error
	err = db.Debug().Model(Kartu{}).Where("no_kartu = ?", uid).Take(&u).Error
	if err != nil {
		return &Kartu{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Kartu{}, errors.New("Kartu tidak ditemukan silahkan aktivasi ulang")
	}
	return u, err
}
