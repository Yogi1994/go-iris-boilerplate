package controllers

import (
	"fmt"
	// "runtime/debug"
	"MVC-golang/models"

	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
)

func GetUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	user, err := models.GetUserById(id)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ApiResource(false, err, "Failed"))
	} else {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(true, user, "Successful operation"))
	}
}

func GetAllUsers(ctx iris.Context) {
	offset := ctx.URLParamIntDefault("offset", 1)
	limit := ctx.URLParamIntDefault("limit", 10)
	name := ctx.URLParam("name")
	orderBy := ctx.URLParam("orderBy")

	users := models.GetAllUsers(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(true, users, "OK"))
}

func UpdateUser(ctx iris.Context) {
	aul := new(models.UpdateUserJson)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
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
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(ApiResource(false, err, "Operation failed"))
		} else {
			id, _ := ctx.Params().GetInt("id")
			uid := uint(id)

			u, u_err := models.UpdateUser(aul, uid)

			if u.ID == 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(ApiResource(false, u_err, "Operation failed"))
			} else {
				ctx.StatusCode(iris.StatusOK)
				ctx.JSON(ApiResource(true, u, "Successful"))
			}
		}
	}
}
