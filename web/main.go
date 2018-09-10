package main

import (
	"fmt"
	"time"

	"goplays/web/app/controllers"
	_ "goplays/web/app/mail"
	"goplays/web/app/service"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

const VERSION = "2.0.1"

func main() {
	//service.Init()
	service.InitMgo()

	beego.AppConfig.Set("version", VERSION)
	runmode := beego.AppConfig.String("runmode")
	if runmode == "dev" {
		beego.SetLevel(beego.LevelDebug)
	} else {
		beego.SetLevel(beego.LevelInformational)
		beego.SetLogger("file", `{"filename":"`+beego.AppConfig.String("log_file")+`"}`)
		beego.BeeLogger.DelLogger("console")
	}

	// 记录启动时间
	beego.AppConfig.Set("up_time", fmt.Sprintf("%d", time.Now().Unix()))

	beego.AddFuncMap("i18n", i18n.Tr)

	/*
		beego.Router("/", &controllers.MainController{}, "*:Index")
		beego.Router("/login", &controllers.MainController{}, "*:Login")
		beego.Router("/logout", &controllers.MainController{}, "*:Logout")
		beego.Router("/profile", &controllers.MainController{}, "*:Profile")
		beego.Router("/regist", &controllers.MainController{}, "*:Regist")
		beego.Router("/servers", &controllers.MainController{}, "*:Servers")
		beego.Router("/files", &controllers.MainController{}, "*:Files")

		beego.AutoRouter(&controllers.UserController{})
		beego.AutoRouter(&controllers.RoleController{})
		beego.AutoRouter(&controllers.MailTplController{})
		beego.AutoRouter(&controllers.MainController{})
		beego.AutoRouter(&controllers.PlayerController{})
		beego.AutoRouter(&controllers.LoggerController{})
		beego.AutoRouter(&controllers.AgencyController{})
		beego.AutoRouter(&controllers.ChartsController{})

		beego.ErrorController(&controllers.ErrorController{})

		beego.SetStaticPath("/assets", "assets")
		beego.SetStaticPath("/contract", "views/main/contract.html")
		beego.SetStaticPath("/rules", "views/main/rules.html")
		beego.SetStaticPath("/download", "views/main/download.html")
		//beego.SetStaticPath("/headimag", "headimag")
		beego.SetStaticPath("/poster", "views/main/poster.html")
	*/

	//namespace
	namespace := beego.AppConfig.String("namespace")
	beego.Trace("namespace: ", namespace)

	//初始化 namespace
	ns :=
		beego.NewNamespace("/"+namespace,

			beego.NSRouter("/", &controllers.MainController{}, "*:Index"),
			beego.NSRouter("/login", &controllers.MainController{}, "*:Login"),
			beego.NSRouter("/logout", &controllers.MainController{}, "*:Logout"),
			beego.NSRouter("/profile", &controllers.MainController{}, "*:Profile"),
			beego.NSRouter("/regist", &controllers.MainController{}, "*:Regist"),
			beego.NSRouter("/servers", &controllers.MainController{}, "*:Servers"),
			beego.NSRouter("/files", &controllers.MainController{}, "*:Files"),

			beego.NSAutoRouter(&controllers.UserController{}),
			beego.NSAutoRouter(&controllers.RoleController{}),
			beego.NSAutoRouter(&controllers.MailTplController{}),
			beego.NSAutoRouter(&controllers.MainController{}),
			beego.NSAutoRouter(&controllers.PlayerController{}),
			beego.NSAutoRouter(&controllers.LoggerController{}),
			beego.NSAutoRouter(&controllers.AgencyController{}),
			beego.NSAutoRouter(&controllers.ChartsController{}),
		)
	//注册 namespace
	beego.AddNamespace(ns)

	beego.ErrorController(&controllers.ErrorController{})

	beego.SetStaticPath("/"+namespace+"/assets", "assets")
	beego.SetStaticPath("/"+namespace+"/contract", "views/main/contract.html")
	beego.SetStaticPath("/"+namespace+"/termsofservice", "views/main/termsofservice.html")
	beego.SetStaticPath("/"+namespace+"/privacypolicy", "views/main/privacypolicy.html")
	beego.SetStaticPath("/"+namespace+"/agencyAgreement", "views/main/agencyAgreement.html")
	//beego.SetStaticPath("/admin/rules", "views/main/rules.html")
	//beego.SetStaticPath("/admin/download", "views/main/download.html")
	//beego.SetStaticPath("/admin/headimag", "headimag")
	beego.SetStaticPath("/"+namespace+"/poster", "views/main/poster.html")

	beego.Run()
}
