package disaster

import "time"

type DisasterEvent struct {
    ID          string    `json:"id"`
    Type        string    `json:"type"`
    Title       string    `json:"title"`
    Magnitude   float64   `json:"magnitude"`
    Latitude    float64   `json:"latitude"`
    Longitude   float64   `json:"longitude"`
    Location    string    `json:"location"`
    Severity    string    `json:"severity"`
    Time        time.Time `json:"time"`
    Source      string    `json:"source"`
}