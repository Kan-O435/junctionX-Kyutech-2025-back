package routes

import (
    "junctionx2025back/internal/api/handlers/satellite"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Satellite Game Backend is running!",
        })
    })
    
    // API v1 group
    v1 := r.Group("/api/v1")
    {
        // 衛星関連
        satelliteGroup := v1.Group("/satellite")
        {
            satelliteGroup.GET("/:id/orbit", satellite.GetOrbit)
            satelliteGroup.POST("/:id/maneuver", satellite.ExecuteManeuver)
            satelliteGroup.GET("/:id/status", satellite.GetStatus)
        }
        
        // デブリ脅威取得（仮実装）
        v1.GET("/mission/debris/:id/threats", func(c *gin.Context) {
            missionID := c.Param("id")
            c.JSON(200, gin.H{
                "mission_id": missionID,
                "threats": []gin.H{
                    {
                        "id": "debris_001",
                        "name": "Rocket Fragment",
                        "distance": 2.5,
                        "time_to_impact": 300,
                        "danger_level": 7,
                    },
                    {
                        "id": "debris_002", 
                        "name": "Satellite Fragment",
                        "distance": 8.1,
                        "time_to_impact": 450,
                        "danger_level": 4,
                    },
                },
                "message": "Sample debris threats",
            })
        })
    }
}