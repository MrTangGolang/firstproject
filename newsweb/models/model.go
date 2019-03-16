package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)


type User struct {
	Id        int
	UserName  string	`orm:"unique"`   //用户名全剧唯一，不能重复  unique
	Pwd       string
	Articles  []*Article  `orm:"rel(m2m)"`

}



type Article struct {
	Id        int       `orm:"Pk;auto"`     			//文章id  pk 主键   auto自增只适用整形
	Title     string	`orm:"size(100)"`   			//文章标题 字符长度100
	Content   string	`orm:"size(500)"`   			//文章内容 字符长度500
	Time      time.Time `orm:"type(datetime);auto_now"` //文章添加时间  auto_now添加时间 auto_now_add修改时间
	ReadCount int		`orm:"default(0)"`				//阅读量  创建完阅读量为0，设置默认0
	Image     string 	`orm:"null"`					//图片 正常orm是不准为空的，不想加图片时候设置允许为空null
	ArticleType *ArticleType `orm:"rel(fk);on_delete(set_null);null"`	//为文章类型设置外键;删除外键设置外键为空；允许外键为空
	Users      []*User	`orm:"reverse(many)"`
}


type ArticleType struct {
	Id        int
	TypeName  string	 `orm:"size(100)"`			//文章类型名
	Articles  []*Article `orm:"reverse(many)"`			//存储的文章
}




func init()  {
	//注册数据库
	orm.RegisterDataBase("default","mysql","root:linran1989a@tcp(127.0.0.1:3306)/newsweb?charset=utf8")

	//注册表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))

	//运行起来
	orm.RunSyncdb("default",false ,true)
}
