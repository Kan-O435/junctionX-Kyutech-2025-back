package mission

import (
    "time"
    "junctionx2025back/internal/models/disaster"
)

type Mission struct {
    ID           string                `json:"id"`
    CallSign     string                `json:"call_sign"`
    Disaster     disaster.DisasterEvent `json:"disaster"`
    Status       string                `json:"status"`
    Priority     string                `json:"priority"`
    CreatedAt    time.Time             `json:"created_at"`
    ChatHistory  []ChatMessage         `json:"chat_history"`
    OperatorInfo OperatorInfo          `json:"operator_info"`
}

type ChatMessage struct {
    ID          string    `json:"id"`
    Role        string    `json:"role"`
    Content     string    `json:"content"`
    Timestamp   time.Time `json:"timestamp"`
    MessageType string    `json:"message_type"`
    Urgent      bool      `json:"urgent"`
}

type OperatorInfo struct {
    Name     string `json:"name"`
    Rank     string `json:"rank"`
    Station  string `json:"station"`
    Shift    string `json:"shift"`
}