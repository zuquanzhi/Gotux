package routes

import (
	"gotux/controllers"
	"gotux/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// API 路由组
	api := r.Group("/api")
	{
		// 公开路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// 公开访问图片信息(通过UUID)
		api.GET("/i/:uuid", controllers.GetImageByUUID)

		// 随机图片 API
		api.GET("/random", controllers.GetRandomImage)           // 返回JSON
		api.GET("/random/image", controllers.ServeRandomImage)   // 直接返回图片
		api.GET("/random/redirect", controllers.RedirectRandomImage) // 重定向到图片

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户相关
			user := authorized.Group("/user")
			{
				user.GET("/profile", controllers.GetProfile)
				user.PUT("/profile", controllers.UpdateProfile)
				user.POST("/change-password", controllers.ChangePassword)
				user.GET("/stats", controllers.GetStats)
				user.GET("/settings", controllers.GetSettings)
				user.PUT("/settings", controllers.UpdateSettings)
			}

			// 图片相关
			image := authorized.Group("/images")
			{
				image.POST("/upload", controllers.UploadImage)
				image.GET("", controllers.GetImages)
				image.GET("/:id", controllers.GetImageDetail)
				image.PUT("/:id", controllers.UpdateImage)
				image.DELETE("/:id", controllers.DeleteImage)
				image.POST("/batch-delete", controllers.BatchDeleteImages)
				image.GET("/:id/links", controllers.GetImageLinks)
			}

			// 管理员路由
			admin := authorized.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", controllers.GetAllUsers)
				admin.PUT("/users/:id/status", controllers.UpdateUserStatus)
				admin.GET("/images", controllers.GetAllImagesAdmin)
				admin.GET("/stats", controllers.GetSystemStats)
			}
		}
	}

	// 直接提供图片文件(通过UUID)
	r.GET("/i/:uuid", controllers.ServeImageByUUID)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
