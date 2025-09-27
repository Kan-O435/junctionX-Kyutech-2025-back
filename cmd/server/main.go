// cmd/server/main.go - CORS修正済み完全版

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// 必要な構造体を直接定義
type ChatMessage struct {
	ID          string    `json:"id"`
	Role        string    `json:"role"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
	MessageType string    `json:"message_type"`
	Urgent      bool      `json:"urgent"`
}

type OperatorInfo struct {
	Name    string `json:"name"`
	Rank    string `json:"rank"`
	Station string `json:"station"`
	Shift   string `json:"shift"`
}

type DisasterEvent struct {
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	Severity    string    `json:"severity"`
	Magnitude   float64   `json:"magnitude"`
	Timestamp   time.Time `json:"timestamp"`
	Coordinates []float64 `json:"coordinates"`
	Source      string    `json:"source"`
}

type Mission struct {
	ID           string         `json:"id"`
	CallSign     string         `json:"call_sign"`
	Disaster     DisasterEvent  `json:"disaster"`
	Status       string         `json:"status"`
	Priority     string         `json:"priority"`
	CreatedAt    time.Time      `json:"created_at"`
	ChatHistory  []ChatMessage  `json:"chat_history"`
	OperatorInfo OperatorInfo   `json:"operator_info"`
}

// DemoGroundControlService - デモ用サービス
type DemoGroundControlService struct {
	missions     map[string]*Mission
	messageCount map[string]int
}

// 固定メッセージ定義
const (
	INITIAL_MESSAGE = `**Ground Control to SEISMIC-56, this is Control Officer Johnson at CONTROL-1. Acknowledge.**

**SEISMIC-56, current time 14:00 Zulu.**

**Sitrep:** Fukuoka, Japan, seismic event confirmed. Magnitude 9.0. Severity level extreme. Widespread infrastructure collapse reported. Multiple districts experiencing critical system failures.

Your mission is immediate deployment to the GYMLABO sector, designated Zone Bravo-2. Check surrounding conditions, report on local damage, and provide situational awareness for further operations.

Do you copy? Over.`

	FOLLOWUP_MESSAGE = `**Ground Control to SEISMIC-56, this is Control Officer Johnson at CONTROL-1.**

**Update:** Sensor anomaly confirmed. Seismic alert was a false detection.

Disregard previous orders. No deployment required. We appreciate your cooperation and rapid readiness. 

CONTROL-1 out.`
)

func NewDemoGroundControlService() *DemoGroundControlService {
	return &DemoGroundControlService{
		missions:     make(map[string]*Mission),
		messageCount: make(map[string]int),
	}
}

func (dgcs *DemoGroundControlService) GetMissions() []*Mission {
	var missions []*Mission
	for _, m := range dgcs.missions {
		missions = append(missions, m)
	}
	return missions
}

func (dgcs *DemoGroundControlService) CreateMissionFromDisaster(disasterEvent DisasterEvent) *Mission {
	// 災害イベントのデフォルト値を確実に設定
	if disasterEvent.Type == "" {
		disasterEvent.Type = "earthquake"
	}
	if disasterEvent.Location == "" {
		disasterEvent.Location = "Fukuoka, Japan - GYMLABO Sector"
	}
	if disasterEvent.Severity == "" {
		disasterEvent.Severity = "extreme"
	}
	if disasterEvent.Magnitude == 0 {
		disasterEvent.Magnitude = 9.0
	}
	if disasterEvent.Source == "" {
		disasterEvent.Source = "Demo Seismic Network"
	}
	if disasterEvent.Timestamp.IsZero() {
		disasterEvent.Timestamp = time.Now()
	}
	if len(disasterEvent.Coordinates) == 0 {
		disasterEvent.Coordinates = []float64{130.4017, 33.5904}
	}

	newMission := &Mission{
		ID:       "SEISMIC-56",
		CallSign: "SEISMIC-56",
		Disaster: disasterEvent, // 確実にdisasterフィールドが設定される
		Status:   "ACTIVE",
		Priority: "PRIORITY ALPHA - CRITICAL",
		CreatedAt: time.Now(),
		ChatHistory: []ChatMessage{},
		OperatorInfo: OperatorInfo{
			Name:    "Johnson",
			Rank:    "Control Officer",
			Station: "CONTROL-1",
			Shift:   "DAY SHIFT",
		},
	}

	// 初期ブリーフィングを追加
	initialBriefing := ChatMessage{
		ID:          fmt.Sprintf("MSG_%d", time.Now().UnixNano()),
		Role:        "assistant",
		Content:     INITIAL_MESSAGE,
		Timestamp:   time.Now(),
		MessageType: "mission_briefing",
		Urgent:      true,
	}
	newMission.ChatHistory = append(newMission.ChatHistory, initialBriefing)

	dgcs.missions[newMission.ID] = newMission
	dgcs.messageCount[newMission.ID] = 0

	log.Printf("Demo Mission Created: %s with disaster: %+v", newMission.CallSign, newMission.Disaster)
	return newMission
}

func (dgcs *DemoGroundControlService) GetMission(id string) (*Mission, bool) {
	m, exists := dgcs.missions[id]
	return m, exists
}

func (dgcs *DemoGroundControlService) ProcessFieldMessage(missionID, userMessage string) ChatMessage {
	m, exists := dgcs.missions[missionID]
	if !exists {
		return ChatMessage{
			ID:          fmt.Sprintf("MSG_ERROR_%d", time.Now().UnixNano()),
			Role:        "assistant",
			Content:     "Ground Control to Field Team. Invalid Mission ID. Please verify call sign. Over.",
			Timestamp:   time.Now(),
			MessageType: "error",
			Urgent:      false,
		}
	}

	// ユーザーメッセージを履歴に追加（ユニークIDを確保）
	userMsg := ChatMessage{
		ID:          fmt.Sprintf("MSG_USER_%d_%d", time.Now().UnixNano(), dgcs.messageCount[missionID]),
		Role:        "user",
		Content:     userMessage,
		Timestamp:   time.Now(),
		MessageType: "field_report",
		Urgent:      false,
	}
	m.ChatHistory = append(m.ChatHistory, userMsg)

	// メッセージカウントを増加
	dgcs.messageCount[missionID]++

	// デモ用固定レスポンス
	var responseContent string
	if dgcs.messageCount[missionID] == 1 {
		responseContent = FOLLOWUP_MESSAGE
		m.Status = "COMPLETED"
		m.Priority = "PRIORITY DELTA - RESOLVED"
	} else {
		responseContent = `**Ground Control to SEISMIC-56.**

Message received and acknowledged. Mission status: COMPLETED. 

Standing by for further instructions.

CONTROL-1 out.`
	}

	// 管制応答を作成（ユニークIDを確保）
	controlMsg := ChatMessage{
		ID:          fmt.Sprintf("MSG_ASSISTANT_%d_%d", time.Now().UnixNano(), dgcs.messageCount[missionID]),
		Role:        "assistant",
		Content:     responseContent,
		Timestamp:   time.Now(),
		MessageType: "control_response",
		Urgent:      dgcs.messageCount[missionID] == 1,
	}
	m.ChatHistory = append(m.ChatHistory, controlMsg)

	log.Printf("Demo Response Sent: Message #%d for %s", dgcs.messageCount[missionID], missionID)
	return controlMsg
}

func (dgcs *DemoGroundControlService) CreateInitialMission() *Mission {
	demoDisaster := DisasterEvent{
		Type:        "earthquake",
		Location:    "Fukuoka, Japan - GYMLABO Sector",
		Severity:    "extreme",
		Magnitude:   9.0,
		Timestamp:   time.Now(),
		Coordinates: []float64{130.4017, 33.5904},
		Source:      "Demo Seismic Network",
	}

	return dgcs.CreateMissionFromDisaster(demoDisaster)
}

func (dgcs *DemoGroundControlService) ResetDemo() {
	dgcs.missions = make(map[string]*Mission)
	dgcs.messageCount = make(map[string]int)
	log.Printf("Demo reset completed")
}

func main() {
	// デモ用サービス作成
	demoService := NewDemoGroundControlService()
	
	// デモ用初期ミッション自動作成
	initialMission := demoService.CreateInitialMission()
	log.Printf("DEMO MODE: Initial mission created - %s", initialMission.CallSign)

	// Gin router の設定
	router := gin.Default()
	
	// CORS設定を強化
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "false")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		log.Printf("Request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("Origin"))
		c.Next()
	})

	// ルート設定
	setupRoutes(router, demoService)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Demo Mission: SEISMIC-56")
	log.Printf("Send any message to trigger false alarm response")
	
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(router *gin.Engine, demoService *DemoGroundControlService) {
	// API グループを作成
	api := router.Group("/api/v1")

	// OPTIONS プリフライトリクエストを明示的に処理
	api.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	// ミッション関連API
	api.GET("/missions", func(c *gin.Context) {
		missions := demoService.GetMissions()
		
		// レスポンス前にdisasterフィールドが存在することを確認
		for _, mission := range missions {
			if mission.Disaster.Type == "" {
				mission.Disaster.Type = "earthquake"
			}
			if mission.Disaster.Location == "" {
				mission.Disaster.Location = "Unknown Location"
			}
			if mission.Disaster.Severity == "" {
				mission.Disaster.Severity = "unknown"
			}
		}
		
		// フロントエンドが期待する形式に合わせる
		if len(missions) > 0 {
			// 最初のミッションを直接返す（フロントエンドが単一ミッションを期待している場合）
			c.JSON(200, missions[0])
		} else {
			// ミッションがない場合はデフォルトミッションを作成
			defaultMission := demoService.CreateInitialMission()
			c.JSON(200, defaultMission)
		}
	})
	
	api.POST("/missions", func(c *gin.Context) {
		var request DisasterEvent
		
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("JSON binding error: %v", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Creating mission from disaster: %+v", request)

		if request.Timestamp.IsZero() {
			request.Timestamp = time.Now()
		}
		
		if request.Source == "" {
			request.Source = "Frontend Creation"
		}

		mission := demoService.CreateMissionFromDisaster(request)
		
		// ミッションオブジェクトを直接返す（フロントエンドが期待する形式）
		c.JSON(201, mission)
	})
	
	api.POST("/missions/:id/message", func(c *gin.Context) {
		missionID := c.Param("id")
		
		var request struct {
			Message string `json:"message" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Processing message for mission %s: %s", missionID, request.Message)

		response := demoService.ProcessFieldMessage(missionID, request.Message)
		
		log.Printf("Generated response: %+v", response)
		
		// レスポンスメッセージを直接返す（フロントエンドが期待する形式）
		c.JSON(200, response)
	})
	
	api.GET("/missions/:id", func(c *gin.Context) {
		missionID := c.Param("id")
		mission, exists := demoService.GetMission(missionID)
		if !exists {
			c.JSON(404, gin.H{
				"error":      "Mission not found",
				"mission_id": missionID,
			})
			return
		}
		
		// disasterフィールドの存在を確保
		if mission.Disaster.Type == "" {
			mission.Disaster.Type = "earthquake"
		}
		if mission.Disaster.Location == "" {
			mission.Disaster.Location = "Unknown Location"
		}
		if mission.Disaster.Severity == "" {
			mission.Disaster.Severity = "unknown"
		}
		
		// ミッションオブジェクトを直接返す
		c.JSON(200, mission)
	})
	
	api.POST("/missions/demo/reset", func(c *gin.Context) {
		demoService.ResetDemo()
		demoService.CreateInitialMission()
		c.JSON(200, gin.H{
			"status":  "reset",
			"message": "Demo scenario reset and new mission created",
		})
	})

	// 衛星関連API（デモ用モックデータ）
	api.GET("/satellite/:id/orbit", func(c *gin.Context) {
		satelliteID := c.Param("id")
		c.JSON(200, gin.H{
			"satellite_id": satelliteID,
			"orbit": gin.H{
				"altitude": 550.2,
				"inclination": 53.0,
				"period": 95.64,
				"orbital_speed": 7.66,
			},
			"position": gin.H{
				"latitude": 35.6762,
				"longitude": 139.6503,
				"altitude_km": 550.2,
			},
			"timestamp": time.Now(),
		})
	})
	
	api.GET("/satellite/:id/coverage", func(c *gin.Context) {
		satelliteID := c.Param("id")
		c.JSON(200, gin.H{
			"satellite_id": satelliteID,
			"coverage": gin.H{
				"footprint_radius": 2500,
				"elevation_angle": 45.5,
				"visibility": "visible",
			},
			"next_pass": time.Now().Add(time.Hour * 2),
			"timestamp": time.Now(),
		})
	})
	
	api.GET("/satellite/:id/status", func(c *gin.Context) {
		satelliteID := c.Param("id")
		c.JSON(200, gin.H{
			"satellite_id": satelliteID,
			"status": gin.H{
				"operational": true,
				"power": 85.5,
				"health": "nominal",
				"fuel": 75.2,
				"attitude": gin.H{
					"roll": 0.1,
					"pitch": -0.3,
					"yaw": 0.05,
				},
			},
			"timestamp": time.Now(),
		})
	})

	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"mode":    "demo",
			"message": "Disaster Response System operational",
			"demo_info": gin.H{
				"mission_id": "SEISMIC-56",
				"scenario":   "Fukuoka Earthquake Detection",
			},
		})
	})

	// デモ情報
	router.GET("/demo", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"demo_mode": true,
			"scenario":  "Fukuoka Earthquake Response",
			"mission_id": "SEISMIC-56",
			"instructions": gin.H{
				"step_1": "Mission SEISMIC-56 automatically created with earthquake alert",
				"step_2": "Send any message to /api/v1/missions/SEISMIC-56/message",
				"step_3": "System responds with false alarm notification",
				"step_4": "Mission status changes to COMPLETED",
			},
			"endpoints": gin.H{
				"missions":        "/api/v1/missions",
				"send_message":    "/api/v1/missions/SEISMIC-56/message",
				"reset_demo":      "/api/v1/missions/demo/reset",
				"satellite_orbit": "/api/v1/satellite/STARLINK-32713/orbit",
				"satellite_status": "/api/v1/satellite/STARLINK-32713/status",
			},
		})
	})

	// テスト用エンドポイント
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test successful"})
	})

	router.POST("/test", func(c *gin.Context) {
		var body map[string]interface{}
		c.ShouldBindJSON(&body)
		c.JSON(200, gin.H{"received": body, "message": "POST test successful"})
	})
}