package models

import (
	"gotux/config"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.Database.Path), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移 User 和 ImageStats 表
	err = DB.AutoMigrate(&User{}, &ImageStats{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 手动处理 Image 表的迁移
	migrateImageTable()

	log.Println("Database initialized successfully")
}

// migrateImageTable 手动迁移 Image 表
func migrateImageTable() {
	// 检查表是否存在
	if !DB.Migrator().HasTable(&Image{}) {
		// 表不存在,直接创建
		if err := DB.AutoMigrate(&Image{}); err != nil {
			log.Fatal("Failed to create images table:", err)
		}
		log.Println("Created images table")
		return
	}

	// 表已存在,检查是否有 uuid 列
	if !DB.Migrator().HasColumn(&Image{}, "uuid") {
		// 添加 uuid 列(不带 UNIQUE 约束)
		if err := DB.Exec("ALTER TABLE images ADD COLUMN uuid TEXT").Error; err != nil {
			log.Fatal("Failed to add uuid column:", err)
		}
		log.Println("Added uuid column to images table")
	}

	// 为没有 UUID 的图片生成 UUID
	migrateExistingImages()

	// 检查是否有唯一索引
	if !DB.Migrator().HasIndex(&Image{}, "idx_images_uuid") {
		// 创建唯一索引
		if err := DB.Exec("CREATE UNIQUE INDEX idx_images_uuid ON images(uuid)").Error; err != nil {
			log.Fatal("Failed to create unique index on uuid:", err)
		}
		log.Println("Created unique index on uuid column")
	}
}

// migrateExistingImages 为现有图片生成 UUID
func migrateExistingImages() {
	var images []Image
	if err := DB.Where("uuid IS NULL OR uuid = ''").Find(&images).Error; err != nil {
		log.Println("Warning: Failed to query images for UUID migration:", err)
		return
	}

	if len(images) == 0 {
		return
	}

	log.Printf("Migrating %d images to add UUIDs...\n", len(images))

	for i := range images {
		if images[i].UUID == "" {
			// 触发 BeforeCreate hook 来生成 UUID
			images[i].BeforeCreate(DB)
			if err := DB.Model(&images[i]).Update("uuid", images[i].UUID).Error; err != nil {
				log.Printf("Warning: Failed to update UUID for image %d: %v\n", images[i].ID, err)
			}
		}
	}

	log.Println("UUID migration completed")
}
