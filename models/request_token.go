package models

import (
	"fmt"
	"time"

	"Go-iris-boilerplate/database"
	"Go-iris-boilerplate/utils"

	"github.com/jinzhu/gorm"
)

type REQUEST_TYPES string

const (
	Signup         REQUEST_TYPES = "signup"
	ForgotPassword REQUEST_TYPES = "forgotPassword"
)

type REQUEST_STATUS string

const (
	expired REQUEST_STATUS = "expired"
	pending REQUEST_STATUS = "pending"
	done    REQUEST_STATUS = "done"
)

type RequestToken struct {
	gorm.Model

	Token  string         `gorm:"index;unique;not null;VARCHAR(40);"`
	Type   REQUEST_TYPES  `gorm:"type:string;not null"`
	Status REQUEST_STATUS `gorm:"type:string;not null;default:'pending'"`
	Expiry time.Time      `gorm:"not null"`
	UserID uint
}

type RequestTokenJson struct {
	Token  string         `json:"token" validate:"required;unique;VARCHAR(40)"`
	Type   REQUEST_TYPES  `json:"type" validate:"type:string;required"`
	Status REQUEST_STATUS `json:"status" validate:"type:string;required;default:'pending'"`
	Expiry time.Time      `json:"expiry" validate:"required"`
	UserID uint           `json:"user_id" validate:"required"`
}

func RequestTokenExpired() {

}

func CreateRequestToken(aul *RequestTokenJson) (request_token *RequestToken) {
	token := utils.RandASCIIBytes(40)

	request_token = new(RequestToken)
	request_token.Token = token
	request_token.Type = aul.Type
	request_token.UserID = aul.UserID
	if err := database.DB.Create(request_token).Error; err != nil {
		fmt.Printf("Create Request token Error %s", err)
	}

	return request_token
}

/**
 * Get the RequestToken record by id
 * @method GetRequestByToken
 * @param  {[type]}       user  *RequestToken [description]
 */
func GetRequestByToken(token string) *RequestToken {
	request_token := new(RequestToken)
	request_token.Token = token

	if err := database.DB.Where("token = ?", token).First(request_token).Error; err != nil {
		fmt.Printf("GetRequestByToken:%s", err)
	}
	return request_token
}

/**
 */
func DeleteRequestTokenByToken(token string) (err error) {
	request_token := GetRequestByToken(token)

	if err := database.DB.Delete(request_token).Error; err != nil {
		fmt.Printf("DeleteRequestTokenByTokenErr:%s", err)
	}
	return
}
