package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

type User struct {
	Id int
	UserName string
	Passwd string
	Articles[]*Article `orm:"rel(m2m)"`
}

type Article struct {
	Id int `orm:"pk:auto"`
	Title string `orm:"size(20)"`
	Content string `orm:"size(500)"`
	Img string `orm:"size(50);null"`
	Time time.Time `orm:"type(datetime);auto_now_add"`
	Count int `orm:"default(0)"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users[]*User `orm:"reverse(many)"`
}

type ArticleType struct{
	Id int
	TypeName string `orm:"size(20)"`
	Articles[]*Article  `orm:"reverse(many)"`

}


	func init()  {

	orm.RegisterDataBase("default","mysql","mysql:xiuchao16@tcp(182.61.50.167:3306)/c2c")
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Article))
		orm.RegisterModel(new(ArticleType))
	orm.RunSyncdb("default",false,true)

}
