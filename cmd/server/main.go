package main

import (
	"junctionx2025back/internal/api/routes"
	"junctionx2025back/internal/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
    log.Println("🚀 Starting Satellite Game Backend...")
    
    // 設定読み込み
    cfg := config.Load()
    log.Printf("📡 Environment: %s", cfg.Environment)
    
    // Ginエンジン初期化
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    r := gin.Default()
    
    // CORS設定（開発用）
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })
    
    // ルート設定
    routes.SetupRoutes(r)
    
    // サーバー起動
    log.Printf("🌍 Server starting on http://0.0.0.0:%s", cfg.Port)
    log.Println("📋 Available endpoints:")
    log.Println("  GET  /health")
    log.Println("  GET  /api/v1/satellite/{id}/orbit")
    log.Println("  POST /api/v1/satellite/{id}/maneuver")
    log.Println("  GET  /api/v1/mission/debris/{id}/threats")
    log.Println("  GET  /api/v1/mission/debris/list")
    log.Println("  GET  /api/v1/mission/debris/stats")
    
    // Cloud Run では 0.0.0.0:$PORT でのリッスンが必須
    if err := r.Run("0.0.0.0:" + cfg.Port); err != nil {
        log.Fatal("❌ Failed to start server:", err)
    }
}