package satellite

import (
    "time"
    "junctionx2025back/internal/models/common"
)

type Satellite struct {
    ID          string         `json:"id"`
    Name        string         `json:"name"`
    NoradID     string         `json:"norad_id"`
    PlayerID    string         `json:"player_id"`
    State       SatelliteState `json:"state"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
}

type SatelliteState struct {
    Position     common.Vector3D `json:"position"`      // ECI座標(km)
    Velocity     common.Vector3D `json:"velocity"`      // ECI速度(km/s)
    Attitude     Attitude        `json:"attitude"`      // 姿勢
    Fuel         float64         `json:"fuel"`          // 燃料残量(kg)
    Power        float64         `json:"power"`         // 電力レベル(%)
    Health       string          `json:"health"`        // "healthy", "damaged", "critical"
    LastUpdate   time.Time       `json:"last_update"`
}

type Attitude struct {
    Roll  float64 `json:"roll"`   // ロール角(度)
    Pitch float64 `json:"pitch"`  // ピッチ角(度)
    Yaw   float64 `json:"yaw"`    // ヨー角(度)
}

type OrbitElements struct {
    SemiMajorAxis    float64   `json:"semi_major_axis"`    // 軌道長半径(km)
    Eccentricity     float64   `json:"eccentricity"`       // 離心率
    Inclination      float64   `json:"inclination"`        // 軌道傾斜角(度)
    RAAN             float64   `json:"raan"`               // 昇交点赤経(度)
    ArgumentPerigee  float64   `json:"argument_perigee"`   // 近地点引数(度)
    TrueAnomaly      float64   `json:"true_anomaly"`       // 真近点角(度)
    Epoch            time.Time `json:"epoch"`              // 軌道要素時刻
}

type Maneuver struct {
    ID           string          `json:"id"`
    SatelliteID  string          `json:"satellite_id"`
    PlayerID     string          `json:"player_id"`
    ThrustVector common.Vector3D `json:"thrust_vector"`  // 推力ベクトル(N)
    Duration     float64         `json:"duration"`       // 継続時間(秒)
    StartTime    time.Time       `json:"start_time"`
    Status       string          `json:"status"`         // "planned", "executing", "completed", "failed"
    FuelConsumed float64         `json:"fuel_consumed"`  // 消費燃料(kg)
}