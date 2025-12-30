package model

import "time"

type Entry struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"` // file vs dir
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
}
