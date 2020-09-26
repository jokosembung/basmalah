package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User_mobile struct {
	ID          uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Seq_id      uint32    `gorm:"" json:"seq_id"`
	Nickname    string    `gorm:"size:255;not null;unique" json:"nickname"`
	Username    string    `gorm:"size:100;not null;unique" json:"username"`
	Password    string    `gorm:"size:100;not null;" json:"psw"`
	Dev_id      string    `gorm:"size:255;not null;" json:"dev_id"`
	App_ver     string    `gorm:"size:50;not null;" json:"version"`
	Status      int       `gorm:"size:1;not null;" json:"status"`
	Customer_id uint32    `gorm:"" json:"customer_id"`
	Kartu_id    int32     `gorm:"" json:"kartu_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Device_list struct {
	Status int `gorm:"" json:"status"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User_mobile) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User_mobile) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User_mobile) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if err := checkmail.ValidateFormat(u.Username); err != nil {
			return errors.New("Invalid Username")
		}

		return nil
	case "login":
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Dev_id == "" {
			return errors.New("Required device_id")
		}

		return nil

	default:
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if err := checkmail.ValidateFormat(u.Username); err != nil {
			return errors.New("Invalid Username")
		}
		return nil
	}
}

func (u *User_mobile) SaveUser(db *gorm.DB) (*User_mobile, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User_mobile{}, err
	}
	return u, nil
}

func (u *User_mobile) FindAllUsers(db *gorm.DB) (*[]User_mobile, error) {
	var err error
	users := []User_mobile{}
	err = db.Debug().Model(&User_mobile{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User_mobile{}, err
	}
	return &users, err
}

func (u *User_mobile) validateToken(db *gorm.DB, tokenid uint32) (bool, error) {
	var err error
	um := &User_mobile{}
	err = db.Debug().Model(User_mobile{}).Where("seq_id = ?", tokenid).Take(&um).Error
	if err != nil {
		return false, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return false, errors.New("Not valid ID")
	}
	return true, nil
}

func (u *User_mobile) FindUserByUsername(db *gorm.DB, uid string) (*User_mobile, error) {
	var err error
	um := &User_mobile{}
	err = db.Debug().Model(User_mobile{}).Where("username = ?", uid).Take(&um).Error
	if err != nil {
		return &User_mobile{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User_mobile{}, errors.New("User Not Found")
	}
	return um, err
}

func (u *User_mobile) UpdateUsermobile(db *gorm.DB, uid uint32) error {

	db = db.Debug().Model(&User_mobile{}).Where("seq_id = ?", uid).Take(&User_mobile{}).UpdateColumns(
		map[string]interface{}{
			"app_ver":   u.App_ver,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (u *User_mobile) UpdateDevice(db *gorm.DB, uid string) error {

	db = db.Debug().Model(&Device_list{}).Where("dev_id = ?", u.Dev_id).Take(&Device_list{}).UpdateColumns(
		map[string]interface{}{
			"status": 4,
		},
	)
	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (u *User_mobile) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User_mobile{}).Where("id = ?", uid).Take(&User_mobile{}).Delete(&User_mobile{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
