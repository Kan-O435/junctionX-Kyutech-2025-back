package satellite

import (
	"time"

	"junctionx2025back/internal/models/common"
)

// 衛星メタ情報（リアルタイム映像用）
type SatelliteInfo struct {
	ID              string
	Name            string
	Type            string
	Resolution      float64
	UpdateInterval  time.Duration
	Coverage        string
	Status          string
	Position        common.Vector3D
	Capabilities    []string
}

// 位置・ズーム等の要求
type Location struct {
	Latitude  float64
	Longitude float64
	Zoom      int
}

type VideoRequest struct {
	Latitude           float64
	Longitude          float64
	Zoom               int
	RequiredResolution float64
	SatelliteID        string
}

type SpectralBand struct {
	Name       string
	Wavelength string
	Purpose    string
}

type VideoSize struct {
	Width  int
	Height int
}

type VideoData struct {
	VideoURL     string
	ThumbnailURL string
	StreamURL    string
	Format       string
	Codec        string
	Bitrate      string
	FrameRate    int
	Duration     int
	Size         VideoSize
	Bands        []SpectralBand
}

type QualityMetrics struct {
	OverallQuality     float64
	CloudCoverage      float64
	AtmosphericClarity float64
	SunAngle           float64
	SignalStrength     float64
	ViewingAngle       float64
}

type VideoResponse struct {
	VideoID       string
	SatelliteID   string
	SatelliteName string
	Location      Location
	VideoData     VideoData
	CaptureTime   time.Time
	Resolution    float64
	Quality       QualityMetrics
	NextUpdate    time.Time
	Status        string
}

type StreamRequest struct {
	SatelliteID string
}

type StreamQuality struct {
	Resolution string
	Bitrate    string
	FrameRate  int
	Format     string
}

type StreamResponse struct {
	StreamID         string
	StreamURL        string
	VideoURL         string
	Status           string
	StartTime        time.Time
	ExpectedDuration time.Duration
	Quality          StreamQuality
}

type VideoRecord struct {
	Timestamp    time.Time
	VideoURL     string
	ThumbnailURL string
	SatelliteID  string
	Quality      float64
}

type SatelliteView struct {
	SatelliteID  string
	VideoURL     string
	ThumbnailURL string
	Resolution   float64
	UpdateTime   time.Time
	Status       string
}

type MultiViewRequest struct {
	SatelliteIDs []string
	Latitude     float64
	Longitude    float64
	Zoom         int
}

type MultiViewResponse struct {
	Location    Location
	Views       []SatelliteView
	SyncTime    time.Time
	TotalViews  int
	Status      string
}


