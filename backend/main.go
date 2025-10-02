package main

import (
	"gotux/config"
	"gotux/models"
	"gotux/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库
	models.InitDB()

		// 创建默认管理员账户
		models.CreateDefaultAdmin()

		// 自动修正所有管理员配额为无限制（0）
		if err := models.FixAdminQuota(); err != nil {
			log.Println("自动修正管理员配额失败:", err)
		}

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.Default()

	// 配置 CORS - 开发模式允许所有来源
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Image-UUID", "X-Image-ID"},
		AllowCredentials: false, // AllowAllOrigins 时必须设为 false
	}))

	// 静态文件服务 - 用于访问上传的图片
	r.Static("/uploads", "./uploads")

	// 注册路由
	routes.SetupRoutes(r)

	// 启动服务器
	port := config.AppConfig.Server.Port
	log.Printf("Server is running on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
