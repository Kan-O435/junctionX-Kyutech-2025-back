package satellite

import (
    "fmt"
    "time"

    "junctionx2025back/internal/models/common"
    "junctionx2025back/internal/models/satellite"
)

type OrbitService struct{}

func NewOrbitService() *OrbitService {
    return &OrbitService{}
}

type OrbitState struct {
    Position common.Vector3D
    Velocity common.Vector3D
    Elements satellite.OrbitElements
}

// 現在の軌道状態を取得
func (s *OrbitService) GetCurrentOrbit(satelliteID string) (*OrbitState, error) {
    state := &OrbitState{
        Position: common.Vector3D{X: 6800.0, Y: 0.0, Z: 0.0},
        Velocity: common.Vector3D{X: 0.0, Y: 7.66, Z: 0.0},
        Elements: satellite.OrbitElements{
            SemiMajorAxis:   6800.0,
            Eccentricity:    0.001,
            Inclination:     51.6,
            RAAN:            0.0,
            ArgumentPerigee: 0.0,
            TrueAnomaly:     0.0,
            Epoch:           time.Now(),
        },
    }
    return state, nil
}

// 軌道変更を実行
func (s *OrbitService) ExecuteManeuver(maneuver satellite.Maneuver) (map[string]interface{}, error) {
    currentState, err := s.GetCurrentOrbit(maneuver.SatelliteID)
    if err != nil {
        return nil, err
    }

    deltaV := calculateDeltaV(maneuver.ThrustVector, maneuver.Duration)
    newVelocity := currentState.Velocity.Add(deltaV)
    fuelConsumed := calculateFuelConsumption(maneuver.ThrustVector.Magnitude(), maneuver.Duration)
    newElements := calculateNewOrbitElements(currentState.Position, newVelocity)

    result := map[string]interface{}{
        "success":        true,
        "maneuver_id":    generateManeuverID(),
        "new_velocity":   newVelocity,
        "new_altitude":   calculateAltitude(currentState.Position),
        "fuel_consumed":  fuelConsumed,
        "fuel_remaining": 100.0 - fuelConsumed,
        "delta_v":        deltaV.Magnitude() * 1000,
        "orbit_elements": newElements,
    }
    return result, nil
}

// ΔV 計算
func calculateDeltaV(thrustVector common.Vector3D, duration float64) common.Vector3D {
    satelliteMass := 1000.0
    acceleration := thrustVector.Scale(1.0 / satelliteMass)
    deltaV := acceleration.Scale(duration / 1000.0)
    return deltaV
}

// 燃料消費
func calculateFuelConsumption(thrustMagnitude, duration float64) float64 {
    consumptionRate := 0.01
    return thrustMagnitude * duration * consumptionRate
}

// 新しい軌道要素計算
func calculateNewOrbitElements(position, velocity common.Vector3D) satellite.OrbitElements {
    mu := 398600.4418
    r := position.Magnitude()
    v := velocity.Magnitude()
    energy := (v*v)/2 - mu/r
    semiMajorAxis := -mu / (2 * energy)

    return satellite.OrbitElements{
        SemiMajorAxis:   semiMajorAxis,
        Eccentricity:    0.001,
        Inclination:     51.6,
        RAAN:            0.0,
        ArgumentPerigee: 0.0,
        TrueAnomaly:     0.0,
        Epoch:           time.Now(),
    }
}

// 高度計算
func calculateAltitude(position common.Vector3D) float64 {
    earthRadius := 6371.0
    return position.Magnitude() - earthRadius
}

// Maneuver ID 生成
func generateManeuverID() string {
    return fmt.Sprintf("maneuver_%d", time.Now().Unix())
}