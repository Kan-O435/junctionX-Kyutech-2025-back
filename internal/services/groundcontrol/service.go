// internal/services/groundcontrol/demo_service.go
package groundcontrol

import (
	"fmt"
	"log"
	"strings"
	"time"

	"junctionx2025back/internal/models/disaster"
	"junctionx2025back/internal/models/mission"
)

// DemoGroundControlService - デモ用の固定レスポンスサービス
type DemoGroundControlService struct {
	missions    map[string]*mission.Mission
	messageCount map[string]int // ミッション毎のメッセージカウント
}

func NewDemoGroundControlService() *DemoGroundControlService {
	return &DemoGroundControlService{
		missions:     make(map[string]*mission.Mission),
		messageCount: make(map[string]int),
	}
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

// GetMissions - ミッション一覧取得
func (dgcs *DemoGroundControlService) GetMissions() []*mission.Mission {
	var missions []*mission.Mission
	for _, m := range dgcs.missions {
		missions = append(missions, m)
	}
	return missions
}

// CreateMissionFromDisaster - 災害からミッション生成（デモ用）
func (dgcs *DemoGroundControlService) CreateMissionFromDisaster(disasterEvent disaster.DisasterEvent) *mission.Mission {
	newMission := &mission.Mission{
		ID:       "SEISMIC-56", // 固定ID
		CallSign: "SEISMIC-56",
		Disaster: disasterEvent,
		Status:   "ACTIVE",
		Priority: "PRIORITY ALPHA - CRITICAL",
		CreatedAt: time.Now(),
		ChatHistory: []mission.ChatMessage{},
		OperatorInfo: mission.OperatorInfo{
			Name:    "Johnson",
			Rank:    "Control Officer",
			Station: "CONTROL-1",
			Shift:   "DAY SHIFT",
		},
	}

	// 初期ブリーフィングを追加
	initialBriefing := mission.ChatMessage{
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

	log.Printf("🎯 Demo Mission Created: %s", newMission.CallSign)
	return newMission
}

// GetMission - 特定ミッション取得
func (dgcs *DemoGroundControlService) GetMission(id string) (*mission.Mission, bool) {
	m, exists := dgcs.missions[id]
	return m, exists
}

// ProcessFieldMessage - フィールドメッセージ処理（デモ用固定レスポンス）
func (dgcs *DemoGroundControlService) ProcessFieldMessage(missionID, userMessage string) mission.ChatMessage {
	m, exists := dgcs.missions[missionID]
	if !exists {
		return mission.ChatMessage{
			ID:          fmt.Sprintf("MSG_%d", time.Now().UnixNano()),
			Role:        "assistant",
			Content:     "Ground Control to Field Team. Invalid Mission ID. Please verify call sign. Over.",
			Timestamp:   time.Now(),
			MessageType: "error",
			Urgent:      false,
		}
	}

	// ユーザーメッセージを履歴に追加
	userMsg := mission.ChatMessage{
		ID:          fmt.Sprintf("MSG_%d", time.Now().UnixNano()),
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
		// 最初のメッセージに対してはフォローアップメッセージ
		responseContent = FOLLOWUP_MESSAGE
		
		// ミッションステータスを更新
		m.Status = "COMPLETED"
		m.Priority = "PRIORITY DELTA - RESOLVED"
	} else {
		// 2回目以降は標準的な応答
		responseContent = `**Ground Control to SEISMIC-56.**

Message received and acknowledged. Mission status: COMPLETED. 

Standing by for further instructions.

CONTROL-1 out.`
	}

	// 管制応答を作成
	controlMsg := mission.ChatMessage{
		ID:          fmt.Sprintf("MSG_%d", time.Now().UnixNano()),
		Role:        "assistant",
		Content:     responseContent,
		Timestamp:   time.Now(),
		MessageType: "control_response",
		Urgent:      dgcs.messageCount[missionID] == 1, // 最初のレスポンスのみ緊急
	}
	m.ChatHistory = append(m.ChatHistory, controlMsg)

	log.Printf("📡 Demo Response Sent: Message #%d for %s", dgcs.messageCount[missionID], missionID)
	return controlMsg
}

// CreateInitialMission - デモ開始用の初期ミッション作成
func (dgcs *DemoGroundControlService) CreateInitialMission() *mission.Mission {
	// デモ用災害イベント
	demoDisaster := disaster.DisasterEvent{
		Type:        "earthquake",
		Location:    "Fukuoka, Japan - GYMLABO Sector",
		Severity:    "extreme",
		Magnitude:   9.0,
		Timestamp:   time.Now(),
		Coordinates: []float64{130.4017, 33.5904}, // 福岡の座標
		Source:      "Demo Seismic Network",
	}

	return dgcs.CreateMissionFromDisaster(demoDisaster)
}

// ResetDemo - デモリセット（必要に応じて）
func (dgcs *DemoGroundControlService) ResetDemo() {
	dgcs.missions = make(map[string]*mission.Mission)
	dgcs.messageCount = make(map[string]int)
	log.Printf("🔄 Demo reset completed")
}