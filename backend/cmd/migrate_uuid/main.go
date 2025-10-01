package main

import (
	"fmt"
	"gotux/config"
	"gotux/models"
	"log"

	"github.com/google/uuid"
)

func main() {
	// 加载配置
	config.InitConfig()

	// 连接数据库
	models.InitDB()

	log.Println("开始迁移数据库...")

	// 为现有的图片生成UUID
	var images []models.Image
	if err := models.DB.Find(&images).Error; err != nil {
		log.Fatal("查询图片失败:", err)
	}

	log.Printf("找到 %d 张图片需要生成UUID\n", len(images))

	for i := range images {
		if images[i].UUID == "" {
			images[i].UUID = uuid.New().String()
			if err := models.DB.Save(&images[i]).Error; err != nil {
				log.Printf("更新图片 %d UUID失败: %v\n", images[i].ID, err)
			} else {
				log.Printf("图片 %d: 生成UUID %s\n", images[i].ID, images[i].UUID)
			}
		}
	}

	log.Println("数据库迁移完成!")
	fmt.Println("\n重要提示:")
	fmt.Println("1. UUID字段已添加到所有图片")
	fmt.Println("2. 现在可以使用 /i/:uuid 访问图片")
	fmt.Println("3. 旧的 /uploads/:path 路径仍然可用")
}
