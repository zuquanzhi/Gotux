package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gotux/config"
	"gotux/middleware"
	"gotux/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析表单
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "表单解析失败"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有上传文件"})
		return
	}

	var uploadedImages []models.Image
	var errors []string

	for _, file := range files {
		// 检查文件大小
		if file.Size > config.AppConfig.Upload.MaxSize {
			errors = append(errors, fmt.Sprintf("%s: 文件大小超过限制", file.Filename))
			continue
		}

		// 检查文件类型
		if !isAllowedFileType(file.Header.Get("Content-Type")) {
			errors = append(errors, fmt.Sprintf("%s: 不支持的文件类型", file.Filename))
			continue
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: 文件打开失败", file.Filename))
			continue
		}

		// 计算文件哈希
		hash := md5.New()
		if _, err := io.Copy(hash, src); err != nil {
			src.Close()
			errors = append(errors, fmt.Sprintf("%s: 哈希计算失败", file.Filename))
			continue
		}
		hashStr := hex.EncodeToString(hash.Sum(nil))
		src.Close()

		// 检查是否已存在相同文件
		existingImage, err := models.GetImageByHash(hashStr, userID)
		if err == nil {
			uploadedImages = append(uploadedImages, *existingImage)
			continue
		}

		// 生成唯一文件名
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		// 按日期组织文件夹
		dateFolder := time.Now().Format("2006/01/02")
		fullPath := filepath.Join(config.AppConfig.Upload.StoragePath, dateFolder)

		// 创建目录
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			errors = append(errors, fmt.Sprintf("%s: 创建目录失败", file.Filename))
			continue
		}

		// 保存文件
		filePath := filepath.Join(fullPath, newFileName)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			errors = append(errors, fmt.Sprintf("%s: 文件保存失败", file.Filename))
			continue
		}

		// 获取图片尺寸
		width, height := getImageDimensions(filePath)

		// 创建数据库记录
		image := models.Image{
			UserID:       userID,
			FileName:     newFileName,
			OriginalName: file.Filename,
			FilePath:     filepath.Join(dateFolder, newFileName),
			FileSize:     file.Size,
			MimeType:     file.Header.Get("Content-Type"),
			Width:        width,
			Height:       height,
			Hash:         hashStr,
			IsPublic:     true,
		}

		if err := models.CreateImage(&image); err != nil {
			os.Remove(filePath) // 删除已保存的文件
			errors = append(errors, fmt.Sprintf("%s: 数据库保存失败", file.Filename))
			continue
		}

		uploadedImages = append(uploadedImages, image)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("成功上传 %d 个文件", len(uploadedImages)),
		"images":  uploadedImages,
		"errors":  errors,
	})
}

// GetImages 获取图片列表
func GetImages(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var images []models.Image
	var total int64
	var err error

	if keyword != "" {
		images, total, err = models.SearchImages(userID, keyword, page, pageSize)
	} else {
		images, total, err = models.GetImagesByUserID(userID, page, pageSize)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取图片列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images":    images,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetImageDetail 获取图片详情
func GetImageDetail(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	image, err := models.GetImageByID(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 检查权限
	if image.UserID != userID {
		user, _ := middleware.GetUser(c)
		if !user.IsAdmin() {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问该图片"})
			return
		}
	}

	// 增加访问次数
	models.IncrementViewCount(uint(imageID))

	c.JSON(http.StatusOK, image)
}

// UpdateImage 更新图片信息
func UpdateImage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	image, err := models.GetImageByID(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 检查权限
	if image.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限编辑该图片"})
		return
	}

	var req struct {
		Description string `json:"description"`
		Tags        string `json:"tags"`
		IsPublic    *bool  `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	image.Description = req.Description
	image.Tags = req.Tags
	if req.IsPublic != nil {
		image.IsPublic = *req.IsPublic
	}

	if err := image.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"image":   image,
	})
}

// DeleteImage 删除图片
func DeleteImage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	image, err := models.GetImageByID(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 检查权限
	if image.UserID != userID {
		user, _ := middleware.GetUser(c)
		if !user.IsAdmin() {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除该图片"})
			return
		}
	}

	// 删除文件
	fullPath := filepath.Join(config.AppConfig.Upload.StoragePath, image.FilePath)
	os.Remove(fullPath)

	// 删除数据库记录
	if err := image.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// BatchDeleteImages 批量删除图片
func BatchDeleteImages(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req struct {
		ImageIDs []uint `json:"image_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	deletedCount := 0
	for _, imageID := range req.ImageIDs {
		image, err := models.GetImageByID(imageID)
		if err != nil {
			continue
		}

		// 检查权限
		if image.UserID != userID {
			user, _ := middleware.GetUser(c)
			if !user.IsAdmin() {
				continue
			}
		}

		// 删除文件
		fullPath := filepath.Join(config.AppConfig.Upload.StoragePath, image.FilePath)
		os.Remove(fullPath)

		// 删除数据库记录
		if err := image.Delete(); err == nil {
			deletedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("成功删除 %d 张图片", deletedCount),
		"deleted_count": deletedCount,
	})
}

// GetImageLinks 获取图片链接（多种格式）
func GetImageLinks(c *gin.Context) {
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	image, err := models.GetImageByID(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 构建完整URL
	baseURL := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	imageURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, baseURL, image.FilePath)

	// 生成各种格式的链接
	links := map[string]string{
		"url":                imageURL,
		"html":               fmt.Sprintf(`<img src="%s" alt="%s" />`, imageURL, image.OriginalName),
		"markdown":           fmt.Sprintf(`![%s](%s)`, image.OriginalName, imageURL),
		"bbcode":             fmt.Sprintf(`[img]%s[/img]`, imageURL),
		"markdown_with_link": fmt.Sprintf(`[![%s](%s)](%s)`, image.OriginalName, imageURL, imageURL),
	}

	c.JSON(http.StatusOK, gin.H{
		"image": image,
		"links": links,
	})
}

// 辅助函数

func isAllowedFileType(mimeType string) bool {
	for _, allowed := range config.AppConfig.Upload.AllowedTypes {
		if allowed == mimeType {
			return true
		}
	}
	return false
}

func getImageDimensions(filePath string) (int, int) {
	img, err := imaging.Open(filePath)
	if err != nil {
		return 0, 0
	}
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

// GetStats 获取统计信息
func GetStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取图片总数
	var imageCount int64
	models.DB.Model(&models.Image{}).Where("user_id = ?", userID).Count(&imageCount)

	// 获取存储使用量
	storageUsed, _ := models.GetUserStorageUsed(userID)

	// 获取总访问量
	var totalViews int64
	models.DB.Table("image_stats").
		Joins("JOIN images ON images.id = image_stats.image_id").
		Where("images.user_id = ?", userID).
		Select("COALESCE(SUM(image_stats.view_count), 0)").
		Scan(&totalViews)

	c.JSON(http.StatusOK, gin.H{
		"image_count":  imageCount,
		"storage_used": storageUsed,
		"total_views":  totalViews,
	})
}
