package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type ConfSync struct {
	ID            uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Br_id         string `json:"br_id"`
	Br_code       string `gorm:"size:255;not null;unique" json:"br_code"`
	Local_cust_id int32  `gorm:"size:100;not null;unique" json:"local_cust_id"`
	Row_limit     string `json:"row_limit"`
}

func (u *ConfSync) FindConfig(db *gorm.DB, brcode string) (*ConfSync, error) {
	var err error
	err = db.Debug().Model(ConfSync{}).Where("br_code = ?", brcode).Take(&u).Error
	if err != nil {
		return &ConfSync{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &ConfSync{}, errors.New("Kartu tidak ditemukan silahkan aktivasi ulang")
	}
	return u, err
}

func (u *ConfSync) UpdateConfigField(db *gorm.DB, brcode string, updateConf ConfSync) error {

	db = db.Debug().Model(&ConfSync{}).Where("br_code = ?", brcode).Updates(updateConf)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
