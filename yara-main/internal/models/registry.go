package models

import (
	"time"
)

// PolicyInfo 策略信息
type PolicyInfo struct {
	Path       string                 `json:"path"`
	Data       map[string]interface{} `json:"data"`
	Metadata   map[string]interface{} `json:"metadata"`
	UpdateTime time.Time              `json:"update_time"`
} 