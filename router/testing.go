package router

import (
	"encoding/base64"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"ssh-api/common"
)

type testingBody struct {
	Host       string `json:"host" validate:"required,hostname|ip"`
	Port       uint   `json:"port" validate:"required,numeric"`
	Username   string `json:"username" validate:"required,alphanum"`
	Password   string `json:"password" validate:"required_without=PrivateKey"`
	PrivateKey string `json:"private_key" validate:"required_without=Password,omitempty,base64"`
	Passphrase string `json:"passphrase"`
}

func (app *Router) TestingRoute(ctx iris.Context) {
	var body testingBody
	ctx.ReadJSON(&body)
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		ctx.JSON(iris.Map{
			"error": 1,
			"msg":   err.Error(),
		})
		return
	}
	data, err := base64.StdEncoding.DecodeString(body.PrivateKey)
	if err != nil {
		ctx.JSON(iris.Map{
			"error": 1,
			"msg":   err.Error(),
		})
		return
	}
	client, err := app.Client.Testing(common.TestingOption{
		Host:     body.Host,
		Port:     body.Port,
		Username: body.Username,
		Key:      data,
	})
	if err != nil {
		ctx.JSON(iris.Map{
			"error": 1,
			"msg":   err.Error(),
		})
		return
	}
	defer client.Close()
	ctx.JSON(iris.Map{
		"error": 0,
		"msg":   "ok",
	})
}
