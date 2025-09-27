package mission

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    
    "junctionx2025back/internal/models/disaster"
    "junctionx2025back/internal/services/groundcontrol"
)

type MissionHandler struct {
    groundControl *groundcontrol.GroundControlService
    upgrader      websocket.Upgrader
}

func NewMissionHandler(groundControl *groundcontrol.GroundControlService) *MissionHandler {
    return &MissionHandler{
        groundControl: groundControl,
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool { return true },
        },
    }
}

func (h *MissionHandler) GetMissions(c *gin.Context) {
    missions := h.groundControl.GetMissions()
    c.JSON(200, gin.H{
        "missions": missions,
        "total":    len(missions),
    })
}

func (h *MissionHandler) CreateMission(c *gin.Context) {
    var disaster disaster.DisasterEvent
    if err := c.ShouldBindJSON(&disaster); err != nil {
        c.JSON(400, gin.H{"error": "Invalid disaster data"})
        return
    }
    
    mission := h.groundControl.CreateMissionFromDisaster(disaster)
    c.JSON(201, mission)
}

func (h *MissionHandler) GetMission(c *gin.Context) {
    missionID := c.Param("id")
    mission, exists := h.groundControl.GetMission(missionID)
    if !exists {
        c.JSON(404, gin.H{"error": "Mission not found"})
        return
    }
    c.JSON(200, mission)
}

func (h *MissionHandler) SendMessage(c *gin.Context) {
    missionID := c.Param("id")
    
    var request struct {
        Message string `json:"message" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": "Message required"})
        return
    }
    
    response := h.groundControl.ProcessFieldMessage(missionID, request.Message)
    c.JSON(200, response)
}

func (h *MissionHandler) WebSocketChat(c *gin.Context) {
    missionID := c.Param("id")
    
    conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    defer conn.Close()
    
    // Send existing history
    if mission, exists := h.groundControl.GetMission(missionID); exists {
        for _, msg := range mission.ChatHistory {
            conn.WriteJSON(msg)
        }
    }
    
    // Handle real-time messages
    for {
        var input struct {
            Message string `json:"message"`
        }
        
        if err := conn.ReadJSON(&input); err != nil {
            break
        }
        
        response := h.groundControl.ProcessFieldMessage(missionID, input.Message)
        conn.WriteJSON(response)
    }
}
