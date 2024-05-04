package router

import (
	"ginchat/docs"
	"ginchat/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	//登录首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	//跳转到注册页面
	r.GET("/toRegister", service.ToRegister)
	//登录后跳转到聊天首页
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", service.Chat)

	//查找好友
	r.POST("/searchFriends", service.SearchFriends)

	//用户模块
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	r.POST("/attach/upload", service.Upload)

	//添加好友
	r.POST("/contact/addfriend", service.AddFriend)
	return r
}
