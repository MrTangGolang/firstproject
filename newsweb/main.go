package main

import (
	_ "newsweb/routers"
	"github.com/astaxie/beego"
	_"newsweb/models"
)

func main() {
	//视图函数第三步    在beego.Run之前
	//将前端第PrePage函数和下面第PrePageIndex函数建立关系
	//参数（前端第函数名，下面第函数名没有括号）
	beego.AddFuncMap("PrePage",PrePageIndex)
	beego.AddFuncMap("NextPage",NextPageIndex)


	beego.Run()
}



//创建视图函数第二步  定义一个实现视图的函数
//上一页  参数（当前页 类型int）返回值前一页类型
func PrePageIndex(pageindex int)int {
	PrePage:=pageindex-1
	//当我们点击上一页，发现可以出现0.-1.-2.-3...
	if PrePage<1{
		PrePage=1
	}
	return PrePage
}
//下一页  参数（下一页  类型int）返回值下一个类型
//当我们点击下一页，发现可以出现NextPage>pagecount的情况 所以我们要用另一种方法创建视图函数
func NextPageIndex(pageindex int,pagecount float64)int  {
	NextPage:=pageindex+1

	if NextPage>int(pagecount){
		NextPage=int(pagecount)
	}
	return NextPage

}

