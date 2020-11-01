package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)
var (
	bindAddress = ":9999"
)

func main()  {
	r := gin.Default()
	// 创建基于cookie的存储引擎，secret11111 参数是用于加密的密钥
	//store := cookie.NewStore([]byte("secret11111"))


	// 初始化基于redis的存储引擎
	// 参数说明：
	//    第1个参数 - redis最大的空闲连接数
	//    第2个参数 - 数通信协议tcp或者udp
	//    第3个参数 - redis地址, 格式，host:port
	//    第4个参数 - redis密码
	//    第5个参数 - session加密密钥
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))


	r.Use(sessions.Sessions("mysession",store))

	r.Static("/static","static")
	r.LoadHTMLGlob("views/**/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200,"home/index",gin.H{})
	})

	r.GET("/user", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil{
			count = 0
		}else{
			count = v.(int)
			count ++
		}
		session.Set("count",count)
		session.Save()
		c.JSON(200,gin.H{
			"user":h.GetUsers(),
			"count":count,
		})
	})

	r.GET("/wsPage",wsPage)

	//开启管理
	go h.run()
	r.Run(bindAddress)
}
