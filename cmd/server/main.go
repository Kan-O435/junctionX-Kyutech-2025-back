package main

import (
    "log"
    "junctionx2025back/internal/api/routes"
    "junctionx2025back/internal/config"
    
    "github.com/gin-gonic/gin"
)

func main() {
    log.Println("🛰️ Starting Disaster Response Ground Control System...")
    
    // 設定読み込み
    cfg := config.Load()
    log.Printf("📡 Environment: %s", cfg.Environment)
    
    // Ginモード設定
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // ルート設定（CORSも含めて全てroutes.SetupRoutesで処理）
    r := routes.SetupRoutes(cfg)
    
    // サーバー起動
    log.Printf("🌍 Server starting on http://localhost:%s", cfg.Port)
    log.Println("📋 Available endpoints:")
    log.Println("  GET  /health")
    log.Println("  GET  /api/v1/missions")
    log.Println("  POST /api/v1/missions")
    log.Println("  GET  /api/v1/missions/:id")
    log.Println("  POST /api/v1/missions/:id/message")
    log.Println("  GET  /api/v1/satellite/:id/orbit")
    log.Println("  GET  /api/v1/satellite/:id/status")
    log.Println("  GET  /api/v1/satellite/:id/coverage")
    log.Println("  GET  /api/v1/ws/missions/:id")
    log.Println("  GET  /api/docs")
    log.Println("  GET  /debug/routes")
    
    if err := r.Run(":" + cfg.Port); err != nil {
        log.Fatal("❌ Failed to start server:", err)
    }
}