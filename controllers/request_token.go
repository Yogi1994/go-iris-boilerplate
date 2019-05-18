package controllers

import (
	"MVC-golang/models"
	"MVC-golang/utils"

	"github.com/kataras/iris"
)

type RequestJson struct {
	token string
}

func ForgotPassword(ctx iris.Context) {

	email := ctx.Params().GetString("email")

	user, _ := models.GetUserByEmail(email)
	println(email)
	if user == nil {
		// TODO:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ApiResource(false, nil, "operation failed"))
	} else {
		request_token := new(models.RequestTokenJson)
		request_token.Type = models.ForgotPassword
		request_token.UserID = user.ID
		req_token := models.CreateRequestToken(request_token)

		if req_token == nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(ApiResource(false, nil, "operation failed"))
		} else {
			utils.SendEmail("Reset Password",
				"templates/forgot_password.html",
				map[string]string{
					"FirstName": user.FirstName,
					"link":      "token=" + req_token.Token,
				},
				user.Email,
			)
			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(ApiResource(true, nil, "Success"))
		}
	}
}

func ResetPassword(ctx iris.Context) {

	token := ctx.Params().GetString("token")
	request_token := models.GetRequestByToken(token)
	if request_token == nil {
		if request_token.Type == models.ForgotPassword {
			new_password := ctx.Params().GetString("new_password")
			confirm_password := ctx.Params().GetString("confirm_password")
			if new_password == confirm_password {
				// email := ctx.Params().GetString("email")
				user, err := models.UpdateUserPassword(new_password, request_token.UserID)
				if user == nil {
					ctx.StatusCode(iris.StatusBadRequest)
					ctx.JSON(ApiResource(false, err, "operation failed"))
				} else {
					ctx.StatusCode(iris.StatusOK)
					ctx.JSON(ApiResource(true, user, "Successful operation"))
				}
			} else {
				// password mismatch
			}
		} else {
			// Bad request
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(ApiResource(false, nil, "operation failed"))
		}
	} else {
		// Invalid token
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ApiResource(false, nil, "Invalid token"))
	}

}

func VerifyEmail(ctx iris.Context) {

	aul := new(RequestJson)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ApiResource(false, nil, "Invalid token"))
	} else {
		println(aul.token)
		request_token := models.GetRequestByToken(aul.token)
		println(request_token.Token)
		if request_token != nil {
			if request_token.Type == models.Signup {
				user, err := models.MarkUserVerified(request_token.UserID)
				if user == nil {
					ctx.StatusCode(iris.StatusBadRequest)
					ctx.JSON(ApiResource(false, err, "operation failed"))
				} else {
					// remove token after the suessful operation
					if err := models.DeleteRequestTokenByToken(aul.token); err == nil {
						ctx.StatusCode(iris.StatusOK)
						ctx.JSON(ApiResource(true, user, "Successful operation"))
					} else {
						ctx.StatusCode(iris.StatusBadRequest)
						ctx.JSON(ApiResource(false, err, "operation failed"))
					}
				}
			} else {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(ApiResource(false, nil, "invalid token"))
			}
		} else {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(ApiResource(false, nil, "Invalid token"))
		}
	}
}
