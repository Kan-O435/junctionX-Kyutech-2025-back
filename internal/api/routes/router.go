package routes

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	disasterHandlers "junctionx2025back/internal/api/handlers/disaster"
	missionHandler "junctionx2025back/internal/api/handlers/mission"
	videoHandlers "junctionx2025back/internal/api/handlers/satellite/video"
	"junctionx2025back/internal/config"
	"junctionx2025back/internal/services/groundcontrol"
)

// helper parsing with fallback
func parseFloat(s string, fallback float64) float64 {
    if v, err := strconv.ParseFloat(s, 64); err == nil {
        return v
    }
    return fallback
}

func parseInt(s string, fallback int) int {
    if v, err := strconv.Atoi(s); err == nil {
        return v
    }
    return fallback
}

// SetupRoutes configures all routes for the disaster response system
func SetupRoutes(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS configuration for frontend integration (no external dependency)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

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
				// Align with FE OrbitResponse shape
				c.JSON(200, gin.H{
					"satellite_id":  satelliteId,
					"timestamp":     time.Now(),
					"altitude":      550 + rand.Float64()*50,      // km
					"orbital_speed": 7.66 + rand.Float64()*0.5,   // km/s
					"position": gin.H{"x": rand.Float64()*1000 - 500, "y": rand.Float64()*1000 - 500, "z": rand.Float64()*1000 - 500},
					"velocity": gin.H{"x": rand.Float64()*2 - 1, "y": rand.Float64()*2 - 1, "z": rand.Float64()*2 - 1},
				})
			})

			satelliteGroup.GET("/:id/status", func(c *gin.Context) {
				satelliteId := c.Param("id")
				// Align with FE StatusResponse shape
				attitude := gin.H{"roll": rand.Float64()*10 - 5, "pitch": rand.Float64()*10 - 5, "yaw": rand.Float64()*10 - 5}
				c.JSON(200, gin.H{
					"satellite_id": satelliteId,
					"status": gin.H{
						"position":    gin.H{"x": rand.Float64()*1000 - 500, "y": rand.Float64()*1000 - 500, "z": rand.Float64()*1000 - 500},
						"velocity":    gin.H{"x": rand.Float64()*2 - 1, "y": rand.Float64()*2 - 1, "z": rand.Float64()*2 - 1},
						"attitude":    attitude,
						"fuel":        100*rand.Float64(),
						"power":       50 + rand.Float64()*50, // percent
						"health":      "operational",
						"last_update": time.Now(),
					},
				})
			})

			satelliteGroup.GET("/:id/coverage", func(c *gin.Context) {
				satelliteId := c.Param("id")
				c.JSON(200, gin.H{
					"satellite_id":  satelliteId,
					"visibility":    "excellent",
					"next_pass":     time.Now().Add(time.Hour * 2),
					"current_position": gin.H{
						"latitude":    35.6762 + (rand.Float64()-0.5)*10,
						"longitude":   139.6503 + (rand.Float64()-0.5)*10,
						"altitude_km": 550 + rand.Float64()*50,
					},
				})
			})

			// GET /api/v1/satellite/available
			satelliteGroup.GET("/available", func(c *gin.Context) {
				satellites := []gin.H{
					{
						"id":              "himawari8",
						"name":            "Himawari-8",
						"type":            "Geostationary Weather",
						"resolution":      1000.0,
						"update_interval": "10m",
						"coverage":        "Asia-Pacific",
						"status":          "active",
						"capabilities":    []string{"visible", "infrared", "water_vapor", "realtime"},
					},
					{
						"id":              "goes16",
						"name":            "GOES-16",
						"type":            "Geostationary Weather",
						"resolution":      500.0,
						"update_interval": "15m",
						"coverage":        "Americas",
						"status":          "active",
						"capabilities":    []string{"visible", "infrared", "lightning", "realtime"},
					},
					{
						"id":              "terra",
						"name":            "Terra",
						"type":            "Earth Observation",
						"resolution":      250.0,
						"update_interval": "1h",
						"coverage":        "Global",
						"status":          "active",
						"capabilities":    []string{"visible", "infrared", "thermal", "multispectral"},
					},
					{
						"id":              "landsat8",
						"name":            "Landsat 8",
						"type":            "Earth Observation",
						"resolution":      30.0,
						"update_interval": "24h",
						"coverage":        "Global",
						"status":          "active",
						"capabilities":    []string{"visible", "infrared", "thermal", "high_resolution"},
					},
					{
						"id":              "worldview3",
						"name":            "WorldView-3",
						"type":            "Commercial High-Resolution",
						"resolution":      0.31,
						"update_interval": "48h",
						"coverage":        "On-demand",
						"status":          "active",
						"capabilities":    []string{"visible", "infrared", "ultra_high_resolution"},
					},
				}
				c.JSON(200, gin.H{
					"satellites": satellites,
					"total":      len(satellites),
					"message":    "Available satellites for video streaming",
				})
			})
		}

		disasterGroup := v1.Group("/disaster")
		{
			disasterGroup.GET("/fires", disasterHandlers.GetFires)
			disasterGroup.GET("/fires/number1", disasterHandlers.GetFiresNumber1)
			disasterGroup.GET("/fires/active", disasterHandlers.GetActiveFires)
			disasterGroup.GET("/fires/global", disasterHandlers.GetGlobalFires)
			disasterGroup.GET("/fires/historical", disasterHandlers.GetHistoricalFires)
			disasterGroup.GET("/fires/area", disasterHandlers.GetFiresByArea)
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

		// Realtime video endpoint (returns dynamic URLs based on query)
		v1.GET("/satellite/video/realtime", videoHandlers.GetRealtimeVideo)

		// Serve dynamic thumbnail by redirecting to a fetched image (seeded by lat/lon/time)
		v1.GET("/satellite/:id/video/thumb", func(c *gin.Context) {
			sid := c.Param("id")
			lat := c.DefaultQuery("lat", "0")
			lon := c.DefaultQuery("lon", "0")
			seed := sid + "_" + lat + "_" + lon + "_" + time.Now().Format("200601021504")
			// Use picsum as a lightweight image source
			c.Redirect(302, "https://picsum.photos/seed/"+seed+"/1280/720.jpg")
		})

		// Serve video stream by redirecting to a public sample based on satellite id (placeholder stream)
		v1.GET("/satellite/:id/video/stream", func(c *gin.Context) {
			// Could switch by :id in future
			c.Redirect(302, "https://www.w3schools.com/html/mov_bbb.mp4")
		})

		// Debris threats endpoint for SatellitePanel
		v1.GET("/mission/debris/:id/threats", func(c *gin.Context) {
			missionID := c.Param("id")
			// シンプルなモック: 位置と速度を含む
			c.JSON(200, gin.H{
				"mission_id": missionID,
				"threats": []gin.H{
					{
						"id": "debris_001",
						"name": "Rocket Fragment",
						"position": gin.H{"x": 1200.0, "y": -350.0, "z": 540.0},
						"velocity": gin.H{"x": 7.5, "y": -1.2, "z": 0.8},
						"danger_level": 7,
						"collision_probability": 0.62,
						"time_to_closest": int((10 * time.Minute).Milliseconds()),
						"closest_distance": 2.5,
						"detected_at": time.Now(),
					},
					{
						"id": "debris_002",
						"name": "Satellite Fragment",
						"position": gin.H{"x": -800.0, "y": 220.0, "z": -150.0},
						"velocity": gin.H{"x": 6.8, "y": 0.9, "z": -0.4},
						"danger_level": 4,
						"collision_probability": 0.18,
						"time_to_closest": int((25 * time.Minute).Milliseconds()),
						"closest_distance": 8.1,
						"detected_at": time.Now(),
					},
				},
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
				"environment":    os.Getenv("GIN_MODE"),
				"port":           ":8080",
				"cors":           "enabled",
				"websocket":      "enabled",
				"ground_control": "active",
				"satellite_mock": "enabled",
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

	// Simple media aliases to avoid 404 in mock response
	// These endpoints redirect to public sample assets
	r.GET("/media/sample.mp4", func(c *gin.Context) {
		c.Redirect(302, "https://www.w3schools.com/html/mov_bbb.mp4")
	})
	r.GET("/media/sample.jpg", func(c *gin.Context) {
		c.Redirect(302, "https://images.unsplash.com/photo-1454789548928-9efd52dc4031?q=80&w=1600&auto=format&fit=crop")
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
	log.Printf("  GET  /api/v1/satellite/available")
	log.Printf("  GET  /api/v1/satellite/video/realtime")
	log.Printf("  GET  /api/docs")

	return r
}
