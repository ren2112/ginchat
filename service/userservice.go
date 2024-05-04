package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetIUserList
// @Summary 所有用户
// @Tags 用户模块
// Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": "查询成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param repassword formData string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")
	fmt.Println("repassword:", repassword)
	user.LoginTime = time.Now()
	user.HeartbeatTime = time.Now()
	user.LogOutTime = time.Now()

	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名或者密码不能为空！",
			"data":    user,
		})
		return
	}

	salt := fmt.Sprintf("%06d", rand.Int31())

	//检验名字是否被使用
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名已存在!",
			"data":    user,
		})
		return
	}

	//验证与输入的密码是否一致
	if repassword != password {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
			"data":    user,
		})
		return
	}

	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "新增用户成功",
		"data":    user,
	})
}

// FindUserByNameAndPwd
// @Summary 校验用户名字和密码
// @Tags 用户模块
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}

	//name := c.Query("name")
	//password := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	//获取用户密码，因为接下来要对输入的密码和数据库的密码进行校验（使用md5里面写好的解密方法）
	user := models.FindUserByName(name)
	//若不存在用户
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "该用户不存在！",
			"data":    data,
		})
		return
	}
	//开始校验密码
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名或密码错误",
			"data":    data,
		})
		return
	}

	//注意password都是明文
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0, //0成功 -1失败
		"message": "登陆成功",
		"data":    data,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id formData string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	//从url的query里面获取id，由于是string，需要转换为uint
	id, _ := strconv.Atoi(c.Request.FormValue("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    -1,
		"message": "删除用户成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	//从url的postform里面获取id，由于是string，需要转换为uint
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)

	//获取新的名字和密码
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	//对需要校验的字段进行校验
	_, err := govalidator.ValidateStruct(user)
	//校验失败
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "修改参数不匹配",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "更新用户成功",
			"data":    user,
		})
	}
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	//将http更新为websocket
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//当函数调用完毕前一刻走defer，为了websocket能够关闭
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(ws)
	//调用下面获取订阅频道的数据函数
	MsgHandler(ws, c)
}

// 获取订阅频道的数据函数
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	//使用我们在systeminit里面封装好的订阅函数，获取订阅频道的数据
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// 查找朋友
func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(id))
	//c.JSON(200, gin.H{
	//	"code":    0,
	//	"message": "查询好友成功",
	//	"data":    users,
	//})
	utils.RespOKList(c.Writer, users, len(users))
}

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetId, _ := strconv.Atoi(c.Request.FormValue("targetName"))
	fmt.Println(userId, targetId)
	code := models.AddFriend(uint(userId), uint(targetId))
	if code == 0 {
		utils.RespOK(c.Writer, code, "添加好友成功")
	} else {
		utils.RespFail(c.Writer, "添加失败")
	}
}
