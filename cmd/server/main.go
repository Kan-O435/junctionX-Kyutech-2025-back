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
    
    // Ginモード設定
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    // ルート設定（エンジン生成をルータに委譲）
    r := routes.SetupRoutes(cfg)
    
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