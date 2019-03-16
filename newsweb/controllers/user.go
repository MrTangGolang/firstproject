package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsweb/models"
)

type UserController struct {
	beego.Controller
}


//展示登陆页面
func (this *UserController)ShowLogin()  {
	//获取cookie数据 ，获取加密的秘文（key=userName=enc)
	dec:=this.Ctx.GetCookie("userName")
	//解密 enc是切片-dec切片-userName切片
	userName,_:=base64.StdEncoding.DecodeString(dec)

	//如果userName不等于空，说明设置了cookie值，我们将数据传递给前端页面
	//前端页面需要用value接收一下
	if string(userName)!=""{
		//非空说明记住密码，如果空，清除cookie数据
		this.Data["userName"]=string(userName)
		this.Data["checked"]="checked"
	}else{
		this.Data["userName"]=""
		this.Data["checked"]=""
	}



	this.TplName="login.html"
}


//展示注册页面
func (this *UserController)ShowRegister()  {
	this.TplName="register.html"
}


//注册业务
func (this *UserController)HandleRegister()  {
	//1、接收前端数据
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	//2、数据校验
	if userName==""||password==""{
		beego.Error("用户名或者密码不能为空")
		this.TplName="register.html"
		return
	}
	//3、操作数据 插入数据（将用户名和密码插入数据库）
	//1、获取orm对象
	o:=orm.NewOrm()
	//2、获取插入对象
	var user models.User
	//3、为插入对象赋值
	user.UserName=userName
	user.Pwd=password
	//4、插入数据
	_,err:=o.Insert(&user)
	if err!=nil{
		beego.Error("插入数据失败",err)
		this.TplName="register.html"
		return
	}

	//4、返回数据
	//this.Ctx.WriteString("注册成功")
	//this.TplName="login.html"
	this.Redirect("/login",302)
}


//登陆业务
func (this *UserController)HandleLogin()  {
	//1、获取前端数据
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	//2、校验数据
	if userName==""||password==""{
		beego.Error("用户名或密码不能为空")
		this.TplName="login"
		return
	}
	//3、操作数据  查询数据（将前端的用户名和密码和数据的用户名密码对比是否一致）
	//1、获取orm对象
	o:=orm.NewOrm()
	//2、获取查询对象
	var user models.User
	//3、给查询对象赋值
	user.UserName=userName
	user.Pwd=password
	//4、查询数据
	err:=o.Read(&user,"UserName")
	if err!=nil{
		beego.Error("用户名不存在，请重新输入")
		this.TplName="login"
		return
	}
	if user.Pwd!=password{
		beego.Error("用户名与密码不符，请重新输入")
		this.TplName="login"
		return
	}

//--------------------------------------------------------------------------
	//获取是否记住用户名
	remember:=this.GetString("remember")//此函数返回的是on
	if remember=="on"{
		//应用base64加密处理cookie不能存储中文的问题
		//此函数参数是切片 返回一个加密的秘文enc
		enc:=base64.StdEncoding.EncodeToString([]byte(userName))
		//设置cookie（"key值"，value，时间）
		this.Ctx.SetCookie("userName",enc,3600*1)//value为用户名，时间为1小时失效
	}else{
		this.Ctx.SetCookie("userName",userName,-1)//删除cookie  -1删除  0永久保留
	}
	//下次获取cookie的时候是显示登陆页面的时候
//------------------------------------------------------------------------------




	//4、返回数据
	//this.TplName="index.html"
	//登陆成功后设置session 用于左上角头像欢迎用户...
	//在展示文章列表页面获取数据
	this.SetSession("userName",userName)
	this.Redirect("/articlelist",302)
}


//处理退出登陆业务
func (this *UserController)Logout()  {
	//删除session
	this.DelSession("userName")
	//返回数据
	this.Redirect("/login",302)

}
