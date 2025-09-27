package routes

import (
    "log"
    "os"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"

    missionHandler "junctionx2025back/internal/api/handlers/mission"
    groundcontrol "junctionx2025back/internal/services/groundcontrol"
)

// SetupRoutes は渡された Gin エンジンにルートを登録します
func SetupRoutes(r *gin.Engine) *gin.Engine {
    // CORS設定
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // ヘルスチェック
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":    "healthy",
            "timestamp": time.Now(),
            "service":   "Junction X 2025 - Ground Control System",
        })
    })

    // GroundControlサービス初期化
    log.Printf("🔧 Initializing Ground Control Service...")
    gcService := groundcontrol.NewGroundControlService("JunctionX-2025")
    missionH := missionHandler.NewMissionHandler(gcService)
    log.Printf("✅ Ground Control Service initialized")

    // API v1 グループ
    v1 := r.Group("/api/v1")
    {
        missionsGroup := v1.Group("/missions")
        {
            missionsGroup.GET("", missionH.GetMissions)
            missionsGroup.POST("", missionH.CreateMission)
            missionsGroup.GET("/:id", missionH.GetMission)
            missionsGroup.POST("/:id/message", missionH.SendMessage)
        }

        v1.GET("/ws/missions/:id", missionH.WebSocketChat)
    }

    // 管理者向けルート
    adminGroup := r.Group("/admin")
    {
        adminGroup.GET("/status", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "server_status":   "operational",
                "ground_control":  "online",
                "gemini_api":      "connected",
                "active_missions": 0,
                "uptime":          time.Since(time.Now()),
            })
        })

        adminGroup.GET("/logs", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "logs": []string{
                    "Ground Control System initialized",
                    "Gemini API connected successfully",
                    "Mission handler ready for deployment",
                    "All systems operational",
                },
            })
        })

        adminGroup.GET("/config", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "environment":   os.Getenv("GIN_MODE"),
                "port":          ":8080",
                "cors":          "enabled",
                "websocket":     "enabled",
                "ground_control": "active",
            })
        })
    }

    // デバッグ用ルート
    r.GET("/debug/routes", func(c *gin.Context) {
        routes := []string{}
        for _, route := range r.Routes() {
            routes = append(routes, route.Method+" "+route.Path)
        }
        c.JSON(200, gin.H{
            "registered_routes": routes,
            "total_routes":      len(routes),
        })
    })

    log.Printf("🛰️ Ground Control System routes setup complete")
    return r
}
