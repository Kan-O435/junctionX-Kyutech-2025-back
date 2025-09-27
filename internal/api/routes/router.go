package routes

import (
	"github.com/gin-gonic/gin"
	disasterHandlers "junctionx2025back/internal/api/handlers/disaster"
	"junctionx2025back/internal/api/handlers/satellite"
	satelliteVideoHandlers "junctionx2025back/internal/api/handlers/satellite/video"
)

func SetupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Satellite Game Backend is running!",
		})
	})

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// 衛星関連（既存）
		satelliteGroup := v1.Group("/satellite")
		{
			// 既存の衛星制御機能
			satelliteGroup.GET("/:id/orbit", satellite.GetOrbit)
			satelliteGroup.POST("/:id/maneuver", satellite.ExecuteManeuver)
			satelliteGroup.GET("/:id/status", satellite.GetStatus)

			// ===== 新規追加：衛星映像機能 =====

			// 利用可能な衛星一覧
			satelliteGroup.GET("/available", satelliteVideoHandlers.GetAvailableSatellites)

			// リアルタイム映像取得
			satelliteGroup.GET("/video/realtime", satelliteVideoHandlers.GetRealtimeVideo)

			// 映像履歴取得
			satelliteGroup.GET("/video/history", satelliteVideoHandlers.GetVideoHistory)

			// 複数衛星同時観測
			satelliteGroup.POST("/video/multi-view", satelliteVideoHandlers.GetMultiSatelliteView)

			// ライブストリーミング開始
			satelliteGroup.POST("/video/stream/start", satelliteVideoHandlers.StartLiveStream)

			// ライブストリーミング停止
			satelliteGroup.POST("/video/stream/stop", satelliteVideoHandlers.StopLiveStream)

			// 特定衛星の詳細情報
			satelliteGroup.GET("/:id/info", satelliteVideoHandlers.GetSatelliteInfo)

			// 衛星の現在位置と視野範囲
			satelliteGroup.GET("/:id/coverage", satelliteVideoHandlers.GetSatelliteCoverage)

			// 指定地点の観測可能衛星一覧
			satelliteGroup.GET("/coverage/location", satelliteVideoHandlers.GetLocationCoverage)
		}

		// ===== 新規追加：自然災害監視機能 =====
		disasterGroup := v1.Group("/disaster")
		{
			// アクティブな災害一覧
			disasterGroup.GET("/active", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"disasters": []gin.H{
						{
							"id":        "earthquake_001",
							"type":      "earthquake",
							"magnitude": 7.2,
							"location": gin.H{
								"latitude":  35.6762,
								"longitude": 139.6503,
								"country":   "Japan",
							},
							"severity": "high",
							"status":   "active",
						},
						{
							"id":        "typhoon_001",
							"type":      "typhoon",
							"magnitude": 4.0,
							"location": gin.H{
								"latitude":  26.0,
								"longitude": 140.0,
								"country":   "Japan",
							},
							"severity": "critical",
							"status":   "active",
						},
					},
					"total":   2,
					"message": "Active disasters detected",
				})
			})

			// 災害詳細情報
			disasterGroup.GET("/:id", func(c *gin.Context) {
				disasterID := c.Param("id")
				c.JSON(200, gin.H{
					"disaster_id": disasterID,
					"type":        "earthquake",
					"magnitude":   7.2,
					"location": gin.H{
						"latitude":  35.6762,
						"longitude": 139.6503,
						"depth":     30.0,
					},
					"satellite_observations": []gin.H{
						{
							"satellite_id": "himawari8",
							"image_url":    "/api/v1/satellite/himawari8/disaster/" + disasterID,
							"capture_time": "2025-09-27T10:00:00Z",
							"quality":      0.95,
						},
					},
					"message": "Disaster details with satellite imagery",
				})
			})

			// 災害地域のリアルタイム衛星映像
			disasterGroup.GET("/:id/video", satelliteVideoHandlers.GetDisasterVideo)

			//火災地域の緯度経度情報
			disasterGroup.GET("/fires", disasterHandlers.GetFires)
			disasterGroup.GET("/fires/number1", disasterHandlers.GetFiresNumber1)
			disasterGroup.GET("/fires/active", disasterHandlers.GetActiveFires)
			disasterGroup.GET("/fires/global", disasterHandlers.GetGlobalFires)
			disasterGroup.GET("/fires/historical", disasterHandlers.GetHistoricalFires)
			disasterGroup.GET("/fires/area", disasterHandlers.GetFiresByArea)

		}

		// デブリ脅威取得（既存）
		v1.GET("/mission/debris/:id/threats", func(c *gin.Context) {
			missionID := c.Param("id")
			c.JSON(200, gin.H{
				"mission_id": missionID,
				"threats": []gin.H{
					{
						"id":             "debris_001",
						"name":           "Rocket Fragment",
						"distance":       2.5,
						"time_to_impact": 300,
						"danger_level":   7,
					},
					{
						"id":             "debris_002",
						"name":           "Satellite Fragment",
						"distance":       8.1,
						"time_to_impact": 450,
						"danger_level":   4,
					},
				},
				"message": "Sample debris threats",
			})
		})

		// ===== 新規追加：WebSocketエンドポイント =====

		// リアルタイム災害通知用WebSocket
		v1.GET("/ws/disaster", func(c *gin.Context) {
			// WebSocket接続処理（実装は別途必要）
			c.JSON(200, gin.H{
				"message":  "WebSocket endpoint for disaster notifications",
				"endpoint": "ws://localhost:8080/api/v1/ws/disaster",
			})
		})

		// リアルタイム衛星映像ストリーミング用WebSocket
		v1.GET("/ws/satellite/stream", func(c *gin.Context) {
			// WebSocket接続処理（実装は別途必要）
			c.JSON(200, gin.H{
				"message":  "WebSocket endpoint for satellite video streaming",
				"endpoint": "ws://localhost:8080/api/v1/ws/satellite/stream",
			})
		})

		// ===== 新規追加：ファイル配信エンドポイント =====

		// 衛星画像ファイル配信
		v1.Static("/files/satellite", "./data/satellite_images")

		// 災害画像ファイル配信
		v1.Static("/files/disaster", "./data/disaster_images")

		// サムネイル画像配信
		v1.Static("/files/thumbnails", "./data/thumbnails")
	}

	// ===== 新規追加：管理者用エンドポイント =====
	adminGroup := r.Group("/admin")
	{
		// 衛星ステータス一覧
		adminGroup.GET("/satellites", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"satellites": []gin.H{
					{
						"id":           "himawari8",
						"name":         "Himawari-8",
						"status":       "operational",
						"last_contact": "2025-09-27T10:00:00Z",
						"data_quality": 0.95,
					},
					{
						"id":           "goes16",
						"name":         "GOES-16",
						"status":       "operational",
						"last_contact": "2025-09-27T09:55:00Z",
						"data_quality": 0.92,
					},
				},
				"total": 2,
			})
		})

		// システム統計
		adminGroup.GET("/stats", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"total_requests": 1245,
				"active_streams": 3,
				"satellite_health": gin.H{
					"operational": 7,
					"maintenance": 1,
					"offline":     0,
				},
				"uptime":       "99.8%",
				"last_updated": "2025-09-27T10:00:00Z",
			})
		})
	}
}
