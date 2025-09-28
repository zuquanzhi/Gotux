package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gotux/database"
	"gotux/handlers"
)

type App struct {
	DBManager *database.DatabaseManager
}

func (app *App) initDB() {
	// 创建uploads目录
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}

	// 初始化数据库管理器
	app.DBManager = database.NewDatabaseManager()
	
	// 初始化主数据库
	err := app.DBManager.InitMainDB()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := &App{}
	app.initDB()

	// 创建处理器
	authHandler := &handlers.AuthHandler{DBManager: app.DBManager}
	imageHandler := &handlers.ImageHandler{DBManager: app.DBManager}

	// 上传文件服务 - 支持用户专属目录
	http.HandleFunc("/uploads/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "."+r.URL.Path)
	})

	// API路由处理
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// CORS处理
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch {
		case r.URL.Path == "/api/auth/register":
			authHandler.Register(w, r)
		case r.URL.Path == "/api/auth/login":
			authHandler.Login(w, r)
		case r.URL.Path == "/api/auth/change-password":
			authHandler.ChangePassword(w, r)
		case r.URL.Path == "/api/images/upload":
			imageHandler.UploadImage(w, r)
		case r.URL.Path == "/api/images/my":
			imageHandler.GetMyImages(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	// 静态文件服务
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 使用相对路径
		frontendDir := "../frontend"

		if r.URL.Path == "/" {
			indexPath := filepath.Join(frontendDir, "index.html")
			http.ServeFile(w, r, indexPath)
		} else if r.URL.Path == "/register" {
			registerPath := filepath.Join(frontendDir, "register.html")
			http.ServeFile(w, r, registerPath)
		} else if r.URL.Path == "/upload" {
			uploadPath := filepath.Join(frontendDir, "upload.html")
			http.ServeFile(w, r, uploadPath)
		} else if r.URL.Path == "/dashboard" {
			dashboardPath := filepath.Join(frontendDir, "dashboard.html")
			http.ServeFile(w, r, dashboardPath)
		} else {
			http.NotFound(w, r)
		}
	})

	log.Println("Server starting on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
