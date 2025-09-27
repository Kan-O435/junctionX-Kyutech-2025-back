package services

import (
    "math"
    "math/rand"
    "time"
    "junctionx2025back/internal/models/common"
    "junctionx2025back/internal/models/mission"
)

type DebrisService struct{}

func NewDebrisService() *DebrisService {
    return &DebrisService{}
}

// デブリ脅威データを生成（シミュレーション用）
func (s *DebrisService) GenerateDebrisThreats(missionID string) []mission.DebrisThreat {
    threats := make([]mission.DebrisThreat, 0, 8)
    earthRadius := 6371.0 // km
    rand.Seed(time.Now().UnixNano())

    for i := 0; i < 8; i++ {
        // 地球周辺の軌道にランダムに配置
        altitude := 400 + rand.Float64()*600 // 400-1000km高度
        angle := (float64(i) / 8.0) * 2 * math.Pi
        radius := earthRadius + altitude

        x := math.Cos(angle) * radius
        y := (rand.Float64() - 0.5) * 200 // ±100kmの高さ変動
        z := math.Sin(angle) * radius

        // 軌道速度（km/s）
        orbitalSpeed := math.Sqrt(398600 / radius) // 軌道速度の簡易計算
        vx := -math.Sin(angle) * orbitalSpeed
        vy := (rand.Float64() - 0.5) * 0.5
        vz := math.Cos(angle) * orbitalSpeed

        // 危険度に基づく衝突確率の計算
        dangerLevel := rand.Intn(10) + 1
        collisionProbability := float64(dangerLevel) / 10.0 * rand.Float64() * 0.8

        // 最接近時間（1-2時間の範囲）
        timeToClosest := time.Duration(rand.Float64()*7200000) * time.Millisecond

        // 最接近距離（1-50kmの範囲）
        closestDistance := rand.Float64()*50 + 1

        threat := mission.DebrisThreat{
            ID:                  missionID + "_debris_" + string(rune('0'+i)),
            NoradID:            "NORAD-" + string(rune('0'+10000+i)),
            Name:               "デブリ " + string(rune('0'+i+1)),
            Position:           common.Vector3D{X: x, Y: y, Z: z},
            Velocity:           common.Vector3D{X: vx, Y: vy, Z: vz},
            Size:               rand.Float64()*5 + 0.5, // 0.5-5.5m
            Mass:               rand.Float64()*1000 + 100, // 100-1100kg
            DangerLevel:        dangerLevel,
            TimeToClosest:      timeToClosest,
            ClosestDistance:    closestDistance,
            CollisionProbability: collisionProbability,
            DetectedAt:         time.Now(),
        }

        threats = append(threats, threat)
    }

    return threats
}

// デブリ統計情報を生成
func (s *DebrisService) GetDebrisStats(threats []mission.DebrisThreat) map[string]interface{} {
    if len(threats) == 0 {
        return map[string]interface{}{
            "total_threats": 0,
            "high_risk_count": 0,
            "collision_probability_avg": 0.0,
            "avg_danger_level": 0.0,
        }
    }

    highRiskCount := 0
    totalCollisionProb := 0.0
    totalDangerLevel := 0

    for _, threat := range threats {
        if threat.DangerLevel >= 7 {
            highRiskCount++
        }
        totalCollisionProb += threat.CollisionProbability
        totalDangerLevel += threat.DangerLevel
    }

    return map[string]interface{}{
        "total_threats": len(threats),
        "high_risk_count": highRiskCount,
        "collision_probability_avg": totalCollisionProb / float64(len(threats)),
        "avg_danger_level": float64(totalDangerLevel) / float64(len(threats)),
    }
}
