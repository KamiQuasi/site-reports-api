package sitereports

import "encoding/json"

// Report struct
type Report struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	StartTime []byte   `json:"starttime"`
	EndTime   []byte   `json:"endtime"`
	Type      string   `json:"type"`
	Results   []Result `json:"results"`
}

// Result struct
type Result struct {
	TimeStamp []byte          `json:"timestamp"`
	Type      string          `json:"type"`
	Notes     string          `json:"notes"`
	Data      json.RawMessage `json:"data"`
}
