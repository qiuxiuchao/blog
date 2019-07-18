package routers

import (
	"classOne/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
    //beego.Router("/", &controllers.MainController{})
	beego.InsertFilter("/Article/*",beego.BeforeRouter,FilterFunc)
	beego.Router("/register", &controllers.RegController{},"get:ShowReg;post:HandleReg")
	beego.Router("/", &controllers.LoginController{},"get:ShowLogin;post:HandleLogin")

	beego.Router("/Article/ShowArticle", &controllers.ArticleController{},"get:ShowArticleList")

	beego.Router("/Article/DeleteArticle", &controllers.ArticleController{},"get:HandleDelete")

	beego.Router("/Article/ShowContent/:id", &controllers.ArticleController{},"get:ShowContent")

	beego.Router("/Article/AddArticle", &controllers.ArticleController{},"get:ShowAddArticleList;post:HandleAddArticle")

	beego.Router("/Article/UpdateArticle", &controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")

	beego.Router("/Article/AddArticleType", &controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")

	beego.Router("/Logout", &controllers.ArticleController{},"get:Logout")

	beego.Router("/SendMail", &controllers.ArticleController{},"get:SendMail")
}

var FilterFunc = func(ctx *context.Context) {
	userName:=ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302,"/")
	}
}
