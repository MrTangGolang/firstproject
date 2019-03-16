package routers

import (
	//"context"  编辑器自动导的包是错误的  需要自己手动导beego的context
	//"github.com/astaxie/beego/context"
	"github.com/astaxie/beego"
	"newsweb/controllers"
)

func init() {
	//在执行控制器之前，找到路由之后添加过滤器 所有article 开头的都要过滤
	//                     参数（"路径前缀/所有， 过滤器位置，过滤器函数名
	//beego.InsertFilter("/article/*",beego.BeforeExec,funcFilter)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	beego.Router("/articlelist",&controllers.ArticleController{},"get:ShowArticleList")
    beego.Router("//addarticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/articledetail",&controllers.ArticleController{},"get:ShowArticleDetail")
    beego.Router("/articleupdate",&controllers.ArticleController{},"get:ShowArticleUpdate;post:HandleArticleUpdate")
    beego.Router("/articledelete",&controllers.ArticleController{},"get:DeletArticcle")
    beego.Router("/addtype",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
	beego.Router("/deletetype",&controllers.ArticleController{},"get:HandleDeleteType")
    beego.Router("/logout",&controllers.UserController{},"get:Logout")
}

////创建过滤器函数
//var funcFilter = func(ctx*context.Context){
//	//登陆校验   从session中获取用户名
//	userName:=ctx.Input.Session("userName")
//	//为空，没有登陆 返回导登陆页
//	if userName==nil{
//		//注意ctx里导redirect函数参数相反（状态码，"路径"）
//		ctx.Redirect(302,"/login")
//	}
//}