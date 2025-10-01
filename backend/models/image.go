package models

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	FileName     string         `gorm:"not null" json:"file_name"`
	OriginalName string         `gorm:"not null" json:"original_name"`
	FilePath     string         `gorm:"not null" json:"file_path"`
	FileSize     int64          `json:"file_size"`
	MimeType     string         `json:"mime_type"`
	Width        int            `json:"width"`
	Height       int            `json:"height"`
	Hash         string         `gorm:"index" json:"hash"`
	Description  string         `json:"description"`
	Tags         string         `json:"tags"` // 逗号分隔的标签
	IsPublic     bool           `gorm:"default:true" json:"is_public"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Stats        *ImageStats    `gorm:"foreignKey:ImageID" json:"stats,omitempty"`
}

type ImageStats struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ImageID   uint      `gorm:"uniqueIndex;not null" json:"image_id"`
	ViewCount int64     `gorm:"default:0" json:"view_count"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateImage 创建图片记录
func CreateImage(image *Image) error {
	return DB.Create(image).Error
}

// GetImageByID 根据ID获取图片
func GetImageByID(id uint) (*Image, error) {
	var image Image
	if err := DB.Preload("User").Preload("Stats").First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImagesByUserID 获取用户的所有图片
func GetImagesByUserID(userID uint, page, pageSize int) ([]Image, int64, error) {
	var images []Image
	var total int64

	query := DB.Model(&Image{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Stats").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&images).Error; err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

// GetAllImages 获取所有图片（管理员用）
func GetAllImages(page, pageSize int) ([]Image, int64, error) {
	var images []Image
	var total int64

	if err := DB.Model(&Image{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := DB.Preload("User").Preload("Stats").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&images).Error; err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

// UpdateImage 更新图片信息
func (i *Image) Update() error {
	return DB.Save(i).Error
}

// DeleteImage 删除图片
func (i *Image) Delete() error {
	return DB.Delete(i).Error
}

// SearchImages 搜索图片
func SearchImages(userID uint, keyword string, page, pageSize int) ([]Image, int64, error) {
	var images []Image
	var total int64

	query := DB.Model(&Image{}).Where("user_id = ?", userID)

	if keyword != "" {
		query = query.Where("original_name LIKE ? OR description LIKE ? OR tags LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Stats").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&images).Error; err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

// IncrementViewCount 增加访问次数
func IncrementViewCount(imageID uint) error {
	var stats ImageStats

	// 尝试查找现有统计记录
	err := DB.Where("image_id = ?", imageID).First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新记录
			stats = ImageStats{
				ImageID:   imageID,
				ViewCount: 1,
			}
			return DB.Create(&stats).Error
		}
		return err
	}

	// 更新统计
	return DB.Model(&stats).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetUserStorageUsed 获取用户已使用的存储空间
func GetUserStorageUsed(userID uint) (int64, error) {
	var total int64
	err := DB.Model(&Image{}).Where("user_id = ?", userID).Select("COALESCE(SUM(file_size), 0)").Scan(&total).Error
	return total, err
}

// GetImageByHash 根据哈希值查找图片（用于去重）
func GetImageByHash(hash string, userID uint) (*Image, error) {
	var image Image
	err := DB.Where("hash = ? AND user_id = ?", hash, userID).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}
