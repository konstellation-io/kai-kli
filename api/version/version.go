package version

import "time"

type Version struct {
	Tag          string    `json:"tag"`
	CreationDate time.Time `json:"creationDate"`
	Status       string    `json:"status"`
}
