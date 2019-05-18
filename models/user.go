package models

import (
	"Go-iris-boilerplate/database"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Email       string `gorm:"unique;VARCHAR(254)"`
	FirstName   string `gorm:"not null VARCHAR(120)"`
	LastName    string `gorm:"VARCHAR(120) default:''"`
	Password    string `gorm:"not null VARCHAR(191)"`
	IsVerified  bool   `gorm:"default:false"`
	IsSuperUser bool   `gorm:"default:false"`
}

type UserJson struct {
	Email     string `json:"email" validate:"required,gte=5,lte=254"`
	Password  string `json:"password" validate:"required,gte=6,lte=191"`
	FirstName string `json:"first_name" validate:"required,gte=2,lte=120"`
	LastName  string `json:"last_name" validate:"lte=120"`
}

type ShowUserJson struct {
	ID        uint
	Email     string `json:"email" validate:"required,gte=5,lte=254"`
	FirstName string `json:"first_name" validate:"required,gte=2,lte=120"`
	LastName  string `json:"last_name" validate:"lte=120"`
}

type UpdateUserJson struct {
	FirstName string `json:"first_name" validate:"gte=2,lte=120"`
	LastName  string `json:"last_name" validate:"lte=120"`
}

/**
 * Get the user record by id
 * @method GetUserById
 * @param  {[type]}       user  *User [description]
 */
func GetUserById(id uint) (user *User, err error) {
	user = new(User)
	user.ID = id
	if err := database.DB.Find(&user, id).Error; err != nil {
		fmt.Printf("GetUserByIdErr:%s", err)
	}
	usr := new(ShowUserJson)

	usr.Email = user.Email
	usr.FirstName = user.FirstName
	usr.LastName = user.LastName

	return
}

/**
 * Get the user record by email
 * @method GetUserById
 * @param  {[type]}       user  *User [description]
 */
func GetUserByEmail(email string) (user *User, err error) {
	user = new(User)
	user.Email = email

	if err := database.DB.First(user).Error; err != nil {
		fmt.Printf("GetUserByEmailErr:%s", err)
	}
	return
}

/**
 * CreateUser
 * @method CreateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateUser(aul *UserJson) (user *User, err error) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = new(User)
	user.Email = aul.Email
	user.Password = string(hash)
	user.FirstName = aul.FirstName
	user.LastName = aul.LastName

	err = database.DB.Create(user).Error
	if err != nil {
		fmt.Printf("CreateUserErr:%s", err)
	}

	return
}

func UpdateUser(uj *UpdateUserJson, id uint) (user *User, err error) {
	user, _ = GetUserById(id)

	user.FirstName = uj.FirstName
	user.LastName = uj.LastName

	err = database.DB.Model(user).Updates(user).Error
	if err != nil {
		fmt.Printf("UpdateUserErr:%s", err)
	}

	return
}

/**
 *
 */
func UpdateUserPassword(password string, id uint) (user *User, err error) {
	user = new(User)
	user.ID = id
	user.Password = password

	err = database.DB.Model(user).Updates(user).Error
	if err != nil {
		fmt.Printf("UpdateUserPasswordErr:%s", err)
	}
	return

}

/**
 *
 */
func MarkUserVerified(id uint) (user *User, err error) {
	user = new(User)
	user.ID = id
	user.IsVerified = true

	err = database.DB.Model(user).Updates(user).Error
	if err != nil {
		fmt.Printf("MarkUserVerifiedErr:%s", err)
	}

	return
}

func GetAllUsers(name, orderBy string, offset, limit int) (usr []*ShowUserJson) {
	users := []*User{}
	if err := database.GetAll(name, orderBy, offset, limit).Find(&users).Error; err != nil {
		fmt.Printf("GetAllUserErr:%s", err)
	}
	for _, element := range users {
		u := new(ShowUserJson)
		u.Email = element.Email
		u.ID = element.ID
		u.FirstName = element.FirstName
		u.LastName = element.LastName
		usr = append(usr, u)
	}
	return
}

/**
 * Verify user login
 * @method UserAdminCheckLogin
 * @param  {[type]}       username string [description]
 */
func UserAdminCheckLogin(email string) User {
	u := User{}
	if err := database.DB.Where("email = ?", email).First(&u).Error; err != nil {
		fmt.Printf("UserAdminCheckLoginErr:%s", err)
	}
	return u
}

func CheckLogin(email, password string) (response string, status bool, msg string) {
	user := UserAdminCheckLogin(email)
	if user.ID == 0 {
		msg = "User does not exist"
		return
	} else {
		if ok := bcrypt.Match(password, user.Password); ok {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)
			// 1 hour expiry
			claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			claims["iat"] = time.Now().Unix()
			token.Claims = claims
			tokenString, err := token.SignedString([]byte("secret"))

			if err != nil {
				msg = err.Error()
				return
			}

			msg = "Logged in successfully"
			response = tokenString
			status = true

			return
		} else {
			msg = "Wrong email or password"
			status = false
			return
		}
	}
}
