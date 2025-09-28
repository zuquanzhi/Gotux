package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gotux/handlers"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB *sql.DB
}

func (app *App) initDB() {
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}

	db, err := sql.Open("sqlite3", "./gotux.db")
	if err != nil {
		log.Fatal(err)
	}

	// 创建用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TEXT DEFAULT (datetime('now'))
	);`

	// 创建图片表
	createImagesTable := `
	CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER DEFAULT 0,
		filename TEXT NOT NULL,
		original_name TEXT NOT NULL,
		size INTEGER NOT NULL,
		mime_type TEXT NOT NULL,
		url TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createImagesTable)
	if err != nil {
		log.Fatal(err)
	}

	app.DB = db
}

func main() {
	app := &App{}
	app.initDB()

	// 初始化处理器
	authHandler := &handlers.AuthHandler{DB: app.DB}
	imageHandler := &handlers.ImageHandler{DB: app.DB}

	// 路由
	// http.HandleFunc("/api/auth/register", authHandler.Register)
	// http.HandleFunc("/api/auth/login", authHandler.Login)
	// http.HandleFunc("/api/images/upload", imageHandler.UploadImage)
	// http.HandleFunc("/api/images/my", imageHandler.GetMyImages)
	http.HandleFunc("/uploads/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "."+r.URL.Path)
	})

	// CORS 中间件
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 根据路径分发请求
		switch {
		case r.URL.Path == "/api/auth/register":
			authHandler.Register(w, r)
		case r.URL.Path == "/api/auth/login":
			authHandler.Login(w, r)
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
		} else if r.URL.Path == "/dashboard" {
			dashboardPath := filepath.Join(frontendDir, "dashboard.html")
			http.ServeFile(w, r, dashboardPath)
		} else {
			http.NotFound(w, r)
		}
	})

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
