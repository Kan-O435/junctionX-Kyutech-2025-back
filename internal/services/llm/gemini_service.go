package llm

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type GeminiService struct {
    apiKey     string
    model      string
    httpClient *http.Client
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

func NewGeminiService(apiKey string) *GeminiService {
    return &GeminiService{
        apiKey:     apiKey,
        model:      "gemini-pro",
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }
}

func (gs *GeminiService) CallAPI(prompt string) string {
    url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", 
        gs.model, gs.apiKey)
    
    request := GeminiRequest{
        Contents: []GeminiContent{
            {
                Parts: []GeminiPart{
                    {Text: prompt},
                },
            },
        },
        GenerationConfig: GenerationConfig{
            Temperature:     0.8,
            TopK:           40,
            TopP:           0.95,
            MaxOutputTokens: 800,
        },
        SafetySettings: []SafetySetting{
            {Category: "HARM_CATEGORY_HARASSMENT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
            {Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
            {Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
            {Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
        },
    }
    
    jsonData, _ := json.Marshal(request)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := gs.httpClient.Do(req)
    if err != nil {
        return "Ground Control to Field Team. API communication failure. Stand by."
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    var geminiResp GeminiResponse
    json.Unmarshal(body, &geminiResp)
    
    if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
        return geminiResp.Candidates[0].Content.Parts[0].Text
    }
    
    return "Ground Control to Field Team. No response from system. Stand by."
}