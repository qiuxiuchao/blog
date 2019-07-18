package main

import (
	_ "classOne/routers"
	"github.com/astaxie/beego"
	_ "classOne/models"
	"strconv"
)

func main() {
	beego.AddFuncMap("ShowPrePage",HandlePrePage)
	beego.AddFuncMap("ShowNextPage",HandleNextPage)
	beego.Run()
}

func HandlePrePage(data int)(string){
	pageIndex := data - 1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}

func HandleNextPage(data int)(int){
	pageIndex := data + 1
	//pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex
}

