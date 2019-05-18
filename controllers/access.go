package controllers

import (
	"Go-iris-boilerplate/models"
	"Go-iris-boilerplate/utils"
	"fmt"

	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
)

type Cred struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserLogin(ctx iris.Context) {
	aul := new(Cred)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ApiResource(false, nil, "Request parameter error"))
	} else {
		if EmailErr := validate.Var(aul.Email, "required,gte=5,lte=254"); EmailErr != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(ApiResource(false, nil, "Email format error"))
		} else if PwdErr := validate.Var(aul.Password, "required,gte=6,lte=191"); PwdErr != nil {

			ctx.JSON(ApiResource(false, nil, "Wrong password format"))
			ctx.StatusCode(iris.StatusBadRequest)
		} else {
			ctx.StatusCode(iris.StatusOK)
			response, status, msg := models.CheckLogin(aul.Email, aul.Password)
			if status == false {
				ctx.StatusCode(iris.StatusBadRequest)
			}

			ctx.JSON(ApiResource(status, response, msg))
		}
	}
}

func UserSignup(ctx iris.Context) {
	aul := new(models.UserJson)

	if err := ctx.ReadJSON(&aul); err != nil {
		// debug.PrintStack()
		fmt.Println(aul)
		ctx.StatusCode(iris.StatusUnauthorized) //TODO: check this status code
		ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(aul)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println()
				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.Type())
				fmt.Println(err.Param())
				fmt.Println()
			}
		} else {
			u, err := models.CreateUser(aul)

			if u.ID == 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(ApiResource(false, err, "operation failed"))
			} else {
				request_token := new(models.RequestTokenJson)
				request_token.Type = models.Signup
				request_token.UserID = u.ID
				req_token := models.CreateRequestToken(request_token)
				if req_token == nil {
					// Failed to create token
					// TODO: needs to create resent token link
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.JSON(ApiResource(false, nil, "Failed"))
					return
				}
				utils.SendEmail("Email Verification",
					"templates/verify_email.html",
					map[string]string{
						"FirstName": u.FirstName,
						"link":      "token=" + req_token.Token,
					},
					u.Email,
				)
				ctx.StatusCode(iris.StatusOK)
				ctx.JSON(ApiResource(true, u, "Successful"))
			}
		}
	}
}
