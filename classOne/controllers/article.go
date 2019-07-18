package controllers

import (
	"bytes"
	"classOne/models"
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"github.com/gomodule/redigo/redis"
	"math"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

//func (this* ArticleController)HandleSelect()  {
//    typeName:=this.GetString("select")
//
//    if typeName == ""{
//		beego.Info("失败")
//		return
//	}
//	o:=orm.NewOrm()
//	var articles []models.Article
//	o.QueryTable("Article").Filter("ArticleType__TypeName",typeName).All(&articles)
//	beego.Info(articles)
//
//}

func (this* ArticleController) ShowArticleList(){

	o:=orm.NewOrm()
	qs := o.QueryTable("Article")
	//var articles []models.Article
	//qs.All(&articles)

	count,err := qs.RelatedSel("ArticleType").Count()
	if err != nil{
		beego.Info("查询错误")
		return
	}

	pageIndex1 := this.GetString("pageIndex")

	pageIndex,err:= strconv.Atoi(pageIndex1)


	if err != nil {
		pageIndex =1
	}

	FirstPage := false
	if pageIndex == 1{
		FirstPage = true
	}
	var types []models.ArticleType

	conn,_:=redis.Dial("tcp","182.61.50.167:6379")
	defer conn.Close()
	conn.Do("AUTH", "xiuchao16")
	buffer,err:=redis.Bytes(conn.Do("get","types"))
	if err != nil{
		beego.Info("redis错误")
	}
	dec := gob.NewDecoder(bytes.NewReader(buffer))
	dec.Decode(&types)


	//if len(types)==0{
	//	o.QueryTable("ArticleType").All(&types)
	//	var buffer bytes.Buffer
	//	enc :=gob.NewEncoder(&buffer)
	//	enc.Encode(types)
	//	_,err=conn.Do("set","types",buffer.Bytes())
	//	if err != nil{
	//		beego.Info("数据库错误")
	//		return
	//	}
	//}


	this.Data["Types"] = types


    pageSize := 2
    start:=pageSize*(pageIndex-1)
	pageCount:=float64(count)/float64(pageSize)
	pageCount1:=math.Ceil(pageCount)
	LastPage := false
	if pageIndex == int(pageCount1){
		LastPage = true
	}

	var articleswithtype []models.Article
	typeName := this.GetString("select")
	if typeName == "" || typeName == "全部"{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articleswithtype)
	}else{
		count,_ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		pageCount=float64(count)/float64(pageSize)
		pageCount1=math.Ceil(pageCount)
		LastPage = false
		if pageIndex == int(pageCount1){
			LastPage = true
		}
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articleswithtype)
	}

	userName:=this.GetSession("userName")
	this.Data["userName"] = userName

	this.Data["LastPage"] = LastPage

	this.Data["typeName"] = typeName

	this.Data["pageIndex"] = pageIndex

	this.Data["FirstPage"] = FirstPage
    this.Data["articles"] = articleswithtype
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount1

	this.Layout="layout.html"

	this.TplName="index.html"
}

func (this* ArticleController) ShowContent(){
	o:=orm.NewOrm()
	id2 := this.GetString(":id")

	id,_:=strconv.Atoi(id2)

	article := models.Article{Id:id}
	err:=o.Read(&article)
	if err != nil{
		beego.Info("数据为空")
		return
	}
	article.Count+=1

	m2m:=o.QueryM2M(&article,"Users")
	userName:=this.GetSession("userName")
	user:=models.User{}
	user.UserName = userName.(string)
	o.Read(&user,"userName")
	_,err=m2m.Add(&user)
	if err != nil{
		beego.Info("插入失败")
	}
	o.Update(&article)


	//o.LoadRelated(&article,"Users")
	//o.QueryTable("Article").RelatedSel("User").Filter("Users__User__UserName",userName.(string)).Distinct().Filter("Id",id).One(&article)
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)

	this.Data["users"] = users
	this.Data["articles"] = article

	this.Layout="layout.html"

	this.LayoutSections=make(map[string]string)
	this.LayoutSections["contentHead"]="head.html"

	this.TplName="content.html"
}


func (this* ArticleController) ShowAddArticleList(){
	o:=orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["Types"] = types
	this.Layout="layout.html"
	this.TplName="add.html"
}

func (this* ArticleController) HandleUpdate(){
     name:=this.GetString("articleName")
	content:=this.GetString("content")
	id,_:=this.GetInt("id")


	if name == "" ||  content ==""{
		beego.Info("更新失败")
		return
	}

	o:=orm.NewOrm()
	article := models.Article{Id:id}
	err1 := o.Read(&article)
	if err1 != nil{
		beego.Info("文章不存在")
		return
	}

	article.Title= name
	article.Content= content

	f,h,err:=this.GetFile("uploadname")

	defer f.Close()
	if err != nil{
		beego.Info("上传失败")
	}else{
		ext:=path.Ext(h.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
			beego.Info("格式不正确")
			return
		}
		if h.Size>5000000{
			beego.Info("文件太大")
			return
		}
		fileName:=time.Now().Format("2006-01-02")
		err=this.SaveToFile("uploadname","./static/img/"+fileName+".jpg")
		if err != nil{
			beego.Info("上传失败")
			return
		}
		article.Img="/static/img/"+fileName+".jpg"
	}

	_,err = o.Update(&article)
	if err != nil{
		beego.Info("更新失败")
		return
	}
	this.Redirect("/Article/ShowArticle",302)
}


func (this* ArticleController) ShowUpdate(){
	id:=this.GetString("id")
	if id == ""{
		beego.Info("链接错误")
		return
	}
	o:=orm.NewOrm()
	article := models.Article{}
	id2,_:=strconv.Atoi(id)
	article.Id=id2
	err:=o.Read(&article)
	if err != nil{
		beego.Info("查询错误")
	}
	this.Data["article"] = article
	this.TplName="update.html"


}

func (this* ArticleController) HandleDelete(){

   id,_:=this.GetInt("id")
	o:=orm.NewOrm()
	article := models.Article{Id:id}
	o.Delete(&article)
	this.Redirect("/Article/ShowArticle",302)

}

func (this* ArticleController)HandleAddArticle(){
	artiName := this.GetString("articleName")
	artiContent := this.GetString("content")

	f,h,err:=this.GetFile("uploadname")
	defer f.Close()


	ext:=path.Ext(h.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Info("格式不正确")
		return
	}
	if h.Size>5000000{
		beego.Info("文件太大")
		return
	}
	fileName:=time.Now().Format("2006-01-02")

	err=this.SaveToFile("uploadname","./static/img/"+fileName+".jpg")
	if err != nil{
		beego.Info("上传失败")
		return
	}

	o:=orm.NewOrm()
	article := models.Article{}
	article.Title=artiName
	article.Content=artiContent
	article.Img="/static/img/"+fileName+".jpg"

	typeName:=this.GetString("select")
	if typeName == ""{
		beego.Info("类型数据失败")
		return
	}
	var artiType models.ArticleType
	artiType.TypeName = typeName
	err =o.Read(&artiType,"typeName")
	if err != nil{
		beego.Info("类型错误")
		return
	}
	article.ArticleType = &artiType
	_,err=o.Insert(&article)
	if err != nil{
		beego.Info("插入数据失败")
		return
	}
	this.Redirect("/Article/ShowArticle",302)

}

func (this* ArticleController)ShowAddType()  {
	o:=orm.NewOrm()
	var artiTypes[]models.ArticleType
	_,err:=o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil{
		beego.Info("查询类型失败")
	}
	this.Data["types"] = artiTypes

	this.Layout="layout.html"
	this.TplName="addType.html"
}

func (this* ArticleController)HandleAddType()  {
	typename:=this.GetString("typeName")
	if typename == ""{
		beego.Info("类型数据为空")
		return
	}
	o:=orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName = typename
	_,err:=o.Insert(&artiType)
	if err != nil{
		beego.Info("插入失败")
		return
	}
	this.Redirect("/Article/AddArticleType",302)
}

func (this* ArticleController)Logout()  {
	this.DelSession("userName")
	this.Redirect("/",302)
	return
}

func (this* ArticleController)SendMail(){
	config :=`{"username":"1162750474@qq.com","password":"vlpofbboxrcticgc","host":"smtp.qq.com","port":587}`
	email := utils.NewEMail(config)
	email.From="1162750474@qq.com"
	email.To = []string{"746152132@qq.com"}
	email.Subject="激活验证邮件"
	email.Text="http://localhost:8080/"
	email.HTML="<h1>特别提示</h1>"
	email.Send()
}