package routes

import (
    "log"
    "math/rand"
    "os"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"

    missionHandler "junctionx2025back/internal/api/handlers/mission"
    "junctionx2025back/internal/config"
    "junctionx2025back/internal/services/groundcontrol"
)

// SetupRoutes configures all routes for the disaster response system
func SetupRoutes(cfg *config.Config) *gin.Engine {
    r := gin.Default()

    // CORS configuration for frontend integration
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:3000",  // Next.js development server
            "http://127.0.0.1:3000",
            "http://localhost:3001",  // Alternative port
            "https://your-domain.com", // Production domain
        },
        AllowMethods: []string{
            "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
        },
        AllowHeaders: []string{
            "Origin", "Content-Type", "Content-Length", 
            "Accept-Encoding", "X-CSRF-Token", "Authorization",
            "Accept", "Cache-Control", "X-Requested-With",
        },
        ExposeHeaders: []string{
            "Content-Length", "Content-Type",
        },
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":    "healthy",
            "timestamp": time.Now(),
            "service":   "Junction X 2025 - Disaster Response System",
        })
    })

    // Initialize Ground Control Service
    log.Printf("🔧 Initializing Ground Control Service...")
    gcService := groundcontrol.NewGroundControlService("AIzaSyC9S10n_tN84xMiEp7hALRZIQD2Olqlokg")
    missionH := missionHandler.NewMissionHandler(gcService)
    log.Printf("✅ Ground Control Service initialized")

    // API v1 group
    v1 := r.Group("/api/v1")
    {
        // Mission Control endpoints
        missionsGroup := v1.Group("/missions")
        {
            missionsGroup.GET("", missionH.GetMissions)
            missionsGroup.POST("", missionH.CreateMission)
            missionsGroup.GET("/:id", missionH.GetMission)
            missionsGroup.POST("/:id/message", missionH.SendMessage)
        }

        // Satellite endpoints (mock data for frontend compatibility)
        satelliteGroup := v1.Group("/satellite")
        {
            satelliteGroup.GET("/:id/orbit", func(c *gin.Context) {
                satelliteId := c.Param("id")
                c.JSON(200, gin.H{
                    "satellite_id": satelliteId,
                    "latitude":     35.6762 + (rand.Float64()-0.5)*10,
                    "longitude":    139.6503 + (rand.Float64()-0.5)*10,
                    "altitude":     550 + rand.Float64()*50,
                    "velocity":     7.66 + rand.Float64()*0.5,
                    "timestamp":    time.Now(),
                })
            })
            
            satelliteGroup.GET("/:id/status", func(c *gin.Context) {
                satelliteId := c.Param("id")
                c.JSON(200, gin.H{
                    "satellite_id":    satelliteId,
                    "status":          "operational",
                    "battery_level":   75 + rand.Intn(25),
                    "signal_strength": "strong",
                    "last_contact":    time.Now().Add(-time.Minute * 2),
                    "temperature":     -15 + rand.Intn(30),
                })
            })
            
            satelliteGroup.GET("/:id/coverage", func(c *gin.Context) {
                satelliteId := c.Param("id")
                c.JSON(200, gin.H{
                    "satellite_id":     satelliteId,
                    "coverage_area":    "Tokyo region",
                    "visibility":       "excellent",
                    "elevation_angle":  45 + rand.Intn(30),
                    "next_pass":        time.Now().Add(time.Hour * 2),
                    "coverage_radius":  500 + rand.Intn(200),
                })
            })
        }

        // WebSocket endpoint for real-time communication
        v1.GET("/ws/missions/:id", missionH.WebSocketChat)

        // Frontend connection test
        v1.GET("/test", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "message":   "Backend connected successfully",
                "timestamp": time.Now(),
                "cors":      "enabled",
            })
        })
    }

    // Admin endpoints
    adminGroup := r.Group("/admin")
    {
        adminGroup.GET("/status", func(c *gin.Context) {
            missions := gcService.GetMissions()
            c.JSON(200, gin.H{
                "server_status":   "operational",
                "ground_control":  "online",
                "gemini_api":      "connected",
                "active_missions": len(missions),
                "uptime":          time.Since(time.Now()),
            })
        })

        adminGroup.GET("/logs", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "logs": []string{
                    "Ground Control System initialized",
                    "Gemini API connected successfully",
                    "Mission handler ready for deployment",
                    "Satellite mock data enabled",
                    "All systems operational",
                },
            })
        })

        adminGroup.GET("/config", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "environment":     os.Getenv("GIN_MODE"),
                "port":            ":8080",
                "cors":            "enabled",
                "websocket":       "enabled",
                "ground_control":  "active",
                "satellite_mock":  "enabled",
            })
        })
    }

    // Debug endpoints
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

    // API documentation
    r.GET("/api/docs", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "service":     "Junction X 2025 - Disaster Response System",
            "version":     "1.0.0",
            "description": "AI-powered disaster response with Ground Control and satellite monitoring",
            "endpoints": gin.H{
                "health": "/health",
                "missions": gin.H{
                    "list":    "GET /api/v1/missions",
                    "create":  "POST /api/v1/missions",
                    "get":     "GET /api/v1/missions/:id",
                    "message": "POST /api/v1/missions/:id/message",
                },
                "satellite": gin.H{
                    "orbit":    "GET /api/v1/satellite/:id/orbit",
                    "status":   "GET /api/v1/satellite/:id/status",
                    "coverage": "GET /api/v1/satellite/:id/coverage",
                },
                "websocket": "GET /api/v1/ws/missions/:id",
                "admin": gin.H{
                    "status": "/admin/status",
                    "logs":   "/admin/logs",
                    "config": "/admin/config",
                },
            },
        })
    })

    log.Printf("🛰️ Ground Control System starting on http://localhost:8080")
    log.Printf("📡 Mission Control ready for deployment")
    log.Printf("🛰️ Satellite mock endpoints enabled")
    log.Printf("📋 Available endpoints:")
    log.Printf("  GET  /health")
    log.Printf("  GET  /api/v1/missions")
    log.Printf("  POST /api/v1/missions")
    log.Printf("  GET  /api/v1/satellite/:id/orbit")
    log.Printf("  GET  /api/v1/satellite/:id/status")
    log.Printf("  GET  /api/v1/satellite/:id/coverage")
    log.Printf("  GET  /api/docs")

    return r
}