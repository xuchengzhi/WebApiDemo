package main

import (
	// _ "BeeGoApi/routers"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type Response struct {
	Code int
	Msg  string
	Data string
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Data["json"] = Response{200, "success.beego", "ok"}
	c.ServeJSON()
}

func main() {
	beego.BConfig.Listen.HTTPPort = 80
	beego.Router("/", &MainController{})
	beego.Run()
}
