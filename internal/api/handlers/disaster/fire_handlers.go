package disaster

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	FIRMS_API_BASE_URL   = "https://firms.modaps.eosdis.nasa.gov/api/country/csv"
	FIRMS_API_AREA_URL   = "https://firms.modaps.eosdis.nasa.gov/api/area/csv"
	FIRMS_API_HISTORICAL = "https://firms.modaps.eosdis.nasa.gov/archive/csv"
	MAP_KEY              = "3175ece84a9610a7c8a836dd6a7d245b"
)

type FireData struct {
	Latitude   float64 `json:"latitude"`   // 観測された火災地点の緯度（単位: 度）
	Longitude  float64 `json:"longitude"`  // 観測された火災地点の経度（単位: 度）
	Brightness float64 `json:"brightness"` // MODIS が観測した火災の放射輝度値（火災の強度指標）
	Scan       float64 `json:"scan"`       // 観測ピクセルのスキャン方向の解像度（km単位の幅）
	Track      float64 `json:"track"`      // 観測ピクセルのトラック方向の解像度（km単位の高さ）
	AcqDate    string  `json:"acq_date"`   // 観測日（UTC）
	AcqTime    string  `json:"acq_time"`   // 観測時刻（UTC、HHMM形式）
	Satellite  string  `json:"satellite"`  // 観測衛星（例: A = Aqua, T = Terra）
	Confidence int     `json:"confidence"` // 検出の信頼度（0〜100、数値が高いほど確実な火災）
	Version    string  `json:"version"`    // FIRMS データセットのバージョン
	FRP        float64 `json:"frp"`        // Fire Radiative Power (火災放射強度、MW単位)
	DayNight   string  `json:"daynight"`   // 観測が昼間か夜間か（D=昼, N=夜）
}

type FireResponse struct {
	Fires   []FireData `json:"fires"`
	Total   int        `json:"total"`
	Message string     `json:"message"`
}

// GET: api/v1/disaster/fires
// 国コードや座標を基にNASA FIRMS APIから火災データを取得しフィルタリングする
func GetFires(c *gin.Context) {
	country := c.DefaultQuery("country", "USA")
	source := "MODIS_NRT" // 常に MODIS 固定
	zahyou := c.DefaultQuery("zahyou", "-125,24,-66,49")
	dayRange := c.DefaultQuery("dayrange", "1")

	url := fmt.Sprintf("%s/%s/%s/%s/%s", FIRMS_API_AREA_URL, MAP_KEY, source, zahyou, dayRange)
	log.Printf("Request URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fire data", "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "NASA FIRMS API returned error", "status": resp.StatusCode, "message": string(body)})
		return
	}

	fires, err := parseCSVResponse(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV data", "message": err.Error()})
		return
	}

	var filteredFires []FireData
	if country == "JPN" {
		for _, fire := range fires {
			if fire.Latitude >= 24 && fire.Latitude <= 46 && fire.Longitude >= 123 && fire.Longitude <= 146 {
				filteredFires = append(filteredFires, fire)
			}
		}
	} else if country == "USA" {
		for _, fire := range fires {
			if fire.Latitude >= 24 && fire.Latitude <= 49 && fire.Longitude >= -125 && fire.Longitude <= -66 {
				filteredFires = append(filteredFires, fire)
			}
		}
	} else {
		filteredFires = fires
	}

	response := FireResponse{
		Fires:   filteredFires,
		Total:   len(filteredFires),
		Message: fmt.Sprintf("Fire data retrieved for %s (total global fires: %d)", country, len(fires)),
	}
	c.JSON(http.StatusOK, response)
}

// GET: api/v1/disaster/fires/number1
// GetFires関数中のbrightnessが最も高い地点を出力
func GetFiresNumber1(c *gin.Context) {
	country := c.DefaultQuery("country", "USA")
	source := "MODIS_NRT" // 常に MODIS 固定
	zahyou := c.DefaultQuery("zahyou", "-125,24,-66,49")
	dayRange := c.DefaultQuery("dayrange", "1")

	url := fmt.Sprintf("%s/%s/%s/%s/%s", FIRMS_API_AREA_URL, MAP_KEY, source, zahyou, dayRange)
	log.Printf("Request URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fire data", "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "NASA FIRMS API returned error", "status": resp.StatusCode, "message": string(body)})
		return
	}

	fires, err := parseCSVResponse(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV data", "message": err.Error()})
		return
	}

	// データが空の場合は早期リターン
	if len(fires) == 0 {
		response := FireResponse{
			Fires:   []FireData{},
			Total:   0,
			Message: fmt.Sprintf("No fire data found for %s", country),
		}
		c.JSON(http.StatusOK, response)
		return
	}

	var filteredFires []FireData
	if country == "JPN" {
		for _, fire := range fires {
			if fire.Latitude >= 24 && fire.Latitude <= 46 && fire.Longitude >= 123 && fire.Longitude <= 146 {
				filteredFires = append(filteredFires, fire)
			}
		}
	} else if country == "USA" {
		for _, fire := range fires {
			if fire.Latitude >= 24 && fire.Latitude <= 49 && fire.Longitude >= -125 && fire.Longitude <= -66 {
				filteredFires = append(filteredFires, fire)
			}
		}
	} else {
		filteredFires = fires
	}

	// フィルタリング後のデータが空の場合は早期リターン
	if len(filteredFires) == 0 {
		response := FireResponse{
			Fires:   []FireData{},
			Total:   0,
			Message: fmt.Sprintf("No fire data found for %s after filtering", country),
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// brightnessが最も高い火災データのみを抽出
	maxBrightness := filteredFires[0].Brightness
	var highestBrightnessFires []FireData

	// 最高brightnessを見つける
	for _, fire := range filteredFires {
		if fire.Brightness > maxBrightness {
			maxBrightness = fire.Brightness
		}
	}

	// 最高brightnessの火災データのみを抽出
	for _, fire := range filteredFires {
		if fire.Brightness == maxBrightness {
			highestBrightnessFires = append(highestBrightnessFires, fire)
		}
	}

	response := FireResponse{
		Fires: highestBrightnessFires,
		Total: len(highestBrightnessFires),
		Message: fmt.Sprintf("Highest brightness fire data retrieved for %s (brightness: %.2f, total fires before filtering: %d)",
			country, maxBrightness, len(fires)),
	}
	c.JSON(http.StatusOK, response)
}

// GET: api/v1/disaster/fires/active
// 指定した信頼度以上の火災データをフィルタリングして返す関数
func GetActiveFires(c *gin.Context) {
	country := c.DefaultQuery("country", "JPN")
	minConfidence := c.DefaultQuery("confidence", "20")
	source := "MODIS_NRT"

	url := fmt.Sprintf("%s/%s/%s/%s/1", FIRMS_API_BASE_URL, MAP_KEY, source, country)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active fire data", "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	allFires, err := parseCSVResponse(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV data", "message": err.Error()})
		return
	}

	minConf, _ := strconv.Atoi(minConfidence)
	activeFires := make([]FireData, 0)
	for _, fire := range allFires {
		if fire.Confidence >= minConf {
			activeFires = append(activeFires, fire)
		}
	}

	response := FireResponse{
		Fires:   activeFires,
		Total:   len(activeFires),
		Message: fmt.Sprintf("Active fires with confidence >= %s for %s", minConfidence, country),
	}
	c.JSON(http.StatusOK, response)
}

// GET: api/v1/disaster/fires/area
// 指定範囲の座標に基づいてNASA FIRMS APIから火災データを取得する
func GetFiresByArea(c *gin.Context) {
	latMin := c.DefaultQuery("lat_min", "24")
	lonMin := c.DefaultQuery("lon_min", "123")
	latMax := c.DefaultQuery("lat_max", "46")
	lonMax := c.DefaultQuery("lon_max", "146")
	source := "MODIS_NRT"
	dayRange := c.DefaultQuery("dayrange", "3")

	areaCoords := fmt.Sprintf("%s,%s,%s,%s", lonMin, latMin, lonMax, latMax)
	url := fmt.Sprintf("%s/%s/%s/%s/%s", FIRMS_API_AREA_URL, MAP_KEY, source, areaCoords, dayRange)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch area fire data", "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	fires, err := parseCSVResponse(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV data", "message": err.Error()})
		return
	}

	response := FireResponse{
		Fires:   fires,
		Total:   len(fires),
		Message: fmt.Sprintf("Fire data retrieved for area (%s,%s) to (%s,%s)", latMin, lonMin, latMax, lonMax),
	}
	c.JSON(http.StatusOK, response)
}

// GET: api/v1/disaster/fires/global
// 世界全体の火災データを取得し日本周辺に限定して返す
func GetGlobalFires(c *gin.Context) {
	source := "MODIS_NRT"
	dayRange := c.DefaultQuery("dayrange", "1")

	url := fmt.Sprintf("%s/%s/%s/world/%s", FIRMS_API_AREA_URL, MAP_KEY, source, dayRange)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch global fire data", "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	fires, err := parseCSVResponse(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV data", "message": err.Error()})
		return
	}

	// 日本周辺だけフィルタリング
	japanFires := make([]FireData, 0)
	for _, fire := range fires {
		if fire.Latitude >= 24 && fire.Latitude <= 46 && fire.Longitude >= 123 && fire.Longitude <= 146 {
			japanFires = append(japanFires, fire)
		}
	}

	response := FireResponse{
		Fires:   japanFires,
		Total:   len(japanFires),
		Message: fmt.Sprintf("Global fire data filtered for Japan region (%d total fires worldwide)", len(fires)),
	}
	c.JSON(http.StatusOK, response)
}

// GET: api/v1/disaster/fires/historical
// 指定した日付範囲に含まれる火災データを抽出する関数
func GetHistoricalFires(c *gin.Context) {
	country := c.DefaultQuery("country", "JPN")
	source := "MODIS_NRT"
	dayRange := c.DefaultQuery("dayrange", "7")

	url := fmt.Sprintf("%s/%s/%s/%s/%s", FIRMS_API_BASE_URL, MAP_KEY, source, country, dayRange)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch historical fire data", "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	fires, err := parseCSVResponse(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV data", "message": err.Error()})
		return
	}

	response := FireResponse{
		Fires:   fires,
		Total:   len(fires),
		Message: fmt.Sprintf("Historical fire data retrieved for %s (last %s days)", country, dayRange),
	}
	c.JSON(http.StatusOK, response)
}

// API から火災データを取得し CSV をパースして返す関数
func parseCSVResponse(body io.Reader) ([]FireData, error) {
	reader := csv.NewReader(body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return []FireData{}, nil
	}

	// デバッグ用: 最初の数行をログ出力
	log.Printf("CSV Header: %v", records[0])
	if len(records) > 1 {
		log.Printf("First data row: %v", records[1])
	}

	startIdx := 0
	if strings.Contains(strings.ToLower(records[0][0]), "latitude") {
		startIdx = 1
	}

	fires := make([]FireData, 0, len(records)-startIdx)
	for i := startIdx; i < len(records); i++ {
		record := records[i]
		fire, err := parseFireRecord(record)
		if err == nil {
			fires = append(fires, fire)
		}
	}
	return fires, nil
}

// API 応答の CSV を FireData 構造体のスライスに変換する関数
func parseFireRecord(record []string) (FireData, error) {
	if len(record) < 13 {
		return FireData{}, fmt.Errorf("incomplete record: %v", record)
	}

	latitude, _ := strconv.ParseFloat(record[0], 64)
	longitude, _ := strconv.ParseFloat(record[1], 64)
	brightness, _ := strconv.ParseFloat(record[2], 64)
	scan, _ := strconv.ParseFloat(record[3], 64)
	track, _ := strconv.ParseFloat(record[4], 64)
	confidence, _ := strconv.Atoi(record[8])
	frp, _ := strconv.ParseFloat(record[10], 64)

	return FireData{
		Latitude:   latitude,
		Longitude:  longitude,
		Brightness: brightness,
		Scan:       scan,
		Track:      track,
		AcqDate:    record[5],
		AcqTime:    record[6],
		Satellite:  record[7],
		Confidence: confidence,
		Version:    record[9],
		FRP:        frp,
		DayNight:   record[11],
	}, nil
}
