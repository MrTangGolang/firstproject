package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"newsweb/models"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

//展示文章列表页面
func (this *ArticleController) ShowArticleList() {
	//获取session数据
	userName:=this.GetSession("userName")
	//如果session数据为空，说明没有登陆，那数据指向登陆界面
	//没有登陆，就别访问我们的界面，使用session时，需要在aap中开启session
	if userName==nil{
		this.Redirect("/login",302)
		return
	}
	//将session数据传递给页面   因为上面的userName是接口类型，前端要的是字符串类型，所以要断言
	this.Data["userName"]=userName.(string)



	//获数据库所有文章 展示给页面
	//查询数据
	//获取orm 对象
	o:=orm.NewOrm()
	//获取查询对象  所有文章用切片存储
	var articles []models.Article
	//查询表   指定查询Article这张表   返回一个表
	qs:=o.QueryTable("Article")
	//查询表内的所有文章 获取分页数据时只获取当前页显示的文章，所以不再查询所有
	//qs.All(&articles)
	//将数据传递给视图
	this.Data["article"]=articles
	//------------------------以上没有问题------------------------



	//获取总记录数和总页数
	//获取总记录数
	count,_:=qs.Count()


	//获取总页数
	pagesize:=int64(2)            //每页显示数量
	pagecount:=float64(count)/float64(pagesize)	  //总页数=总记录数/每页显示数量  整型/整型 省略小数  所以我们转为浮点型

	//因为总页数只能为整数， 按照实际情况向上取整math.Ceil  向下取整math.Floor()
	pagecount=math.Ceil(pagecount)

	//将总记录数总页数传递一个页面
	this.Data["count"]=count
	this.Data["pagecount"]=pagecount

	//----------------------------以上没有问题----------------------------

	//实现首页和末页
	//获取前端首页数据
	pageindex,err:=this.GetInt("pageindex")
	if err!=nil{
		pageindex=1
	}
	//获取分页数据
	start:=pagesize*(int64(pageindex)-1)
	//查询部分数据qs.Limit(页面容量，从哪里开始)
	//RelateSel 一对多关系表查询中，用来指定另外一张表多函数
	qs.Limit(pagesize,start).RelatedSel("ArticleType").All(&articles)
	//将数据传递给视图
	//this.Data["article"]=articles
	//----------------------------以上没有问题------------------------------

	//将分页数据传递给视图
	//this.Data["pageindex"]=pageindex

	//---------------------以上是分页实现，以下是实现分类下拉条-------------------


	//根据下拉框获取分类文章
	//获取数据
	typename:=this.GetString("select")
	//把选中的类型显示在下拉框中，将数据传递给前端页面
	this.Data["typename"]=typename

	//查询的文章，想要找类型对应的文章，相当于文章关联的类型表
	//查询部分数据qs.Limit(页面容量，从哪里开始).指定关联表RelatedSel("指定表的表名").过滤器Filter("指定表的名称__按字段查"，查询字段的值）.All(所有的文章）
	qs.Limit(pagesize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typename).All(&articles)
	//将数据传递给视图
	this.Data["article"]=articles
	//将分页数据传递给视图
	this.Data["pageindex"]=pageindex















	//实现分类下拉条
	//创建存储容器
	var articletypes []models.ArticleType
	//查询所有分类（指定查询表）（放到存储容器中）
	o.QueryTable("ArticleType").All(&articletypes)
	//将数据传递给视图
	this.Data["articletypes"]=articletypes


	//指定模版页
	this.Layout="layout.html"

	//指定视图
	this.TplName="index.html"


}

//展示文章添加页面
func (this *ArticleController)ShowAddArticle()  {
	//获取文章分类信息
	//获取orm对象
	o:=orm.NewOrm()
	//创建存储容器
	var articletypes []models.ArticleType
	//获取文章分类并存入容器   (指定查询的表）(将获取的所有数据存入容器）
	o.QueryTable("ArticleType").All(&articletypes)
	//将数据返回给视图
	this.Data["articletypes"]=articletypes
	//-------------------------------------------------------------------




	this.TplName="add.html"
}

//处理文章添加业务
func (this *ArticleController)HandleAddArticle()  {
	//1、获取数据
	//获取文章数据
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	//获取图片数据
	file,head,err:=this.GetFile("uploadname")
	//文件流 需要关闭
	defer file.Close()


	//------获取文章类型----------------------------------
	typename:=this.GetString("select")
	//--------------------------------------------------

	//2、校验数据
	//校验文章数据
	if articleName==""||content==""{
		this.Data["errmsg"]="文章标题或文章内容不能为空，请重新输入"
		this.TplName="add.html"
		return
	}
	//校验图片数据
	if err!=nil{
		this.Data["errmsg"]="添加图片失败，请重新上传"
		this.TplName="add.html"
		return
	}

	//校验图片大小
	if head.Size>500000{
		this.Data["errmsg"]="图片尺寸过大，请重新上传"
		this.TplName="add.html"
		return
	}
	//校验图片格式
	fileExt:=path.Ext(head.Filename)
	if fileExt!=".jpg"&&fileExt!=".png"&&fileExt!=".jpeg"{
		this.Data["errmsg"]="图片格式不正确，请重新上传"
		this.TplName="add.html"
		return
	}
	//防止图片名重复
	filename:=time.Now().Format("2006-01-02-15:04:05")+fileExt
	//保存图片
	this.SaveToFile("uploadname","./static/image/"+filename)



	//3、操作数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取插入对象
	var article models.Article
	//为插入对象赋值
	article.Title=articleName
	article.Content=content
	article.Image="/static/image/"+filename


	//--------------根据类型名称获取类型对象--------------------
	var articletype models.ArticleType
	articletype.TypeName=typename
	o.Read(&articletype,"TypeName")
	//----为插入对象分类  直接赋值发现类型不匹配
	//article.ArticleType=typename
	article.ArticleType=&articletype
	//-----------------------------------------------------------



	//插入数据
	_,err=o.Insert(&article)
	if err!=nil{
		this.Data["errmsg"]="添加文章失败，请重新添加"
		this.TplName="add.html"
		return
	}




	//4、返回数据
	this.Redirect("/articlelist",302)
}


//展示文章详情
func (this *ArticleController)ShowArticleDetail()  {
	//获取文章ID数据
	articleId,err:=this.GetInt("id")
	//校验数据
	if err!=nil{
		this.Data["errmsg"]="获取文章路失败"
		this.TplName="index.html"
		return
	}

	//查询数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取查询对象
	var article models.Article
	//为查询对象赋值
	article.Id=articleId
	//查询
	err=o.Read(&article,"Id")
	if err!=nil{
		this.Data["errmsg"]="获取文章路失败"
		this.TplName="index.html"
		return
	}

	//-------------------------------------------------
	//获取article对象，上面已经获取完
	//获取多对多操作对象 (操作对象，"多对多的关联表字段"
	m2m:=o.QueryM2M(&article,"Users")
	//获取要插入的数据
	var user models.User
	//从session中获取userName 接口类型
	userName:=this.GetSession("userName")
	//给插入对象赋值
	user.UserName=userName.(string)//断言
	//查询userName是否存在，之前的userName存在于session中
	o.Read(&user,"UserName")
	//插入多对多关系
	m2m.Add(user)
	//加载关系   //在文章详情页展示 文章里面查用户  第一种方法
	o.LoadRelated(&article,"Users")


	//第二种多对多关系
	var users []models.User
	//o.QueryTable("指定查询表"),Filter("User里的Articles字段 双下划线  表名  双下划线  表中字段"，"指定查询字段").去重函数.All(所有用户)
	o.QueryTable("User").Filter("Articles__Article__Id",articleId).Distinct().All(&users)
	//将数据传递给视图
	this.Data["users"]=users
	//-----------------------------------------------------




	//返回数据
	this.Data["article"]=article
	this.TplName="content.html"

}


//展示编辑文章文章详情业务
func (this *ArticleController)ShowArticleUpdate()  {
	//获取数据
	articleid,err:=this.GetInt("id")


	//获取错误信息
	errmsg:=this.GetString("errmsg")
	if errmsg!=""{
		this.Data["errmsg"]=errmsg
	}



	//校验数据
	if err!=nil{
		//这里用渲染不合适
		//this.Data["errmsg"]="访问路径错误"
		//this.TplName="index.html"
		//return
		errmsg:="访问路径错误"
		this.Redirect("/articlelist?=errmsg"+errmsg,302)
		return
	}

	//查询数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取查询对象
	var article models.Article
	//为查询对象赋值
	article.Id=articleid
	//查询数据
	o.Read(&article,"Id")

	//返回数据
	this.Data["article"]=article
	this.TplName="update.html"

}


//处理文章编辑业务
func (this *ArticleController)HandleArticleUpdate()  {
	//获取数据
    articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	file,head,err:=this.GetFile("uploadname")
	articleid,err2:=this.GetInt("id")
	//校验数据
	if articleName==""||content==""||err!=nil||err2!=nil{
		errmsg:="文章标题或内容不能为空"
		this.Redirect("/articleupdate?id="+strconv.Itoa(articleid)+"&errmsg="+errmsg,302)
		return
	}
	defer file.Close()
	//校验图片尺寸
	if head.Size>50000{
		errmsg:="图片尺寸过大，请重新上传"
		this.Redirect("/articleupdate?id="+strconv.Itoa(articleid)+"&errmsg="+errmsg,302)
		return
	}
	//校验图片格式
	fileExt:=path.Ext(head.Filename)
	if fileExt!=".jpg"&&fileExt!=".png"&&fileExt!=".jpeg"{
		errmsg:="图片格式错误，请重新上传"
		this.Redirect("/articleupdate?id="+strconv.Itoa(articleid)+"&errmsg="+errmsg,302)
		return
	}
	//防止图片名重复
	filename:=time.Now().Format("2006-01-02-15:04:05")+fileExt
	this.SaveToFile("uploadname","./static/image"+filename)

	//更新数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取更新对象
	var article models.Article
	//给更新对象赋值
	article.Id=articleid
	//更新
	err=o.Read(&article,"Id")
	if err!=nil{
		errmsg:="更新文章不存在"
		this.Redirect("/articleupdate?id="+strconv.Itoa(articleid)+"&errmsg="+errmsg,302)
		return
	}
	article.Title=articleName
	article.Content=content
	article.Image="/static/image"+filename
	o.Update(&article)


	//返回数据
	//this.Data["article"]=article
	//this.TplName="update.html"
	this.Redirect("/articlelist",302)
}


//处理删除业务
func (this *ArticleController)DeletArticcle()  {
	//获取数据
	articleid,err:=this.GetInt("id")
	//校验数据
	if err!=nil{
		beego.Error("访问路径错误")
		this.Redirect("/articlelist",302)
		return
	}
	//删除数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取删除对象
	var article models.Article
	//为删除对象赋值
	article.Id=articleid
	//删除数据
	_,err=o.Delete(&article)
	if err!=nil{
		beego.Error("删除数据失败")
		this.Redirect("/articlelist",302)
		return
	}

	//返回数据
	this.Redirect("/articlelist",302)

}


//展示文章分类页面
func (this *ArticleController)ShowAddType()  {
	//获取所有类型，并展示
	//获取orm对象
	o:=orm.NewOrm()
	//创建存储容器
	var articletypes []models.ArticleType
	//指定查询表
	qs:=o.QueryTable("ArticleType")
	//查询所有类型
	qs.All(&articletypes)
	//返回数据给视图
	this.Data["articletypes"]=articletypes


	this.TplName="addType.html"

}


//处理文章分类业务
func (this *ArticleController)HandleAddType()  {
	//获取数据
	typeName:=this.GetString("typeName")
	//校验数据
	if typeName==""{
		this.Data["errmsg"]="分类标题内容不能为空"
		this.Redirect("/addtype",302)
		return
	}
	//插入数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取插入对象
	var articletype models.ArticleType
	//为插入对象赋值
	articletype.TypeName=typeName
	//插入数据
	_,err:=o.Insert(&articletype)
	if err!=nil{
		this.Data["errmsg"]="添加文章分类失败，请重新添加"
		this.Redirect("/addtype",302)
		return
	}

	//返回数据
	this.Redirect("/addtype",302)

}


//删除文章分类业务
func (this *ArticleController)HandleDeleteType()  {
	//获取数据
	articleid,err:=this.GetInt("id")

	//校验数据
	if err!=nil{
		this.Data["errmsg"]="您所要删除的文章不存在"
		this.Redirect("/addtype",302)
		return
	}
	//删除数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取删除对象
	var articletype models.ArticleType
	//为删除对象赋值
	articletype.Id=articleid
	//删除数据
	_,err=o.Delete(&articletype)
	if err!=nil{
		this.Data["errmsg"]="删除文章失败"
		this.Redirect("/addtype",302)
		return
	}
	//返回数据
	this.Redirect("/addtype",302)

}
