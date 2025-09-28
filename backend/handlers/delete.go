package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gotux/utils"
)

func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DeleteImage called: Method=%s, URL=%s\n", r.Method, r.URL.Path)
	
	if r.Method != "DELETE" {
		fmt.Printf("Invalid method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从JWT获取用户ID
	fmt.Printf("Getting authorization header...\n")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		fmt.Printf("No authorization header or invalid format\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "需要登录才能删除图片",
		})
		return
	}

	tokenString := authHeader[7:]
	fmt.Printf("Parsing JWT token...\n")
	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		fmt.Printf("JWT parsing failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "无效的登录令牌",
		})
		return
	}

	userID := claims.UserID
	fmt.Printf("User ID: %d\n", userID)

	// 从URL路径中提取图片ID
	fmt.Printf("Extracting image ID from URL path: %s\n", r.URL.Path)
	pathParts := strings.Split(r.URL.Path, "/")
	fmt.Printf("Path parts: %v\n", pathParts)
	if len(pathParts) < 4 {
		fmt.Printf("Invalid path parts length: %d\n", len(pathParts))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "无效的图片ID",
		})
		return
	}

	imageIDStr := pathParts[3] // /api/images/{id}
	fmt.Printf("Image ID string: %s\n", imageIDStr)
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Failed to parse image ID: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "无效的图片ID格式",
		})
		return
	}
	fmt.Printf("Parsed image ID: %d\n", imageID)

	// 获取用户专属数据库
	userDB, err := h.DBManager.GetUserDB(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "无法访问用户数据库",
		})
		return
	}
	defer userDB.Close()

	// 查询图片信息
	var filename, originalName string
	err = userDB.QueryRow("SELECT filename, original_name FROM images WHERE id = ?", imageID).Scan(&filename, &originalName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "图片不存在",
		})
		return
	}

	// 删除数据库记录
	_, err = userDB.Exec("DELETE FROM images WHERE id = ?", imageID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "删除数据库记录失败",
		})
		return
	}

	// 删除物理文件
	filePath := filepath.Join("uploads", "user_"+strconv.FormatInt(userID, 10), filename)
	if err := os.Remove(filePath); err != nil {
		// 即使文件删除失败，也不返回错误，因为数据库记录已经删除
		// 这可能是因为文件已经被手动删除或不存在
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("Image deletion successful for ID: %d\n", imageID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "图片删除成功",
	})
}