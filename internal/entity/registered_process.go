package entity

import "time"

type RegisteredProcess struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	Type       string    `json:"type"`
	Image      string    `json:"image"`
	UploadDate time.Time `json:"uploadDate"`
	Owner      string    `json:"owner"`
	Status     string    `json:"status"`
}
