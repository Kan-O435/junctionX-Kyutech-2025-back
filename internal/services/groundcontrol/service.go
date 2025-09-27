package groundcontrol

import (
    "fmt"
    "log"
    "strings"
    "time"
    
    "junctionx2025back/internal/models/disaster"
    "junctionx2025back/internal/models/mission"
    "junctionx2025back/internal/services/llm"
)

type GroundControlService struct {
    geminiService *llm.GeminiService
    missions      map[string]*mission.Mission
}

func NewGroundControlService(apiKey string) *GroundControlService {
    return &GroundControlService{
        geminiService: llm.NewGeminiService(apiKey),
        missions:      make(map[string]*mission.Mission),
    }
}

func (gcs *GroundControlService) getGroundControlPrompt() string {
    return `You are a Ground Control Officer at a Disaster Response Command Center. Respond with the following characteristics:

【Communication Style】
- Use military/aviation control terminology
- Be concise and direct with clear points
- Use radio communication phrases: "Ground Control to Field Team", "Copy that", "Roger"
- ALWAYS use call signs (e.g., "Field Team Alpha")

【Required Terminology】
- "Ground Control to [team], acknowledge"
- "Copy that", "Roger", "Negative", "Stand by"
- "Proceed", "Hold position", "Sitrep"
- "Priority Alpha", "Code Red"

【Response Format Example】
"Ground Control to Field Team Alpha, this is Control Officer Johnson.
Current time 15:42 Zulu. Your mission is...
Do you copy? Over."

Respond as a professional control officer with military precision.`
}

// GetMissions - Handler compatible method name
func (gcs *GroundControlService) GetMissions() []*mission.Mission {
    var missions []*mission.Mission
    for _, m := range gcs.missions {
        missions = append(missions, m)
    }
    return missions
}

// CreateMissionFromDisaster - Handler compatible method name
func (gcs *GroundControlService) CreateMissionFromDisaster(disasterEvent disaster.DisasterEvent) *mission.Mission {
    newMission := &mission.Mission{
        ID:       fmt.Sprintf("MISSION_%d", time.Now().Unix()),
        CallSign: gcs.generateCallSign(disasterEvent),
        Disaster: disasterEvent,
        Status:   "ACTIVE",
        Priority: gcs.calculatePriority(disasterEvent),
        CreatedAt: time.Now(),
        ChatHistory: []mission.ChatMessage{},
        OperatorInfo: gcs.generateOperator(),
    }
    
    briefing := gcs.generateInitialBriefing(newMission)
    newMission.ChatHistory = append(newMission.ChatHistory, briefing)
    
    gcs.missions[newMission.ID] = newMission
    
    log.Printf("🎯 Mission Deployed: %s [%s]", newMission.CallSign, newMission.Priority)
    return newMission
}

func (gcs *GroundControlService) GetMission(id string) (*mission.Mission, bool) {
    m, exists := gcs.missions[id]
    return m, exists
}

func (gcs *GroundControlService) ProcessFieldMessage(missionID, userMessage string) mission.ChatMessage {
    m, exists := gcs.missions[missionID]
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
    
    // Add user message
    userMsg := mission.ChatMessage{
        ID:          fmt.Sprintf("MSG_%d", time.Now().UnixNano()),
        Role:        "user",
        Content:     userMessage,
        Timestamp:   time.Now(),
        MessageType: gcs.detectMessageType(userMessage),
        Urgent:      gcs.isUrgent(userMessage),
    }
    m.ChatHistory = append(m.ChatHistory, userMsg)
    
    // Update priority if urgent
    if userMsg.Urgent && !strings.Contains(m.Priority, "CRITICAL") {
        m.Priority = "PRIORITY ALPHA - CRITICAL"
    }
    
    // Generate response - FIXED: Changed prompt to responsePrompt
    responsePrompt := fmt.Sprintf(`%s

【Current Mission】
- Call Sign: %s
- Disaster: %s (%s severity)
- Location: %s
- Magnitude: %.1f
- Priority: %s
- Operator: %s (%s)

【Recent Communications】
%s

【Field Report】
"%s"

Respond as Ground Control with specific instructions and military precision.`,
        gcs.getGroundControlPrompt(),
        m.CallSign,
        m.Disaster.Type,
        m.Disaster.Severity,
        m.Disaster.Location,
        m.Disaster.Magnitude,
        m.Priority,
        m.OperatorInfo.Name,
        m.OperatorInfo.Rank,
        gcs.getRecentCommsHistory(m),
        userMessage,
    )
    
    // Use CallAPI method (correct method name)
    response := gcs.geminiService.CallAPI(responsePrompt)
    if response == "" {
        response = "Ground Control to Field Team. Communication system experiencing interference. Switching to backup channel. Standby. Over."
        log.Printf("❌ Gemini API Error: Empty response")
    }
    
    controlMsg := mission.ChatMessage{
        ID:          fmt.Sprintf("MSG_%d", time.Now().UnixNano()),
        Role:        "assistant",
        Content:     response,
        Timestamp:   time.Now(),
        MessageType: "control_response",
        Urgent:      userMsg.Urgent,
    }
    m.ChatHistory = append(m.ChatHistory, controlMsg)
    
    return controlMsg
}

// Helper functions
func (gcs *GroundControlService) generateCallSign(disasterEvent disaster.DisasterEvent) string {
    callSigns := map[string]string{
        "earthquake": "SEISMIC",
        "wildfire":   "FIREBIRD", 
        "tsunami":    "TIDAL",
        "volcano":    "MAGMA",
        "hurricane":  "CYCLONE",
        "flood":      "DELUGE",
        "tornado":    "VORTEX",
    }
    
    prefix := callSigns[disasterEvent.Type]
    if prefix == "" {
        prefix = "ALPHA"
    }
    
    suffix := time.Now().Unix() % 100
    return fmt.Sprintf("%s-%02d", prefix, suffix)
}

func (gcs *GroundControlService) calculatePriority(disasterEvent disaster.DisasterEvent) string {
    switch disasterEvent.Severity {
    case "extreme":
        return "PRIORITY ALPHA - CRITICAL"
    case "critical", "high":
        return "PRIORITY BRAVO - HIGH"
    case "medium":
        return "PRIORITY CHARLIE - MEDIUM"
    default:
        return "PRIORITY DELTA - ROUTINE"
    }
}

func (gcs *GroundControlService) generateOperator() mission.OperatorInfo {
    operators := []string{"Johnson", "Smith", "Williams", "Brown", "Davis", "Wilson"}
    ranks := []string{"Control Officer", "Senior Controller", "Flight Director", "Mission Commander"}
    stations := []string{"CONTROL-1", "CONTROL-2", "BACKUP-1"}
    
    now := time.Now().Unix()
    return mission.OperatorInfo{
        Name:    operators[now%int64(len(operators))],
        Rank:    ranks[now%int64(len(ranks))],
        Station: stations[now%int64(len(stations))],
        Shift:   gcs.getCurrentShift(),
    }
}

func (gcs *GroundControlService) getCurrentShift() string {
    hour := time.Now().Hour()
    switch {
    case hour >= 6 && hour < 14:
        return "DAY SHIFT"
    case hour >= 14 && hour < 22:
        return "EVENING SHIFT"
    default:
        return "NIGHT SHIFT"
    }
}

func (gcs *GroundControlService) generateInitialBriefing(m *mission.Mission) mission.ChatMessage {
    // FIXED: Changed prompt to responsePrompt
    responsePrompt := fmt.Sprintf(`%s

【Mission Initialization】
- Call Sign: %s
- Disaster: %s at %s
- Magnitude: %.1f
- Severity: %s
- Priority: %s
- Operator: %s (%s)
- Station: %s

Generate initial mission briefing with military precision and radio protocol.`,
        gcs.getGroundControlPrompt(),
        m.CallSign,
        m.Disaster.Type,
        m.Disaster.Location,
        m.Disaster.Magnitude,
        m.Disaster.Severity,
        m.Priority,
        m.OperatorInfo.Name,
        m.OperatorInfo.Rank,
        m.OperatorInfo.Station,
    )
    
    // Use CallAPI method (correct method name)
    response := gcs.geminiService.CallAPI(responsePrompt)
    if response == "" {
        response = fmt.Sprintf(
            "Ground Control to %s. This is %s at %s. Mission initiated for %s event at %s. Severity level: %s. %s. Standing by for field reports. How copy? Over.",
            m.CallSign,
            m.OperatorInfo.Name,
            m.OperatorInfo.Station,
            m.Disaster.Type,
            m.Disaster.Location,
            m.Disaster.Severity,
            m.Priority,
        )
        log.Printf("❌ Gemini API Error in briefing: Empty response")
    }
    
    return mission.ChatMessage{
        ID:          fmt.Sprintf("MSG_%d", time.Now().UnixNano()),
        Role:        "assistant",
        Content:     response,
        Timestamp:   time.Now(),
        MessageType: "mission_briefing",
        Urgent:      strings.Contains(m.Priority, "ALPHA"),
    }
}

func (gcs *GroundControlService) getRecentCommsHistory(m *mission.Mission) string {
    if len(m.ChatHistory) <= 1 {
        return "(No previous communications)"
    }
    
    var history strings.Builder
    recent := m.ChatHistory
    if len(recent) > 4 {
        recent = recent[len(recent)-4:]
    }
    
    for _, msg := range recent {
        sender := "Field Team"
        if msg.Role == "assistant" {
            sender = "Ground Control"
        }
        urgentFlag := ""
        if msg.Urgent {
            urgentFlag = " [URGENT]"
        }
        history.WriteString(fmt.Sprintf("- %s%s: %s\n", sender, urgentFlag, msg.Content))
    }
    
    return history.String()
}

func (gcs *GroundControlService) detectMessageType(message string) string {
    lower := strings.ToLower(message)
    
    if strings.Contains(lower, "photo") || strings.Contains(lower, "image") || strings.Contains(lower, "picture") {
        return "photo_report"
    } else if strings.Contains(lower, "location") || strings.Contains(lower, "position") || strings.Contains(lower, "coordinates") {
        return "location_report"
    } else if strings.Contains(lower, "complete") || strings.Contains(lower, "finished") || strings.Contains(lower, "done") {
        return "completion_report"
    } else if gcs.isUrgent(message) {
        return "emergency_report"
    } else if strings.Contains(lower, "status") || strings.Contains(lower, "update") {
        return "status_report"
    } else if strings.Contains(lower, "?") || strings.Contains(lower, "how") || strings.Contains(lower, "what") {
        return "field_question"
    }
    return "general_communication"
}

func (gcs *GroundControlService) isUrgent(message string) bool {
    urgentKeywords := []string{
        "emergency", "urgent", "help", "danger", "injured", "trapped",
        "collapse", "fire", "explosion", "medical", "critical", "casualties",
        "immediate", "asap", "now", "panic", "rescue", "bleeding",
        "code red", "mayday", "sos", "distress",
    }
    
    lower := strings.ToLower(message)
    for _, keyword := range urgentKeywords {
        if strings.Contains(lower, keyword) {
            return true
        }
    }
    return false
}