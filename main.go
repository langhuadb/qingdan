package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
var(
	db *gorm.DB
)
type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"`
}
func initMySQL()(err error){
	dsn:="root:123@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	db,err = gorm.Open("mysql",dsn)
	if err != nil {
		return
	}
	return  db.DB().Ping()
}
func main() {
	//创建数据库
	//sql:create database bubble
	//连接数据库
	err:=initMySQL()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//模型绑定
	db.AutoMigrate(&Todo{})
	r := gin.Default()
	r.Static("/static","static")
	// 告诉gin框架去哪里找模板文件
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK,"index.html",nil)
	})

	//v1
	v1Group := r.Group("v1")
	{
		//待办事项
		//添加
		v1Group.POST("/todo", func(c *gin.Context) {
			var todo Todo
			c.BindJSON(&todo)
			if err = db.Create(&todo).Error;err != nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
			}else {
				c.JSON(http.StatusOK,todo)
			}
		})
		//查看
		//查看所有待办事项
		v1Group.GET("/todo", func(c  *gin.Context) {
			//查询todo这个表里所有数据
			var todolist []Todo
			if err = db.Find(&todolist).Error;err!=nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
			}else {
				c.JSON(http.StatusOK,todolist)
			}
		})
		//查看某一个待办事项
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})
		//修改
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id,ok:= c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK,gin.H{"error":"无效的id"})
				return
			}
			var todo Todo
			if err = db.Where("id=?",id).First(&todo).Error;err != nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
				return
			}
			c.BindJSON(&todo)
			if err = db.Save(&todo).Error;err != nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
			}else {
				c.JSON(http.StatusOK,todo)
			}
		})
		//删除
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id,ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK,gin.H{"error":"无效的id"})
				return
			}
			if err = db.Where("id=?",id).Delete(Todo{}).Error;err!=nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
			}else {
				c.JSON(http.StatusOK,gin.H{id:"deleted"})
			}
		})
	}

	r.Run(":9000")
}
