package controllers

import (
	"classOne/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type RegController struct {
	beego.Controller
}

type LoginController struct {
	beego.Controller
}

func (this* RegController) ShowReg(){
	this.TplName="register.html"
}

func (this* LoginController) ShowLogin(){
	name:=this.Ctx.GetCookie("userName")
	if name !=""{
		passwd:=this.Ctx.GetCookie("passwd")
		this.Data["userName"] =name
		this.Data["passwd"] = passwd
		this.Data["check"] = "checked"
	}

	this.TplName="login.html"
}

func (this* LoginController) HandleLogin(){
	name := this.GetString("userName")
	passwd := this.GetString("password")
	if name=="" || passwd=="" {
		beego.Info("不能为空")
		this.TplName="register.html"
		this.TplName="login.html"
		return
	}

	o:=orm.NewOrm()
	user := models.User{}
	user.UserName=name
	err:=o.Read(&user,"userName")
	if err != nil {
		beego.Info("登录失败")
		this.TplName="login.html"
		return
	}
	if user.Passwd != passwd{
		beego.Info("登录失败")
		this.TplName="login.html"
		return
	}

	check:=this.GetString("remember")

   if check =="on"{
	   this.Ctx.SetCookie("userName",name,time.Second*3600)
	   this.Ctx.SetCookie("passwd",passwd,time.Second*3600)
   }else{
	   this.Ctx.SetCookie("userName","s",-1)
	   this.Ctx.SetCookie("passwd","s",-1)
   }

   this.SetSession("userName",name)

	this.Redirect("/Article/ShowArticle",302)
}

func (this* RegController) HandleReg(){
	name := this.GetString("userName")
	passwd := this.GetString("password")

	if name=="" || passwd=="" {
		beego.Info("不能为空")
		this.TplName="register.html"
		return
	}
	o:=orm.NewOrm()
	user := models.User{}
	user.UserName=name
	user.Passwd = passwd
	_,err:=o.Insert(&user)
    if err != nil{
		beego.Info("失败")
	}


	this.Redirect("/",302)
}

