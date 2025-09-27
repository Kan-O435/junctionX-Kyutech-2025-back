package llm

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
		"log"
    "net/http"
    "time"
)

type GeminiService struct {
    apiKey     string
    model      string
    httpClient *http.Client
		baseURL    string
}

type GeminiRequest struct {
    Contents         []GeminiContent   `json:"contents"`
    GenerationConfig GenerationConfig `json:"generationConfig"`
    SafetySettings   []SafetySetting  `json:"safetySettings"`
}

type GeminiContent struct {
    Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
    Text string `json:"text"`
}

type GenerationConfig struct {
    Temperature     float64 `json:"temperature"`
    TopK           int     `json:"topK"`
    TopP           float64 `json:"topP"`
    MaxOutputTokens int    `json:"maxOutputTokens"`
}

type SafetySetting struct {
    Category  string `json:"category"`
    Threshold string `json:"threshold"`
}

type GeminiResponse struct {
    Candidates []struct {
        Content struct {
            Parts []struct {
                Text string `json:"text"`
            } `json:"parts"`
        } `json:"content"`
        FinishReason string `json:"finishReason"`
    } `json:"candidates"`
}

// internal/services/llm/gemini_service.go
func NewGeminiService(apiKey string) *GeminiService {
    return &GeminiService{
        apiKey:     apiKey,
        httpClient: &http.Client{Timeout: 30 * time.Second},
        baseURL:    "https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent",
    }
}

func (gs *GeminiService) CallAPI(prompt string) string {
    url := fmt.Sprintf("%s?key=%s", gs.baseURL, gs.apiKey)
    
    // デバッグ情報を追加
    log.Printf("=== Gemini API Debug ===")
    log.Printf("URL: %s", url)
    log.Printf("Prompt length: %d", len(prompt))
    log.Printf("Prompt preview: %.100s", prompt)
    
    request := GeminiRequest{
        Contents: []GeminiContent{
            {
                Parts: []GeminiPart{
                    {Text: prompt},
                },
            },
        },
        GenerationConfig: GenerationConfig{
            Temperature:     0.7,
            TopK:           30,
            TopP:           0.7,
            MaxOutputTokens: 750,
        },
        SafetySettings: []SafetySetting{
            {Category: "HARM_CATEGORY_HARASSMENT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
            {Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
            {Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
            {Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
        },
    }
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        log.Printf("JSON marshal error: %v", err)
        return "Ground Control to Field Team. Data encoding error. Stand by."
    }
    
    // リクエストボディをログ出力（デバッグ用）
    log.Printf("Request JSON: %s", string(jsonData))
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Request creation error: %v", err)
        return "Ground Control to Field Team. Request creation error. Stand by."
    }
    req.Header.Set("Content-Type", "application/json")
    
    log.Printf("Calling Gemini API...")
    resp, err := gs.httpClient.Do(req)
    if err != nil {
        log.Printf("HTTP request error: %v", err)
        return "Ground Control to Field Team. API communication failure. Stand by."
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Response read error: %v", err)
        return "Ground Control to Field Team. Response read error. Stand by."
    }
    
    // レスポンス詳細をログ出力
    log.Printf("Gemini API status: %d", resp.StatusCode)
    log.Printf("Response body: %s", string(body))
    
    if resp.StatusCode != 200 {
        log.Printf("Gemini API error response: %s", string(body))
        return "Ground Control to Field Team. API status error. Stand by."
    }
    
    var geminiResp GeminiResponse
    if err := json.Unmarshal(body, &geminiResp); err != nil {
        log.Printf("JSON unmarshal error: %v", err)
        log.Printf("Raw response: %s", string(body))
        return "Ground Control to Field Team. Response parsing error. Stand by."
    }
    
    // レスポンス構造をログ出力
    log.Printf("Candidates count: %d", len(geminiResp.Candidates))
    if len(geminiResp.Candidates) > 0 {
        log.Printf("Parts count: %d", len(geminiResp.Candidates[0].Content.Parts))
        log.Printf("Finish reason: %s", geminiResp.Candidates[0].FinishReason)
    }
    
    if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
        response := geminiResp.Candidates[0].Content.Parts[0].Text
        log.Printf("Gemini API success: %s", response)
        return response
    }
    
    log.Printf("No valid response from Gemini API")
    return "Ground Control to Field Team. No response from system. Stand by."
}