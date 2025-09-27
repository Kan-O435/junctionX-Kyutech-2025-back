// internal/api/handlers/satellite/video/video_handlers.go
package video

import (
    "net/http"
    "strconv"
    "time"
    
    "github.com/gin-gonic/gin"
)

// 利用可能な衛星一覧
func GetAvailableSatellites(c *gin.Context) {
    satellites := []gin.H{
        {
            "id": "himawari8",
            "name": "Himawari-8",
            "type": "Geostationary Weather",
            "resolution": 1000.0,
            "update_interval": "10m",
            "coverage": "Asia-Pacific",
            "status": "active",
            "capabilities": []string{"visible", "infrared", "water_vapor", "realtime"},
        },
        {
            "id": "goes16",
            "name": "GOES-16",
            "type": "Geostationary Weather", 
            "resolution": 500.0,
            "update_interval": "15m",
            "coverage": "Americas",
            "status": "active",
            "capabilities": []string{"visible", "infrared", "lightning", "realtime"},
        },
        {
            "id": "terra",
            "name": "Terra",
            "type": "Earth Observation",
            "resolution": 250.0,
            "update_interval": "1h",
            "coverage": "Global",
            "status": "active",
            "capabilities": []string{"visible", "infrared", "thermal", "multispectral"},
        },
        {
            "id": "landsat8",
            "name": "Landsat 8",
            "type": "Earth Observation",
            "resolution": 30.0,
            "update_interval": "24h",
            "coverage": "Global",
            "status": "active", 
            "capabilities": []string{"visible", "infrared", "thermal", "high_resolution"},
        },
        {
            "id": "worldview3",
            "name": "WorldView-3",
            "type": "Commercial High-Resolution",
            "resolution": 0.31,
            "update_interval": "48h",
            "coverage": "On-demand",
            "status": "active",
            "capabilities": []string{"visible", "infrared", "ultra_high_resolution"},
        },
    }
    
    c.JSON(http.StatusOK, gin.H{
        "satellites": satellites,
        "total": len(satellites),
        "message": "Available satellites for video streaming",
    })
}

// リアルタイム映像取得
func GetRealtimeVideo(c *gin.Context) {
    // クエリパラメータ取得
    latStr := c.Query("latitude")
    lonStr := c.Query("longitude")
    zoomStr := c.Query("zoom")
    resolutionStr := c.Query("required_resolution")
    preferSatellite := c.Query("prefer_satellite")
    
    // パラメータ検証
    lat, err := strconv.ParseFloat(latStr, 64)
    if err != nil || lat < -90 || lat > 90 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid latitude. Must be between -90 and 90",
            "received": latStr,
        })
        return
    }
    
    lon, err := strconv.ParseFloat(lonStr, 64)
    if err != nil || lon < -180 || lon > 180 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid longitude. Must be between -180 and 180", 
            "received": lonStr,
        })
        return
    }
    
    zoom, err := strconv.Atoi(zoomStr)
    if err != nil || zoom < 1 || zoom > 20 {
        zoom = 10 // デフォルト値
    }
    
    var requiredResolution float64
    if resolutionStr != "" {
        requiredResolution, _ = strconv.ParseFloat(resolutionStr, 64)
    }
    
    // 最適な衛星を選択
    selectedSatellite := selectBestSatellite(lat, lon, requiredResolution, preferSatellite)
    
    // レスポンス作成
    videoResponse := gin.H{
        "video_id": "video_" + selectedSatellite["id"].(string) + "_" + strconv.FormatInt(time.Now().Unix(), 10),
        "satellite_id": selectedSatellite["id"],
        "satellite_name": selectedSatellite["name"],
        "location": gin.H{
            "latitude": lat,
            "longitude": lon,
            "zoom": zoom,
        },
        "video_data": gin.H{
            "video_url": "/api/v1/satellite/" + selectedSatellite["id"].(string) + "/video/stream?lat=" + latStr + "&lon=" + lonStr + "&zoom=" + zoomStr,
            "thumbnail_url": "/api/v1/satellite/" + selectedSatellite["id"].(string) + "/video/thumb?lat=" + latStr + "&lon=" + lonStr,
            "stream_url": "wss://api.example.com/satellite/" + selectedSatellite["id"].(string) + "/stream?lat=" + latStr + "&lon=" + lonStr,
            "format": "mp4",
            "codec": "h264",
            "bitrate": "2000kbps",
            "frame_rate": 30,
            "duration": 0, // ライブストリーム
            "size": gin.H{
                "width": 1920,
                "height": 1080,
            },
        },
        "capture_time": time.Now().Format(time.RFC3339),
        "resolution": selectedSatellite["resolution"],
        "quality": gin.H{
            "overall_quality": 0.85,
            "cloud_coverage": calculateCloudCoverage(lat, lon),
            "atmospheric_clarity": 0.9,
            "sun_angle": calculateSunAngle(lat, lon),
            "signal_strength": 0.95,
            "viewing_angle": 0.8,
        },
        "next_update": time.Now().Add(15 * time.Minute).Format(time.RFC3339),
        "status": "streaming",
    }
    
    c.JSON(http.StatusOK, videoResponse)
}

// 映像履歴取得
func GetVideoHistory(c *gin.Context) {
    latStr := c.Query("latitude")
    lonStr := c.Query("longitude")
    hoursStr := c.Query("hours")
    
    lat, err := strconv.ParseFloat(latStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
        return
    }
    
    lon, err := strconv.ParseFloat(lonStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
        return
    }
    
    hours, err := strconv.Atoi(hoursStr)
    if err != nil || hours < 1 {
        hours = 24 // デフォルト24時間
    }
    
    var records []gin.H
    for i := 0; i < hours; i++ {
        pastTime := time.Now().Add(-time.Hour * time.Duration(i))
        
        record := gin.H{
            "timestamp": pastTime.Format(time.RFC3339),
            "video_url": "/api/v1/satellite/history/video?lat=" + latStr + "&lon=" + lonStr + "&time=" + strconv.FormatInt(pastTime.Unix(), 10),
            "thumbnail_url": "/api/v1/satellite/history/thumb?lat=" + latStr + "&lon=" + lonStr + "&time=" + strconv.FormatInt(pastTime.Unix(), 10),
            "satellite_id": "himawari8",
            "quality": 0.8 - float64(i)*0.02, // 時間が経つほど品質低下
        }
        
        records = append(records, record)
    }
    
    c.JSON(http.StatusOK, gin.H{
        "location": gin.H{
            "latitude": lat,
            "longitude": lon,
        },
        "records": records,
        "total_records": len(records),
        "time_range_hours": hours,
        "message": "Video history retrieved successfully",
    })
}

// 複数衛星同時観測
func GetMultiSatelliteView(c *gin.Context) {
    var request struct {
        Latitude     float64  `json:"latitude" binding:"required"`
        Longitude    float64  `json:"longitude" binding:"required"`
        Zoom         int      `json:"zoom"`
        SatelliteIDs []string `json:"satellite_ids" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    var views []gin.H
    for _, satelliteID := range request.SatelliteIDs {
        view := gin.H{
            "satellite_id": satelliteID,
            "video_url": "/api/v1/satellite/" + satelliteID + "/video?lat=" + strconv.FormatFloat(request.Latitude, 'f', 4, 64) + "&lon=" + strconv.FormatFloat(request.Longitude, 'f', 4, 64),
            "thumbnail_url": "/api/v1/satellite/" + satelliteID + "/thumb?lat=" + strconv.FormatFloat(request.Latitude, 'f', 4, 64) + "&lon=" + strconv.FormatFloat(request.Longitude, 'f', 4, 64),
            "resolution": getSatelliteResolution(satelliteID),
            "update_time": time.Now().Format(time.RFC3339),
            "status": "available",
        }
        views = append(views, view)
    }
    
    c.JSON(http.StatusOK, gin.H{
        "location": gin.H{
            "latitude": request.Latitude,
            "longitude": request.Longitude,
            "zoom": request.Zoom,
        },
        "views": views,
        "sync_time": time.Now().Format(time.RFC3339),
        "total_views": len(views),
        "status": "synchronized",
    })
}

// ライブストリーミング開始
func StartLiveStream(c *gin.Context) {
    var request struct {
        SatelliteID      string  `json:"satellite_id" binding:"required"`
        Latitude         float64 `json:"latitude" binding:"required"`
        Longitude        float64 `json:"longitude" binding:"required"`
        DurationMinutes  int     `json:"duration_minutes"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if request.DurationMinutes <= 0 {
        request.DurationMinutes = 120 // デフォルト2時間
    }
    
    streamID := "stream_" + request.SatelliteID + "_" + strconv.FormatInt(time.Now().Unix(), 10)
    
    c.JSON(http.StatusOK, gin.H{
        "stream_id": streamID,
        "stream_url": "wss://api.example.com/satellite/stream/" + streamID,
        "video_url": "https://api.example.com/satellite/video/" + streamID + ".m3u8",
        "status": "starting",
        "start_time": time.Now().Format(time.RFC3339),
        "expected_duration_minutes": request.DurationMinutes,
        "quality": gin.H{
            "resolution": "1920x1080",
            "bitrate": "2000kbps",
            "frame_rate": 30,
            "format": "HLS",
        },
        "satellite": gin.H{
            "id": request.SatelliteID,
            "name": getSatelliteName(request.SatelliteID),
        },
    })
}

// ライブストリーミング停止
func StopLiveStream(c *gin.Context) {
    var request struct {
        StreamID string `json:"stream_id" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "stream_id": request.StreamID,
        "status": "stopped",
        "stop_time": time.Now().Format(time.RFC3339),
        "message": "Live stream stopped successfully",
    })
}

// 衛星詳細情報
func GetSatelliteInfo(c *gin.Context) {
    satelliteID := c.Param("id")
    
    satelliteInfo := getSatelliteDetails(satelliteID)
    if satelliteInfo == nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Satellite not found",
            "satellite_id": satelliteID,
        })
        return
    }
    
    c.JSON(http.StatusOK, satelliteInfo)
}

// 衛星カバレッジ情報
func GetSatelliteCoverage(c *gin.Context) {
    satelliteID := c.Param("id")
    
    c.JSON(http.StatusOK, gin.H{
        "satellite_id": satelliteID,
        "coverage_area": gin.H{
            "type": "circle",
            "center": gin.H{
                "latitude": 0,
                "longitude": 140.7,
            },
            "radius_km": 10000,
        },
        "current_position": gin.H{
            "latitude": 0,
            "longitude": 140.7,
            "altitude_km": 35786,
        },
        "next_pass": time.Now().Add(90 * time.Minute).Format(time.RFC3339),
        "visibility": "excellent",
    })
}

// 指定地点の観測可能衛星
func GetLocationCoverage(c *gin.Context) {
    latStr := c.Query("latitude")
    lonStr := c.Query("longitude")
    
    lat, err := strconv.ParseFloat(latStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
        return
    }
    
    lon, err := strconv.ParseFloat(lonStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
        return
    }
    
    availableSatellites := []gin.H{
        {
            "satellite_id": "himawari8",
            "satellite_name": "Himawari-8",
            "visibility": "excellent",
            "elevation_angle": 45.2,
            "next_pass": time.Now().Add(15 * time.Minute).Format(time.RFC3339),
            "resolution_at_location": 1000.0,
        },
        {
            "satellite_id": "terra",
            "satellite_name": "Terra",
            "visibility": "good",
            "elevation_angle": 32.1,
            "next_pass": time.Now().Add(75 * time.Minute).Format(time.RFC3339),
            "resolution_at_location": 250.0,
        },
    }
    
    c.JSON(http.StatusOK, gin.H{
        "location": gin.H{
            "latitude": lat,
            "longitude": lon,
        },
        "available_satellites": availableSatellites,
        "total_satellites": len(availableSatellites),
        "best_satellite": "himawari8",
        "updated_at": time.Now().Format(time.RFC3339),
    })
}

// 災害地域の映像
func GetDisasterVideo(c *gin.Context) {
    disasterID := c.Param("id")
    
    c.JSON(http.StatusOK, gin.H{
        "disaster_id": disasterID,
        "video_streams": []gin.H{
            {
                "satellite_id": "himawari8",
                "video_url": "/api/v1/satellite/himawari8/disaster/" + disasterID,
                "type": "infrared",
                "last_update": time.Now().Format(time.RFC3339),
                "quality": 0.95,
            },
            {
                "satellite_id": "terra",
                "video_url": "/api/v1/satellite/terra/disaster/" + disasterID,
                "type": "visible",
                "last_update": time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
                "quality": 0.88,
            },
        },
        "disaster_info": gin.H{
            "type": "earthquake",
            "severity": "high",
            "location": gin.H{
                "latitude": 35.6762,
                "longitude": 139.6503,
            },
        },
        "message": "Real-time satellite video of disaster area",
    })
}

// ===== ヘルパー関数 =====

func selectBestSatellite(lat, lon, requiredResolution float64, preferSatellite string) gin.H {
    // 簡易的な衛星選択ロジック
    if preferSatellite != "" {
        return getSatelliteDetails(preferSatellite).(gin.H)
    }
    
    // アジア太平洋地域ならHimawari-8
    if lat >= -60 && lat <= 60 && lon >= 80 && lon <= 200 {
        return gin.H{
            "id": "himawari8",
            "name": "Himawari-8",
            "resolution": 1000.0,
        }
    }
    
    // デフォルトはTerra
    return gin.H{
        "id": "terra",
        "name": "Terra",
        "resolution": 250.0,
    }
}

func calculateCloudCoverage(lat, lon float64) float64 {
    // 簡易的な雲量計算
    if lat >= -30 && lat <= 30 {
        return 0.6 // 熱帯
    }
    return 0.3 // その他
}

func calculateSunAngle(lat, lon float64) float64 {
    hour := float64(time.Now().Hour())
    localNoon := 12.0 + (lon / 15.0)
    timeDiff := hour - localNoon
    if timeDiff < 0 {
        timeDiff = -timeDiff
    }
    return 1.0 - (timeDiff / 12.0)
}

func getSatelliteResolution(satelliteID string) float64 {
    resolutions := map[string]float64{
        "himawari8": 1000.0,
        "goes16": 500.0,
        "terra": 250.0,
        "landsat8": 30.0,
        "worldview3": 0.31,
    }
    if res, exists := resolutions[satelliteID]; exists {
        return res
    }
    return 1000.0
}

func getSatelliteName(satelliteID string) string {
    names := map[string]string{
        "himawari8": "Himawari-8",
        "goes16": "GOES-16",
        "terra": "Terra",
        "landsat8": "Landsat 8",
        "worldview3": "WorldView-3",
    }
    if name, exists := names[satelliteID]; exists {
        return name
    }
    return "Unknown Satellite"
}

func getSatelliteDetails(satelliteID string) interface{} {
    satellites := map[string]gin.H{
        "himawari8": {
            "id": "himawari8",
            "name": "Himawari-8",
            "type": "Geostationary Weather",
            "resolution": 1000.0,
            "status": "operational",
            "orbit_type": "geostationary",
            "altitude_km": 35786,
            "launch_date": "2014-10-07",
            "next_pass": time.Now().Add(10 * time.Minute).Format(time.RFC3339),
        },
        "terra": {
            "id": "terra",
            "name": "Terra",
            "type": "Earth Observation",
            "resolution": 250.0,
            "status": "operational",
            "orbit_type": "polar",
            "altitude_km": 705,
            "launch_date": "1999-12-18",
            "next_pass": time.Now().Add(90 * time.Minute).Format(time.RFC3339),
        },
        "goes16": {
            "id": "goes16",
            "name": "GOES-16",
            "type": "Geostationary Weather",
            "resolution": 500.0,
            "status": "operational",
            "orbit_type": "geostationary",
            "altitude_km": 35786,
            "launch_date": "2016-11-19",
            "next_pass": time.Now().Add(15 * time.Minute).Format(time.RFC3339),
        },
        "landsat8": {
            "id": "landsat8",
            "name": "Landsat 8",
            "type": "Earth Observation",
            "resolution": 30.0,
            "status": "operational",
            "orbit_type": "polar",
            "altitude_km": 705,
            "launch_date": "2013-02-11",
            "next_pass": time.Now().Add(16 * time.Hour).Format(time.RFC3339),
        },
        "worldview3": {
            "id": "worldview3",
            "name": "WorldView-3",
            "type": "Commercial High-Resolution",
            "resolution": 0.31,
            "status": "operational",
            "orbit_type": "polar",
            "altitude_km": 617,
            "launch_date": "2014-08-13",
            "next_pass": time.Now().Add(2 * time.Hour).Format(time.RFC3339),
        },
    }
    
    if satellite, exists := satellites[satelliteID]; exists {
        return satellite
    }
    return nil
}