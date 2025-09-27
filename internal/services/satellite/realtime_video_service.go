// internal/services/satellite/realtime_video_service.go
package satellite

import (
    "context"
    "fmt"
    "time"
    
    "junctionx2025back/internal/models/satellite"
    "junctionx2025back/internal/models/common"
)

type RealtimeVideoService struct {
    satellites map[string]*satellite.SatelliteInfo
}

func NewRealtimeVideoService() *RealtimeVideoService {
    return &RealtimeVideoService{
        satellites: make(map[string]*satellite.SatelliteInfo),
    }
}

// 利用可能な衛星一覧取得
func (s *RealtimeVideoService) GetAvailableSatellites(ctx context.Context) ([]satellite.SatelliteInfo, error) {
    satellites := []satellite.SatelliteInfo{
        // 静止気象衛星
        {
            ID:           "himawari8",
            Name:         "Himawari-8",
            Type:         "Geostationary Weather",
            Resolution:   1000.0, // 1km
            UpdateInterval: time.Minute * 10,
            Coverage:     "Asia-Pacific",
            Status:       "active",
            Position: common.Vector3D{
                X: 35786, // 静止軌道高度
                Y: 0,
                Z: 0,
            },
            Capabilities: []string{"visible", "infrared", "water_vapor", "realtime"},
        },
        {
            ID:           "goes16",
            Name:         "GOES-16",
            Type:         "Geostationary Weather",
            Resolution:   500.0,
            UpdateInterval: time.Minute * 15,
            Coverage:     "Americas",
            Status:       "active",
            Position: common.Vector3D{
                X: 35786,
                Y: -7500,
                Z: 0,
            },
            Capabilities: []string{"visible", "infrared", "lightning", "realtime"},
        },
        {
            ID:           "goes17",
            Name:         "GOES-17",
            Type:         "Geostationary Weather",
            Resolution:   500.0,
            UpdateInterval: time.Minute * 15,
            Coverage:     "Pacific",
            Status:       "active",
            Position: common.Vector3D{
                X: 35786,
                Y: 12000,
                Z: 0,
            },
            Capabilities: []string{"visible", "infrared", "realtime"},
        },
        // 極軌道衛星
        {
            ID:           "terra",
            Name:         "Terra",
            Type:         "Earth Observation",
            Resolution:   250.0,
            UpdateInterval: time.Hour * 1,
            Coverage:     "Global",
            Status:       "active",
            Position: common.Vector3D{
                X: 705,
                Y: 0,
                Z: 0,
            },
            Capabilities: []string{"visible", "infrared", "thermal", "multispectral"},
        },
        {
            ID:           "aqua",
            Name:         "Aqua",
            Type:         "Earth Observation",
            Resolution:   250.0,
            UpdateInterval: time.Hour * 1,
            Coverage:     "Global",
            Status:       "active",
            Position: common.Vector3D{
                X: 705,
                Y: 1000,
                Z: 0,
            },
            Capabilities: []string{"visible", "infrared", "thermal", "ocean_color"},
        },
        {
            ID:           "landsat8",
            Name:         "Landsat 8",
            Type:         "Earth Observation",
            Resolution:   30.0,
            UpdateInterval: time.Hour * 24, // 16日周期
            Coverage:     "Global",
            Status:       "active",
            Position: common.Vector3D{
                X: 705,
                Y: 2000,
                Z: 0,
            },
            Capabilities: []string{"visible", "infrared", "thermal", "multispectral", "high_resolution"},
        },
        {
            ID:           "sentinel2",
            Name:         "Sentinel-2A",
            Type:         "Earth Observation",
            Resolution:   10.0,
            UpdateInterval: time.Hour * 12,
            Coverage:     "Global",
            Status:       "active",
            Position: common.Vector3D{
                X: 786,
                Y: 3000,
                Z: 0,
            },
            Capabilities: []string{"visible", "infrared", "multispectral", "high_resolution"},
        },
        // 商用高解像度衛星
        {
            ID:           "worldview3",
            Name:         "WorldView-3",
            Type:         "Commercial High-Resolution",
            Resolution:   0.31, // 31cm
            UpdateInterval: time.Hour * 48,
            Coverage:     "On-demand",
            Status:       "active",
            Position: common.Vector3D{
                X: 617,
                Y: 4000,
                Z: 0,
            },
            Capabilities: []string{"visible", "infrared", "ultra_high_resolution", "on_demand"},
        },
    }
    
    return satellites, nil
}

// 指定位置のリアルタイム映像取得
func (s *RealtimeVideoService) GetRealtimeVideo(ctx context.Context, req satellite.VideoRequest) (*satellite.VideoResponse, error) {
    // 衛星選択
    selectedSatellite, err := s.selectBestSatellite(req)
    if err != nil {
        return nil, err
    }
    
    // 映像データ生成
    videoData := s.generateVideoData(selectedSatellite, req)
    
    response := &satellite.VideoResponse{
        VideoID:     fmt.Sprintf("video_%s_%d", selectedSatellite.ID, time.Now().Unix()),
        SatelliteID: selectedSatellite.ID,
        SatelliteName: selectedSatellite.Name,
        Location: satellite.Location{
            Latitude:  req.Latitude,
            Longitude: req.Longitude,
            Zoom:      req.Zoom,
        },
        VideoData:    videoData,
        CaptureTime:  time.Now(),
        Resolution:   selectedSatellite.Resolution,
        Quality:      s.calculateVideoQuality(selectedSatellite, req),
        NextUpdate:   time.Now().Add(selectedSatellite.UpdateInterval),
        Status:       "streaming",
    }
    
    return response, nil
}

// 最適な衛星選択
func (s *RealtimeVideoService) selectBestSatellite(req satellite.VideoRequest) (*satellite.SatelliteInfo, error) {
    satellites, _ := s.GetAvailableSatellites(context.Background())
    
    var bestSatellite *satellite.SatelliteInfo
    bestScore := 0.0
    
    for _, sat := range satellites {
        score := s.calculateSatelliteScore(sat, req)
        if score > bestScore {
            bestScore = score
            bestSatellite = &sat
        }
    }
    
    if bestSatellite == nil {
        return nil, fmt.Errorf("no suitable satellite found")
    }
    
    return bestSatellite, nil
}

// 衛星スコア計算
func (s *RealtimeVideoService) calculateSatelliteScore(sat satellite.SatelliteInfo, req satellite.VideoRequest) float64 {
    score := 0.0
    
    // 解像度スコア
    if req.RequiredResolution > 0 {
        if sat.Resolution <= req.RequiredResolution {
            score += 30.0
        } else {
            score += 30.0 * (req.RequiredResolution / sat.Resolution)
        }
    } else {
        score += 20.0
    }
    
    // 更新頻度スコア
    updateMinutes := sat.UpdateInterval.Minutes()
    if updateMinutes <= 15 {
        score += 25.0
    } else if updateMinutes <= 60 {
        score += 20.0
    } else {
        score += 10.0
    }
    
    // カバレッジスコア
    if s.isInCoverage(sat, req.Latitude, req.Longitude) {
        score += 25.0
    }
    
    // リアルタイム対応スコア
    for _, capability := range sat.Capabilities {
        if capability == "realtime" {
            score += 20.0
        }
    }
    
    return score
}

// カバレッジ判定
func (s *RealtimeVideoService) isInCoverage(sat satellite.SatelliteInfo, lat, lon float64) bool {
    switch sat.Coverage {
    case "Asia-Pacific":
        return lat >= -60 && lat <= 60 && lon >= 80 && lon <= 200
    case "Americas":
        return lat >= -60 && lat <= 60 && lon >= -180 && lon <= -30
    case "Pacific":
        return lat >= -60 && lat <= 60 && lon >= 120 && lon <= -120
    case "Global", "On-demand":
        return true
    default:
        return true
    }
}

// 映像データ生成
func (s *RealtimeVideoService) generateVideoData(sat *satellite.SatelliteInfo, req satellite.VideoRequest) satellite.VideoData {
    return satellite.VideoData{
        VideoURL:     fmt.Sprintf("/api/v1/satellite/%s/video/stream?lat=%.4f&lon=%.4f&zoom=%d", sat.ID, req.Latitude, req.Longitude, req.Zoom),
        ThumbnailURL: fmt.Sprintf("/api/v1/satellite/%s/video/thumb?lat=%.4f&lon=%.4f", sat.ID, req.Latitude, req.Longitude),
        StreamURL:    fmt.Sprintf("wss://api.example.com/satellite/%s/stream?lat=%.4f&lon=%.4f", sat.ID, req.Latitude, req.Longitude),
        Format:       "mp4",
        Codec:        "h264",
        Bitrate:      "2000kbps",
        FrameRate:    30,
        Duration:     0, // ライブストリーム
        Size: satellite.VideoSize{
            Width:  1920,
            Height: 1080,
        },
        Bands: s.getAvailableBands(sat),
    }
}

// 利用可能なバンド取得
func (s *RealtimeVideoService) getAvailableBands(sat *satellite.SatelliteInfo) []satellite.SpectralBand {
    bands := []satellite.SpectralBand{}
    
    for _, capability := range sat.Capabilities {
        switch capability {
        case "visible":
            bands = append(bands, satellite.SpectralBand{
                Name:       "Visible",
                Wavelength: "0.4-0.7µm",
                Purpose:    "True color imaging",
            })
        case "infrared":
            bands = append(bands, satellite.SpectralBand{
                Name:       "Near Infrared",
                Wavelength: "0.7-1.4µm",
                Purpose:    "Vegetation analysis",
            })
        case "thermal":
            bands = append(bands, satellite.SpectralBand{
                Name:       "Thermal Infrared",
                Wavelength: "8-12µm",
                Purpose:    "Temperature measurement",
            })
        case "water_vapor":
            bands = append(bands, satellite.SpectralBand{
                Name:       "Water Vapor",
                Wavelength: "6.2µm",
                Purpose:    "Atmospheric moisture",
            })
        }
    }
    
    return bands
}

// 映像品質計算
func (s *RealtimeVideoService) calculateVideoQuality(sat *satellite.SatelliteInfo, req satellite.VideoRequest) satellite.QualityMetrics {
    // 雲量をランダムに生成（実際のAPIでは気象データから取得）
    cloudCoverage := s.generateCloudCoverage(req.Latitude, req.Longitude)
    
    // 大気透明度
    atmosphericClarity := 1.0 - (cloudCoverage * 0.3)
    
    // 太陽角度（時刻による）
    sunAngle := s.calculateSunAngle(req.Latitude, req.Longitude, time.Now())
    
    // 総合品質スコア
    overallQuality := (atmosphericClarity + sunAngle) / 2.0
    
    return satellite.QualityMetrics{
        OverallQuality:      overallQuality,
        CloudCoverage:       cloudCoverage,
        AtmosphericClarity:  atmosphericClarity,
        SunAngle:           sunAngle,
        SignalStrength:     0.9, // 衛星の信号強度
        ViewingAngle:       s.calculateViewingAngle(sat, req.Latitude, req.Longitude),
    }
}

// 雲量生成（簡易版）
func (s *RealtimeVideoService) generateCloudCoverage(lat, lon float64) float64 {
    // 緯度による雲量の違いを模擬
    if lat >= -30 && lat <= 30 {
        return 0.6 // 熱帯：雲が多い
    } else if lat >= 30 && lat <= 60 || lat >= -60 && lat <= -30 {
        return 0.3 // 温帯：中程度
    } else {
        return 0.1 // 極地：雲が少ない
    }
}

// 太陽角度計算
func (s *RealtimeVideoService) calculateSunAngle(lat, lon float64, t time.Time) float64 {
    // 簡易的な太陽角度計算
    hour := float64(t.Hour())
    localNoon := 12.0 + (lon / 15.0) // 経度から地方時の正午を計算
    
    timeDiff := hour - localNoon
    if timeDiff < 0 {
        timeDiff = -timeDiff
    }
    
    // 正午に近いほど高いスコア
    return 1.0 - (timeDiff / 12.0)
}

// 視野角計算
func (s *RealtimeVideoService) calculateViewingAngle(sat *satellite.SatelliteInfo, lat, lon float64) float64 {
    // 静止衛星の場合、赤道からの距離で視野角を計算
    if sat.Type == "Geostationary Weather" {
        return 1.0 - (lat*lat)/(60*60) // 赤道付近で最適
    }
    return 0.8 // 極軌道衛星は一定
}

// 地点の映像履歴取得
func (s *RealtimeVideoService) GetVideoHistory(ctx context.Context, lat, lon float64, hours int) ([]satellite.VideoRecord, error) {
    records := []satellite.VideoRecord{}
    
    for i := 0; i < hours; i++ {
        pastTime := time.Now().Add(-time.Hour * time.Duration(i))
        
        record := satellite.VideoRecord{
            Timestamp:   pastTime,
            VideoURL:    fmt.Sprintf("/api/v1/satellite/history/video?lat=%.4f&lon=%.4f&time=%d", lat, lon, pastTime.Unix()),
            ThumbnailURL: fmt.Sprintf("/api/v1/satellite/history/thumb?lat=%.4f&lon=%.4f&time=%d", lat, lon, pastTime.Unix()),
            SatelliteID: "himawari8", // 最も頻繁に更新される衛星
            Quality:     0.8 - float64(i)*0.05, // 時間が経つほど品質低下
        }
        
        records = append(records, record)
    }
    
    return records, nil
}

// 複数衛星での同時観測
func (s *RealtimeVideoService) GetMultiSatelliteView(ctx context.Context, req satellite.MultiViewRequest) (*satellite.MultiViewResponse, error) {
    views := []satellite.SatelliteView{}
    
    for _, satelliteID := range req.SatelliteIDs {
        view := satellite.SatelliteView{
            SatelliteID:   satelliteID,
            VideoURL:     fmt.Sprintf("/api/v1/satellite/%s/video?lat=%.4f&lon=%.4f", satelliteID, req.Latitude, req.Longitude),
            ThumbnailURL: fmt.Sprintf("/api/v1/satellite/%s/thumb?lat=%.4f&lon=%.4f", satelliteID, req.Latitude, req.Longitude),
            Resolution:   s.getSatelliteResolution(satelliteID),
            UpdateTime:   time.Now(),
            Status:       "available",
        }
        views = append(views, view)
    }
    
    return &satellite.MultiViewResponse{
        Location: satellite.Location{
            Latitude:  req.Latitude,
            Longitude: req.Longitude,
            Zoom:      req.Zoom,
        },
        Views:        views,
        SyncTime:     time.Now(),
        TotalViews:   len(views),
        Status:       "synchronized",
    }, nil
}

// 衛星解像度取得
func (s *RealtimeVideoService) getSatelliteResolution(satelliteID string) float64 {
    resolutions := map[string]float64{
        "himawari8":   1000.0,
        "goes16":      500.0,
        "goes17":      500.0,
        "terra":       250.0,
        "aqua":        250.0,
        "landsat8":    30.0,
        "sentinel2":   10.0,
        "worldview3":  0.31,
    }
    
    if res, exists := resolutions[satelliteID]; exists {
        return res
    }
    return 1000.0 // デフォルト
}

// ライブストリーミング開始
func (s *RealtimeVideoService) StartLiveStream(ctx context.Context, req satellite.StreamRequest) (*satellite.StreamResponse, error) {
    streamID := fmt.Sprintf("stream_%s_%d", req.SatelliteID, time.Now().Unix())
    
    return &satellite.StreamResponse{
        StreamID:    streamID,
        StreamURL:   fmt.Sprintf("wss://api.example.com/satellite/stream/%s", streamID),
        VideoURL:    fmt.Sprintf("https://api.example.com/satellite/video/%s.m3u8", streamID),
        Status:      "starting",
        StartTime:   time.Now(),
        ExpectedDuration: time.Hour * 2, // 2時間の連続配信
        Quality: satellite.StreamQuality{
            Resolution: "1920x1080",
            Bitrate:    "2000kbps",
            FrameRate:  30,
            Format:     "HLS",
        },
    }, nil
}