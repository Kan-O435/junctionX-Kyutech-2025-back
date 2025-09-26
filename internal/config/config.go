package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port        string
    Environment string
    DatabaseURL string
    RedisURL    string
    JWTSecret   string
    
    // Space-Track API
    SpaceTrackUsername string
    SpaceTrackPassword string
    
    // Physics settings
    PhysicsTickRate     int
    OrbitPropagationStep int
}

func Load() *Config {
    return &Config{
        Port:        getEnv("PORT", "8080"),
        Environment: getEnv("GIN_MODE", "debug"),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/satellite_game?sslmode=disable"),
        RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
        JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
        
        SpaceTrackUsername: getEnv("SPACETRACK_USERNAME", ""),
        SpaceTrackPassword: getEnv("SPACETRACK_PASSWORD", ""),
        
        PhysicsTickRate:     getEnvInt("PHYSICS_TICK_RATE", 10),
        OrbitPropagationStep: getEnvInt("ORBIT_PROPAGATION_STEP", 60),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}