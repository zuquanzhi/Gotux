// 数据库迁移工具
// 使用方法: go run cmd/migrate/main.go
// 这将自动添加新的用户设置字段到数据库

package main

import (
	"fmt"
	"log"
	"os"

	"gotux/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库连接
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./gotux.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	models.DB = db

	// 自动迁移
	fmt.Println("🚀 开始迁移数据库...")
	if err := db.AutoMigrate(&models.User{}, &models.Image{}, &models.ImageStats{}); err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err)
	}

	fmt.Println("✅ 数据库迁移完成！")
	fmt.Println("\n📋 已添加/更新以下字段到 User 表:")
	fmt.Println("   - custom_domain (自定义域名)")
	fmt.Println("   - default_link_format (默认链接格式)")
	fmt.Println("   - enable_watermark (启用水印)")
	fmt.Println("   - watermark_text (水印文字)")
	fmt.Println("   - watermark_position (水印位置)")
	fmt.Println("   - compress_image (压缩图片)")
	fmt.Println("   - compress_quality (压缩质量)")
	fmt.Println("   - max_image_size (最大图片大小)")
	fmt.Println("   - allowed_image_types (允许的图片类型)")
	fmt.Println("   - enable_image_review (图片审核)")
	fmt.Println("   - storage_quota (存储配额)")
	fmt.Println("   - used_storage (已使用存储)")
	fmt.Println("\n💡 提示: 重启后端服务以使更改生效")
}
