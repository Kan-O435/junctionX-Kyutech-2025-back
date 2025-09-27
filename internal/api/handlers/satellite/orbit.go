package satellite

import (
    "net/http"
    "time"
    
    satelliteModel "junctionx2025back/internal/models/satellite"
    "junctionx2025back/internal/models/common"
    
    "github.com/gin-gonic/gin"
)

// GET /api/v1/satellite/:id/orbit
func GetOrbit(c *gin.Context) {
    satelliteID := c.Param("id")
    
    // サンプルデータ
    response := map[string]interface{}{
        "satellite_id": satelliteID,
        "timestamp":    time.Now(),
        "position": common.Vector3D{
            X: 6800.0,
            Y: 0.0,
            Z: 0.0,
        },
        "velocity": common.Vector3D{
            X: 0.0,
            Y: 7.66,
            Z: 0.0,
        },
        "altitude":     429.0,
        "orbital_speed": 7.66,
    }
    
    c.JSON(http.StatusOK, response)
}

// POST /api/v1/satellite/:id/maneuver  
func ExecuteManeuver(c *gin.Context) {
    satelliteID := c.Param("id")
    
    var request struct {
        PlayerID     string          `json:"player_id"`
        ThrustVector common.Vector3D `json:"thrust_vector"`
        Duration     float64         `json:"duration"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // サンプル応答
    result := map[string]interface{}{
        "success":        true,
        "satellite_id":   satelliteID,
        "maneuver_id":    "maneuver_001",
        "fuel_consumed":  5.2,
        "new_altitude":   435.0,
        "delta_v":        request.ThrustVector.Magnitude() * request.Duration / 100,
    }
    
    c.JSON(http.StatusOK, result)
}

func GetStatus(c *gin.Context) {
    satelliteID := c.Param("id")
    
    status := satelliteModel.SatelliteState{
        Position: common.Vector3D{X: 6800.0, Y: 0.0, Z: 0.0},
        Velocity: common.Vector3D{X: 0.0, Y: 7.66, Z: 0.0},
        Fuel:     85.5,
        Health:   "healthy",
    }
    
    c.JSON(http.StatusOK, gin.H{
        "satellite_id": satelliteID,
        "status":       status,
    })
}