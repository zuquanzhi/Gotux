package controllers

import (
	"gotux/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllUsers 获取所有用户（管理员）
func GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := models.GetAllUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetAllImagesAdmin 获取所有图片（管理员）
func GetAllImagesAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	images, total, err := models.GetAllImages(page, pageSize)
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

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active disabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	user.Status = req.Status
	if err := user.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"user":    user,
	})
}

// UpdateUserQuota 更新用户存储配额(管理员)
func UpdateUserQuota(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req struct {
		StorageQuota int64 `json:"storage_quota" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 检查新配额是否小于当前使用量
	storageUsed, _ := models.GetUserStorageUsed(uint(userID))
	if req.StorageQuota > 0 && req.StorageQuota < storageUsed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "新配额不能小于当前使用量",
			"storage_used":  storageUsed,
			"storage_quota": req.StorageQuota,
		})
		return
	}

	user.StorageQuota = req.StorageQuota
	if err := user.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "配额更新成功",
		"user_id":       user.ID,
		"storage_quota": user.StorageQuota,
		"storage_used":  storageUsed,
	})
}

// GetSystemStats 获取系统统计信息
func GetSystemStats(c *gin.Context) {
	var userCount int64
	var imageCount int64
	var totalStorage int64
	var totalViews int64

	models.DB.Model(&models.User{}).Count(&userCount)
	models.DB.Model(&models.Image{}).Count(&imageCount)
	models.DB.Model(&models.Image{}).Select("COALESCE(SUM(file_size), 0)").Scan(&totalStorage)
	models.DB.Model(&models.ImageStats{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)

	c.JSON(http.StatusOK, gin.H{
		"user_count":    userCount,
		"image_count":   imageCount,
		"total_storage": totalStorage,
		"total_views":   totalViews,
	})
}
