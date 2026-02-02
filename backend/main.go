package main

import (
	"log"
	"pr-business/config"
	"pr-business/middlewares"
	"pr-business/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 初始化Redis
	redisClient, err := config.InitRedis(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 中间件
	r.Use(middlewares.CORS())
	r.Use(middlewares.Logger())
	r.Use(middlewares.Recovery())

	// 注册路由
	routes.SetupRoutes(r, db, redisClient, cfg)

	// 启动服务器
	log.Printf("Server starting on port %s...", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
