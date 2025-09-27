package mission

import (
    "time"
    "junctionx2025back/internal/models/common"
)

type DebrisAvoidanceMission struct {
    Base                                    // 基底ミッション構造を埋め込み
    Threats         []DebrisThreat         `json:"threats"`
    AvoidanceActions []AvoidanceAction     `json:"avoidance_actions"`
    Settings        DebrisSettings         `json:"settings"`
}

type Base struct {
    ID          string          `json:"id"`
    Type        string          `json:"type"`
    PlayerID    string          `json:"player_id"`
    SatelliteID string          `json:"satellite_id"`
    Status      string          `json:"status"`      // "active", "completed", "failed"
    Score       int             `json:"score"`
    StartTime   time.Time       `json:"start_time"`
    EndTime     *time.Time      `json:"end_time,omitempty"`
    Progress    float64         `json:"progress"`    // 0.0 - 1.0
}

type DebrisThreat struct {
    ID                  string          `json:"id"`
    NoradID             string          `json:"norad_id"`
    Name                string          `json:"name"`
    Position            common.Vector3D `json:"position"`
    Velocity            common.Vector3D `json:"velocity"`
    Size                float64         `json:"size"`           // メートル
    Mass                float64         `json:"mass"`           // kg
    DangerLevel         int             `json:"danger_level"`   // 1-10
    TimeToClosest       time.Duration   `json:"time_to_closest"`
    ClosestDistance     float64         `json:"closest_distance"` // km
    CollisionProbability float64        `json:"collision_probability"` // 0-1
    DetectedAt          time.Time       `json:"detected_at"`
}

type AvoidanceAction struct {
    ID              string          `json:"id"`
    ThreatID        string          `json:"threat_id"`
    ActionType      string          `json:"action_type"`    // "dodge", "brake", "accelerate"
    ThrustVector    common.Vector3D `json:"thrust_vector"`
    Duration        float64         `json:"duration"`
    ExecutedAt      time.Time       `json:"executed_at"`
    Result          string          `json:"result"`         // "success", "partial", "failed"
    FuelCost        float64         `json:"fuel_cost"`
    SafetyMargin    float64         `json:"safety_margin"`  // 実際の回避距離
}

type DebrisSettings struct {
    MaxThreats          int     `json:"max_threats"`
    ThreatFrequency     float64 `json:"threat_frequency"`     // threats per minute
    CollisionThreshold  float64 `json:"collision_threshold"`  // km
    FuelLimit           float64 `json:"fuel_limit"`           // kg
    TimeLimit           int     `json:"time_limit"`           // seconds
    Difficulty          string  `json:"difficulty"`           // "easy", "normal", "hard"
}