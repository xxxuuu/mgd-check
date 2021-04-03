package web

import (
	"github.com/gin-gonic/gin"
	"mgd-check/core"
	"net/http"
)

func InitWeb() {
	// 初始化gin & 注册路由
	r := gin.Default()

	r.LoadHTMLGlob("web/template/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/register", func(c *gin.Context) {
		info := core.CheckInfo{
			c.PostForm("phone"),
			c.PostForm("password"),
			c.PostForm("description"),
			c.PostForm("country"),
			c.PostForm("province"),
			c.PostForm("city"),
			c.PostForm("address"),
			c.PostForm("latitude"),
			c.PostForm("longitude"),
			c.PostForm("email"),
		}
		core.GetDb().Register(info)
		c.HTML(http.StatusOK, "register-result.html", gin.H{
			"result": "录入成功",
		})
	})

	// default port 8080
	_ = r.Run()
}

